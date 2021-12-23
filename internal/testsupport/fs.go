package testsupport

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TempDir creates a temporary directory and makes sure it gets removed at
// the end of the test.
func TempDir(t *testing.T) string {
	t.Helper()

	prefix := strings.ReplaceAll(t.Name(), string(os.PathSeparator), "_")
	tempDir, err := os.MkdirTemp("", prefix)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Error(err)
		}
	})
	return tempDir
}

// ProjectRoot searches from the current working directory upwards until it
// finds a go.mod file. The directory containing the go.mod file is then
// assumed to be the root of this project.
func ProjectRoot(t *testing.T) string {
	t.Helper()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	cur, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for cur != home {
		goModFile := filepath.Join(cur, "go.mod")
		s, err := os.Stat(goModFile)
		if os.IsNotExist(err) {
			cur = filepath.Dir(cur)
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		if !s.Mode().IsRegular() {
			t.Fatalf("not a file: %s", cur)
		}
		return cur
	}
	t.Fatalf("project root not found: reached: %s", home)
	return ""
}

// PathExists checks if a file or directory exists at path.
func PathExists(t *testing.T, path string) bool {
	t.Helper()

	_, err := os.Stat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	if err != nil {
		t.Errorf("path exists: unexpected error: %v", err)
	}
	return true
}

// DirNotEmpty checks that the directory is not empty.
func DirNotEmpty(t *testing.T, path string) bool {
	if !PathExists(t, path) {
		return false
	}
	es, err := os.ReadDir(path)
	if err != nil {
		t.Errorf("read dir: %v", err)
	}
	return len(es) > 0
}
