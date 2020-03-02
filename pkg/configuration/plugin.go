package configuration

import (
	"github.com/oclaussen/dodo/pkg/configuration"
	"github.com/oclaussen/dodo/pkg/plugin"
	"github.com/oclaussen/dodo/pkg/types"
)

type Configuration struct{}

func init() {
	plugin.RegisterPluginServer(
		configuration.PluginType,
		&configuration.Plugin{Impl: &Configuration{}},
	)
}

func (p *Configuration) GetClientOptions(_ string) (*configuration.ClientOptions, error) {
	return &configuration.ClientOptions{}, nil
}

func (p *Configuration) UpdateConfiguration(backdrop *types.Backdrop) (*types.Backdrop, error) {
	return LoadBackdrop(backdrop.Name)
}

func (p *Configuration) Provision(_ string) error {
	return nil
}
