package scanning

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func RunTrufflehog(repoPath string, outPath string, cliArgs ...string) error {
	// check if bare repo
	bare := false
	_, gitHeadErr := os.Stat(fmt.Sprintf("%s/HEAD", repoPath))
	_, dotGitErr := os.Stat(fmt.Sprintf("%s/.git", repoPath))
	// check if .git directory exists
	if os.IsNotExist(dotGitErr) && gitHeadErr == nil {
		logrus.Debug("scanning trufflehog using --bare")
		bare = true
	}

	args := []string{
		"git",
		fmt.Sprintf("file://%s", repoPath),
		"--json",
		"--no-fail",
	}
	if bare {
		args = append(args, "--bare")
	}
	args = append(args, cliArgs...)
	cmd := exec.Command("trufflehog", args...)

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("could not open %s for writing: %v", outPath, err)
	}

	cmd.Stdout = f
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		logrus.Errorf("failed to start trufflehog: %v", err)
		return err
	}

	if err := cmd.Wait(); err != nil {
		logrus.Errorf("trufflehog failed: %v", err)
		return err
	}

	return nil
}
