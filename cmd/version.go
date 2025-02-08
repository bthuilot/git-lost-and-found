package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Version string = "dev"
	Commit  string = "N/A"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of git-lost-and-found",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		fmt.Printf("git-lost-and-found %s (%s)\n", Version, Commit)
		return nil
	},
}
