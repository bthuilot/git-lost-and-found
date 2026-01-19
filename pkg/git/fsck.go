// Copyright (C) 2024-2026 Bryce Thuilot
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the FSF, either version 3 of the License, or (at your option) any later version.
// See the LICENSE file in the root of this repository for full license text or
// visit: <https://www.gnu.org/licenses/gpl-3.0.html>.

package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"go.uber.org/zap"
)

type DanglingObjects struct {
	Blobs   []*object.Blob
	Commits []*object.Commit
	Trees   []*object.Tree
}

func FindDanglingObjects(r *git.Repository, repoPath string) (DanglingObjects, error) {
	var (
		d    DanglingObjects
		buff = bytes.NewBuffer(nil)
	)
	cmd := exec.Command("git", "-C", repoPath, "fsck", "--lost-found")
	cmd.Stdout = buff
	cmd.Stderr = os.Stderr
	scanner := bufio.NewScanner(buff)
	if err := cmd.Run(); err != nil {
		return d, err
	}

	// Parse out objects
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			split := strings.Split(line, " ")
			if len(split) != 3 {
				zap.L().Warn("unexpected output from git fsck", zap.String("git-fsck-output", line))
				continue
			}
			hash := plumbing.NewHash(split[2])
			switch split[1] {
			case "blob":
				obj, err := r.BlobObject(hash)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting blob object: %s\n", err)
					continue
				}
				d.Blobs = append(d.Blobs, obj)
			case "commit":
				obj, err := r.CommitObject(hash)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting commit object: %s\n", err)
					continue
				}
				d.Commits = append(d.Commits, obj)
			case "tree":
				obj, err := r.TreeObject(hash)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting tree object: %s\n", err)
					continue
				}
				d.Trees = append(d.Trees, obj)
			}
		}
	}
	return d, scanner.Err()
}
