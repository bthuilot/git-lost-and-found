package scanning

import (
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func ExecScanner(dir string, cmdArgs []string) error {
	if len(cmdArgs) == 0 {
		logrus.Warnf("no scanner command provided, exiting")
		return nil
	}

	for i, arg := range cmdArgs {
		cmdArgs[i] = strings.ReplaceAll(arg, "{}", dir)
	}

	logrus.WithFields(logrus.Fields{
		"scanner": cmdArgs[0],
		"args":    cmdArgs[1:],
	}).Debug("executing scanner command")

	// #nosec G204
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir

	return cmd.Run()
}
