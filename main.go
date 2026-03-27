package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/suapapa/gox/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
