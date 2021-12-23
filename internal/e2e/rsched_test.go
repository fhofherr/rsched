package e2e_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fhofherr/rsched/internal/cmd"
	"github.com/fhofherr/rsched/internal/restic"
	"github.com/fhofherr/rsched/internal/testsupport"
	"github.com/stretchr/testify/assert"
)

func TestRschedBackup(t *testing.T) {
	if os.Getenv("RSCHED_TEST_SKIP_E2E") != "" {
		t.Skip("RSCHED_TEST_SKIP_E2E environment variable set")
	}
	prjRoot := testsupport.ProjectRoot(t)
	resticBinary := filepath.Join(prjRoot, "bin", "restic")
	if !testsupport.PathExists(t, resticBinary) {
		t.Skipf("%s does not exist.", resticBinary)
	}

	repo := filepath.Join(testsupport.TempDir(t), "repo")
	rsched := &cmd.RSched{
		Scheduler: &restic.Scheduler{},
	}
	cfg := cmd.Config{
		BackupPath:         prjRoot,
		BackupSchedule:     restic.ScheduleOnce,
		ResticPasswordFile: filepath.Join("testdata", "restic_password_file"),
		ResticRepository:   repo,
		ResticBinary:       resticBinary,
	}
	go rsched.Run(cfg)
	defer rsched.Shutdown()

	assert.Eventually(t, func() bool {
		// TODO: This is a crude way that restic is completed. Maybe we find
		// something better later on.
		return testsupport.DirNotEmpty(t, filepath.Join(repo, "snapshots"))
	}, 15*time.Second, 100*time.Millisecond)
}
