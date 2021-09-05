package config

import (
	"cuelang.org/go/cue"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
)

func BackdropsFromValue(v cue.Value) (map[string]*api.Backdrop, error) {
	return BackdropsFromMap(v)
}

func BackdropsFromMap(v cue.Value) (map[string]*api.Backdrop, error) {
	out := map[string]*api.Backdrop{}

	err := eachInMap(v, func(name string, v cue.Value) error {
		r, err := BackdropFromStruct(name, v)
		if err == nil {
			out[name] = r
		}

		return err

	})

	return out, err
}

func BackdropFromStruct(name string, v cue.Value) (*api.Backdrop, error) {
	out := &api.Backdrop{
		Name:       name,
		Entrypoint: &api.Entrypoint{},
	}

	if p, ok := property(v, "name"); ok {
		if n, err := p.String(); err != nil {
			return nil, err
		} else {
			out.Name = n
		}
	}

	if p, ok := property(v, "container_name"); ok {
		if n, err := p.String(); err != nil {
			return nil, err
		} else {
			out.ContainerName = n
		}
	}

	if p, ok := property(v, "runtime"); ok {
		if n, err := p.String(); err != nil {
			return nil, err
		} else {
			out.Runtime = n
		}
	}

	if p, ok := property(v, "script"); ok {
		if n, err := p.String(); err != nil {
			return nil, err
		} else {
			out.Entrypoint.Script = n
		}
	}

	if p, ok := property(v, "user"); ok {
		if n, err := p.String(); err != nil {
			return nil, err
		} else {
			out.User = n
		}
	}

	if p, ok := property(v, "working_dir"); ok {
		if n, err := p.String(); err != nil {
			return nil, err
		} else {
			out.WorkingDir = n
		}
	}

	if p, ok := property(v, "image"); ok {
		if n, err := p.String(); err == nil {
			out.ImageId = n
		} else if b, err := BuildInfoFromStruct(p); err == nil {
			out.BuildInfo = b
		} else {
			return nil, ErrUnexpectedSpec
		}
	}

	if p, ok := property(v, "build"); ok {
		if n, err := p.String(); err == nil {
			out.ImageId = n
		} else if b, err := BuildInfoFromStruct(p); err == nil {
			out.BuildInfo = b
		} else {
			return nil, ErrUnexpectedSpec
		}
	}

	if p, ok := property(v, "aliases"); ok {
		if as, err := StringListFromValue(p); err != nil {
			return nil, err
		} else {
			out.Aliases = as
		}
	}

	if p, ok := property(v, "capabilities"); ok {
		if as, err := StringListFromValue(p); err != nil {
			return nil, err
		} else {
			out.Capabilities = as
		}
	}

	if p, ok := property(v, "interpreter"); ok {
		if as, err := StringListFromValue(p); err != nil {
			return nil, err
		} else {
			out.Entrypoint.Interpreter = as
		}
	}

	if p, ok := property(v, "environment"); ok {
		if env, err := EnvironmentVariablesFromValue(p); err != nil {
			return nil, err
		} else {
			out.Environment = env
		}
	}

	if p, ok := property(v, "volumes"); ok {
		if vs, err := VolumeMountsFromValue(p); err != nil {
			return nil, err
		} else {
			out.Volumes = vs
		}
	}

	if p, ok := property(v, "ports"); ok {
		if ps, err := PortBindingsFromValue(p); err != nil {
			return nil, err
		} else {
			out.Ports = ps
		}
	}

	return out, nil
}
