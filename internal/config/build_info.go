package config

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/wabenet/dodo-config/pkg/cuetils"
	api "github.com/wabenet/dodo-core/api/core/v1alpha5"
)

func BuildInfoFromStruct(_ string, value cue.Value) (*api.BuildInfo, error) {
	out := &api.BuildInfo{}

	if p, ok, err := cuetils.Extract(value, "name", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "name", err)
	} else if ok {
		out.ImageName = p
	}

	if p, ok, err := cuetils.Extract(value, "builder", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "builder", err)
	} else if ok {
		out.Builder = p
	}

	if p, ok, err := cuetils.Extract(value, "context", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "context", err)
	} else if ok {
		out.Context = p
	}

	if p, ok, err := cuetils.Extract(value, "dockerfile", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "dockerfile", err)
	} else if ok {
		out.Dockerfile = p
	}

	if p, ok, err := cuetils.Extract(value, "steps", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "steps", err)
	} else if ok {
		out.InlineDockerfile = []string{p}
	}

	if p, ok, err := cuetils.Extract(value, "dependencies", cuetils.OneOrMore(cuetils.String)); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "dependencies", err)
	} else if ok {
		out.Dependencies = p
	}

	if p, ok, err := cuetils.Extract(value, "arguments", cuetils.ListOrDict(BuildArgumentFromStruct)); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "arguments", err)
	} else if ok {
		out.Arguments = p
	}

	if p, ok, err := cuetils.Extract(value, "secrets", cuetils.ListOrDict(BuildSecretFromStruct)); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "secrets", err)
	} else if ok {
		out.Secrets = p
	}

	if p, ok, err := cuetils.Extract(value, "ssh_agents", cuetils.ListOrDict(BuildSSHAgentFromStruct)); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "ssh_agents", err)
	} else if ok {
		out.SshAgents = p
	}

	return out, nil
}

func BuildArgumentFromStruct(name string, value cue.Value) (*api.BuildArgument, error) {
	out := &api.BuildArgument{Key: name}

	if p, ok, err := cuetils.Extract(value, "name", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "name", err)
	} else if ok {
		out.Key = p
	}

	if p, ok, err := cuetils.Extract(value, "value", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "value", err)
	} else if ok {
		out.Value = p
	}

	return out, nil
}

func BuildSecretFromStruct(name string, value cue.Value) (*api.BuildSecret, error) {
	out := &api.BuildSecret{Id: name}

	if p, ok, err := cuetils.Extract(value, "id", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "id", err)
	} else if ok {
		out.Id = p
	}

	if p, ok, err := cuetils.Extract(value, "path", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "path", err)
	} else if ok {
		out.Path = p
	}

	return out, nil
}

func BuildSSHAgentFromStruct(name string, value cue.Value) (*api.SshAgent, error) {
	out := &api.SshAgent{Id: name}

	if p, ok, err := cuetils.Extract(value, "path", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "path", err)
	} else if ok {
		out.Id = p
	}

	if p, ok, err := cuetils.Extract(value, "identity_file", cuetils.String); err != nil {
		return nil, fmt.Errorf("invalid config for %s: %w", "identity_file", err)
	} else if ok {
		out.IdentityFile = p
	}

	return out, nil
}
