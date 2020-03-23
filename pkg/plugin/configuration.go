package plugin

import (
	"fmt"

	"github.com/dodo/dodo-config/pkg/decoder"
	cfgtypes "github.com/dodo/dodo-config/pkg/types"
	"github.com/oclaussen/dodo/pkg/configuration"
	"github.com/oclaussen/dodo/pkg/types"
	"github.com/oclaussen/go-gimme/configfiles"
	"github.com/sahilm/fuzzy"
)

type Configuration struct{}

func (p *Configuration) GetClientOptions(_ string) (*configuration.ClientOptions, error) {
	return &configuration.ClientOptions{}, nil
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
				"backdrops": decoder.Map(cfgtypes.NewBackdrop(), &backdrops),
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

	// TODO: this will be ignored by the plugin client
	matches := fuzzy.Find(backdrop.Name, names)
	if len(matches) == 0 {
		return nil, fmt.Errorf("could not find any configuration for backdrop '%s'", backdrop.Name)
	}
	return nil, fmt.Errorf("backdrop '%s' not found, did you mean '%s'?", backdrop.Name, matches[0].Str)
}

func (p *Configuration) Provision(_ string) error {
	return nil
}
