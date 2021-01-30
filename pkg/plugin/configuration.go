package plugin

import (
	"fmt"

	api "github.com/dodo-cli/dodo-core/api/v1alpha1"
	"github.com/dodo-cli/dodo-core/pkg/decoder"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/configuration"
	"github.com/dodo-cli/dodo-core/pkg/types"
	log "github.com/hashicorp/go-hclog"
	"github.com/oclaussen/go-gimme/configfiles"
	"github.com/sahilm/fuzzy"
)

var _ configuration.Configuration = &Configuration{}

type Configuration struct{}

func (p *Configuration) Type() plugin.Type {
	return configuration.Type
}

func (p *Configuration) Init() error {
	return nil
}

func (p *Configuration) PluginInfo() (*api.PluginInfo, error) {
	return &api.PluginInfo{Name: "config"}, nil
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

			d := decoder.New(configFile.Path)
			backdrops := map[string]*api.Backdrop{}
			d.DecodeYaml(configFile.Content, &backdrops, map[string]decoder.Decoding{
				"backdrops": decoder.Map(types.NewBackdrop(), &backdrops),
			})

			for name, b := range backdrops {
				b.Name = name // TODO: This shouldn't be necessary. And not happen here
				log.L().Debug("found backdrop", "name", b.Name)
				result = append(result, b)
			}

			return false
		},
	})

	return result, nil
}
