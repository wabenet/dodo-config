package config

import (
	"cuelang.org/go/cue"
	"github.com/dodo-cli/dodo-config/pkg/cuetils"
	"github.com/dodo-cli/dodo-config/pkg/spec"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/hashicorp/go-multierror"
)

type Config struct {
	Backdrops map[string]*api.Backdrop
	Includes  []string
}

func GetAllBackdrops(filenames ...string) (map[string]*api.Backdrop, error) {
	var errs error
	backdrops := map[string]*api.Backdrop{}

	for _, filename := range filenames {
		value, err := cuetils.ReadYAMLFileWithSpec(spec.CueSpec, filename)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		config, err := ConfigFromValue(value)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		for name, backdrop := range config.Backdrops {
			backdrops[name] = backdrop
		}

		for _, include := range config.Includes {
			included, err := GetAllBackdrops(include)
			if err != nil {
				errs = multierror.Append(errs, err)
				continue
			}

			for name, backdrop := range included {
				backdrops[name] = backdrop
			}
		}
	}

	return backdrops, errs
}

func ConfigFromValue(v cue.Value) (*Config, error) {
	out := &Config{}

	if p, ok := cuetils.Get(v, "backdrops"); ok {
		if bs, err := BackdropsFromValue(p); err != nil {
			return nil, err
		} else {
			out.Backdrops = bs
		}
	}

	if p, ok := cuetils.Get(v, "include"); ok {
		if is, err := IncludesFromValue(p); err != nil {
			return nil, err
		} else {
			out.Includes = is
		}
	}

	return out, nil
}
