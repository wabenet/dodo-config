package command

import (
	"fmt"

	log "github.com/hashicorp/go-hclog"
	"github.com/oclaussen/dodo/pkg/decoder"
	"github.com/oclaussen/dodo/pkg/types"
	"github.com/oclaussen/go-gimme/configfiles"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Config plugin sub commands",
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
			backdrops := map[string]*types.Backdrop{}
			_, err := configfiles.GimmeConfigFiles(&configfiles.Options{
				Name:                      "dodo",
				Extensions:                []string{"yaml", "yml", "json"},
				IncludeWorkingDirectories: true,
				Filter: func(configFile *configfiles.ConfigFile) bool {
					d := decoder.New(configFile.Path)
					d.DecodeYaml(configFile.Content, &backdrops, map[string]decoder.Decoding{
						"backdrops": decoder.Map(types.NewBackdrop(), &backdrops),
					})

					return false
				},
			})

			if err != nil {
				log.L().Error("error finding config files", "error", err)
			}

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
			backdrops := map[string]*types.Backdrop{}
			_, err := configfiles.GimmeConfigFiles(&configfiles.Options{
				FileGlobs:        args,
				UseFileGlobsOnly: true,
				Filter: func(configFile *configfiles.ConfigFile) bool {
					d := decoder.New(configFile.Path)
					d.DecodeYaml(configFile.Content, &backdrops, map[string]decoder.Decoding{
						"backdrops": decoder.Map(types.NewBackdrop(), &backdrops),
					})

					for _, err := range d.Errors() {
						fmt.Println(err)
					}

					return false
				},
			})

			if err != nil {
				log.L().Error("error finding config files", "error", err)
			}

			return nil
		},
	}
}
