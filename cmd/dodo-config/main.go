package main

import (
	"os"

	plugin "github.com/wabenet/dodo-config"
)

func main() {
	os.Exit(plugin.RunMe())
}
