package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func MakeRef(r *git.Repository, name string, commit *object.Commit) error {
	ref := plumbing.NewHashReference(plumbing.ReferenceName(name), commit.Hash)
	return r.Storer.SetReference(ref)
}

func RemoveReferences(r *git.Repository, names []string) error {
	for _, name := range names {
		err := r.Storer.RemoveReference(plumbing.ReferenceName(name))
		if err != nil {
			return err
		}
	}
	return nil
}
