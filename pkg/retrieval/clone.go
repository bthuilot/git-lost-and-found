package retrieval

import (
	"os"

	"github.com/go-git/go-git/v5"
)

func CloneRepo(repoURL string) (*git.Repository, error) {
	// CLone to temp directory
	tmpDir, err := os.MkdirTemp("", "*")
	if err != nil {
		return nil, err
	}
	return git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL: repoURL,
		// Mirror: true,
	})
}

func ExistingRepo(repoPath string) (*git.Repository, error) {
	return git.PlainOpen(repoPath)
}
