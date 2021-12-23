package restic

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// CmdRunner defines the Run method which allows to run an external command.
//
// Run may return an incomplete Error. Callers may enrich it with additional
// information if possible. Additionally Run may return any other error
// if this is necessary.
//
// The purpose of CmdRunner is mostly to allow for easy testing of various
// restic invocations.
type CmdRunner interface {
	Run(cmd *exec.Cmd) error
}

type osRunner struct{}

func (r *osRunner) Run(cmd *exec.Cmd) error {
	var (
		exitErr *exec.ExitError
		stderr  strings.Builder
	)

	cmd.Stderr = &stderr

	err := cmd.Run()
	if errors.As(err, &exitErr) {
		return Error{
			ExitCode: exitErr.ExitCode(),
			Stderr:   stderr.String(),
		}
	}
	return err
}

func joinEnv(env map[string]string) []string {
	res := make([]string, 0, len(env))
	for k, v := range env {
		res = append(res, fmt.Sprintf("%s=%v", k, v))
	}
	return res
}
