package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/fixi2/InfraTrack/internal/cli"
)

type exitCoder interface {
	ExitCode() int
}

func main() {
	root, err := cli.NewRootCommand()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)

		var ec exitCoder
		if errors.As(err, &ec) {
			os.Exit(ec.ExitCode())
		}

		os.Exit(1)
	}
}
