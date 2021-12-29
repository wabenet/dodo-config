package config

import (
	"errors"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/yaml"
	"github.com/dodo-cli/dodo-config/pkg/spec"
	"github.com/dodo-cli/dodo-config/pkg/template"
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
		config, err := ParseConfig(filename)
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

func ParseConfig(filename string) (*Config, error) {
	ctx := cuecontext.New()

	bis := load.Instances([]string{"-"}, &load.Config{
		Stdin: strings.NewReader(spec.CueSpec),
	})

	if len(bis) != 1 {
		return nil, errors.New("expected exactly one instance")
	}

	bi := bis[0]

	if bi.Err != nil {
		return nil, bi.Err
	}

	yamlFile, err := yaml.Extract(filename, nil)
	if err != nil {
		return nil, err
	}

	yamlFile, err = template.TemplateCueAST(yamlFile)
	if err != nil {
		return nil, err
	}

	if err := bi.AddSyntax(yamlFile); err != nil {
		return nil, err
	}

	value := ctx.BuildInstance(bi)
	if err := value.Err(); err != nil {
		return nil, err
	}

	if err := value.Validate(cue.Concrete(true), cue.Final()); err != nil {
		return nil, err
	}

	return ConfigFromValue(value)
}

func ConfigFromValue(v cue.Value) (*Config, error) {
	out := &Config{}

	if p, ok := property(v, "backdrops"); ok {
		if bs, err := BackdropsFromValue(p); err != nil {
			return nil, err
		} else {
			out.Backdrops = bs
		}
	}

	if p, ok := property(v, "include"); ok {
		if is, err := IncludesFromValue(p); err != nil {
			return nil, err
		} else {
			out.Includes = is
		}
	}

	return out, nil
}
