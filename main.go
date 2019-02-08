package main

import (
	"fmt"
	"os"

	"github.com/redmaner/MicroHTTP/micro"
)

func main() {
	args := os.Args

	if len(args) <= 1 {
		showHelp(args)
	}

	m := micro.NewInstanceFromConfig(args[1])
	m.Serve()

}

func showHelp(args []string) {
	fmt.Printf("MicroHTTP version %s\n\nUsage:\n\n\t%s /path/to/config.json\n\n", micro.Version, args[0])
}
