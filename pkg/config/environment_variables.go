package config

import (
	"cuelang.org/go/cue"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/config"
)

func EnvironmentVariablesFromValue(v cue.Value) ([]*api.EnvironmentVariable, error) {
	if out, err := EnvironmentVariablesFromMap(v); err == nil {
		return out, nil
	}

	if out, err := EnvironmentVariablesFromList(v); err == nil {
		return out, nil
	}

	return nil, ErrUnexpectedSpec
}

func EnvironmentVariablesFromMap(v cue.Value) ([]*api.EnvironmentVariable, error) {
	out := []*api.EnvironmentVariable{}

	err := eachInMap(v, func(name string, v cue.Value) error {
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

	err := eachInList(v, func(v cue.Value) error {
		r, err := EnvironmentVariableFromValue("", v)
		if err == nil {
			out = append(out, r)
		}

		return err
	})

	return out, err
}

func EnvironmentVariableFromValue(name string, v cue.Value) (*api.EnvironmentVariable, error) {
	if out, err := EnvironmentVariableFromString(name, v); err == nil {
		return out, err
	}

	if out, err := EnvironmentVariableFromStruct(name, v); err == nil {
		return out, err
	}

	return nil, ErrUnexpectedSpec
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

	if p, ok := property(v, "name"); ok {
		if v, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Key = v
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
