package config

import (
	"cuelang.org/go/cue"
	"strings"
)

func StringFromValue(v cue.Value) (string, error) {
	return v.String()
}

func StringListFromValue(v cue.Value) ([]string, error) {
	out := []string{}

	err := eachInList(v, func(v cue.Value) error {
		str, err := StringFromValue(v)
		if err == nil {
			out = append(out, str)
		}

		return err
	})

	return out, err
}

func property(v cue.Value, name string) (cue.Value, bool) {
	p := v.LookupPath(cue.MakePath(cue.Str(name)))
	return p, p.Exists()
}

func eachInList(v cue.Value, f func(cue.Value) error) error {
	iter, err := v.List()
	if err != nil {
		return err
	}

	for iter.Next() {
		if err := f(iter.Value()); err != nil {
			return err
		}
	}

	return nil
}

func eachInMap(v cue.Value, f func(string, cue.Value) error) error {
	iter, err := v.Fields()
	if err != nil {
		return err
	}

	for iter.Next() {
		// CUE selector is supposed an unambigous path selector in the
		// map, not just the key. So it might be quoted.
		// FIXME: We simply trim the quotes here, to get the map key from
		// the selector, which is kinda hacky and will probably cause
		// trouble later
		name := strings.Trim(iter.Selector().String(), `"`)

		if err := f(name, iter.Value()); err != nil {
			return err
		}
	}

	return nil
}
