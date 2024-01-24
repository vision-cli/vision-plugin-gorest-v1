package cmd

import (
	"github.com/spf13/cobra"

	"github.com/vision-cli/vision-plugin-gorest-v1/cmd/generate"
	"github.com/vision-cli/vision-plugin-gorest-v1/cmd/info"
	"github.com/vision-cli/vision-plugin-gorest-v1/cmd/initialise"
	"github.com/vision-cli/vision-plugin-gorest-v1/cmd/version"
)

func init() {
	rootCmd.AddCommand(initialise.InitCmd)
	rootCmd.AddCommand(info.InfoCmd)
	rootCmd.AddCommand(version.VersionCmd)
	rootCmd.AddCommand(generate.GenerateCmd)
}

var rootCmd = &cobra.Command{
	Use:                "gorest",
	FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
}

func Execute() {
	rootCmd.Execute()
}
