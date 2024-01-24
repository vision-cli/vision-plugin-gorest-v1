package generate

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/vision-cli/vision-plugin-gorest-v1/cmd/model"
)

const DEFAULT_OUTPUT_DIR = "./services"
const TEMPLATE_DIR = "_template"

//go:embed all:_template
var templateFiles embed.FS

var GenerateCmd = &cobra.Command{
	Use:   "generate <output services folder> <config file>",
	Short: "generate a REST service",
	Long:  "generate a golang REST service in the services directory using the vision.json configuration",
	Args:  cobra.NoArgs,
	RunE:  generateAndCheck,
}

type success struct {
	Success bool `json:"success"`
}

type convertConfig struct {
	PluginConfig model.PluginConfig `json:"config"`
	GoVersion    string
}

type SinglePluginConfig struct {
	PluginName string
	ModuleName string
	Command    string
}

type SingleConvertConfig struct {
	PluginConfig SinglePluginConfig
	GoVersion    string
}

// wraps the run function to determine a success or failed response
func generateAndCheck(cmd *cobra.Command, args []string) error {
	jEnc := json.NewEncoder(os.Stdout)

	err := run(cmd, args)
	if err != nil {
		_ = jEnc.Encode(success{Success: false})
		return fmt.Errorf("generating template: %w", err)
	}

	err = jEnc.Encode(success{Success: true})
	if err != nil {
		return fmt.Errorf("encoding JSON response: %w", err)
	}

	return nil
}

func run(cmd *cobra.Command, args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	outputPath := DEFAULT_OUTPUT_DIR
	vPath := filepath.Join(wd, "vision.json")

	// if outputPath dir does not exist, create dir
	_, err = os.Stat(outputPath)
	if os.IsNotExist(err) {
		os.MkdirAll(outputPath, os.ModePerm)
	} else if err != nil {
		return fmt.Errorf("searching for output dir: %w", err)
	}

	vj, err := openVisionJson(vPath)
	if err != nil {
		return fmt.Errorf("opening vision.json: %w", err)
	}

	for _, m := range vj.PluginConfig.ModuleNames {
		// convert struct to use correct JSON tag
		var convConf SingleConvertConfig
		convConf.PluginConfig.PluginName = vj.PluginConfig.PluginName
		convConf.PluginConfig.ModuleName = m
		convConf.PluginConfig.Command = vj.PluginConfig.Command
		convConf.GoVersion = vj.GoVersion

		parts := strings.Split(m, "/")
		singlePath := filepath.Join(outputPath, filepath.Join(parts[2:]...))
		err = walkDirAndClone(wd, singlePath, &convConf)
		if err != nil {
			return fmt.Errorf("walking dir and cloning: %w", err)
		}
		err = execGoModTidy(singlePath)
		if err != nil {
			return fmt.Errorf("executing go mod tidy: %w", err)
		}
	}
	return nil
}

func execGoModTidy(outputPath string) error {
	if outputPath == "." {
		outputPath = ""
	}
	c := exec.Command("go", "mod", "tidy")
	c.Dir = outputPath
	_, err := c.Output()
	if err != nil {
		return fmt.Errorf("running 'go mod tidy': %w", err)
	}
	return nil
}

func walkDirAndClone(wd, outputPath string, vj *SingleConvertConfig) error {
	return fs.WalkDir(templateFiles, TEMPLATE_DIR, func(path string, d fs.DirEntry, err error) error {
		newPath := filepath.Join(wd, outputPath, strings.TrimPrefix(path, filepath.Join(TEMPLATE_DIR, "/")))

		switch {
		case path == TEMPLATE_DIR: // skip the top level template dir
			return nil
		case d.IsDir(): // if it is a dir then create it
			return cloneDir(newPath)
		case filepath.Ext(newPath) == ".tmpl":
			err := cloneExecTmpl(path, newPath, vj)
			if err != nil {
				return fmt.Errorf("cloning template files: %w", err)
			}

			return nil
		default:
			cloneFile(path, newPath)
			if err != nil {
				return fmt.Errorf("cloning files: %w", err)
			}
			return nil
		}
	})
}

func openVisionJson(vPath string) (*convertConfig, error) {
	f, err := os.OpenFile(vPath, os.O_RDWR, 0444)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading bytes: %w", err)
	}

	var jsonData model.PluginData
	if err = json.Unmarshal(b, &jsonData); err != nil {
		return nil, fmt.Errorf("unmarshalling json: %w", err)
	}

	gv, err := getLatestGoVersion()
	if err != nil {
		return nil, fmt.Errorf("getting latest Go version: %w", err)
	}

	// convert struct to use correct JSON tag
	var convConf convertConfig
	convConf.PluginConfig = jsonData.PluginConfig
	convConf.GoVersion = gv

	return &convConf, nil
}

// if path is a directory, just copy it
func cloneDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// if file isn't template file, just copy it
func cloneFile(src, dst string) error {
	fsrc, err := templateFiles.Open(src)
	if err != nil {
		return fmt.Errorf("opening from templateFiles: %w", err)
	}
	defer fsrc.Close()
	fdst, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if os.IsExist(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("[clone] opening from clone: %w", err)
	}
	defer fdst.Close()
	_, err = io.Copy(fdst, fsrc)
	return err
}

func cloneExecTmpl(src, dst string, vj *SingleConvertConfig) error {
	// open file and read it
	trimmedNewPath := strings.TrimSuffix(dst, filepath.Ext(dst))
	err := cloneFile(src, trimmedNewPath)
	if err != nil {
		return fmt.Errorf("cloning file: %w", err)
	}
	f, err := os.OpenFile(trimmedNewPath, os.O_RDWR, 0444) // only enable reading mode as we do not need to write anything
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}
	err = f.Truncate(0)
	if err != nil {
		return fmt.Errorf("truncating: %w", err)
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("seeking: %w", err)
	}

	tmplEx, err := template.New("templateFile").Parse(string(b))
	if err != nil {
		return fmt.Errorf("creating template file: %w", err)
	}

	return tmplEx.Execute(f, vj)
}

func getLatestGoVersion() (string, error) {
	cmd := "curl 'https://go.dev/VERSION?m=text'"
	b, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", fmt.Errorf("curling Go version: %w", err)
	}

	goVersion := string(b)[2:8]

	return goVersion, nil
}
