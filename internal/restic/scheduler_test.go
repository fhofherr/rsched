package restic_test

import (
	"context"
	"testing"
	"time"

	"github.com/fhofherr/rsched/internal/restic"
	"github.com/stretchr/testify/assert"
)

func TestScheduler_ScheduleBackup(t *testing.T) {
	t.Run("schedule backup once", func(t *testing.T) {
		called := make(chan struct{})
		s := &restic.Scheduler{
			BackupFunc: func(ctx context.Context, path string, os ...restic.Option) error {
				close(called)
				return nil
			},
		}
		defer s.Shutdown()

		err := s.ScheduleBackup(restic.ScheduleOnce, "/some/path")
		if !assert.NoError(t, err) {
			return
		}
		select {
		case <-called:
			return
		case <-time.After(10 * time.Millisecond):
			t.Error("BackupFunc not called within 10ms")
		}
	})

	t.Run("schedule backup regularly", func(t *testing.T) {
		s := &restic.Scheduler{
			BackupFunc: func(ctx context.Context, path string, os ...restic.Option) error {
				return nil
			},
		}
		defer s.Shutdown()

		err := s.ScheduleBackup("@hourly", "/some/path")
		if !assert.NoError(t, err) {
			return
		}
		restic.AssertSchedulerHasSingleJob(t, s)
	})

	t.Run("invalid cron schedule", func(t *testing.T) {
		s := &restic.Scheduler{
			BackupFunc: func(ctx context.Context, path string, os ...restic.Option) error {
				return nil
			},
		}
		defer s.Shutdown()

		err := s.ScheduleBackup("invalid", "/some/path")
		assert.Error(t, err)
	})
}

func TestScheduler_Shutdown(t *testing.T) {
	ready := make(chan struct{})
	shutdownDly := 10 * time.Millisecond

	s := &restic.Scheduler{
		BackupFunc: func(ctx context.Context, path string, os ...restic.Option) error {
			close(ready)
			<-ctx.Done()
			<-time.After(shutdownDly)
			return ctx.Err()
		},
	}

	if err := s.ScheduleBackup(restic.ScheduleOnce, "/some/path"); !assert.NoError(t, err) {
		return
	}

	select {
	case <-ready:
		break
	case <-time.After(10 * time.Millisecond):
		t.Error("Test not ready within 10ms")
	}

	start := time.Now()
	s.Shutdown()
	end := time.Now()

	assert.GreaterOrEqual(t, end.Sub(start), shutdownDly)
}
