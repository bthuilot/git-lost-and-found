package git

import (
	"crypto/rand"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	memoryStorage "github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMakeRef(t *testing.T) {
	type testCase struct {
		name string
		// TODO: eventually test the expected error
		doesErr bool
		refName string
		hashStr string
	}

	tests := []testCase{
		{
			name:    "valid hash",
			refName: "refs/heads/test",
			hashStr: "1234567890123456789012345678901234567890",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gitRepo := &git.Repository{
				Storer: memoryStorage.NewStorage(),
			}

			treeHashStr := make([]byte, 40)
			_, _ = rand.Read(treeHashStr)

			commit := &object.Commit{
				Hash: plumbing.NewHash(test.hashStr),
				Author: object.Signature{
					Name:  "test",
					Email: "test@test.com",
					When:  time.Now(),
				},
				Committer: object.Signature{
					Name:  "test",
					Email: "test@test.com",
					When:  time.Now(),
				},
				Message:      "this is my test commit message",
				TreeHash:     plumbing.NewHash(string(treeHashStr)),
				ParentHashes: nil,
				Encoding:     "UTF-8",
			}
			err := MakeRef(gitRepo, test.refName, commit)
			if test.doesErr {
				assert.Errorf(t, err, "expected an error while creating ref")
				return
			} else {
				assert.NoError(t, err, "expected no error while creating ref")
			}

			// check if the ref was created
			ref, err := gitRepo.Reference(plumbing.ReferenceName(test.refName), true)
			assert.NoErrorf(t, err, "expected no error while retrieving created ref")
			if ref == nil {
				t.Errorf("expected ref to be created")
				return
			}

			assert.Equal(t, test.hashStr, ref.Hash().String())
		})
	}
}
