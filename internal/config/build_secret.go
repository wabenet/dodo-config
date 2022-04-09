package config

import (
	"cuelang.org/go/cue"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/hashicorp/go-multierror"
)

func BuildSecretsFromValue(v cue.Value) ([]*api.BuildSecret, error) {
	var errs error

	if out, err := BuildSecretsFromMap(v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	if out, err := BuildSecretsFromList(v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	return nil, errs
}

func BuildSecretsFromMap(v cue.Value) ([]*api.BuildSecret, error) {
	out := []*api.BuildSecret{}

	err := eachInMap(v, func(name string, v cue.Value) error {
		r, err := BuildSecretFromValue(name, v)
		if err == nil {
			out = append(out, r)
		}

		return err

	})

	return out, err
}

func BuildSecretsFromList(v cue.Value) ([]*api.BuildSecret, error) {
	out := []*api.BuildSecret{}

	err := eachInList(v, func(v cue.Value) error {
		r, err := BuildSecretFromValue("", v)
		if err == nil {
			out = append(out, r)
		}

		return err
	})

	return out, err
}

func BuildSecretFromValue(name string, v cue.Value) (*api.BuildSecret, error) {
	var errs error

	if out, err := BuildSecretFromStruct(name, v); err == nil {
		return out, err
	} else {
		errs = multierror.Append(errs, err)
	}

	return nil, errs
}

func BuildSecretFromStruct(name string, v cue.Value) (*api.BuildSecret, error) {
	out := &api.BuildSecret{Id: name}

	if p, ok := property(v, "id"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Id = n
		}
	}

	if p, ok := property(v, "path"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Path = n
		}
	}

	return out, nil
}
