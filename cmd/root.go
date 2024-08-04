package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-scanner",
	Short: "git-scanner searches *all* commits in a git repository for credentials",
	Long:  `git-scanner searches *all* commits in a git repository for credentials, including orphaned commits.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
