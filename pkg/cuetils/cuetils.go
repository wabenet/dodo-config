package cuetils

import (
	"errors"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/encoding/yaml"
	"github.com/wabenet/dodo-config/pkg/template"
)

func Get(v cue.Value, property string) (cue.Value, bool) {
	p := v.LookupPath(cue.MakePath(cue.Str(property)))
	return p, p.Exists()
}

func IterList(v cue.Value, f func(cue.Value) error) error {
	iter, err := v.List()
	if err != nil {
		return err
	}

	for iter.Next() {
		if err := f(iter.Value()); err != nil {
			return err
		}
	}

	return nil
}

func IterMap(v cue.Value, f func(string, cue.Value) error) error {
	iter, err := v.Fields()
	if err != nil {
		return err
	}

	for iter.Next() {
		// CUE selector is supposed an unambigous path selector in the
		// map, not just the key. So it might be quoted.
		// FIXME: We simply trim the quotes here, to get the map key from
		// the selector, which is kinda hacky and will probably cause
		// trouble later
		name := strings.Trim(iter.Selector().String(), `"`)

		if err := f(name, iter.Value()); err != nil {
			return err
		}
	}

	return nil
}

func ReadYAMLFileWithSpec(spec string, filename string) (cue.Value, error) {
	ctx := cuecontext.New()

	bis := load.Instances([]string{"-"}, &load.Config{
		Stdin: strings.NewReader(spec),
	})

	if len(bis) != 1 {
		return cue.Value{}, errors.New("expected exactly one instance")
	}

	bi := bis[0]

	if bi.Err != nil {
		return cue.Value{}, bi.Err
	}

	yamlFile, err := yaml.Extract(filename, nil)
	if err != nil {
		return cue.Value{}, err
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
		return cue.Value{}, err
	}

	if err := bi.AddSyntax(yamlFile); err != nil {
		return cue.Value{}, err
	}

	value := ctx.BuildInstance(bi)
	if err := value.Err(); err != nil {
		return cue.Value{}, err
	}

	if err := value.Validate(cue.Concrete(true), cue.Final()); err != nil {
		return cue.Value{}, err
	}

	return value, nil
}
