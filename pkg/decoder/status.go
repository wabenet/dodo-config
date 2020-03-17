package decoder

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Status struct {
	file   string
	path   []string
	errors []error
}

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

func (s *Status) LoadYaml(decode Decoder, content []byte) {
	var mapType map[interface{}]interface{}
	if err := yaml.Unmarshal(content, &mapType); err != nil {
		s.Error("invalid yaml")
		return
	}
        s.Run("", decode, mapType)
}
