package includes

import (
	_ "embed"
	"fmt"

	"cuelang.org/go/cue"
	"github.com/hashicorp/go-multierror"
	"github.com/wabenet/dodo-config/pkg/cuetils"
)

//go:embed includes.cue
var CueSpec string

type include struct {
	file string
}

func includeFromStruct(_ string, v cue.Value) (*include, error) {
	out := &include{}

	if p, ok, err := cuetils.Extract(v, "file", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "file", err)
	} else if ok {
		out.file = p
	}

	return out, nil
}

func ResolveIncludes(filenames ...string) ([]string, error) {
	var errs error
	resolved := filenames

	for _, filename := range filenames {
		v, err := cuetils.ReadYAMLFileWithSpec(CueSpec, filename)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		p, ok, err := cuetils.Extract(v, "include", cuetils.List(includeFromStruct))
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		if !ok {
			continue
		}

		for _, include := range p {
			incs, err := ResolveIncludes(include.file)
			if err != nil {
				errs = multierror.Append(errs, err)
				continue
			}

			resolved = append(resolved, incs...)
		}
	}

	return resolved, errs
}
