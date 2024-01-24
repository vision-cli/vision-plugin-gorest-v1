package model

import (
	"io"
	"os"
)

var Out io.Writer = os.Stdout

const PluginSemVer = "v0.0.1"
const PluginName = "vision-plugin-plugin-v1"
const PluginCommand = "plugin"
const PluginOutputDir = "."

type PluginConfig struct {
	PluginName  string   `json:"plugin_name"`
	ModuleNames []string `json:"module_names"`
	Command     string   `json:"command"`
}

type PluginData struct {
	PluginConfig PluginConfig `json:"gorest"`
}
