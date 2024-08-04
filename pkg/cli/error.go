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
