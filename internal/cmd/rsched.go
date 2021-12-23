package cmd

import (
	"log"

	"github.com/fhofherr/rsched/internal/restic"
)

// RSched implements the rsched command.
type RSched struct {
	Scheduler ResticScheduler
}

// Run executes rsched based on the passed config.
func (r *RSched) Run(cfg Config) {
	env := Environ()

	if cfg.ResticRepository != "" && env[restic.EnvResticRepository] == "" {
		env[restic.EnvResticRepository] = cfg.ResticRepository
	}
	if cfg.ResticPasswordFile != "" && env[restic.EnvResticPasswordFile] == "" {
		env[restic.EnvResticPasswordFile] = cfg.ResticPasswordFile
	}
	if cfg.BackupSchedule != "" {
		r.scheduleBackup(cfg, env)
	}
}

// Shutdown performs a graceful shutdown of rsched.
func (r *RSched) Shutdown() {
	panic("not implemented")
}

func (r *RSched) scheduleBackup(cfg Config, env map[string]string) {
	err := r.Scheduler.ScheduleBackup(cfg.BackupSchedule, cfg.BackupPath, restic.WithEnv(env))
	if err != nil {
		log.Printf("Failed to schedule backup: %v", err)
	}
}

// ResticScheduler represents the actual restic scheduler.
type ResticScheduler interface {
	ScheduleBackup(schedule, path string, os ...restic.Option) error
}
