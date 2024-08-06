package processor

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
)

func ProcessCommit(commit *object.Commit, blobCache map[plumbing.Hash]struct{}, args GitleaksArgs) ([]GitleaksResult, error) {
	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	exclusions := map[string][]string{
		"Go":      {"vendor/", "go.sum", "go.mod", "*.pb.go", "*.pb.gw.go"},
		"NodeJS":  {"node_modules/", "*.lock"},
		"Python":  {"__pycache__/", "*.pyc", "*.pyo", "*.pyd"},
		"Java":    {"target/", "build/", "*.class", "*.war", "*.ear", "*.jar"},
		"C#":      {"obj/", "bin/", "*.dll", "*.exe", "*.csproj", "*.sln"},
		"Ruby":    {".cache/", "*.rb", ".bundle/"},
		"Generic": {"README.md", "LICENSE", "CHANGELOG.md", "CONTRIBUTING.md", "CODE_OF_CONDUCT.md", "SECURITY.md", "PULL_REQUEST_TEMPLATE.md", "ISSUE_TEMPLATE.md"},
		"Test":    {"*.test", "*.spec", "*.mock", "*.stub", "*.fake", "*.mocks", "*.stubs", "*.fakes", "*.dummies"},
		"Common":  {"test/", "tests/", "node_modules/", "*.log", "*.tmp", "*.bak", "*.swp", "*.swo"},
	}

	shouldSkip := func(filename string) bool {
		for _, exclusions := range exclusions {
			for _, exclusion := range exclusions {
				if strings.Contains(exclusion, "/") {
					if strings.HasPrefix(filename, exclusion) {
						return true
					}
				} else {
					match, _ := filepath.Match(exclusion, filepath.Base(filename))
					if match {
						return true
					}
				}
			}
		}
		return strings.Contains(filepath.Base(filename), "test")
	}

	var results []GitleaksResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	err = tree.Files().ForEach(func(f *object.File) error {
		if shouldSkip(f.Name) {
			logrus.Debugf("Skipping file: %s\n", f.Name)
			return nil
		}

		wg.Add(1)
		go func(f *object.File) {
			defer wg.Done()
			blob := f.Blob

			// Check if we've already processed this blob
			mu.Lock()
			if _, exists := blobCache[blob.Hash]; exists {
				mu.Unlock()
				return
			}
			blobCache[blob.Hash] = struct{}{}
			mu.Unlock()

			content, err := blob.Reader()
			if err != nil {
				logrus.Errorf("Failed to read blob: %v", err)
				return
			}
			defer content.Close()

			var raw string
			raw, err = f.Contents()
			if err != nil {
				logrus.Errorf("Failed to read contents: %v", err)
				return
			}

			found, err := scanContent(commit, f.Name, content, args)
			if err != nil {
				logrus.Errorf("Failed to scan content: %v", err)
				return
			}
			// for all res in found, add raw to res.RawFile
			// add res to results
			// add res to uniqueSecrets
			// add res to enriched
			// add enriched to results

			enriched := make([]GitleaksResult, 0)
			for _, res := range found {
				// add original f.Blob file as string to res.RawFile
				res.RawFile = raw
				enriched = append(enriched, res)
			}

			mu.Lock()
			results = append(results, enriched...)
			mu.Unlock()
		}(f)

		return nil
	})

	wg.Wait()

	return results, err
}

func scanContent(commit *object.Commit, fileName string, content io.Reader, args GitleaksArgs) (results []GitleaksResult, err error) {
	logrus.Debugf("Scanning content for file: %s\n", fileName)
	stdOut, stdErr, err := runGitleaksScan(content, args)
	if err != nil {
		logrus.Errorf("Failed to run Gitleaks scan: %v", err)
		stdErrOutput, _ := io.ReadAll(stdErr)
		logrus.Errorf("Stderr: %s", string(stdErrOutput))
		return
	}

	var r []GitleaksResult
	if err = json.NewDecoder(stdOut).Decode(&r); err != nil {
		logrus.Errorf("Failed to unmarshal Gitleaks result: %v", err)
		return
	}

	if len(r) == 0 {
		return
	}

	for _, result := range r {
		logrus.Debugf("Gitleaks scan result for file %s: %s\n", fileName, result.Secret)
		result.Commit = commit.Hash.String()
		result.File = fileName
		result.Author = commit.Author.String()
		results = append(results, result)
	}

	return
}
func runGitleaksScan(content io.Reader, args GitleaksArgs) (io.Reader, io.Reader, error) {
	logrus.Debugf("Starting Gitleaks scan")

	cmd := exec.Command("gitleaks", "detect", "--pipe", "--report-format", "json", "--report-path", "/dev/stdout", "--exit-code", "0")
	if args.Config != "" {
		cmd.Args = append(cmd.Args, "--config", args.Config)
	}
	cmd.Stdin = content

	outBuf := bytes.NewBuffer(nil)
	errBuf := bytes.NewBuffer(nil)

	cmd.Stdout = outBuf
	cmd.Stderr = errBuf

	if err := cmd.Start(); err != nil {
		logrus.Errorf("Failed to start cmd: %v", err)
		return nil, errBuf, err
	}

	if err := cmd.Wait(); err != nil {
		logrus.Errorf("Cmd failed: %v", err)
		return outBuf, errBuf, err
	}

	return outBuf, errBuf, nil
}
