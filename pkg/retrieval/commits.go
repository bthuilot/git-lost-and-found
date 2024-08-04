package retrieval

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func LookupAllCommits(r *git.Repository) ([]*object.Commit, error) {
	cIter, err := r.CommitObjects()
	if err != nil {
		return nil, err
	}

	var commits []*object.Commit

	err = cIter.ForEach(func(commit *object.Commit) error {
		commits = append(commits, commit)
		return nil
	})
	return commits, err
}
