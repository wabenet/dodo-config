package config

import (
	"github.com/dodo-cli/dodo-config/pkg/cuetils"
	"github.com/dodo-cli/dodo-config/pkg/includes"
	"github.com/dodo-cli/dodo-config/pkg/spec"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/hashicorp/go-multierror"
)

func GetAllBackdrops(filenames ...string) (map[string]*api.Backdrop, error) {
	var errs error
	backdrops := map[string]*api.Backdrop{}

	resolved, err := includes.ResolveIncludes(filenames...)
	if err != nil {
		errs = multierror.Append(errs, err)
		return backdrops, errs
	}

	for _, filename := range resolved {
		value, err := cuetils.ReadYAMLFileWithSpec(spec.CueSpec, filename)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		p, ok := cuetils.Get(value, "backdrops")
		if !ok {
			continue
		}

		bs, err := BackdropsFromValue(p)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		for name, backdrop := range bs {
			backdrops[name] = backdrop
		}
	}

	return backdrops, errs
}
