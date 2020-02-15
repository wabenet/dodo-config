package configuration

import (
	"github.com/hashicorp/go-plugin"
	"github.com/oclaussen/dodo/pkg/plugin/configuration"
	"github.com/oclaussen/dodo/pkg/types"
)

type Configuration struct{}

func NewPlugin() plugin.Plugin {
	return &configuration.Plugin{Impl: &Configuration{}}
}

func (p *Configuration) GetClientOptions(_ string) (*configuration.ClientOptions, error) {
	return &configuration.ClientOptions{}, nil
}

func (p *Configuration) UpdateConfiguration(backdrop *types.Backdrop) (*types.Backdrop, error) {
	conf, err := LoadBackdrop(backdrop.Name)
	if err != nil {
		return nil, err
	}
	backdrop.Merge(conf)
	return backdrop, nil
}

func (p *Configuration) Provision(_ string) error {
	return nil
}
