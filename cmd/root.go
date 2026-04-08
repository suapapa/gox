package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/suapapa/gox/pkg/runner"
)

// Execute runs the root command.
func Execute() error {
	var generate bool
	var update bool
	var verbose bool

	var rootCmd = &cobra.Command{
		Use:   "gox [gox-flags] <package>[@version] [args...]",
		Short: "Go version of npx",
		Long: `gox (Go Execute) is a tool that allows you to run Go packages as commands without having to install them globally. 
It automatically downloads the specified package, compiles it into a local cache, and executes it.`,
		Example: `  gox suapapa/gox --help
  gox github.com/suapapa/gox --help
  gox -v golang.org/x/tools/cmd/goimports@latest -w main.go
  gox --generate foo/bar@v1.2.3
  gox -u golang.org/x/tools/cmd/goimports@latest -w main.go`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pkgWithVersion := args[0]
			pkgArgs := args[1:]

			// Initialize the runner with default options.
			r, err := runner.NewRunner(
				runner.WithGenerate(generate),
				runner.WithUpdate(update),
				runner.WithVerbose(verbose),
			)
			if err != nil {
				return fmt.Errorf("initialization error: %w", err)
			}

			// Execute using background context.
			err = r.Run(context.Background(), pkgWithVersion, pkgArgs)
			if err != nil {
				// If the command itself failed, don't show the help message.
				cmd.SilenceUsage = true
			}
			return err
		},
	}

	rootCmd.Flags().BoolVarP(&generate, "generate", "g", false, "run go generate before building")
	rootCmd.Flags().BoolVarP(&update, "update", "u", false, "force reinstall/update even if the package is already cached")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "print version, cache path, binary path, and install/update steps")
	rootCmd.Flags().SetInterspersed(false)

	return rootCmd.Execute()
}
