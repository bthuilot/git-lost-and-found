package git

import (
	"os"

	"github.com/go-git/go-git/v5"
)

func CloneRepo(repoURL string) (*git.Repository, string, error) {
	// CLone to temp directory
	tmpDir, err := os.MkdirTemp("", "*")
	if err != nil {
		return nil, "", err
	}
	repo, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL:    repoURL,
		Mirror: true,
	})
	return repo, tmpDir, err
}

func ExistingRepo(repoPath string) (*git.Repository, error) {
	return git.PlainOpen(repoPath)
}
