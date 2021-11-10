package configuration

import (
	"fmt"

	"github.com/dodo-cli/dodo-config/pkg/config"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/configuration"
	"github.com/oclaussen/go-gimme/configfiles"
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

	filenames := []string{}
	configfiles.GimmeConfigFiles(&configfiles.Options{
		Name:                      "dodo",
		Extensions:                []string{"yaml", "yml", "json"},
		IncludeWorkingDirectories: true,
		Filter: func(configFile *configfiles.ConfigFile) bool {
			filenames = append(filenames, configFile.Path)
			return false
		},
	})

	backdrops, err := config.GetAllBackdrops(filenames...)
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
