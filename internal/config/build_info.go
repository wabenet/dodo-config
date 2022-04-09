package config

import (
	"cuelang.org/go/cue"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
)

func BuildInfoFromStruct(v cue.Value) (*api.BuildInfo, error) {
	out := &api.BuildInfo{}

	if p, ok := property(v, "name"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.ImageName = n
		}
	}

	if p, ok := property(v, "builder"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Builder = n
		}
	}

	if p, ok := property(v, "context"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Context = n
		}
	}

	if p, ok := property(v, "dockerfile"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Dockerfile = n
		}
	}

	if p, ok := property(v, "steps"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.InlineDockerfile = []string{n}
		}
	}

	if p, ok := property(v, "dependencies"); ok {
		if deps, err := StringListFromValue(p); err != nil {
			return nil, err
		} else {
			out.Dependencies = deps
		}
	}

	if p, ok := property(v, "arguments"); ok {
		if bas, err := BuildArgumentsFromValue(p); err != nil {
			return nil, err
		} else {
			out.Arguments = bas
		}
	}

	if p, ok := property(v, "secrets"); ok {
		if bss, err := BuildSecretsFromValue(p); err != nil {
			return nil, err
		} else {
			out.Secrets = bss
		}
	}

	if p, ok := property(v, "ssh_agents"); ok {
		if bsa, err := BuildSSHAgentsFromValue(p); err != nil {
			return nil, err
		} else {
			out.SshAgents = bsa
		}
	}

	return out, nil
}
