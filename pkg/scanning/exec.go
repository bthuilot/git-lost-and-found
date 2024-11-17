package scanning

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

func ExecScanner(dir string, cmdArgs []string) error {
	if len(cmdArgs) == 0 {
		return errors.New("no scanner command provided")
	}

	for i, arg := range cmdArgs {
		cmdArgs[i] = strings.ReplaceAll(arg, "{}", dir)
	}

	logrus.WithFields(logrus.Fields{
		"scanner": cmdArgs[0],
		"args":    cmdArgs[1:],
	}).Debug("executing scanner command")

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir

	return cmd.Run()
}
