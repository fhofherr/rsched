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
		name        string
		args        []string
		expectedCfg cmd.Config
		assertErr   assert.ErrorAssertionFunc
	}{
		{
			name: "Default config",
			expectedCfg: cmd.Config{
				BackupSchedule: "@hourly",
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
			assert.Equal(t, tt.expectedCfg, cfg)
		})
	}
}
