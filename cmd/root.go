package cmd

import (
	"fmt"
	"github.com/bthuilot/git-lost-and-found/pkg/git"
	gogit "github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
)

var (
	/* Flags */
	// repoURL is the URL of the git repository to scan
	repoURL string
	// repoPath is the path to the git repository to scan
	repoPath string
	// bare is a flag to clone or import the repository as a bare repository
	bare bool
	// keepRefs is a flag to keep refs created for dangling commits
	keepRefs bool
	// cleanup is a flag to remove the cloned repo after scanning
	// NOTE: only valid when --repo-url is set
	cleanup bool
	// logLevel is the log level for the application
	logLevel string
)

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
	scanCmd.PersistentFlags().BoolVarP(&bare, "bare", "b", true, "clone or import the repository as a bare repository")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "log level (debug, info, warn, error, fatal, panic)")
	rootCmd.PersistentFlags().StringVarP(&repoURL, "repo-url", "r", "", "URL of the git repository to scan")
	rootCmd.PersistentFlags().StringVarP(&repoPath, "repo-path", "p", "", "Path to the git repository to scan")
	rootCmd.PersistentFlags().BoolVarP(&keepRefs, "keep-refs", "k", false, "Keep refs created for dangling commits")
	rootCmd.PersistentFlags().BoolVarP(&cleanup, "cleanup", "c", false, "Remove the cloned repository after scanning")
	_ = rootCmd.MarkPersistentFlagFilename("repo-path")
	rootCmd.MarkFlagsMutuallyExclusive("repo-url", "repo-path")
	rootCmd.MarkFlagsOneRequired("repo-url", "repo-path")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
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
