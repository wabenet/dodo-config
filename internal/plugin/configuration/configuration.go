package configuration

import (
	"fmt"

	"github.com/wabenet/dodo-config/internal/config"
	core "github.com/wabenet/dodo-core/pkg/config"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/configuration"
)

const name = "config"

var _ configuration.Configuration = &Configuration{}

type Configuration struct {
	backdrops map[string]configuration.Backdrop
}

func New() *Configuration {
	return &Configuration{}
}

func (p *Configuration) Type() plugin.Type {
	return configuration.Type
}

func (p *Configuration) Metadata() plugin.Metadata {
	return plugin.NewMetadata(configuration.Type, name)
}

func (p *Configuration) Init() (plugin.Config, error) {
	return map[string]string{}, nil
}

func (*Configuration) Cleanup() {}

func (p *Configuration) get() (map[string]configuration.Backdrop, error) {
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

func (p *Configuration) GetBackdrop(name string) (configuration.Backdrop, error) {
	bs, err := p.get()
	if err != nil {
		return configuration.Backdrop{}, err
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

	return configuration.Backdrop{}, fmt.Errorf("could not find any configuration for backdrop '%s'", name)
}

func (p *Configuration) ListBackdrops() ([]configuration.Backdrop, error) {
	result := []configuration.Backdrop{}

	bs, err := p.get()
	if err != nil {
		return nil, err
	}

	for _, b := range bs {
		result = append(result, b)
	}

	return result, nil
}
