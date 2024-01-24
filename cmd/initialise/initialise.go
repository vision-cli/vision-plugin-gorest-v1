package initialise

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	api "github.com/vision-cli/api/v1"
	"github.com/vision-cli/vision-plugin-gorest-v1/cmd/model"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "initialise a project with this plugin",
	Long:  "initialise a project's vision.json file with this plugin's configuration values",
	Args: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return fmt.Errorf("please provide one or more module names in the form github.com/org/domain/subdomain...")
		}

		for _, arg := range args {
			if !strings.HasPrefix(arg, "github.com/") {
				return fmt.Errorf("invalid plugin module name: %s", arg)
			}
		}

		return nil
	},
	RunE: runCommand,
}

func runCommand(cmd *cobra.Command, args []string) error {
	jEnc := json.NewEncoder(model.Out)

	pd := model.PluginConfig{
		PluginName:  "vision-plugin-gorest-v1",
		ModuleNames: args,
		Command:     "gorest",
	}

	err := jEnc.Encode(api.Init{
		Config:  pd,
		Success: true,
	})

	if err != nil {
		return fmt.Errorf("encoding json: %w", err)
	}
	return nil
}
