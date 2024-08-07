package cmd

import (
	"fmt"
	"github.com/bthuilot/git-scanner/pkg/cli"
	"github.com/bthuilot/git-scanner/pkg/git"
	"github.com/bthuilot/git-scanner/pkg/scanning"
	gogit "github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var (
	// Used for flags.
	repoURL       string
	repoPath      string
	outputPath    string
	gitleakConfig string
	keepRefs      bool
)

func init() {
	scanCmd.PersistentFlags().StringVarP(&repoURL, "repo-url", "r", "", "URL of the git repository to scan")
	scanCmd.PersistentFlags().StringVarP(&repoPath, "repo-path", "p", "", "Path to the git repository to scan")
	scanCmd.PersistentFlags().StringVarP(&outputPath, "output", "o", "", "Path to the output directory")
	scanCmd.PersistentFlags().StringVarP(&gitleakConfig, "gitleaks-config", "c", "", "Path to the gitleaks config file")
	scanCmd.PersistentFlags().BoolVarP(&keepRefs, "keep-refs", "k", false, "Keep refs created for dangling commits")
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
			output = os.Stdout
			err    error
		)
		if outputPath != "" && outputPath != "-" {
			output, err = os.Create(outputPath)
			defer output.Close()
			if err != nil {
				cli.ErrorExit(err)
			}
		}

		// clone repo or import existing repo
		r, dir, err := getGitRepository()
		if err != nil {
			cli.ErrorExit(err)
		}

		logrus.Debug("Cloned or imported repository")

		danglingObjs, err := git.FindDanglingObjects(r, dir)
		if err != nil {
			cli.ErrorExit(err)
		}

		var createdRefs []string
		for _, c := range danglingObjs.Commits {
			logrus.Infof("Dangling commit: %s", c.Hash.String())
			ref := fmt.Sprintf("refs/dangling/%s", c.Hash.String())
			if err = git.MakeRef(r, ref, c); err != nil {
				logrus.Warnf("Failed to create ref for dangling commit: %s", c.Hash.String())
				continue
			}
			createdRefs = append(createdRefs, ref)
		}

		var gitleaksArgs []string
		if gitleakConfig != "" {
			gitleaksArgs = append(gitleaksArgs, "-c", gitleakConfig)
		}

		if err = scanning.RunGitleaks(dir, outputPath, gitleaksArgs...); err != nil {
			cli.ErrorExit(err)
		}

		if !keepRefs {
			if err = git.RemoveReferences(r, createdRefs); err != nil {
				cli.ErrorExit(err)
			}
		}

	},
}

func getGitRepository() (*gogit.Repository, string, error) {
	var (
		r   *gogit.Repository
		dir string = repoPath
		err error
	)
	if repoURL != "" {
		r, dir, err = git.CloneRepo(repoURL)
		if err != nil {
			return nil, "", err
		}
		logrus.Infof("Cloned repo: %s", repoURL)
	} else {
		r, err = git.ExistingRepo(repoPath)
		if err != nil {
			return nil, "", err
		}
		logrus.Infof("Using existing repo: %s", repoPath)
	}
	return r, dir, nil
}
