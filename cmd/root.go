package cmd

import (
	"fmt"
	"github.com/bthuilot/git-scanner/pkg/git"
	gogit "github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	repoURL  string
	repoPath string
	keepRefs bool

	logLevel string
	bare     bool
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

func getGitRepository() (*gogit.Repository, string, error) {
	var (
		r   *gogit.Repository
		dir string = repoPath
		err error
	)
	if repoURL != "" {
		r, dir, err = git.CloneRepository(repoURL, bare)
		if err != nil {
			return nil, "", err
		}
		logrus.Infof("Cloned repo: %s", repoURL)
	} else {
		r, err = git.ImportRepository(repoPath)
		if err != nil {
			return nil, "", err
		}
		logrus.Infof("Using existing repo: %s", repoPath)
	}
	return r, dir, nil
}
