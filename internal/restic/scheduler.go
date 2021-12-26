package restic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/robfig/cron/v3"
)

// ScheduleOnce is a special schedule which signals that a certain job should
// be executed only once.
const ScheduleOnce = "once"

// Scheduler takes care of scheduling the various restic commands and handling
// graceful shutdown.
//
// Schedule Argument
//
// Some methods of Scheduler expect a schedule argument of type string. The
// value of schedule may either be a valid cron expresssion as supported by
// github.com/robfig/cron or the special value "once". Passing "once" leads
// to the respective function being executed immediately in a separate go
// routine.
type Scheduler struct {
	// The function that is called whenever it is time to create a backup.
	// Defaults to Backup.
	BackupFunc func(ctx context.Context, path string, os ...Option) error

	once       sync.Once
	cron       *cron.Cron
	sempaphore chan struct{}
	shutdown   chan struct{}
}

// ScheduleBackup ensures the BackupFunc is being called according to schedule.
//
// See the documentation of the Scheduler type for the definition of schedule.
func (s *Scheduler) ScheduleBackup(schedule, path string, os ...Option) error {
	return s.scheduleFunc(schedule, func(ctx context.Context) {
		log.Println("Beginning backup")
		if err := s.BackupFunc(ctx, path, os...); err != nil {
			var rErr Error

			log.Printf("Error during backup: %v", err)
			if errors.As(err, &rErr) && len(rErr.Stderr) > 0 {
				log.Printf("Restic stderr: %s", string(rErr.Stderr))
			}
			return
		}
		log.Println("Backup successfully completed")
	})
}

// Run starts the Scheduler in the calling go routine.
func (s *Scheduler) Run() {
	s.init()
	s.cron.Run()
}

// Shutdown performs a graceful shutdown of the Scheduler.
//
// Any currently running jobs get notified of the imminent shutdown. Currently
// running invocations of restic may get a signal sent using kill.
//
// Shutdown blocks the calling go routine until all currently running tasks
// are completed.
func (s *Scheduler) Shutdown() {
	s.init() // Call init to ensure s.shutdown exists even if nothing was scheduled

	close(s.shutdown)
	s.cron.Stop()

	// Wait for all jobs to finish by acquiring the semaphore. This only works
	// as long as the semaphore can be acquired only once. Should this change
	// we need to acquire the semaphore as often as the sempaphore is buffered.
	s.acquireSemaphore(context.Background())
}

func (s *Scheduler) scheduleFunc(schedule string, f func(context.Context)) error {
	s.init()

	log.Printf("Adding job with schedule %q", schedule)
	job := s.newJob(f)
	if schedule == ScheduleOnce {
		go job()
		return nil
	}
	if _, err := s.cron.AddFunc(schedule, job); err != nil {
		return fmt.Errorf("add cron entry: %v", err)
	}
	return nil
}

func (s *Scheduler) init() {
	s.once.Do(func() {
		s.shutdown = make(chan struct{})
		// Initialize the semaphore with 1 as we do not want to have more than
		// one job running at any time.
		s.sempaphore = make(chan struct{}, 1)
		s.cron = cron.New()

		if s.BackupFunc == nil {
			s.BackupFunc = Backup
		}
	})
}

// notifyShutdown creates a context using ctx that is gets canceled once
// Shutdown is called.
func (s *Scheduler) notifyShutdown(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		select {
		case <-s.shutdown:
			cancel()
		case <-ctx.Done():
			// Abort, some other part of the code called cancel
		}
	}()

	return ctx, cancel
}

// newJob wraps a function f to be notified of scheduler shutdown and
// acquire and release the semaphore.
func (s *Scheduler) newJob(f func(context.Context)) func() {
	return func() {
		ctx, cancel := s.notifyShutdown(context.Background())
		defer cancel()

		if !s.acquireSemaphore(ctx) {
			return
		}
		defer s.releaseSempaphore()

		f(ctx)
	}
}

func (s *Scheduler) acquireSemaphore(ctx context.Context) bool {
	select {
	case s.sempaphore <- struct{}{}:
		return true
	case <-ctx.Done():
		return false
	}
}

func (s *Scheduler) releaseSempaphore() {
	<-s.sempaphore
}
