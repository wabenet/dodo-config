package config

import (
	"cuelang.org/go/cue"
	"github.com/dodo-cli/dodo-config/pkg/cuetils"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/hashicorp/go-multierror"
)

func BuildSSHAgentsFromValue(v cue.Value) ([]*api.SshAgent, error) {
	var errs error

	if out, err := BuildSSHAgentsFromMap(v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	if out, err := BuildSSHAgentsFromList(v); err == nil {
		return out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	return nil, errs
}

func BuildSSHAgentsFromMap(v cue.Value) ([]*api.SshAgent, error) {
	out := []*api.SshAgent{}

	err := cuetils.IterMap(v, func(name string, v cue.Value) error {
		r, err := BuildSSHAgentFromValue(name, v)
		if err == nil {
			out = append(out, r)
		}

		return err

	})

	return out, err
}

func BuildSSHAgentsFromList(v cue.Value) ([]*api.SshAgent, error) {
	out := []*api.SshAgent{}

	err := cuetils.IterList(v, func(v cue.Value) error {
		r, err := BuildSSHAgentFromValue("", v)
		if err == nil {
			out = append(out, r)
		}

		return err
	})

	return out, err
}

func BuildSSHAgentFromValue(name string, v cue.Value) (*api.SshAgent, error) {
	var errs error

	if out, err := BuildSSHAgentFromStruct(name, v); err == nil {
		return out, err
	} else {
		errs = multierror.Append(errs, err)
	}

	return nil, errs
}

func BuildSSHAgentFromStruct(name string, v cue.Value) (*api.SshAgent, error) {
	out := &api.SshAgent{Id: name}

	if p, ok := cuetils.Get(v, "path"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Id = n
		}
	}

	if p, ok := cuetils.Get(v, "identity_file"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.IdentityFile = n
		}
	}

	return out, nil
}
