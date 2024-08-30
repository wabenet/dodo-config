package config

import (
	"github.com/hashicorp/go-multierror"
	"github.com/wabenet/dodo-config/pkg/cuetils"
	"github.com/wabenet/dodo-config/pkg/includes"
	"github.com/wabenet/dodo-config/pkg/spec"
	api "github.com/wabenet/dodo-core/api/core/v1alpha5"
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

		if p, ok, err := cuetils.Extract(value, "backdrops", cuetils.Map(BackdropFromStruct)); err != nil {
			errs = multierror.Append(errs, err)
			continue
		} else if ok {
			for name, backdrop := range p {
				backdrops[name] = backdrop
			}
		}
	}

	return backdrops, errs
}
