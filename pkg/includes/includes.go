package includes

import (
	_ "embed"

	"cuelang.org/go/cue"
	"github.com/dodo-cli/dodo-config/pkg/cuetils"
	"github.com/hashicorp/go-multierror"
)

//go:embed includes.cue
var CueSpec string

func ResolveIncludes(filenames ...string) ([]string, error) {
	var errs error
	resolved := filenames

	for _, filename := range filenames {
		v, err := cuetils.ReadYAMLFileWithSpec(CueSpec, filename)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		p, ok := cuetils.Get(v, "include")
		if !ok {
			continue
		}

		includes, err := includesFromValue(p)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		for _, include := range includes {
			incs, err := ResolveIncludes(include)
			if err != nil {
				errs = multierror.Append(errs, err)
				continue
			}

			resolved = append(resolved, incs...)
		}
	}

	return resolved, errs
}

func includesFromValue(v cue.Value) ([]string, error) {
	var errs error

	if out, err := includesFromList(v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	return nil, errs
}

func includesFromList(v cue.Value) ([]string, error) {
	out := []string{}

	err := cuetils.IterList(v, func(v cue.Value) error {
		if p, ok := cuetils.Get(v, "file"); ok {
			f, err := p.String()
			if err == nil {
				out = append(out, f)
			}

			return err
		}

		return nil
	})

	return out, err
}
