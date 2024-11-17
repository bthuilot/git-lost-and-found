package cmd

import (
	"fmt"
	"github.com/bthuilot/git-lost-and-found/pkg/git"
	"github.com/bthuilot/git-lost-and-found/pkg/scanning"
	"github.com/sirupsen/logrus"
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

func init() {
	scanCmd.Flags().BoolVarP(&bare, "bare", "b", true, "clone or import the repository as a bare repository")
	scanCmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "log level (debug, info, warn, error, fatal, panic)")
	scanCmd.Flags().StringVarP(&repoURL, "repo-url", "r", "", "URL of the git repository to scan")
	scanCmd.Flags().StringVarP(&repoPath, "repo-path", "p", "", "Path to the git repository to scan")
	scanCmd.Flags().BoolVarP(&keepRefs, "keep-refs", "k", false, "Keep refs created for dangling commits")
	scanCmd.Flags().BoolVarP(&cleanup, "cleanup", "c", false, "Remove the cloned repository after scanning")
	_ = scanCmd.MarkFlagFilename("repo-path")
	scanCmd.MarkFlagsMutuallyExclusive("repo-url", "repo-path")
	scanCmd.MarkFlagsOneRequired("repo-url", "repo-path")

	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Retrieve hanging commits and then run a scanner",
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
