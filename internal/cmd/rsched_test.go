package cmd_test

import (
	"testing"

	"github.com/fhofherr/rsched/internal/cmd"
	"github.com/fhofherr/rsched/internal/restic"
	"github.com/stretchr/testify/mock"
)

func TestRSched_Run(t *testing.T) {
	type testCase struct {
		name string
		cfg  cmd.Config
		mock func(t *testing.T, tt *testCase)

		// Set during test execution
		Scheduler *cmd.MockResticScheduler
	}

	tests := []testCase{
		{
			name: "print version",
			cfg: cmd.Config{
				PrintVersion: true,
			},
		},
		{
			name: "backup only",
			cfg: cmd.Config{
				BackupPath:         "/",
				BackupSchedule:     "@hourly",
				ResticPasswordFile: "/path/to/password-file",
				ResticRepository:   "/path/to/repository",
				ResticBinary:       "/path/to/restic",
			},
			mock: func(t *testing.T, tt *testCase) {
				env := cmd.Environ()
				env[restic.EnvResticRepository] = tt.cfg.ResticRepository
				env[restic.EnvResticPasswordFile] = tt.cfg.ResticPasswordFile

				tt.Scheduler.
					On(
						"ScheduleBackup",
						tt.cfg.BackupSchedule,
						tt.cfg.BackupPath,
						mock.MatchedBy(
							restic.MatchOptions(
								t,
								restic.WithEnv(env),
								restic.WithBinary(tt.cfg.ResticBinary),
							),
						),
					).Return(nil)
				tt.Scheduler.On("Run").Return()
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.Scheduler = &cmd.MockResticScheduler{}
			tt.Scheduler.Test(t)

			if tt.mock != nil {
				tt.mock(t, &tt)
			}

			rsched := &cmd.RSched{
				Scheduler: tt.Scheduler,
			}
			rsched.Run(tt.cfg)

			tt.Scheduler.AssertExpectations(t)
		})
	}
}
