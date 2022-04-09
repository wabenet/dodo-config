package config_test

import (
	"github.com/dodo-cli/dodo-config/internal/config"
	"github.com/dodo-cli/dodo-config/pkg/cuetils"
	"github.com/dodo-cli/dodo-config/pkg/spec"
)

const (
	TestConfig = "test/dodo.yaml"
)

func ParseTestConfig() (*config.Config, error) {
	value, err := cuetils.ReadYAMLFileWithSpec(spec.CueSpec, TestConfig)
	if err != nil {
		return nil, err
	}

	return config.ConfigFromValue(value)
}
