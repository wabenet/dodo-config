package main

import (
	"os"

	"github.com/dodo-cli/dodo-config/pkg/plugin"
)

func main() {
	os.Exit(plugin.RunMe())
}
