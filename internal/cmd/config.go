package cmd

import (
	"flag"
	"fmt"

	"github.com/peterbourgon/ff/v3"
)

// Config contains the configuration for the rsched command. The individual
// values can be either set using command line flags or environment variables.
type Config struct {
	BackupSchedule string
}

// TODO implement validation of config values

// LoadConfig loads a new Config from the environment and command line flags.
func LoadConfig(args []string) (Config, error) {
	var cfg Config

	fs := flag.NewFlagSet("rsched", flag.ContinueOnError)
	fs.StringVar(&cfg.BackupSchedule, "backup-schedule", "@hourly", "Interval in which backups should be taken.")

	err := ff.Parse(fs, args, ff.WithEnvVarPrefix("RSCHED_"))
	if err != nil {
		return cfg, fmt.Errorf("parse config: %v", err)
	}

	return cfg, nil
}
