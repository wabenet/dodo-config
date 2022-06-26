package config

import (
	"cuelang.org/go/cue"
	"github.com/wabenet/dodo-config/pkg/cuetils"
	api "github.com/wabenet/dodo-core/api/v1alpha3"
)

func BuildInfoFromStruct(v cue.Value) (*api.BuildInfo, error) {
	out := &api.BuildInfo{}

	if p, ok := cuetils.Get(v, "name"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.ImageName = n
		}
	}

	if p, ok := cuetils.Get(v, "builder"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Builder = n
		}
	}

	if p, ok := cuetils.Get(v, "context"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Context = n
		}
	}

	if p, ok := cuetils.Get(v, "dockerfile"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Dockerfile = n
		}
	}

	if p, ok := cuetils.Get(v, "steps"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.InlineDockerfile = []string{n}
		}
	}

	if p, ok := cuetils.Get(v, "dependencies"); ok {
		if deps, err := StringListFromValue(p); err != nil {
			return nil, err
		} else {
			out.Dependencies = deps
		}
	}

	if p, ok := cuetils.Get(v, "arguments"); ok {
		if bas, err := BuildArgumentsFromValue(p); err != nil {
			return nil, err
		} else {
			out.Arguments = bas
		}
	}

	if p, ok := cuetils.Get(v, "secrets"); ok {
		if bss, err := BuildSecretsFromValue(p); err != nil {
			return nil, err
		} else {
			out.Secrets = bss
		}
	}

	if p, ok := cuetils.Get(v, "ssh_agents"); ok {
		if bsa, err := BuildSSHAgentsFromValue(p); err != nil {
			return nil, err
		} else {
			out.SshAgents = bsa
		}
	}

	return out, nil
}
