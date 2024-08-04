package retrieval

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func RetrieveAllCommits() ([]object.Commit, error) {

	repo, err := git.PlainOpen("/tmp/DVWA")

	if err != nil {
		return nil, err
	}

	ref, err := repo.Head()

	if err != nil {
		return nil, err
	}

	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash()})

	if err != nil {
		return nil, err
	}

	var commits []object.Commit

	err = commitIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, *c)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return commits, nil
}
