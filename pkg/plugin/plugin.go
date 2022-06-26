package plugin

import (
	"os"

	commandimpl "github.com/wabenet/dodo-config/internal/plugin/command"
	configimpl "github.com/wabenet/dodo-config/internal/plugin/configuration"
	"github.com/wabenet/dodo-core/pkg/plugin"
	"github.com/wabenet/dodo-core/pkg/plugin/command"
	"github.com/wabenet/dodo-core/pkg/plugin/configuration"
)

func RunMe() int {
	m := plugin.Init()

	if os.Getenv(plugin.MagicCookieKey) == plugin.MagicCookieValue {
		m.ServePlugins(NewConfiguration())

		return 0
	} else {
		if err := NewCommand().GetCobraCommand().Execute(); err != nil {
			return 1
		}

		return 0
	}
}

func IncludeMe(m plugin.Manager) {
	m.IncludePlugins(NewConfiguration(), NewCommand())
}

func NewCommand() command.Command {
	return commandimpl.New()
}

func NewConfiguration() configuration.Configuration {
	return configimpl.New()
}
