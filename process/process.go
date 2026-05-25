package process

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/akira-io/onyx/osinfo"
	"github.com/akira-io/onyx/shell"
)

// ErrEmptyPath is returned when Relaunch receives no application path.
var ErrEmptyPath = errors.New("process: application path required")

// Relaunch starts a fresh instance of the application at the given path. The
// caller quits the current process afterwards. On macOS the path is the
// `.app` bundle, on Windows the executable, and on Linux the binary.
func Relaunch(applicationPath string) error {
	if applicationPath == "" {
		return ErrEmptyPath
	}
	platform := osinfo.Current()
	switch {
	case platform.IsDarwin():
		opener, err := shell.NewResolver().Lookups([]string{"open", "/usr/bin/open"}).Resolve()
		if err != nil {
			return fmt.Errorf("process: locate open: %w", err)
		}
		return start(opener, "-n", applicationPath)
	case platform.IsWindows():
		return start("cmd", "/c", "start", "", applicationPath)
	default:
		return start(applicationPath)
	}
}

func start(name string, args ...string) error {
	if err := exec.Command(name, args...).Start(); err != nil {
		return fmt.Errorf("process: relaunch: %w", err)
	}
	return nil
}
