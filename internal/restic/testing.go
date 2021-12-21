package restic

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCase represents a test case for a call to restic.
type TestCase struct {
	Name       string
	Repo       string
	Password   string
	BackupPath string

	// Put any additional options in here. The Run method makes sure this gets
	// additionally filled with an WithRunner option pointing to a mock
	// runner for external commands.
	Options     []Option
	Invocations func(t *testing.T, tt *TestCase) []ExpectedInvocation
	Perform     func(t *testing.T, tt *TestCase)
	Assert      func(t *testing.T, tt *TestCase)
}

// Run executes the current test case.
func (tt *TestCase) Run(t *testing.T) {
	runner := &TestCmdRunner{
		T: t,
	}
	if tt.Invocations != nil {
		runner.Invocations = append(runner.Invocations, tt.Invocations(t, tt)...)
	}
	tt.Options = append(tt.Options, WithCmdRunner(runner))

	tt.Perform(t, tt)
	if tt.Assert != nil {
		tt.Assert(t, tt)
	}
	runner.AssertComplete()
}

// ExpectedInvocation represents an invocation of the fake restic executable.
type ExpectedInvocation struct {
	Args []string
	Env  map[string]string
	Code int
}

// TestCmdRunner is a runner that converts the cmd passed to Run to an
// Invocation and compares that to its internal list of invocations. If
// the cmd does not match the expectation a test error is recorded.
type TestCmdRunner struct {
	Invocations []ExpectedInvocation
	T           *testing.T

	pos int
}

// Run checks if cmd matches the next invocation in Invocations. If that
// Invocation has an exit code other than zero an Error is returned.
func (r *TestCmdRunner) Run(cmd *exec.Cmd) error {
	if r.pos >= len(r.Invocations) {
		r.T.Fatalf("More calls than invocations: %d >= %d", r.pos, len(r.Invocations))
	}
	inv := r.Invocations[r.pos]
	r.pos++

	var expectedEnv []string // nolint: prealloc
	for k, v := range inv.Env {
		expectedEnv = append(expectedEnv, fmt.Sprintf("%s=%s", k, v))
	}

	assert := assert.New(r.T)
	assert.ElementsMatch(inv.Args, cmd.Args, "Arguments don't match")
	assert.ElementsMatch(cmd.Env, expectedEnv, "Environment does not match")

	if inv.Code > 0 {
		return Error{
			ExitCode: inv.Code,
		}
	}

	return nil
}

// AssertComplete asserts that all expected restic invocations were actually made.
func (r *TestCmdRunner) AssertComplete() {
	if r.pos != len(r.Invocations) {
		r.T.Errorf("Expected restic to be called %d times; was only called %d times", len(r.Invocations), r.pos)
	}
}
