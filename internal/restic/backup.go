package restic

import (
	"context"
	"fmt"
	"log"
	"os/exec"
)

// Backup calls restic backup to create a backup.
//
// repo defines the location of the restic repository. It needs to be in the
// format defined in the restic documentation
// (https://restic.readthedocs.io/en/stable/030_preparing_a_new_repo.html).
// If the repository is not yet initialized Backup initializes a new restic
// repository. The value of pw is used as the encryption password of the
// repository.
//
// Additional options can be passed using opts. Any options not relevant for
// creating a backup are silently ignored.
func Backup(ctx context.Context, path string, os ...Option) error {
	var opts options

	if err := opts.Apply(os); err != nil {
		return fmt.Errorf("backup options: %v", err)
	}

	if !repoInitialized(ctx, opts) {
		if err := initializeRepo(ctx, opts); err != nil {
			return err // no need to wrap, error already handled by initializeRepo
		}
	}

	cmd := exec.CommandContext(ctx, opts.Restic, "backup", path)
	cmd.Env = joinEnv(opts.Env)

	if err := opts.Runner.Run(cmd); err != nil {
		if err == ctx.Err() {
			return err
		}
		if rErr, ok := asError(err); ok {
			rErr.Command = "backup"
			return rErr
		}
		return fmt.Errorf("restic backup: %v", err)
	}

	return nil
}

func repoInitialized(ctx context.Context, opts options) bool {
	cmd := exec.CommandContext(ctx, opts.Restic, "snapshots")
	cmd.Env = joinEnv(opts.Env)

	if err := opts.Runner.Run(cmd); err != nil {
		log.Printf("repo initialization check encountered error: %v", err)
		return false
	}
	return true
}

func initializeRepo(ctx context.Context, opts options) error {
	cmd := exec.CommandContext(ctx, opts.Restic, "init")
	cmd.Env = joinEnv(opts.Env)

	if err := opts.Runner.Run(cmd); err != nil {
		if err == ctx.Err() {
			return err
		}
		if rErr, ok := asError(err); ok {
			rErr.Command = "init"
			return rErr
		}
		return fmt.Errorf("initialize restic repository: %v", err)
	}
	return nil
}
