package main

import (
	"os"

	"github.com/dodo/dodo-config/pkg/command"
	config "github.com/dodo/dodo-config/pkg/plugin"
	log "github.com/hashicorp/go-hclog"
	"github.com/oclaussen/dodo/pkg/appconfig"
	dodo "github.com/oclaussen/dodo/pkg/plugin"
)

func main() {
	if os.Getenv(dodo.MagicCookieKey) == dodo.MagicCookieValue {
		dodo.ServePlugins(&config.Configuration{})
	} else {
		log.SetDefault(log.New(appconfig.GetLoggerOptions()))
		cmd := command.NewCommand()
		if err := cmd.Execute(); err != nil {
			os.Exit(1)
		}
	}
}
