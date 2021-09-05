package config_test

import (
	"testing"

	"github.com/dodo-cli/dodo-config/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestBuildInfo(t *testing.T) {
	config, err := config.ParseConfig("test/dodo.yaml")

	assert.Nil(t, err)

	backdrop, ok := config["test_build"]
	assert.True(t, ok)

        build := backdrop.BuildInfo
	assert.NotNil(t, backdrop.BuildInfo)

	assert.Equal(t, "testimage", build.ImageName)
	assert.Equal(t, "/some/path", build.Context)
	assert.Equal(t, "/some/other/path", build.Dockerfile)
	assert.Equal(t, []string{"FROM foo\n"}, build.InlineDockerfile)
}
