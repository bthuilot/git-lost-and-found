// Copyright (C) 2025 Bryce Thuilot
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the FSF, either version 3 of the License, or (at your option) any later version.
// See the LICENSE file in the root of this repository for full license text or
// visit: <https://www.gnu.org/licenses/gpl-3.0.html>.

package scanning

import (
	"os"
	"os/exec"
	"strings"

	"go.uber.org/zap"
)

func ExecScanner(dir string, cmdArgs []string) error {
	if len(cmdArgs) == 0 {
		zap.L().Warn("no scanner command provided, exiting")
		return nil
	}

	for i, arg := range cmdArgs {
		cmdArgs[i] = strings.ReplaceAll(arg, "{}", dir)
	}

	zap.L().
		Info("executing scanner command", zap.String("scanner", cmdArgs[0]), zap.Strings("args", cmdArgs[1:]))

	// #nosec G204
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dir

	return cmd.Run()
}
