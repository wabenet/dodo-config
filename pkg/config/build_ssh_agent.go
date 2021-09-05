package config

import (
	"cuelang.org/go/cue"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
)

func BuildSSHAgentsFromValue(v cue.Value) ([]*api.SshAgent, error) {
	if out, err := BuildSSHAgentsFromMap(v); err == nil {
		return out, nil
	}

	if out, err := BuildSSHAgentsFromList(v); err == nil {
		return out, nil
	}

	return nil, ErrUnexpectedSpec
}

func BuildSSHAgentsFromMap(v cue.Value) ([]*api.SshAgent, error) {
	out := []*api.SshAgent{}

	err := eachInMap(v, func(name string, v cue.Value) error {
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

	err := eachInList(v, func(v cue.Value) error {
		r, err := BuildSSHAgentFromValue("", v)
		if err == nil {
			out = append(out, r)
		}

		return err
	})

	return out, err
}

func BuildSSHAgentFromValue(name string, v cue.Value) (*api.SshAgent, error) {
	if out, err := BuildSSHAgentFromStruct(name, v); err == nil {
		return out, err
	}

	return nil, ErrUnexpectedSpec
}

func BuildSSHAgentFromStruct(name string, v cue.Value) (*api.SshAgent, error) {
	out := &api.SshAgent{Id: name}

	if p, ok := property(v, "path"); ok {
		if n, err := p.String(); err != nil {
			return nil, err
		} else {
			out.Id = n
		}
	}

	if p, ok := property(v, "identity_file"); ok {
		if n, err := p.String(); err != nil {
			return nil, err
		} else {
			out.IdentityFile = n
		}
	}

	return out, nil
}
