package info

import (
	_ "embed"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	api "github.com/vision-cli/api/v1"
)

//go:embed info.txt
var infoOutput string

var InfoCmd = &cobra.Command{
	Use:   "info",
	Short: "return info about go-rest-server-plugin}",
	Long:  "return detailed information about the go-rest-server-plugin plugin",
	Run: func(cmd *cobra.Command, args []string) {
		json.NewEncoder(os.Stdout).Encode(api.Info{
			ShortDescription: "go-rest-server-plugin",
			LongDescription:  infoOutput,
		})
	},
}
