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

func ProcessCommit(commit *object.Commit, blobCache map[plumbing.Hash]struct{}, args GitleaksArgs) ([]GitleaksResult, error) {
	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	var results []GitleaksResult

	err = tree.Files().ForEach(func(f *object.File) error {
		blob := f.Blob

		// Check if we've already processed this blob
		if _, exists := blobCache[blob.Hash]; !exists {
			content, err := blob.Reader()
			if err != nil {
				return err
			}
			defer content.Close()

			// add blob to cache
			blobCache[blob.Hash] = struct{}{}
			found, err := scanContent(commit, f.Name, content, args)
			if err != nil {
				return err
			}
			results = append(results, found...)
		}
		return nil
	})

	return results, err
}

func scanContent(commit *object.Commit, fileName string, content io.Reader, args GitleaksArgs) (results []GitleaksResult, err error) {
	logrus.Infof("Scanning content for file: %s\n", fileName)
	stdOut, stdErr, err := runGitleaksScan(content, args)
	if err != nil {
		logrus.Errorf("Failed to run Gitleaks scan: %v", err)
		stdErrOutput, _ := io.ReadAll(stdErr)
		logrus.Errorf("Stderr: %s", string(stdErrOutput))
		return nil, err
	}

	var r []GitleaksResult
	if err = json.NewDecoder(stdOut).Decode(&r); err != nil {
		logrus.Errorf("Failed to unmarshal Gitleaks result: %v", err)
		return nil, err
	}

	for _, result := range r {
		logrus.Debugf("Gitleaks scan result for file %s: %s\n", fileName, result.Secret)
		result.Commit = commit.Hash.String()
		result.File = fileName
		result.Author = commit.Author.String()
		results = append(results, result)
	}

	return
	// logrus.Infof("Gitleaks scan result for file %s: %s\n", fileName, result)
}

func runGitleaksScan(content io.Reader, args GitleaksArgs) (io.Reader, io.Reader, error) {
	logrus.Info("Starting Gitleaks scan")

	cmd := exec.Command("gitleaks", "detect", "--pipe", "--report-format", "json", "--report-path", "/dev/stdout", "--exit-code", "0")
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

	return outBuf, errBuf, nil
}
