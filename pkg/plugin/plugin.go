package plugin

import (
	"github.com/oclaussen/dodo/pkg/command"
	"github.com/oclaussen/dodo/pkg/configuration"
	"github.com/oclaussen/dodo/pkg/plugin"
)

func RegisterPlugin() {
	plugin.RegisterPluginServer(
		command.PluginType,
		&command.Plugin{Impl: &Commands{}},
	)
	plugin.RegisterPluginServer(
		configuration.PluginType,
		&configuration.Plugin{Impl: &Configuration{}},
	)
}
