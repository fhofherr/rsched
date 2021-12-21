package restic_test

import (
	"context"
	"testing"

	"github.com/fhofherr/rsched/internal/restic"
	"github.com/stretchr/testify/assert"
)

func TestBackup(t *testing.T) {
	tests := []restic.TestCase{
		{
			Name:       "invalid environment",
			BackupPath: "/never/backed/up",
			Perform: func(t *testing.T, tt *restic.TestCase) {
				err := restic.Backup(context.Background(), tt.BackupPath, tt.Options...)
				assert.Error(t, err)
			},
		},
		{
			Name:       "repository not yet initialized",
			Repo:       "/path/to/repository",
			Password:   "super secret",
			BackupPath: "/some/path",
			Invocations: func(t *testing.T, tt *restic.TestCase) []restic.ExpectedInvocation {
				return []restic.ExpectedInvocation{
					{
						Args: []string{"restic", "snapshots"},
						Env: map[string]string{
							"RESTIC_REPOSITORY": tt.Repo,
							"RESTIC_PASSWORD":   tt.Password,
						},
						// Restic returns non-zero exit code if repo not initialized or something else went wrong.
						// TODO: we may wan't to parse restic's output in the future.
						Code: 1,
					},
					{
						Args: []string{"restic", "init"},
						Env: map[string]string{
							"RESTIC_REPOSITORY": tt.Repo,
							"RESTIC_PASSWORD":   tt.Password,
						},
					},
					{
						Args: []string{"restic", "backup", tt.BackupPath},
						Env: map[string]string{
							"RESTIC_REPOSITORY": tt.Repo,
							"RESTIC_PASSWORD":   tt.Password,
						},
					},
				}
			},
			Perform: func(t *testing.T, tt *restic.TestCase) {
				tt.Options = append(tt.Options, restic.WithRepository(tt.Repo), restic.WithPassword(tt.Password))
				err := restic.Backup(context.Background(), tt.BackupPath, tt.Options...)
				assert.NoError(t, err)
			},
		},
		{
			Name:       "error during repo initialization",
			Repo:       "/path/to/repository",
			Password:   "super secret",
			BackupPath: "/some/path",
			Invocations: func(t *testing.T, tt *restic.TestCase) []restic.ExpectedInvocation {
				return []restic.ExpectedInvocation{
					{
						Args: []string{"restic", "snapshots"},
						Env: map[string]string{
							"RESTIC_REPOSITORY": tt.Repo,
							"RESTIC_PASSWORD":   tt.Password,
						},
						Code: 1,
					},
					{
						Args: []string{"restic", "init"},
						Env: map[string]string{
							"RESTIC_REPOSITORY": tt.Repo,
							"RESTIC_PASSWORD":   tt.Password,
						},
						Code: 1,
					},
				}
			},
			Perform: func(t *testing.T, tt *restic.TestCase) {
				tt.Options = append(tt.Options, restic.WithRepository(tt.Repo), restic.WithPassword(tt.Password))
				err := restic.Backup(context.Background(), tt.BackupPath, tt.Options...)
				assert.ErrorIs(t, err, restic.Error{
					Command:  "init",
					ExitCode: 1,
				})
			},
		},
		{
			Name:       "repository already initialized",
			Repo:       "/other/path/to/repository",
			Password:   "even more secret",
			BackupPath: "/more/important/data",
			Invocations: func(t *testing.T, tt *restic.TestCase) []restic.ExpectedInvocation {
				return []restic.ExpectedInvocation{
					{
						Args: []string{"restic", "snapshots"},
						Env: map[string]string{
							"RESTIC_REPOSITORY": tt.Repo,
							"RESTIC_PASSWORD":   tt.Password,
						},
					},
					{
						Args: []string{"restic", "backup", tt.BackupPath},
						Env: map[string]string{
							"RESTIC_REPOSITORY": tt.Repo,
							"RESTIC_PASSWORD":   tt.Password,
						},
					},
				}
			},
			Perform: func(t *testing.T, tt *restic.TestCase) {
				tt.Options = append(tt.Options, restic.WithRepository(tt.Repo), restic.WithPassword(tt.Password))
				err := restic.Backup(context.Background(), tt.BackupPath, tt.Options...)
				assert.NoError(t, err)
			},
		},
		{
			Name:       "fatal error during backup",
			Repo:       "/other/path/to/repository",
			Password:   "even more secret",
			BackupPath: "/more/important/data",
			Invocations: func(t *testing.T, tt *restic.TestCase) []restic.ExpectedInvocation {
				return []restic.ExpectedInvocation{
					{
						Args: []string{"restic", "snapshots"},
						Env: map[string]string{
							"RESTIC_REPOSITORY": tt.Repo,
							"RESTIC_PASSWORD":   tt.Password,
						},
					},
					{
						Args: []string{"restic", "backup", tt.BackupPath},
						Env: map[string]string{
							"RESTIC_REPOSITORY": tt.Repo,
							"RESTIC_PASSWORD":   tt.Password,
						},
						Code: 1,
					},
				}
			},
			Perform: func(t *testing.T, tt *restic.TestCase) {
				tt.Options = append(tt.Options, restic.WithRepository(tt.Repo), restic.WithPassword(tt.Password))
				err := restic.Backup(context.Background(), tt.BackupPath, tt.Options...)
				assert.ErrorIs(t, err, restic.Error{
					Command:  "backup",
					ExitCode: 1,
				})
			},
		},
		{
			Name:       "incomplete backup",
			Repo:       "/yet/another/repository",
			Password:   "secret secret secret",
			BackupPath: "/really/important/data",
			Invocations: func(t *testing.T, tt *restic.TestCase) []restic.ExpectedInvocation {
				return []restic.ExpectedInvocation{
					{
						Args: []string{"restic", "snapshots"},
						Env: map[string]string{
							"RESTIC_REPOSITORY": tt.Repo,
							"RESTIC_PASSWORD":   tt.Password,
						},
					},
					{
						Args: []string{"restic", "backup", tt.BackupPath},
						Env: map[string]string{
							"RESTIC_REPOSITORY": tt.Repo,
							"RESTIC_PASSWORD":   tt.Password,
						},
						Code: 3,
					},
				}
			},
			Perform: func(t *testing.T, tt *restic.TestCase) {
				tt.Options = append(tt.Options, restic.WithRepository(tt.Repo), restic.WithPassword(tt.Password))
				err := restic.Backup(context.Background(), tt.BackupPath, tt.Options...)

				assert.ErrorIs(t, err, restic.Error{
					Command:  "backup",
					ExitCode: 3,
				})
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, tt.Run)
	}
}
