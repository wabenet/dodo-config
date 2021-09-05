package config_test

import (
	"testing"

	"github.com/dodo-cli/dodo-config/pkg/config"
	"github.com/stretchr/testify/assert"
)

func ignore_TestInclude(t *testing.T) {
	config, err := config.ParseConfig("test/dodo.yaml")
	assert.Nil(t, err)

	_, ok := config["test2"]

	assert.True(t, ok)
}
