package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"cuelang.org/go/cue/errors"
	"github.com/dodo-cli/dodo-config/internal/config"
	api "github.com/dodo-cli/dodo-core/api/v1alpha2"
	core "github.com/dodo-cli/dodo-core/pkg/config"
	"github.com/dodo-cli/dodo-core/pkg/plugin"
	"github.com/dodo-cli/dodo-core/pkg/plugin/command"
	log "github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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
			backdrops, err := config.GetAllBackdrops(core.GetConfigFiles()...)
			if err != nil {
				log.L().Error(err.Error())
			}

			p := getPrettyPrinter()

			for name := range backdrops {
				p.Fprintf(os.Stdout, "%s\n", name) // TODO filename
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
			p := getPrettyPrinter()

			_, err := config.GetAllBackdrops(args...)
			if err == nil {
				p.Fprintf(os.Stdout, "configuration is valid!\n")

				return nil
			}

			cwd, _ := os.Getwd()
			w := &bytes.Buffer{}

			errors.Print(w, err, &errors.Config{
				Format: func(w io.Writer, format string, args ...interface{}) {
					p.Fprintf(w, format, args...)
				},
				Cwd:     cwd,
				ToSlash: false,
			})

			fmt.Fprintf(os.Stdout, string(w.Bytes()))

			return nil
		},
	}
}

func getPrettyPrinter() *message.Printer {
	loc := os.Getenv("LC_ALL")
	if loc == "" {
		loc = os.Getenv("LANG")
	}
	loc = strings.Split(loc, ".")[0]

	return message.NewPrinter(language.Make(loc))
}
