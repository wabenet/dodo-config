package cuetils

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/yaml"
	"github.com/wabenet/dodo-config/pkg/template"
)

func ReadYAMLFileWithSpec(spec, filename string) (cue.Value, error) {
	ctx := cuecontext.New()

	instances := load.Instances([]string{"-"}, &load.Config{
		Stdin: strings.NewReader(spec),
	})

	if len(instances) != 1 {
		return cue.Value{}, fmt.Errorf("expected exactly one instance")
	}

	instance := instances[0]

	if instance.Err != nil {
		return cue.Value{}, instance.Err
	}

	yamlFile, err := yaml.Extract(filename, nil)
	if err != nil {
		return cue.Value{}, fmt.Errorf("get yaml contents of %s: %w", filename, err)
	}

	// We template the source on the AST level - still not entirely sure
	// that's a good idea.
	// We could template the values after building with CUE, but that will
	// proably mess too much with CUEs validation engine.
	// We could template the source files directly before parsing, but I'm a
	// big fan if our the sources are always valid yaml and can be validated
	// as such - otherwise they can get really messy and it kinda defeats
	// the point of using CUE in the first place.

	yamlFile, err = template.TemplateCueAST(yamlFile)
	if err != nil {
		return cue.Value{}, fmt.Errorf("templating error in %s: %w", filename, err)
	}

	if err := instance.AddSyntax(yamlFile); err != nil {
		return cue.Value{}, fmt.Errorf("cue syntax error in %s: %w", filename, err)
	}

	value := ctx.BuildInstance(instance)
	if err := value.Err(); err != nil {
		return cue.Value{}, fmt.Errorf("cue build error in %s: %w", filename, err)
	}

	if err := value.Validate(cue.Concrete(true), cue.Final()); err != nil {
		return cue.Value{}, fmt.Errorf("cue validation error in %s: %w", filename, err)
	}

	return value, nil
}
