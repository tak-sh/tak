package main

import (
	"fmt"
	"github.com/tak-sh/tak/pkg/cli"
	"os"
)

var version string

func main() {
	err := cli.New(version).Run(os.Args)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}
