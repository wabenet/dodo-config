package config

import (
	"cuelang.org/go/cue"
	"github.com/dodo-cli/dodo-config/pkg/cuetils"
	"github.com/hashicorp/go-multierror"
)

func IncludesFromValue(v cue.Value) ([]string, error) {
	var errs error

	if out, err := IncludesFromList(v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	return nil, errs
}

func IncludesFromList(v cue.Value) ([]string, error) {
	out := []string{}

	err := cuetils.IterList(v, func(v cue.Value) error {
		if p, ok := cuetils.Get(v, "file"); ok {
			f, err := StringFromValue(p)
			if err == nil {
				out = append(out, f)
			}

			return err
		}

		return nil
	})

	return out, err
}
