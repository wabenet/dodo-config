package main

import (
	config "github.com/dodo/dodo-config/pkg/plugin"
	dodo "github.com/oclaussen/dodo/pkg/plugin"
)

func main() {
	config.RegisterPlugin()
	dodo.ServePlugins()
}
