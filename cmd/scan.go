package cmd

import (
	"os"

	"github.com/bthuilot/git-scanner/pkg/cli"
	"github.com/bthuilot/git-scanner/pkg/processor"
	"github.com/bthuilot/git-scanner/pkg/retrieval"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	repoURL       string
	repoPath      string
	outputPath    string
	gitleakConfig string
)

func init() {
	scanCmd.PersistentFlags().StringVar(&repoURL, "repo-url", "", "URL of the git repository to scan")
	scanCmd.PersistentFlags().StringVar(&repoPath, "repo-path", "", "Path to the git repository to scan")
	scanCmd.PersistentFlags().StringVar(&outputPath, "output", "", "Path to the output directory")
	scanCmd.PersistentFlags().StringVar(&gitleakConfig, "gitleaks-config", "", "Path to the gitleaks config file")
	_ = scanCmd.MarkPersistentFlagFilename("repo-path")
	_ = scanCmd.MarkPersistentFlagFilename("gitleaks-config")
	scanCmd.MarkFlagsMutuallyExclusive("repo-url", "repo-path")
	scanCmd.MarkFlagsOneRequired("repo-url", "repo-path")
	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scan all commits of a git repository",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			output  = os.Stdout
			err     error
			results []processor.GitleaksResult
		)
		if outputPath != "" && outputPath != "-" {
			output, err = os.Create(outputPath)
			defer output.Close()
			if err != nil {
				logrus.Error(err)
				cli.ErrorExit(err)
			}
		}

		r, err := getGitRepository()
		if err != nil {
			logrus.Error(err)
			cli.ErrorExit(err)
		}

		logrus.Debug("Cloned or imported repository")

		commits, err := retrieval.LookupAllCommits(r)
		if err != nil {
			logrus.Error(err)
			cli.ErrorExit(err)
		}
		logrus.Infof("Gathered %d commits for repo %s", len(commits), repoURL)

		// Process all commits and gather results
		blobCache := make(map[plumbing.Hash]struct{})
		uniqueSecrets := make(map[string]processor.GitleaksResult)
		for _, commit := range commits {
			logrus.Debug("Processing commit: ", commit.Hash)
			commitResults, err := processor.ProcessCommit(commit, blobCache, processor.GitleaksArgs{Config: gitleakConfig})
			if err != nil {
				logrus.Errorf("error processing commit %s: %s", commit.String(), err)
				continue
			}
			for _, result := range commitResults {
				if _, exists := uniqueSecrets[result.Secret]; !exists {
					results = append(results, result)
					uniqueSecrets[result.Secret] = result
				}
			}
		}

		logrus.Infof("Processed all commits and found %d secrets for repo %s", len(results), repoURL)

		if len(results) > 0 {
			logrus.Infof("Writing results to %s", outputPath)
			if err := cli.WriteResults(output, uniqueSecrets); err != nil {
				logrus.Error(err)
				cli.ErrorExit(err)
			}
		} else {
			if _, err := os.Stat(outputPath); err == nil {
				if err := os.Remove(outputPath); err != nil {
					logrus.Error(err)
					cli.ErrorExit(err)
				}
			}
		}
		logrus.Info("Scan completed successfully")
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
