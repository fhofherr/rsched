package cmd

import (
	"github.com/fhofherr/rsched/internal/restic"
	"github.com/stretchr/testify/mock"
)

// MockResticScheduler is a mock implementation of the ResticScheduler
// interface used during testing.
type MockResticScheduler struct {
	mock.Mock
}

// ScheduleBackup registers a call to itself and returns the arguments
// it was mocked for.
//
// Since restic.Option is a function which cannot be compared
// restic.MatchOptions needs to be used together with mock.MatchedBy during
// mocking.
func (m *MockResticScheduler) ScheduleBackup(schedule, path string, os ...restic.Option) error {
	// Do not expand os to make this work with restic.MatchOptions
	args := m.Called(schedule, path, os)
	return args.Error(0)
}
