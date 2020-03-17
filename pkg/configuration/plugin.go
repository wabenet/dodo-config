package configuration

import (
	"fmt"

	"github.com/dodo/dodo-config/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/configuration"
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

func (p *Configuration) GetClientOptions(_ string) (*configuration.ClientOptions, error) {
	return &configuration.ClientOptions{}, nil
}

func (p *Configuration) UpdateConfiguration(backdrop *types.Backdrop) (*types.Backdrop, error) {
	ptr, decode := NewConfig()
	configfiles.GimmeConfigFiles(&configfiles.Options{
		Name:                      "dodo",
		Extensions:                []string{"yaml", "yml", "json"},
		IncludeWorkingDirectories: true,
		Filter: func(configFile *configfiles.ConfigFile) bool {
			s := decoder.New(configFile.Path)
			s.LoadYaml(decode, configFile.Content)
                        return false
		},
	})

	// TODO: wtf this cast
	config := *(ptr.(**Config))
	if result, ok := config.Backdrops[backdrop.Name]; ok {
		return result, nil
	}

	// TODO: this will be ignored by the plugin client
	matches := fuzzy.Find(backdrop.Name, names(config))
	if len(matches) == 0 {
		return nil, fmt.Errorf("could not find any configuration for backdrop '%s'", backdrop.Name)
	}
	return nil, fmt.Errorf("backdrop '%s' not found, did you mean '%s'?", backdrop.Name, matches[0].Str)
}

func names(config *Config) []string {
	var result []string
	if config.Backdrops != nil {
		for name := range config.Backdrops {
			result = append(result, name)
		}
	}
	if config.Groups != nil {
		for _, group := range config.Groups {
			result = append(result, names(group)...)
		}
	}
	return result
}

func (p *Configuration) Provision(_ string) error {
	return nil
}
