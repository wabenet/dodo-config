package main

import (
	"github.com/dodo/dodo-config/pkg/command"
	"github.com/dodo/dodo-config/pkg/configuration"
	"github.com/oclaussen/dodo/pkg/plugin"
)

func main() {
	command.RegisterPlugin()
	configuration.RegisterPlugin()
	plugin.ServePlugins()
}
