package clipboard

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/akira-io/onyx/osinfo"
)

var ErrClipboardUnavailable = errors.New("clipboard: no supported backend available")

func Read() (string, error) {
	platform := osinfo.Current()
	switch {
	case platform.IsDarwin():
		return runReader("pbpaste")
	case platform.IsWindows():
		return runReader("powershell", "-NoProfile", "-Command", "Get-Clipboard")
	case platform.IsLinux():
		for _, b := range linuxReaders() {
			out, err := runReader(b.cmd, b.args...)
			if err == nil {
				return out, nil
			}
		}
		return "", ErrClipboardUnavailable
	}
	return "", ErrClipboardUnavailable
}

func Write(text string) error {
	platform := osinfo.Current()
	switch {
	case platform.IsDarwin():
		return runWriter(text, "pbcopy")
	case platform.IsWindows():
		return runWriter(text, "powershell", "-NoProfile", "-Command", "Set-Clipboard")
	case platform.IsLinux():
		for _, b := range linuxWriters() {
			if err := runWriter(text, b.cmd, b.args...); err == nil {
				return nil
			}
		}
		return ErrClipboardUnavailable
	}
	return ErrClipboardUnavailable
}

type backend struct {
	cmd  string
	args []string
}

func linuxReaders() []backend {
	return []backend{
		{"wl-paste", []string{"--no-newline"}},
		{"xclip", []string{"-selection", "clipboard", "-o"}},
		{"xsel", []string{"--clipboard", "--output"}},
	}
}

func linuxWriters() []backend {
	return []backend{
		{"wl-copy", nil},
		{"xclip", []string{"-selection", "clipboard"}},
		{"xsel", []string{"--clipboard", "--input"}},
	}
}

func runReader(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return "", fmt.Errorf("clipboard read via %s: %w", name, err)
	}
	return strings.TrimRight(string(out), "\n"), nil
}

func runWriter(text, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clipboard write via %s: %w", name, err)
	}
	return nil
}
