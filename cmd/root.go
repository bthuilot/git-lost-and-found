package cmd

import (
	"fmt"
	"os"

	"github.com/bthuilot/git-scanner/pkg/retrieval"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-scanner",
	Short: "git-scanner searches *all* commits in a git repository for credentials",
	Long:  `git-scanner searches *all* commits in a git repository for credentials, including orphaned commits.`,
	Run: func(cmd *cobra.Command, args []string) {
		commits, err := retrieval.RetrieveAllCommits()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		logrus.Info("Retrieved all commits")
		for _, commit := range commits {
			logrus.Infof("Commit: %s", commit.Hash)
		}

		// TODO: scan each commit
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
