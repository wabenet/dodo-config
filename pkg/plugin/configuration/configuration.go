package configuration

import (
	"fmt"

	"github.com/dodo-cli/dodo-config/pkg/config"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	core "github.com/dodo-cli/dodo-core/pkg/config"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/configuration"
)

const name = "config"

var _ configuration.Configuration = &Configuration{}

type Configuration struct {
	backdrops map[string]*api.Backdrop
}

func New() *Configuration {
	return &Configuration{}
}

func (p *Configuration) Type() plugin.Type {
	return configuration.Type
}

func (p *Configuration) PluginInfo() *api.PluginInfo {
	return &api.PluginInfo{
		Name: &api.PluginName{Name: name, Type: configuration.Type.String()},
	}
}

func (p *Configuration) Init() (plugin.PluginConfig, error) {
	return map[string]string{}, nil
}

func (p *Configuration) get() (map[string]*api.Backdrop, error) {
	if p.backdrops != nil {
		return p.backdrops, nil
	}

	backdrops, err := config.GetAllBackdrops(core.GetConfigFiles()...)
	if err != nil {
		return nil, err
	}

	p.backdrops = backdrops

	return p.backdrops, nil
}

func (p *Configuration) GetBackdrop(name string) (*api.Backdrop, error) {
	bs, err := p.get()
	if err != nil {
		return nil, err
	}

	for _, b := range bs {
		if b.Name == name {
			return b, nil
		}

		for _, a := range b.Aliases {
			if a == name {
				return b, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find any configuration for backdrop '%s'", name)
}

func (p *Configuration) ListBackdrops() ([]*api.Backdrop, error) {
	result := []*api.Backdrop{}

	bs, err := p.get()
	if err != nil {
		return nil, err
	}

	for _, b := range bs {
		result = append(result, b)
	}

	return result, nil
}
