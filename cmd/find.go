package cmd

import (
	"fmt"
	"github.com/bthuilot/git-scanner/pkg/git"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(findCmd)
}

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find all dangling commits in a repository",
	Long: `Finds all dangling commits in a git repository and creates refs for them.
The refs are created in the format 'refs/dangling/<commit-hash>'.
If -k is not set, the refs will be removed after the scanner command is executed.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		r, dir, err := getGitRepository()
		if err != nil {
			return err
		}

		logrus.WithField("repository_directory", dir).Info("Scanning repository")

		danglingObjs, err := git.FindDanglingObjects(r, dir)
		if err != nil {
			return err
		}

		// TODO: additional support scanning for blobs
		var createdRefs []string
		logrus.Infof("Found %d dangling commits", len(danglingObjs.Commits))
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
				logrus.WithField("created_refs_amt", len(createdRefs)).Debug("removing created refs")
				if err = git.RemoveReferences(r, createdRefs); err != nil {
					logrus.Errorf("Failed to remove created refs: %s", err)
				}
			}()
		}

		return nil
	},
}
