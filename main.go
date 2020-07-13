package main

import (
	"os"

	"github.com/dodo/dodo-config/pkg/command"
	config "github.com/dodo/dodo-config/pkg/plugin"
	dodo "github.com/oclaussen/dodo/pkg/plugin"
)

func main() {
	if os.Getenv(dodo.MagicCookieKey) == dodo.MagicCookieValue {
		config.RegisterPlugin()
		dodo.ServePlugins()
	} else {
		cmd := command.NewCommand()
		if err := cmd.Execute(); err != nil {
			os.Exit(1)
		}
	}
}
