// Copyright (C) 2024-2026 Bryce Thuilot
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the FSF, either version 3 of the License, or (at your option) any later version.
// See the LICENSE file in the root of this repository for full license text or
// visit: <https://www.gnu.org/licenses/gpl-3.0.html>.

package cli

import (
	"fmt"
	"os"
)

func ErrorExit(err error) {
	ErrorMsg(err)
	os.Exit(1)
}

func ErrorMsg(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s", err)
}
