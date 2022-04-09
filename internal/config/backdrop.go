package config

import (
	"cuelang.org/go/cue"
	"github.com/dodo-cli/dodo-config/pkg/cuetils"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/hashicorp/go-multierror"
)

func BackdropsFromValue(v cue.Value) (map[string]*api.Backdrop, error) {
	return BackdropsFromMap(v)
}

func BackdropsFromMap(v cue.Value) (map[string]*api.Backdrop, error) {
	out := map[string]*api.Backdrop{}

	err := cuetils.IterMap(v, func(name string, v cue.Value) error {
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

	if p, ok := cuetils.Get(v, "name"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Name = n
		}
	}

	if p, ok := cuetils.Get(v, "container_name"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.ContainerName = n
		}
	}

	if p, ok := cuetils.Get(v, "runtime"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Runtime = n
		}
	}

	if p, ok := cuetils.Get(v, "script"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.Entrypoint.Script = n
		}
	}

	if p, ok := cuetils.Get(v, "user"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.User = n
		}
	}

	if p, ok := cuetils.Get(v, "working_dir"); ok {
		if n, err := StringFromValue(p); err != nil {
			return nil, err
		} else {
			out.WorkingDir = n
		}
	}

	if p, ok := cuetils.Get(v, "image"); ok {
		if n, b, err := ImageOrBuildInfoFromValue(p); err != nil {
			return nil, err
		} else {
			out.ImageId = n
			out.BuildInfo = b
		}
	}

	if p, ok := cuetils.Get(v, "build"); ok {
		if n, b, err := ImageOrBuildInfoFromValue(p); err != nil {
			return nil, err
		} else {
			out.ImageId = n
			out.BuildInfo = b
		}
	}

	if p, ok := cuetils.Get(v, "aliases"); ok {
		if as, err := StringListFromValue(p); err != nil {
			return nil, err
		} else {
			out.Aliases = as
		}
	}

	if p, ok := cuetils.Get(v, "capabilities"); ok {
		if as, err := StringListFromValue(p); err != nil {
			return nil, err
		} else {
			out.Capabilities = as
		}
	}

	if p, ok := cuetils.Get(v, "interpreter"); ok {
		if as, err := StringListFromValue(p); err != nil {
			return nil, err
		} else {
			out.Entrypoint.Interpreter = as
		}
	}

	if p, ok := cuetils.Get(v, "environment"); ok {
		if env, err := EnvironmentVariablesFromValue(p); err != nil {
			return nil, err
		} else {
			out.Environment = env
		}
	}

	if p, ok := cuetils.Get(v, "volumes"); ok {
		if vs, err := VolumeMountsFromValue(p); err != nil {
			return nil, err
		} else {
			out.Volumes = vs
		}
	}

	if p, ok := cuetils.Get(v, "ports"); ok {
		if ps, err := PortBindingsFromValue(p); err != nil {
			return nil, err
		} else {
			out.Ports = ps
		}
	}

	if p, ok := cuetils.Get(v, "devices"); ok {
		if ds, err := DeviceMappingsFromValue(p); err != nil {
			return nil, err
		} else {
			out.Devices = ds
		}
	}

	return out, nil
}

func ImageOrBuildInfoFromValue(v cue.Value) (string, *api.BuildInfo, error) {
	var errs error

	if out, err := StringFromValue(v); err == nil {
		return out, nil, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	if out, err := BuildInfoFromStruct(v); err == nil {
		return "", out, nil
	} else {
		errs = multierror.Append(errs, err)
	}

	return "", nil, errs
}
