package plugin

import (
	"fmt"

	"github.com/oclaussen/dodo/pkg/configuration"
	"github.com/oclaussen/dodo/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/plugin"
	"github.com/oclaussen/dodo/pkg/types"
	"github.com/oclaussen/go-gimme/configfiles"
	"github.com/sahilm/fuzzy"
)

type Configuration struct{}

func RegisterPlugin() {
	plugin.RegisterPluginServer(
		configuration.PluginType,
		&configuration.Plugin{Impl: &Configuration{}},
	)
}

func (p *Configuration) Init() error {
	return nil
}

func (p *Configuration) UpdateConfiguration(backdrop *types.Backdrop) (*types.Backdrop, error) {
	backdrops := map[string]*types.Backdrop{}
	configfiles.GimmeConfigFiles(&configfiles.Options{
		Name:                      "dodo",
		Extensions:                []string{"yaml", "yml", "json"},
		IncludeWorkingDirectories: true,
		Filter: func(configFile *configfiles.ConfigFile) bool {
			d := decoder.New(configFile.Path)
			d.DecodeYaml(configFile.Content, &backdrops, map[string]decoder.Decoding{
				"backdrops": decoder.Map(types.NewBackdrop(), &backdrops),
			})
			return false
		},
	})

	if result, ok := backdrops[backdrop.Name]; ok {
		return result, nil
	}

	var names []string
	for name := range backdrops {
		names = append(names, name)
	}

	matches := fuzzy.Find(backdrop.Name, names)
	if len(matches) == 0 {
		return nil, fmt.Errorf("could not find any configuration for backdrop '%s'", backdrop.Name)
	}
	return nil, fmt.Errorf("backdrop '%s' not found, did you mean '%s'?", backdrop.Name, matches[0].Str)
}

func (p *Configuration) Provision(_ string) error {
	return nil
}
