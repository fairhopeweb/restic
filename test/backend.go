package test_helper

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/restic/restic"
	"github.com/restic/restic/backend"
	"github.com/restic/restic/backend/local"
	"github.com/restic/restic/server"
)

var TestPassword = "foobar"
var TestCleanup = flag.Bool("test.cleanup", true, "clean up after running tests (remove local backend directory with all content)")
var TestTempDir = flag.String("test.tempdir", "", "use this directory for temporary storage (default: system temp dir)")

func SetupBackend(t testing.TB) *server.Server {
	tempdir, err := ioutil.TempDir(*TestTempDir, "restic-test-")
	OK(t, err)

	// create repository below temp dir
	b, err := local.Create(filepath.Join(tempdir, "repo"))
	OK(t, err)

	// set cache dir
	err = os.Setenv("RESTIC_CACHE", filepath.Join(tempdir, "cache"))
	OK(t, err)

	return server.NewServer(b)
}

func TeardownBackend(t testing.TB, s *server.Server) {
	if !*TestCleanup {
		l := s.Backend().(*local.Local)
		t.Logf("leaving local backend at %s\n", l.Location())
		return
	}

	OK(t, s.Delete())
}

func SetupKey(t testing.TB, s *server.Server, password string) *server.Key {
	k, err := server.CreateKey(s, password)
	OK(t, err)

	return k
}

func SnapshotDir(t testing.TB, server *server.Server, path string, parent backend.ID) *restic.Snapshot {
	arch, err := restic.NewArchiver(server)
	OK(t, err)
	sn, _, err := arch.Snapshot(nil, []string{path}, parent)
	OK(t, err)
	return sn
}