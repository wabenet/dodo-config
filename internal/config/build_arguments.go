package config

import (
	"cuelang.org/go/cue"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/hashicorp/go-multierror"
)

func BuildArgumentsFromValue(v cue.Value) ([]*api.BuildArgument, error) {
	var errs error

	if out, err := BuildArgumentsFromMap(v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	if out, err := BuildArgumentsFromList(v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	return nil, errs
}

func BuildArgumentsFromMap(v cue.Value) ([]*api.BuildArgument, error) {
	out := []*api.BuildArgument{}

	err := eachInMap(v, func(name string, v cue.Value) error {
		r, err := BuildArgumentFromValue(name, v)
		if err == nil {
			out = append(out, r)
		}

		return err

	})

	return out, err
}

func BuildArgumentsFromList(v cue.Value) ([]*api.BuildArgument, error) {
	out := []*api.BuildArgument{}

	err := eachInList(v, func(v cue.Value) error {
		r, err := BuildArgumentFromValue("", v)
		if err == nil {
			out = append(out, r)
		}

		return err
	})

	return out, err
}

func BuildArgumentFromValue(name string, v cue.Value) (*api.BuildArgument, error) {
	var errs error

	if out, err := BuildArgumentFromStruct(name, v); err == nil {
		return out, err
	} else {
		errs = multierror.Append(errs, err)
	}

	return nil, errs
}

func BuildArgumentFromStruct(name string, v cue.Value) (*api.BuildArgument, error) {
	out := &api.BuildArgument{Key: name}

	if p, ok := property(v, "name"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Key = n
		}
	}

	if p, ok := property(v, "value"); ok {
		if v, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Value = v
		}
	}

	return out, nil
}
