package scanning

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func RunGitleaks(repoPath string, outPath string, cliArgs ...string) error {
	args := []string{
		"detect",
		"-s", repoPath,
		"-f", "json",
		"-r", outPath,
		"--exit-code", "0",
	}
	args = append(args, cliArgs...)
	cmd := exec.Command("gitleaks", args...)
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		logrus.Errorf("failed to start gitleaks: %v", err)
		return err
	}

	if err := cmd.Wait(); err != nil {
		logrus.Errorf("gitleaks failed: %v", err)
		return err
	}

	return nil
}
