package main

import (
	_ "github.com/dodo/dodo-config/pkg/configuration"
	"github.com/oclaussen/dodo/pkg/plugin"
)

func main() {
	plugin.ServePlugins()
}
