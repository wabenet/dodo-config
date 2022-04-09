package config

import (
	"cuelang.org/go/cue"
	"github.com/dodo-cli/dodo-config/pkg/cuetils"
)

func StringFromValue(v cue.Value) (string, error) {
	return v.String()
}

func StringListFromValue(v cue.Value) ([]string, error) {
	out := []string{}

	err := cuetils.IterList(v, func(v cue.Value) error {
		str, err := StringFromValue(v)
		if err == nil {
			out = append(out, str)
		}

		return err
	})

	return out, err
}
