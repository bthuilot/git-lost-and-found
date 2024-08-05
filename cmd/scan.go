package cmd

import (
	"github.com/bthuilot/git-scanner/pkg/cli"
	"github.com/bthuilot/git-scanner/pkg/processor"
	"github.com/bthuilot/git-scanner/pkg/retrieval"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	// Used for flags.
	repoURL    string
	repoPath   string
	outputPath string
)

func init() {
	scanCmd.PersistentFlags().StringVar(&repoURL, "repo-url", "", "URL of the git repository to scan")
	scanCmd.PersistentFlags().StringVar(&repoPath, "repo-path", "", "Path to the git repository to scan")
	scanCmd.PersistentFlags().StringVar(&outputPath, "output-path", "", "Path to the output directory")
	_ = scanCmd.MarkPersistentFlagFilename("repo-path")
	scanCmd.MarkFlagsMutuallyExclusive("repo-url", "repo-path")
	scanCmd.MarkFlagsOneRequired("repo-url", "repo-path")
	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan all commits of a git repository",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			output *os.File = os.Stdin
			err    error
		)
		if outputPath != "" && outputPath != "-" {
			output, err = os.Create(outputPath)
			if err != nil {
				logrus.Error(err)
				cli.ErrorExit(err)
			}
		}

		// clone repo or import existing repo
		r, err := getGitRepository()
		if err != nil {
			logrus.Error(err)
			cli.ErrorExit(err)
		}

		commits, err := retrieval.LookupAllCommits(r)
		if err != nil {
			logrus.Error(err)
			cli.ErrorExit(err)
		}
		logrus.Info("Retrieved all commits")

		blobCache := make(map[plumbing.Hash]processor.BlobInfo)
		for _, commit := range commits {
			err = processor.ProcessCommit(commit, output, blobCache)
			if err != nil {
				logrus.Error(err)
				cli.ErrorExit(err)
			}
		}
	},
}

func getGitRepository() (*git.Repository, error) {
	var (
		r   *git.Repository
		err error
	)
	if repoURL != "" {
		r, err = retrieval.CloneRepo(repoURL)
		if err != nil {
			return nil, err
		}
		logrus.Infof("Cloned repo: %s", repoURL)
	} else {
		r, err = retrieval.ExistingRepo(repoPath)
		if err != nil {
			return nil, err
		}
		logrus.Infof("Using existing repo: %s", repoPath)
	}
	return r, nil
}
