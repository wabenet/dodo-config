package main

import (
	"os"

	"github.com/dodo-cli/dodo-config/pkg/command"
	config "github.com/dodo-cli/dodo-config/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/appconfig"
	dodo "github.com/dodo-cli/dodo-core/pkg/plugin"
	log "github.com/hashicorp/go-hclog"
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
