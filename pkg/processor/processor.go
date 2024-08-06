package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
)

type blobMetadata struct {
	commit *object.Commit
	name   string
	path   string
}

const maxBatchSize = 100 // Adjust the batch size as needed

func ProcessCommits(commits []*object.Commit, args GitleaksArgs) ([]GitleaksResult, error) {
	var blobCache = make(map[plumbing.Hash]*blobMetadata)
	var mu sync.Mutex
	tempDir, err := os.MkdirTemp("", "blob-temp")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	logrus.Infof("Processing %d commits", len(commits))

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

	// Collect all unique blobs and their metadata
	for _, commit := range commits {
		tree, err := commit.Tree()
		if err != nil {
			return nil, err
		}

		logrus.Infof("Processing commit: %s", commit.Hash.String())

		err = tree.Files().ForEach(func(f *object.File) error {
			blob := f.Blob

			if shouldSkip(f.Name) {
				logrus.Infof("Skipping file: %s", f.Name)
				return nil
			}
			logrus.Infof("Processing file: %s", f.Name)

			mu.Lock()
			_, exists := blobCache[blob.Hash]
			mu.Unlock()
			if !exists {
				content, err := blob.Reader()
				if err != nil {
					return err
				}
				defer content.Close()

				blobPath := filepath.Join(tempDir, blob.Hash.String())
				outFile, err := os.Create(blobPath)
				if err != nil {
					return err
				}
				_, err = io.Copy(outFile, content)
				if err != nil {
					outFile.Close()
					return err
				}
				outFile.Close()

				mu.Lock()
				blobCache[blob.Hash] = &blobMetadata{
					commit: commit,
					name:   f.Name,
					path:   blobPath,
				}
				mu.Unlock()
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	// Process blobs in batches
	var results []GitleaksResult
	var batch []*blobMetadata
	for _, metadata := range blobCache {
		batch = append(batch, metadata)
		if len(batch) >= maxBatchSize {
			batchResults, err := processBatch(batch, args)
			if err != nil {
				return nil, err
			}
			results = append(results, batchResults...)
			batch = nil // Reset the batch
		}
	}

	// Process any remaining blobs
	if len(batch) > 0 {
		batchResults, err := processBatch(batch, args)
		if err != nil {
			return nil, err
		}
		results = append(results, batchResults...)
	}

	logrus.Infof("Found %d secrets in %d unique blobs", len(results), len(blobCache))
	return results, nil
}

func processBatch(batch []*blobMetadata, args GitleaksArgs) ([]GitleaksResult, error) {
	var combinedBuffer bytes.Buffer
	for _, metadata := range batch {
		combinedBuffer.WriteString(fmt.Sprintf("=== Commit: %s, File: %s ===\n", metadata.commit.Hash, metadata.name))
		blobFile, err := os.Open(metadata.path)
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(&combinedBuffer, blobFile)
		blobFile.Close()
		if err != nil {
			return nil, err
		}
		combinedBuffer.WriteString("\n")
	}

	resRaw, stdErr, err := runGitleaksScan(&combinedBuffer, args)
	if err != nil {
		logrus.Errorf("Failed to run Gitleaks scan: %v", err)
		stdErrOutput, _ := io.ReadAll(stdErr)
		logrus.Errorf("Stderr: %s", string(stdErrOutput))
		return nil, err
	}

	var gitleaksResults []GitleaksResult
	if err = json.NewDecoder(resRaw).Decode(&gitleaksResults); err != nil {
		logrus.Errorf("Failed to unmarshal Gitleaks result: %v", err)
		return nil, err
	}

	var results []GitleaksResult
	for _, result := range gitleaksResults {
		for _, metadata := range batch {
			blobFile, err := os.ReadFile(metadata.path)
			if err != nil {
				logrus.Errorf("Failed to read blob file: %v", err)
				continue
			}
			if bytes.Contains(blobFile, []byte(result.Secret)) {
				results = append(results, GitleaksResult{
					Secret:      result.Secret,
					Commit:      metadata.commit.Hash.String(),
					File:        metadata.name,
					Author:      metadata.commit.Author.String(),
					RawFile:     string(blobFile),
					Description: result.Description,
					RuleID:      result.RuleID,
					Tags:        result.Tags,
				})
			}
		}
	}

	return results, nil
}

func runGitleaksScan(content io.Reader, args GitleaksArgs) (io.Reader, io.Reader, error) {
	logrus.Info("Starting Gitleaks scan")

	// Create a temporary file for the Gitleaks report
	tempFile, err := os.CreateTemp("", "gitleaks-report-*.json")
	if err != nil {
		logrus.Errorf("Failed to create temp file: %v", err)
		return nil, nil, err
	}
	defer os.Remove(tempFile.Name())

	cmd := exec.Command("gitleaks", "detect", "--pipe", "--report-format", "json", "--report-path", tempFile.Name(), "--exit-code", "0")
	if args.Config != "" {
		cmd.Args = append(cmd.Args, "--config", args.Config)
	}
	cmd.Stdin = content

	outBuf, errBuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	cmd.Stdout = outBuf
	cmd.Stderr = errBuf

	if err := cmd.Start(); err != nil {
		logrus.Errorf("Failed to start cmd: %v", err)
		return outBuf, errBuf, err
	}

	if err := cmd.Wait(); err != nil {
		logrus.Errorf("Cmd failed (FINDING): %v", err)
		return outBuf, errBuf, err
	}

	reportData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		logrus.Errorf("Failed to read temp file: %v", err)
		return nil, nil, err
	}
	reportDataReader := bytes.NewReader(reportData)

	return reportDataReader, errBuf, nil
}
