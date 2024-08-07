package scanning

import (
	"github.com/sirupsen/logrus"
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

	if err := cmd.Start(); err != nil {
		logrus.Errorf("Failed to start cmd: %v", err)
		return err
	}

	if err := cmd.Wait(); err != nil {
		logrus.Errorf("Cmd failed: %v", err)
		return err
	}

	return nil
}
