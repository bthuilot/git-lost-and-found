// Copyright (C) 2024-2026 Bryce Thuilot
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the FSF, either version 3 of the License, or (at your option) any later version.
// See the LICENSE file in the root of this repository for full license text or
// visit: <https://www.gnu.org/licenses/gpl-3.0.html>.

package cmd

import (
	"fmt"
	"os"

	"github.com/bthuilot/git-lost-and-found/v2/pkg/git"
	"github.com/bthuilot/git-lost-and-found/v2/pkg/scanning"
	gogit "github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
	findCmd.Flags().
		BoolVarP(&bare, "bare", "b", true, "clone or import the repository as a bare repository")
	// findCmd.Flags().StringVarP(&logLevel, "log-level", "l", "info", "log level (debug, info, warn, error, fatal, panic)")
	// findCmd.Flags().StringVar(&logFormat, "log-format", "text", "log format (text, json)")
	findCmd.Flags().StringVarP(&repoURL, "repo-url", "r", "", "URL of the git repository to scan")
	findCmd.Flags().
		StringVarP(&repoPath, "repo-path", "p", "", "Path to the git repository to scan")
	findCmd.Flags().
		BoolVarP(&keepRefs, "keep-refs", "k", false, "Keep refs created for dangling commits")
	findCmd.Flags().
		BoolVarP(&cleanup, "cleanup", "c", false, "Remove the cloned repository after scanning")
	_ = findCmd.MarkFlagFilename("repo-path")
	findCmd.MarkFlagsMutuallyExclusive("repo-url", "repo-path")
	findCmd.MarkFlagsOneRequired("repo-url", "repo-path")

	rootCmd.AddCommand(findCmd)
}

var findCmd = &cobra.Command{
	Use:   "find [flags] -- [optional command]",
	Short: "Find all hanging commits and add references to them. Optionally run a command once references are created",
	Long: `Retrieve all dangling commits in a git repository and then run a given program in the directory before cleaning up.
A command can be added to the positional arguments prefixed with '--', that will run once the refences are created and before cleanup.
The command will be executed in the repository directory, and any '{}' will be replaced with the directory path in the command.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		zap.L().Info("beginning scan")
		r, dir, cleanup, err := getGitRepository()
		if err != nil {
			return err
		}
		defer cleanup()

		zap.L().Info("Scanning repository", zap.String("repository-directory", dir))

		// TODO: additional support scanning for blobs
		danglingObjs, err := git.FindDanglingObjects(r, dir)
		if err != nil {
			return err
		}
		zap.L().
			Info("dangling commits", zap.Int("dangling-commits-amount", len(danglingObjs.Commits)))

		var createdRefs []string
		for _, c := range danglingObjs.Commits {
			zap.L().
				Debug("found dangling commit", zap.String("dangling-commit-sha", c.Hash.String()))
			ref := fmt.Sprintf("refs/dangling/%s", c.Hash.String())
			if err = git.MakeRef(r, ref, c); err != nil {
				zap.L().
					Error("failed to create ref for dangling commit", zap.Error(err), zap.String("danling-commit-sha", c.Hash.String()))
				continue
			}
			createdRefs = append(createdRefs, ref)
		}

		zap.L().
			Info("created refs for dangling commits", zap.Int("created-refs-amount", len(createdRefs)))
		if !keepRefs {
			defer func() {
				removeErr := git.RemoveReferences(r, createdRefs)
				if removeErr != nil {
					zap.L().Error("failed to remove created refs", zap.Error(removeErr))
				}
			}()
		}

		zap.L().
			Debug("executing scanner", zap.String("directory", dir), zap.Strings("arguments", args))
		if err = scanning.ExecScanner(dir, args); err != nil {
			return err
		}

		return nil
	},
}

func getGitRepository() (*gogit.Repository, string, func(), error) {
	var (
		r   *gogit.Repository
		dir = repoPath
		err error
	)
	cleanupF := func() {}
	if repoURL != "" {
		r, dir, err = git.CloneRepository(repoURL, bare)
		if err != nil {
			return nil, "", cleanupF, err
		}
		zap.L().Info("cloned repository", zap.String("repository-url", repoURL))
		cleanupF = func() {
			if cleanup {
				zap.L().Debug("cleaning up cloned repo")
				if err := os.RemoveAll(dir); err != nil {
					zap.L().Error("failed to remove cloned repo", zap.Error(err))
				}
			}
		}
	} else {
		r, err = git.ImportRepository(repoPath)
		if err != nil {
			return nil, "", cleanupF, err
		}
		zap.L().Info("Using existing cloned repository", zap.String("repository-path", repoPath))
	}
	return r, dir, cleanupF, nil
}
