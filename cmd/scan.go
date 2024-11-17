package cmd

import (
	"fmt"
	"github.com/bthuilot/git-scanner/pkg/git"
	"github.com/bthuilot/git-scanner/pkg/scanning"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Retrieve hanging commits and then run a scanner",
	Long: `Retrieve all dangling commits in a git repository and then run a given program in the directory before cleaning up.
The scanner command must be separated from the git-scanner command with '--'.
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
