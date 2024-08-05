package processor

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sirupsen/logrus"
)

func ProcessCommit(commit *object.Commit, output io.Writer, blobCache map[plumbing.Hash]BlobInfo) error {
	tree, err := commit.Tree()
	if err != nil {
		return err
	}

	return tree.Files().ForEach(func(f *object.File) error {
		blob := f.Blob

		// Check if we've already processed this blob
		if _, exists := blobCache[blob.Hash]; !exists {
			content, err := blob.Reader()
			if err != nil {
				return err
			}
			defer content.Close()
			contentBytes, err := io.ReadAll(content)
			if err != nil {
				return err
			}

			blobCache[blob.Hash] = BlobInfo{
				Hash:    blob.Hash,
				Content: contentBytes,
			}
			scanContent(commit, f.Name, contentBytes, output)
		}
		return nil
	})
}

func scanContent(commit *object.Commit, fileName string, content []byte, output io.Writer) {
	logrus.Infof("Scanning content for file: %s\n", fileName)
	stdOut, stdErr, err := runGitleaksScan(string(content))
	if err != nil {
		logrus.Errorf("Failed to run Gitleaks scan: %v", err)
		logrus.Errorf("Gitleaks scan stdOut: %s", stdOut)
		logrus.Errorf("Gitleaks scan stderr: %s", stdErr)

		// marshall the stdOut into a GitleaksResult struct
		// and print the raw Secret

		var results []GitleaksResult
		err := json.Unmarshal([]byte(stdOut), &results)
		if err != nil {
			logrus.Errorf("Failed to unmarshal Gitleaks result: %v", err)
		}

		for _, result := range results {
			result.Commit = commit.Hash.String()
			result.File = fileName
			result.Author = commit.Author.String()
			logrus.Errorf("Gitleaks scan result for file %s: %s\n", fileName, result.Secret)
		}

		// write results to a file
		output.Write([]byte(stdOut))
		//if err := os.WriteFile(fmt.Sprintf("results/gitleaks_%s_results.json", fileName), []byte(stdOut), 0644); err != nil {
		//	logrus.Errorf("Failed to write Gitleaks results to file: %v", err)
		//}

		return
	}
	// logrus.Infof("Gitleaks scan result for file %s: %s\n", fileName, result)
}

func runGitleaksScan(content string) (string, string, error) {
	logrus.Info("Starting Gitleaks scan")

	cmd := exec.Command("gitleaks", "detect", "--pipe", "--report-format", "json", "--report-path", "/dev/stdout")

	stdin, err := cmd.StdinPipe()
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

	go func() {
		defer stdin.Close()
		_, err := io.WriteString(stdin, content)
		if err != nil {
			logrus.Errorf("Failed to write to stdin: %v", err)
		}
	}()

	if err := cmd.Wait(); err != nil {
		logrus.Errorf("Cmd failed (FINDING): %v", err)
		return outBuf.String(), errBuf.String(), err
	}

	logrus.Debugf("Gitleaks scan stdout:\n%s", outBuf.String())
	logrus.Debugf("Gitleaks scan stderr:\n%s", errBuf.String())
	return outBuf.String(), errBuf.String(), nil
}
