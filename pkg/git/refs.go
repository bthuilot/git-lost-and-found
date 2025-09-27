// Copyright (C) 2025 Bryce Thuilot
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the FSF, either version 3 of the License, or (at your option) any later version.
// See the LICENSE file in the root of this repository for full license text or
// visit: <https://www.gnu.org/licenses/gpl-3.0.html>.

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
