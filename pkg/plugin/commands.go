package plugin

import (
	"github.com/dodo/dodo-config/pkg/decoder"
	"github.com/oclaussen/go-gimme/configfiles"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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
			ptr, decode := NewConfig()
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
			config := *(ptr.(**Config))
			for name, _ := range config.Backdrops {
				// TODO filename
				log.Info(name)
			}
			return nil
		},
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
			_, decode := NewConfig()
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
