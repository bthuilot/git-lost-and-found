package git

import (
	"os"

	"github.com/go-git/go-git/v5"
)

func CloneRepository(repoURL string, bare bool) (*git.Repository, string, error) {
	// CLone to temp directory
	tmpDir, err := os.MkdirTemp("", "*")
	if err != nil {
		return nil, "", err
	}
	repo, err := git.PlainClone(tmpDir, bare, &git.CloneOptions{
		URL:    repoURL,
		Mirror: bare,
	})
	return repo, tmpDir, err
}

func ImportRepository(repoPath string) (*git.Repository, error) {
	return git.PlainOpen(repoPath)
}
