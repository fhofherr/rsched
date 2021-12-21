package restic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions_Apply(t *testing.T) {
	tests := []struct {
		name    string
		options []Option
		errMsg  string
	}{
		{
			name:    "RESTIC_REPOSITORY and RESTIC_PASSWORD present",
			options: []Option{WithRepository("/some/repo"), WithPassword("some password")},
		},
		{
			name: "RESTIC_REPOSITORY_FILE and RESTIC_PASSWORD present",
			options: []Option{
				WithEnv(map[string]string{envResticRepositoryFile: "/path/to/some/file"}),
				WithPassword("some password"),
			},
		},
		{
			name: "RESTIC_REPOSITORY and RESTIC_PASSWORD_FILE present",
			options: []Option{
				WithRepository("/some/repo"),
				WithEnv(map[string]string{envResticPasswordFile: "/path/to/some/file"}),
			},
		},
		{
			name: "RESTIC_REPOSITORY and RESTIC_PASSWORD_COMMAND present",
			options: []Option{
				WithRepository("/some/repo"),
				WithEnv(map[string]string{envResticPasswordCommand: "/path/to/some/command"}),
			},
		},
		{
			name:    "RESTIC_REPOSITORY missing",
			options: []Option{WithPassword("/some/password")},
			errMsg:  "environment requires one of: RESTIC_REPOSITORY, RESTIC_REPOSITORY_FILE",
		},
		{
			name:    "RESTIC_PASSWORD missing",
			options: []Option{WithRepository("/some/repo")},
			errMsg:  "environment requires one of: RESTIC_PASSWORD, RESTIC_PASSWORD_FILE, RESTIC_PASSWORD_COMMAND",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var opts options

			err := opts.Apply(tt.options)
			if tt.errMsg == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, tt.errMsg)
		})
	}
}
