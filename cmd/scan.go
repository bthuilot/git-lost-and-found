package cmd

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/bthuilot/git-scanner/pkg/cli"
	"github.com/bthuilot/git-scanner/pkg/processor"
	"github.com/bthuilot/git-scanner/pkg/reporter"
	"github.com/bthuilot/git-scanner/pkg/retrieval"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
			err error
		)

		if outputPath != "" && outputPath != "-" {
			if repoURL != "" {
				var suffix string
				if strings.Contains(repoURL, "https://github.com") {
					split := strings.Split(repoURL, "/")
					suffix = split[len(split)-2] + "/" + split[len(split)-1]
				} else {
					suffix = ""
				}
				outputPath = outputPath + "/" + suffix
			}
			err := os.MkdirAll(outputPath, os.ModePerm)
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

		// filter commits

		blobCache := make(map[plumbing.Hash]processor.BlobInfo)
		reportsByCommit := make(map[string][]processor.SecretsReport)

		for _, commit := range commits {
			reports, err := processor.ProcessCommit(commit, blobCache)
			if err != nil {
				logrus.Error(err)
				cli.ErrorExit(err)
			}
			reportsByCommit[commit.Hash.String()] = reports
		}

		// create a file system where the folder is the commit and then write the reportjobs blob to the file system
		for commitHash, reports := range reportsByCommit {
			for _, report := range reports {
				if report.FileName != "" {
					err = os.MkdirAll(outputPath+"/rawblobs/"+commitHash, os.ModePerm)
					if err != nil {
						logrus.Errorf("unable to make child commit dir %s", err)
						cli.ErrorExit(err)
					}
					cleanFileName := strings.ReplaceAll(report.FileName, "/", "_")
					file, err := os.Create(outputPath + "/rawblobs/" + commitHash + "/" + cleanFileName)
					if err != nil {
						logrus.Error(err)
						cli.ErrorExit(err)
					}
					_, err = file.Write([]byte(report.RawBlob))
					if err != nil {
						logrus.Error(err)
						cli.ErrorExit(err)
					}
				}
			}
		}

		// run gitleaks
		stdOut, _, err := processor.RunGitleaksScanFileSystem(outputPath + "/rawblobs")
		if err != nil {
			logrus.Errorf("Failed to run Gitleaks scan: %v", err)
			err = os.RemoveAll(outputPath)
			if err != nil {
				logrus.Error(err)
			}
			return
		}

		var results []processor.GitleaksResult
		err = json.Unmarshal([]byte(stdOut), &results)
		if err != nil {
			logrus.Errorf("Failed to unmarshal Gitleaks result: %v", err)
			cli.ErrorExit(err)
		}

		if len(results) == 0 {
			// delete the output directory if no results
			err = os.RemoveAll(outputPath)
			if err != nil {
				logrus.Error(err)
				cli.ErrorExit(err)
			}
			logrus.Info("No results found")
			return
		}

		// hydrate the results with the raw blob
		for i, result := range results {
			file, err := os.ReadFile(outputPath + "/rawblobs/" + result.File)
			if err != nil {
				logrus.Error(err)
				cli.ErrorExit(err)
			}
			results[i].RawBlob = string(file)
		}

		var resFile *os.File
		if outputPath != "" {
			resFile, err = os.Create(outputPath + "/results.json")
		} else {
			resFile, err = os.Create("results.json")
		}
		if err != nil {
			logrus.Error(err)
			cli.ErrorExit(err)
		}

		err = reporter.WriteReport(resFile, results)
		if err != nil {
			logrus.Error(err)
			cli.ErrorExit(err)
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
