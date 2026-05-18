package files

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/akira-io/onyx/osinfo"
)

var (
	ErrPathRequired        = errors.New("files: path is required")
	ErrUnsupportedPlatform = errors.New("files: unsupported platform")
)

func OpenPath(path string) error {
	if path == "" {
		return ErrPathRequired
	}
	platform := osinfo.Current()
	switch {
	case platform.IsDarwin():
		return runCommand("open", path)
	case platform.IsLinux():
		return runCommand("xdg-open", path)
	case platform.IsWindows():
		return runCommand("cmd", "/c", "start", "", path)
	}
	return ErrUnsupportedPlatform
}

func OpenURL(url string) error {
	if url == "" {
		return ErrPathRequired
	}
	platform := osinfo.Current()
	switch {
	case platform.IsDarwin():
		return runCommand("open", url)
	case platform.IsLinux():
		return runCommand("xdg-open", url)
	case platform.IsWindows():
		return runCommand("cmd", "/c", "start", "", url)
	}
	return ErrUnsupportedPlatform
}

func RevealInFileManager(path string) error {
	if path == "" {
		return ErrPathRequired
	}
	platform := osinfo.Current()
	switch {
	case platform.IsDarwin():
		return runCommand("open", "-R", path)
	case platform.IsLinux():
		return runCommand("xdg-open", filepath.Dir(path))
	case platform.IsWindows():
		return runCommand("explorer", "/select,"+path)
	}
	return ErrUnsupportedPlatform
}

func runCommand(name string, args ...string) error {
	if err := exec.Command(name, args...).Run(); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}
