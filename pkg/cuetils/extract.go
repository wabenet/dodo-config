package cuetils

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"github.com/hashicorp/go-multierror"
)

type Extractor[V any] func(string, cue.Value) (V, error)

func Extract[V any](value cue.Value, key string, extract Extractor[V]) (V, bool, error) {
	var out V

	def, _ := value.Default() // FIXME this feels broken...
	p := def.LookupPath(cue.MakePath(cue.Str(key)))
	if !p.Exists() {
		return out, false, nil
	}

	out, err := extract(key, p)
	if err != nil {
		return out, true, err
	}

	return out, true, nil
}

func ParseString[V any](extract func(string) (V, error)) Extractor[V] {
	return func(name string, value cue.Value) (V, error) {
		var out V

		s, err := String(name, value)
		if err != nil {
			return out, err
		}

		out, err = extract(s)
		if err != nil {
			return out, err
		}

		return out, nil
	}
}

func Either[V any](extracts []Extractor[V]) Extractor[V] {
	return func(name string, value cue.Value) (V, error) {
		var (
			out  V
			errs error
		)

		for _, extract := range extracts {
			out, err := extract(name, value)
			if err == nil {
				return out, nil
			}

			errs = multierror.Append(errs, err)
		}

		return out, errs
	}
}

func Map[V any](extract Extractor[V]) Extractor[map[string]V] {
	return func(name string, value cue.Value) (map[string]V, error) {
		out := make(map[string]V)

		iter, err := value.Fields()
		if err != nil {
			return out, fmt.Errorf("value for %s is not a map: %w", name, err)
		}

		for iter.Next() {
			name := strings.Trim(iter.Selector().String(), `"`)

			r, err := extract(name, iter.Value())
			if err != nil {
				return out, err
			}

			out[name] = r
		}

		return out, nil
	}
}

func List[V any](extract Extractor[V]) Extractor[[]V] {
	return func(name string, value cue.Value) ([]V, error) {
		out := []V{}

		iter, err := value.List()
		if err != nil {
			return out, fmt.Errorf("value for %s is not a list: %w", name, err)
		}

		for iter.Next() {
			r, err := extract(name, iter.Value())
			if err != nil {
				return out, err
			}

			out = append(out, r)
		}

		return out, nil
	}
}

func OneOrMore[V any](extract Extractor[V]) Extractor[[]V] {
	return func(name string, value cue.Value) ([]V, error) {
		out := []V{}

		if p, err := extract(name, value); err == nil {
			return append(out, p), nil
		}

		iter, err := value.List()
		if err != nil {
			return out, fmt.Errorf("value for %s is not a list: %w", name, err)
		}

		for iter.Next() {
			str, err := extract(name, iter.Value())
			if err != nil {
				return out, err
			}

			out = append(out, str)
		}

		return out, nil
	}
}

func ListOrDict[V any](extract Extractor[V]) Extractor[[]V] {
	return func(name string, value cue.Value) ([]V, error) {
		var errs error

		out, err := fromMap(name, value, extract)
		if err == nil {
			return out, nil
		}

		errs = multierror.Append(errs, err)

		out, err = fromList(name, value, extract)
		if err == nil {
			return out, nil
		}

		errs = multierror.Append(errs, err)

		return nil, errs
	}
}

func fromMap[V any](name string, value cue.Value, extract Extractor[V]) ([]V, error) {
	out := []V{}

	iter, err := value.Fields()
	if err != nil {
		return out, fmt.Errorf("value for %s is not a map: %w", name, err)
	}

	for iter.Next() {
		name := strings.Trim(iter.Selector().String(), `"`)

		r, err := extract(name, iter.Value())
		if err != nil {
			return out, err
		}

		out = append(out, r)
	}

	return out, nil
}

func fromList[V any](name string, value cue.Value, extract Extractor[V]) ([]V, error) {
	out := []V{}

	iter, err := value.List()
	if err != nil {
		return out, fmt.Errorf("value for %s is not a list: %w", name, err)
	}

	for iter.Next() {
		r, err := extract("", iter.Value())
		if err != nil {
			return out, err
		}

		out = append(out, r)
	}

	return out, nil
}
