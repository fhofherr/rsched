package restic

import (
	"errors"
	"fmt"
	"strings"
)

// Error represents an error that occurred while interacting with restic.
type Error struct {
	Command  string
	ExitCode int
	Stderr   string
}

func (e Error) Error() string {
	var sb strings.Builder

	fmt.Fprint(&sb, "restic")
	if e.Command != "" {
		fmt.Fprintf(&sb, " %s", e.Command)
	}
	if e.ExitCode > 0 {
		fmt.Fprintf(&sb, ": exit code: %d", e.ExitCode)
	}

	return sb.String()
}

func asError(err error) (Error, bool) {
	var rErr Error

	ok := errors.As(err, &rErr)
	return rErr, ok
}
