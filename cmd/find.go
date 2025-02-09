package cmd

import (
	"fmt"
	"github.com/bthuilot/git-lost-and-found/pkg/git"
	"github.com/bthuilot/git-lost-and-found/pkg/scanning"
	gogit "github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
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
)

func init() {
	findCmd.Flags().BoolVarP(&bare, "bare", "b", true, "clone or import the repository as a bare repository")
	// findCmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "log level (debug, info, warn, error, fatal, panic)")
	// findCmd.Flags().StringVar(&logFormat, "log-format", "text", "log format (text, json)")
	findCmd.Flags().StringVarP(&repoURL, "repo-url", "r", "", "URL of the git repository to scan")
	findCmd.Flags().StringVarP(&repoPath, "repo-path", "p", "", "Path to the git repository to scan")
	findCmd.Flags().BoolVarP(&keepRefs, "keep-refs", "k", false, "Keep refs created for dangling commits")
	findCmd.Flags().BoolVarP(&cleanup, "cleanup", "c", false, "Remove the cloned repository after scanning")
	_ = findCmd.MarkFlagFilename("repo-path")
	findCmd.MarkFlagsMutuallyExclusive("repo-url", "repo-path")
	findCmd.MarkFlagsOneRequired("repo-url", "repo-path")

	rootCmd.AddCommand(findCmd)
}

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find all hanging commits, reference them and then run a scanner",
	Long: `Retrieve all dangling commits in a git repository and then run a given program in the directory before cleaning up.
The scanner command must be separated from the git-lost-and-found command with '--'.
The command will be executed in the repository directory, and any '{}' will be replaced with the directory path in the command.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		logrus.Info("beginning scan")
		r, dir, cleanup, err := getGitRepository()
		if err != nil {
			return err
		}
		defer cleanup()

		logrus.WithField("repository_directory", dir).Info("Scanning repository")

		// TODO: additional support scanning for blobs
		danglingObjs, err := git.FindDanglingObjects(r, dir)
		if err != nil {
			return err
		}
		logrus.WithField("dangling_commits_amt", len(danglingObjs.Commits)).Info("dangling commits")

		var createdRefs []string
		for _, c := range danglingObjs.Commits {
			logrus.Debugf("Dangling commit: %s", c.Hash.String())
			ref := fmt.Sprintf("refs/dangling/%s", c.Hash.String())
			if err = git.MakeRef(r, ref, c); err != nil {
				logrus.Warnf("Failed to create ref for dangling commit: %s", c.Hash.String())
				continue
			}
			createdRefs = append(createdRefs, ref)
		}

		logrus.WithField("created_refs_amt", len(createdRefs)).Info("created refs for dangling commits")
		if !keepRefs {
			defer func() {
				removeErr := git.RemoveReferences(r, createdRefs)
				if removeErr != nil {
					logrus.Errorf("Failed to remove created refs: %s", removeErr)
				}
			}()
		}

		logrus.Debug("Executing scanner")
		if err = scanning.ExecScanner(dir, args); err != nil {
			return err
		}

		return nil
	},
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
