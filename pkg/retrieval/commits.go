package retrieval

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func LookupAllCommits(r *git.Repository) ([]*object.Commit, error) {
	logIter, err := r.Log(&git.LogOptions{
		All: true,
	})

	if err != nil {
		return nil, err
	}

	// set of commit hashes found via standard log
	foundCommits := make(map[string]struct{})
	err = logIter.ForEach(func(commit *object.Commit) error {
		foundCommits[commit.Hash.String()] = struct{}{}
		return nil
	})

	cIter, err := r.CommitObjects()
	if err != nil {
		return nil, err
	}
	var shadowCommits []*object.Commit
	err = cIter.ForEach(func(commit *object.Commit) error {
		if _, ok := foundCommits[commit.Hash.String()]; !ok {
			shadowCommits = append(shadowCommits, commit)
		}
		return nil
	})

	return shadowCommits, err
}
