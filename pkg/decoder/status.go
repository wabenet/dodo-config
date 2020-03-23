package decoder

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type (
	Decoder  func(*Status, interface{})
	Producer func() (interface{}, Decoder)

	Status struct {
		file   string
		path   []string
		errors []error
	}
)

func New(filename string) *Status {
	return &Status{
		file:   filename,
		path:   []string{},
		errors: []error{},
	}
}

func (s *Status) Error(msg string) {
	// TODO add file and path to error messages
	s.errors = append(s.errors, errors.New(msg))
}

func (s *Status) Errors() []error {
	return s.errors
}

func (s *Status) Run(name string, decode Decoder, value interface{}) {
	s.path = append(s.path, name)
	decode(s, value)
	s.path = s.path[:len(s.path)-1]
}

func (s *Status) DecodeYaml(content []byte, target interface{}, lookup map[string]Decoder) {
	var mapType map[interface{}]interface{}
	if err := yaml.Unmarshal(content, &mapType); err != nil {
		s.Error("invalid yaml")
		return
	}

	var dummySlice []struct{}
	var dummyItem struct{}

	includeHelper := func() (interface{}, Decoder) {
		return &dummyItem, Keys(map[string]Decoder{
			"text": func(s *Status, config interface{}) {
				var decoded string
				String(&decoded)(s, config)
				s.DecodeYaml([]byte(decoded), target, lookup)
			},
			"file": func(s *Status, config interface{}) {
				var decoded string
				String(&decoded)(s, config)

				bytes, err := readFile(decoded)
				if err != nil {
					s.Error("could not read file")
					return
				}

				sub := New(decoded)
				sub.DecodeYaml(bytes, target, lookup)
				s.errors = append(s.errors, sub.errors...)
			},
		})
	}

	lookup["include"] = Kinds(map[reflect.Kind]Decoder{
		reflect.Map:   Singleton(includeHelper, &dummySlice),
		reflect.Slice: Slice(includeHelper, &dummySlice),
	})

	s.Run("", Keys(lookup), mapType)
}

func readFile(filename string) ([]byte, error) {
	if !filepath.IsAbs(filename) {
		directory, err := os.Getwd()
		if err != nil {
			return []byte{}, err
		}
		filename, err = filepath.Abs(filepath.Join(directory, filename))
		if err != nil {
			return []byte{}, err
		}
	}
	return ioutil.ReadFile(filename)
}
