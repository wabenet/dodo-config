package command

import (
	"fmt"

	"github.com/dodo/dodo-config/pkg/configuration"
	"github.com/dodo/dodo-config/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/command"
	"github.com/oclaussen/dodo/pkg/plugin"
	"github.com/oclaussen/go-gimme/configfiles"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func RegisterPlugin() {
	plugin.RegisterPluginServer(
		command.PluginType,
		&command.Plugin{Impl: &Commands{}},
	)
}

type Commands struct {
	cmds map[string]*cobra.Command
}

func (p *Commands) GetCommands() (map[string]*cobra.Command, error) {
	if len(p.cmds) == 0 {
		p.cmds = map[string]*cobra.Command{
			"list":     NewListCommand(),
			"validate": NewValidateCommand(),
		}
	}
	return p.cmds, nil
}

func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "list",
		Short:                 "List available backdrop configurations",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ptr, decode := configuration.NewConfig()
			configfiles.GimmeConfigFiles(&configfiles.Options{
				Name:                      "dodo",
				Extensions:                []string{"yaml", "yml", "json"},
				IncludeWorkingDirectories: true,
				Filter: func(configFile *configfiles.ConfigFile) bool {
					s := decoder.New(configFile.Path)
					s.LoadYaml(decode, configFile.Content)
                                        return false
				},
			})

			// TODO: wtf this cast
                        config := *(ptr.(**configuration.Config))
			list("", config)
			return nil
		},
	}
}

func list(prefix string, config *configuration.Config) {
	if config.Backdrops != nil {
		for name, _ := range config.Backdrops {
			// TODO filename
			if len(prefix) > 0 {
				log.WithFields(log.Fields{"group": prefix}).Info(name)
			} else {
				log.Info(name)
			}
		}
	}
	if config.Groups != nil {
		for name, group := range config.Groups {
			list(fmt.Sprintf("%s/", name), group)
		}
	}
}

func NewValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "validate",
		Short:                 "Validate configuration files for syntax errors",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, decode := configuration.NewConfig()
			configfiles.GimmeConfigFiles(&configfiles.Options{
				FileGlobs:        args,
				UseFileGlobsOnly: true,
				Filter: func(configFile *configfiles.ConfigFile) bool {
					s := decoder.New(configFile.Path)
					s.LoadYaml(decode, configFile.Content)
					for _, err := range s.Errors() {
						log.Warn(err)
					}
					return false
				},
			})
			return nil
		},
	}
}
