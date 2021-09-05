package config

import (
	"errors"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/yaml"
	"github.com/dodo-cli/dodo-config/pkg/spec"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
)

var ErrUnexpectedSpec = errors.New("Whelp, we don't match the expected specification. Now what?")

func ParseConfig(filename string) (map[string]*api.Backdrop, error) {
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

	if err := addFile(ctx, bi, filename); err != nil {
		return nil, err
	}

	value := ctx.BuildInstance(bi)
	if err := value.Err(); err != nil {
		return nil, err
	}

	if err := value.Validate(cue.Concrete(true)); err != nil {
		return nil, err
	}

	return Config(value)
}

func addFile(ctx *cue.Context, bi *build.Instance, filename string) error {
	yamlFile, err := yaml.Extract(filename, nil)
	if err != nil {
		return err
	}

	if err := bi.AddSyntax(yamlFile); err != nil {
		return err
	}

	value := ctx.BuildFile(yamlFile)
	if err := value.Err(); err != nil {
		return err
	}

	if is, ok := property(value, "include"); ok {
		if err := eachInList(is, func(v cue.Value) error {
			if p, ok := property(v, "file"); ok {
				if f, err := p.String(); err == nil {
					return addFile(ctx, bi, f)
				} else {
					return err
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func Config(v cue.Value) (map[string]*api.Backdrop, error) {
	if p, ok := property(v, "backdrops"); !ok {
		return map[string]*api.Backdrop{}, ErrUnexpectedSpec
	} else {
		return BackdropsFromValue(p)
	}
}
