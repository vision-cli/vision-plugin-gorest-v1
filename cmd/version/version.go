package version

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	api "github.com/vision-cli/api/v1"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "the plugin version",
	Long:  "get the version for the vision-plugin-gorest-v1 plugin",
	Run: func(cmd *cobra.Command, args []string) {
		json.NewEncoder(os.Stdout).Encode(
			api.Version{
				SemVer: "v0.0.1",
			})
	},
}
