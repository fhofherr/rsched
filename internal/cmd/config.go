package cmd

import (
	"flag"
	"fmt"

	"github.com/peterbourgon/ff/v3"
)

// Config contains the configuration for the rsched command. The individual
// values can be either set using command line flags or environment variables.
type Config struct {
	PrintVersion       bool
	BackupPath         string
	BackupSchedule     string
	ResticPasswordFile string
	ResticRepository   string
	ResticBinary       string
}

// LoadConfig loads a new Config from the environment and command line flags.
func LoadConfig(args []string) (Config, error) {
	var cfg Config

	fs := flag.NewFlagSet("rsched", flag.ContinueOnError)
	fs.BoolVar(&cfg.PrintVersion, "v", false, "Print version and exit")
	fs.StringVar(&cfg.BackupSchedule, "backup-schedule", "@hourly", "Interval in which backups should be taken.")
	fs.StringVar(&cfg.BackupPath, "restic-backup-path", "/", "Directory to backup.")
	fs.StringVar(
		&cfg.ResticPasswordFile,
		"restic-password-file",
		"",
		`Path to a file containing the restic repository password.

The value of this flag is ignored if the RESTIC_PASSWORD_FILE environment
variable is set.
`)
	fs.StringVar(
		&cfg.ResticRepository,
		"restic-repository",
		"",
		`Location of the restic repository.

See the restic documentation for valid values
(https://restic.readthedocs.io/en/stable/030_preparing_a_new_repo.html).

The value of this flag is ignored if the RESTIC_REPOSITORY environment
variable is set.
`)

	fs.StringVar(&cfg.ResticBinary, "restic-binary", "", "Path to the restic binary")

	err := ff.Parse(fs, args, ff.WithEnvVarPrefix("RSCHED"))
	if err != nil {
		return cfg, fmt.Errorf("parse config: %v", err)
	}

	return cfg, nil
}
