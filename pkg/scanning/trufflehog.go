package scanning

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func RunTrufflehog(repoPath string, outPath string, cliArgs ...string) error {
	args := []string{
		"git",
		fmt.Sprintf("file://%s", repoPath),
		"--json",
		"--no-fail",
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
