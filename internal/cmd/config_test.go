package cmd_test

import (
	"testing"

	"github.com/fhofherr/rsched/internal/cmd"
	"github.com/stretchr/testify/assert"
)

// We do not fake the environment in this test but rather assume ff does this
// correctly.
func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		assertCfg func(t *testing.T, actual cmd.Config)
		assertErr assert.ErrorAssertionFunc
	}{
		{
			name: "Default config",
			assertCfg: func(t *testing.T, actual cmd.Config) {
				expected := cmd.Config{
					BackupSchedule: "@hourly",
					BackupPath:     "/",
				}
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "Print version",
			args: []string{"-v"},
			assertCfg: func(t *testing.T, actual cmd.Config) {
				assert.True(t, actual.PrintVersion)
			},
		},
		{
			name: "Pass backup path",
			args: []string{"-restic-backup-path", "/path/to/backup"},
			assertCfg: func(t *testing.T, actual cmd.Config) {
				assert.Equal(t, "/path/to/backup", actual.BackupPath)
			},
		},
		{
			name: "Pass password file",
			args: []string{"-restic-password-file", "/path/to/password/file"},
			assertCfg: func(t *testing.T, actual cmd.Config) {
				assert.Equal(t, "/path/to/password/file", actual.ResticPasswordFile)
			},
		},
		{
			name: "Pass restic repository",
			args: []string{"-restic-repository", "/path/to/restic/repository"},
			assertCfg: func(t *testing.T, actual cmd.Config) {
				assert.Equal(t, "/path/to/restic/repository", actual.ResticRepository)
			},
		},
		{
			name: "Pass restic binary",
			args: []string{"-restic-binary", "/path/to/restic"},
			assertCfg: func(t *testing.T, actual cmd.Config) {
				assert.Equal(t, "/path/to/restic", actual.ResticBinary)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.assertErr == nil {
				tt.assertErr = assert.NoError
			}
			cfg, err := cmd.LoadConfig(tt.args)
			tt.assertErr(t, err)
			tt.assertCfg(t, cfg)
		})
	}
}
