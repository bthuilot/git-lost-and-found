package git

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
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
				logrus.Warn("Unexpected output from git fsck: ", line)
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
