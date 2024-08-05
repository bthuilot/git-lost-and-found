package processor

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
)

const maxWorkers = 8

type job struct {
	file *object.File
}

func ProcessCommit(commit *object.Commit, blobCache map[plumbing.Hash]BlobInfo) (reports []SecretsReport, err error) {
	if commit == nil {
		return nil, nil
	}
	tree, err := commit.Tree()
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	jobCh := make(chan job, maxWorkers)
	resultCh := make(chan SecretsReport)

	// Worker pool
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(jobCh, resultCh, blobCache, &wg, &mu)
	}

	go func() {
		tree.Files().ForEach(func(f *object.File) error {
			jobCh <- job{file: f}
			return nil
		})
		close(jobCh)
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for report := range resultCh {
		reports = append(reports, report)
	}

	return
}

func worker(jobCh <-chan job, resultCh chan<- SecretsReport, blobCache map[plumbing.Hash]BlobInfo, wg *sync.WaitGroup, mu *sync.Mutex) {

	exclusions := map[string][]string{
		"Go":      {"**/vendor/**", "**/go.sum", "**/go.mod", "**/*.pb.go", "**/*.pb.gw.go"},
		"NodeJS":  {"**/node_modules/**", "**/*.lock"},
		"Python":  {"**/__pycache__/**", "**/*.pyc", "**/*.pyo", "**/*.pyd"},
		"Java":    {"**/target/**", "**/build/**", "**/*.class", "**/*.war", "**/*.ear", "**/*.jar"},
		"C#":      {"**/obj/**", "**/bin/**", "**/*.dll", "**/*.exe", "**/*.csproj", "**/*.sln"},
		"Ruby":    {"**/.cache/**", "**/*.rb", "**/.bundle/**"},
		"Generic": {"**/README.md", "**/LICENSE", "**/CHANGELOG.md", "**/CONTRIBUTING.md", "**/CODE_OF_CONDUCT.md", "**/SECURITY.md", "**/PULL_REQUEST_TEMPLATE.md", "**/ISSUE_TEMPLATE.md"},
		"Test":    {"**/*.test", "**/*.spec", "**/*.mock", "**/*.stub", "**/*.fake", "**/*.mocks", "**/*.stubs", "**/*.fakes", "**/*.dummies"},
	}

	defer wg.Done()
	for j := range jobCh {
		blob := j.file.Blob
		var shouldSkip bool
		for _, exclusions := range exclusions {
			for _, exclusion := range exclusions {
				match, _ := filepath.Match(exclusion, j.file.Name)
				if match {
					shouldSkip = true
					break
				}
			}
		}
		if shouldSkip {
			logrus.Debugf("Skipping file: %s", j.file.Name)
			continue
		}

		mu.Lock()
		_, exists := blobCache[blob.Hash]
		mu.Unlock()

		if !exists {
			content, err := blob.Reader()
			if err != nil {
				logrus.Error(err)
				continue
			}
			defer content.Close()
			contentBytes, err := io.ReadAll(content)
			if err != nil {
				logrus.Error(err)
				continue
			}

			mu.Lock()
			blobCache[blob.Hash] = BlobInfo{
				Hash:    blob.Hash,
				Content: contentBytes,
			}
			mu.Unlock()

			report, err := prepareScanJob(j.file.Name, blob.Hash.String(), contentBytes)
			if err != nil {
				logrus.Error(err)
				continue
			}

			resultCh <- report
		}
	}
}

func prepareScanJob(fileName, blobHash string, content []byte) (report SecretsReport, err error) {
	logrus.Debugf("Scanning content for file: %s\n", fileName)

	return SecretsReport{
		BlobHash: blobHash,
		FileName: fileName,
		RawBlob:  string(content),
		Results:  []GitleaksResult{},
	}, nil

	// stdOut, stdErr, err := runGitleaksScan(string(content))

	// if err != nil {
	// 	logrus.Debugf("Gitleaks scan stdOut: %s", stdOut)
	// 	logrus.Debugf("Gitleaks scan stderr: %s", stdErr)

	// 	var results []GitleaksResult
	// 	err = json.Unmarshal([]byte(stdOut), &results)
	// 	if err != nil {
	// 		logrus.Errorf("Failed to unmarshal Gitleaks result: %v", err)
	// 		return SecretsReport{}, err
	// 	}

	// 	for _, result := range results {
	// 		logrus.Errorf("Gitleaks scan result for file %s: %s\n", fileName, result.Secret)
	// 	}

	// 	report.BlobHash = blobHash
	// 	report.Results = results
	// 	report.FileName = fileName
	// 	report.RawBlob = string(content)
	// 	return report, nil
	// }
	// return SecretsReport{}, nil
}

func RunGitleaksScanFileSystem(scanPath string) (string, string, error) {
	logrus.Infof("Starting Gitleaks scan")

	cmd := exec.Command("gitleaks", "detect", "--no-git", "--report-format", "json", "--report-path", "/dev/stdout")
	cmd.Dir = scanPath

	_, err := cmd.StdinPipe()
	if err != nil {
		logrus.Errorf("Failed to get stdin pipe: %v", err)
		return "", "", err
	}

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err := cmd.Start(); err != nil {
		logrus.Errorf("Failed to start cmd: %v", err)
		return "", "", err
	}

	if err := cmd.Wait(); err != nil {
		logrus.Warnf("Cmd failed (FINDING): %v", err)
		return outBuf.String(), errBuf.String(), nil
	}

	return outBuf.String(), errBuf.String(), fmt.Errorf("gitleaks scan failed to run successfully")
}

// func RunGitleaksScanJob(content string) (string, string, error) {
// 	logrus.Debugf("Starting Gitleaks scan")

// 	cmd := exec.Command("gitleaks", "detect", "--pipe", "--report-format", "json", "--report-path", "/dev/stdout")

// 	stdin, err := cmd.StdinPipe()
// 	if err != nil {
// 		logrus.Errorf("Failed to get stdin pipe: %v", err)
// 		return "", "", err
// 	}

// 	var outBuf, errBuf bytes.Buffer
// 	cmd.Stdout = &outBuf
// 	cmd.Stderr = &errBuf

// 	if err := cmd.Start(); err != nil {
// 		logrus.Errorf("Failed to start cmd: %v", err)
// 		return "", "", err
// 	}

// 	go func() {
// 		defer stdin.Close()
// 		_, err := io.WriteString(stdin, content)
// 		if err != nil {
// 			logrus.Errorf("Failed to write to stdin: %v", err)
// 		}
// 	}()

// 	if err := cmd.Wait(); err != nil {
// 		logrus.Errorf("Cmd failed (FINDING): %v", err)
// 		return outBuf.String(), errBuf.String(), err
// 	}

// 	logrus.Debugf("Gitleaks scan stdout:\n%s", outBuf.String())
// 	logrus.Debugf("Gitleaks scan stderr:\n%s", errBuf.String())
// 	return outBuf.String(), errBuf.String(), nil
// }
