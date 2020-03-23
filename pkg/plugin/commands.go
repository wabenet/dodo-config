package plugin

import (
	"github.com/dodo/dodo-config/pkg/decoder"
	cfgtypes "github.com/dodo/dodo-config/pkg/types"
	"github.com/oclaussen/dodo/pkg/types"
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
			backdrops := map[string]*types.Backdrop{}
			configfiles.GimmeConfigFiles(&configfiles.Options{
				Name:                      "dodo",
				Extensions:                []string{"yaml", "yml", "json"},
				IncludeWorkingDirectories: true,
				Filter: func(configFile *configfiles.ConfigFile) bool {
					s := decoder.New(configFile.Path)
					s.DecodeYaml(configFile.Content, &backdrops, map[string]decoder.Decoder{
						"backdrops": decoder.Map(cfgtypes.NewBackdrop(), &backdrops),
					})
					return false
				},
			})

			for name, _ := range backdrops {
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
			backdrops := map[string]*types.Backdrop{}
			configfiles.GimmeConfigFiles(&configfiles.Options{
				FileGlobs:        args,
				UseFileGlobsOnly: true,
				Filter: func(configFile *configfiles.ConfigFile) bool {
					s := decoder.New(configFile.Path)
					s.DecodeYaml(configFile.Content, &backdrops, map[string]decoder.Decoder{
						"backdrops": decoder.Map(cfgtypes.NewBackdrop(), &backdrops),
					})
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
