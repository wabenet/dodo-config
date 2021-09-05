package plugin

import (
	"fmt"

	"github.com/dodo-cli/dodo-config/pkg/config"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/configuration"
	log "github.com/hashicorp/go-hclog"
	"github.com/oclaussen/go-gimme/configfiles"
	"github.com/sahilm/fuzzy"
)

const name = "config"

var _ configuration.Configuration = &Configuration{}

type Configuration struct{}

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

func (*Configuration) Init() (plugin.PluginConfig, error) {
	return map[string]string{}, nil
}

func (p *Configuration) GetBackdrop(alias string) (*api.Backdrop, error) {
	backdrops, err := p.ListBackdrops()
	if err != nil {
		return nil, err
	}

	if result, err := findBackdrop(backdrops, alias); err == nil {
		return result, nil
	}

	names := []string{}
	for _, b := range backdrops {
		names = append(names, b.Name)
		names = append(names, b.Aliases...)
	}

	matches := fuzzy.Find(alias, names)
	if len(matches) == 0 {
		return nil, fmt.Errorf("could not find any configuration for backdrop '%s'", alias)
	}
	return nil, fmt.Errorf("backdrop '%s' not found, did you mean '%s'?", alias, matches[0].Str)
}

func findBackdrop(backdrops []*api.Backdrop, name string) (*api.Backdrop, error) {
	for _, b := range backdrops {
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
	configfiles.GimmeConfigFiles(&configfiles.Options{
		Name:                      "dodo",
		Extensions:                []string{"yaml", "yml", "json"},
		IncludeWorkingDirectories: true,
		Filter: func(configFile *configfiles.ConfigFile) bool {
			log.L().Debug("checking config file", "path", configFile.Path)

			backdrops, err := config.ParseConfig(configFile.Path)
			if err != nil {
				log.L().Error(err.Error())
			}

			for _, backdrop := range backdrops {
				log.L().Debug("found backdrop", "name", backdrop.Name)
				result = append(result, backdrop)
			}

			return false
		},
	})

	return result, nil
}
