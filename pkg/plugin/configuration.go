package plugin

import (
	"fmt"
	"reflect"

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

	var names []string
	for name := range config.Backdrops {
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

type Config struct {
	Backdrops map[string]*types.Backdrop
}

func NewConfig() (interface{}, decoder.Decoder) {
	target := &Config{Backdrops: map[string]*types.Backdrop{}}
	return &target, func(s *decoder.Status, config interface{}) {
		DoTheThing(&target.Backdrops, "backdrops", decoder.Map(cfgtypes.NewBackdrop(), &target.Backdrops))(s, config)
	}
}

func DoTheThing(target interface{}, key string, decode decoder.Decoder) decoder.Decoder {
	dummy := []interface{}{}
	return decoder.Keys(map[string]decoder.Decoder{
		key: decode,
		"include": decoder.Kinds(map[reflect.Kind]decoder.Decoder{
			reflect.Map:   decoder.Singleton(NewThing(target, key, decode), &dummy),
			reflect.Slice: decoder.Slice(NewThing(target, key, decode), &dummy),
		}),
	})
}

func NewThing(target interface{}, key string, decode decoder.Decoder) decoder.Producer {
	return func() (interface{}, decoder.Decoder) {
		var dummy struct{}
		return &dummy, decoder.Keys(map[string]decoder.Decoder{
			"text": decoder.IncludeText(DoTheThing(target, key, decode)),
			"file": decoder.IncludeFile(DoTheThing(target, key, decode)),
		})
	}
}
