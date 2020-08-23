package main

import (
	"os"

	"github.com/dodo-cli/dodo-config/plugin"
)

func main() {
	os.Exit(plugin.RunMe())
}
