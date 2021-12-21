package restic

import (
	"errors"
	"fmt"
)

// Error represents an error that occurred while interacting with restic.
type Error struct {
	Command  string
	ExitCode int

	// Set err to transport the error that lead to Error. This is mainly for
	// logging and debugging purposes. Values that are provided by the error
	// and have a corresponding Field in this struct need to be set.
	Err error
}

func (e Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("restic %s: %v", e.Command, e.Err)
	}
	return fmt.Sprintf("restic %s: %d", e.Command, e.ExitCode)
}

func asError(err error) (Error, bool) {
	var rErr Error

	ok := errors.As(err, &rErr)
	return rErr, ok
}
