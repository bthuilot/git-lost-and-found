package cmd

import (
	"os"

	"github.com/bthuilot/git-scanner/pkg/cli"
	"github.com/bthuilot/git-scanner/pkg/processor"
	"github.com/bthuilot/git-scanner/pkg/retrieval"
	"github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	repoURL       string
	repoPath      string
	outputPath    string
	gitleakConfig string
	numWorkers    int = 4 // Number of concurrent workers, can be configured
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
		logrus.Infof("Scanning repository %s", repoURL)

		// clone repo or import existing repo
		r, err := getGitRepository()
		if err != nil {
			logrus.Error(err)
			cli.ErrorExit(err)
		}

		logrus.Info("Cloned or imported repository")

		commits, err := retrieval.LookupAllCommits(r)
		logrus.Infof("Retrieved %d commits", len(commits))
		for _, commit := range commits {
			logrus.Infof("Commit: %s", commit.Hash)
		}
		os.Exit(0)
		if err != nil {
			logrus.Error(err)
			cli.ErrorExit(err)
		}
		logrus.Info("Retrieved all commits")

		// Process all commits and gather results
		results, err = processor.ProcessCommits(commits, processor.GitleaksArgs{Config: gitleakConfig})
		if err != nil {
			logrus.Error(err)
			cli.ErrorExit(err)
		}

		logrus.Infof("Processed %d commits", len(commits))
		logrus.Infof("Found %d secrets", len(results))

		uniqueSecrets := make(map[string]processor.GitleaksResult)
		for _, result := range results {
			if _, exists := uniqueSecrets[result.Secret]; !exists {
				results = append(results, result)
				uniqueSecrets[result.Secret] = result
			}
		}

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
