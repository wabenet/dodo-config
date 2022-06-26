package config

import (
	"cuelang.org/go/cue"
	"github.com/hashicorp/go-multierror"
	"github.com/wabenet/dodo-config/pkg/cuetils"
	api "github.com/wabenet/dodo-core/api/v1alpha3"
	"github.com/wabenet/dodo-core/pkg/config"
)

func EnvironmentVariablesFromValue(v cue.Value) ([]*api.EnvironmentVariable, error) {
	var errs error

	if out, err := EnvironmentVariablesFromMap(v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	if out, err := EnvironmentVariablesFromList(v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	return nil, errs
}

func EnvironmentVariablesFromMap(v cue.Value) ([]*api.EnvironmentVariable, error) {
	out := []*api.EnvironmentVariable{}

	err := cuetils.IterMap(v, func(name string, v cue.Value) error {
		r, err := EnvironmentVariableFromValue(name, v)
		if err == nil {
			out = append(out, r)
		}

		return err

	})

	return out, err
}

func EnvironmentVariablesFromList(v cue.Value) ([]*api.EnvironmentVariable, error) {
	out := []*api.EnvironmentVariable{}

	err := cuetils.IterList(v, func(v cue.Value) error {
		r, err := EnvironmentVariableFromValue("", v)
		if err == nil {
			out = append(out, r)
		}

		return err
	})

	return out, err
}

func EnvironmentVariableFromValue(name string, v cue.Value) (*api.EnvironmentVariable, error) {
	var errs error

	if out, err := EnvironmentVariableFromString(name, v); err == nil {
		return out, err
	} else {
		errs = multierror.Append(errs, err)
	}

	if out, err := EnvironmentVariableFromStruct(name, v); err == nil {
		return out, err
	} else {
		errs = multierror.Append(errs, err)
	}

	return nil, errs
}

func EnvironmentVariableFromString(_ string, v cue.Value) (*api.EnvironmentVariable, error) {
	s, err := StringFromValue(v)
	if err != nil {
		return nil, err
	}

	return config.ParseEnvironmentVariable(s)
}

func EnvironmentVariableFromStruct(name string, v cue.Value) (*api.EnvironmentVariable, error) {
	out := &api.EnvironmentVariable{Key: name}

	if p, ok := cuetils.Get(v, "name"); ok {
		if v, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Key = v
		}
	}

	if p, ok := cuetils.Get(v, "value"); ok {
		if v, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Value = v
		}
	}

	return out, nil
}
