package main

import (
	"github.com/dodo/dodo-config/pkg/configuration"
	"github.com/hashicorp/go-plugin"
	dodo "github.com/oclaussen/dodo/pkg/plugin"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(new(log.JSONFormatter))
	plugin.Serve(&plugin.ServeConfig{
		GRPCServer:      plugin.DefaultGRPCServer,
		HandshakeConfig: dodo.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			dodo.Configuration: configuration.NewPlugin(),
		},
	})
}
