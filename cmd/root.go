package cmd

import (
	"fmt"
	"github.com/bthuilot/git-lost-and-found/pkg/git"
	gogit "github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "git-lost-and-found",
	Short: "git-lost-and-found will find all dangling commits in a git repository and create refs for them.",
	Long: `git-lost-and-found will find all dangling commits in a git repository and create refs for them.
This allows for scanners that use 'git log' to search blob data to not miss any changes.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			return err
		}
		logrus.SetLevel(level)

		return nil
	},
}

func init() {
	rootCmd.SetErrPrefix("ERROR: ")
}

func getGitRepository() (*gogit.Repository, string, func(), error) {
	var (
		r   *gogit.Repository
		dir string = repoPath
		err error
	)
	cleanupF := func() {}
	if repoURL != "" {
		r, dir, err = git.CloneRepository(repoURL, bare)
		if err != nil {
			return nil, "", cleanupF, err
		}
		logrus.Infof("Cloned repo: %s", repoURL)
		cleanupF = func() {
			if cleanup {
				logrus.Debug("cleaning up cloned repo")
				if err := os.RemoveAll(dir); err != nil {
					logrus.Errorf("Failed to remove cloned repo: %s", err)
				}
			}
		}
	} else {
		r, err = git.ImportRepository(repoPath)
		if err != nil {
			return nil, "", cleanupF, err
		}
		logrus.Infof("Using existing repo: %s", repoPath)
	}
	return r, dir, cleanupF, nil
}
