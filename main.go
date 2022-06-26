package main

import (
	"os"

	"github.com/wabenet/dodo-config/pkg/plugin"
)

func main() {
	os.Exit(plugin.RunMe())
}
