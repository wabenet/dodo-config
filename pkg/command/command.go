package command

import (
	"fmt"

	"github.com/dodo-cli/dodo-config/pkg/config"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/command"
	log "github.com/hashicorp/go-hclog"
	"github.com/oclaussen/go-gimme/configfiles"
	"github.com/spf13/cobra"
)

const name = "config"

var _ command.Command = &Command{}

type Command struct {
	cmd *cobra.Command
}

func New() *Command {
	return &Command{cmd: NewCommand()}
}

func (p *Command) Type() plugin.Type {
	return command.Type
}

func (p *Command) PluginInfo() *api.PluginInfo {
	return &api.PluginInfo{
		Name: &api.PluginName{Name: name, Type: command.Type.String()},
	}
}

func (*Command) Init() (plugin.PluginConfig, error) {
	return map[string]string{}, nil
}

func (p *Command) GetCobraCommand() *cobra.Command {
	return p.cmd
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: "Config plugin subcommands",
	}

	cmd.AddCommand(NewListCommand())
	cmd.AddCommand(NewValidateCommand())

	return cmd
}

func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:                   "list",
		Short:                 "List available backdrop configurations",
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE: func(cmd *cobra.Command, args []string) error {
			backdrops := map[string]*api.Backdrop{}
			configfiles.GimmeConfigFiles(&configfiles.Options{
				Name:                      "dodo",
				Extensions:                []string{"yaml", "yml", "json"},
				IncludeWorkingDirectories: true,
				Filter: func(configFile *configfiles.ConfigFile) bool {
					backdrops, err := config.ParseConfig(configFile.Path)
					if err != nil {
						log.L().Error(err.Error())
					}

					for name, backdrop := range backdrops {
						backdrops[name] = backdrop // TODO: check for duplicates
					}

					return false
				},
			})

			for name := range backdrops {
				// TODO filename
				fmt.Println(name)
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
			configfiles.GimmeConfigFiles(&configfiles.Options{
				FileGlobs:        args,
				UseFileGlobsOnly: true,
				Filter: func(configFile *configfiles.ConfigFile) bool {
					if _, err := config.ParseConfig(configFile.Path); err != nil {
						fmt.Println(err)
					}

					return false
				},
			})

			return nil
		},
	}
}
