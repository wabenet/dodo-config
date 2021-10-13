package config

import (
	"cuelang.org/go/cue"
)

func IncludesFromValue(v cue.Value) ([]string, error) {
	if out, err := IncludesFromList(v); err == nil {
		return out, nil
	}

	return nil, ErrUnexpectedSpec
}

func IncludesFromList(v cue.Value) ([]string, error) {
	out := []string{}

	err := eachInList(v, func(v cue.Value) error {
		if p, ok := property(v, "file"); ok {
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
