package plugin

import (
	"os"

	"github.com/dodo-cli/dodo-config/pkg/plugin/command"
	"github.com/dodo-cli/dodo-config/pkg/plugin/configuration"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
)

func RunMe() int {
	m := plugin.Init()

	if os.Getenv(plugin.MagicCookieKey) == plugin.MagicCookieValue {
		m.ServePlugins(configuration.New())

		return 0
	} else {
		if err := command.New().GetCobraCommand().Execute(); err != nil {
			return 1
		}

		return 0
	}
}

func IncludeMe(m plugin.Manager) {
	m.IncludePlugins(configuration.New(), command.New())
}
