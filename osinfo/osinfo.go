package osinfo

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Platform struct {
	identifier string
}

const (
	darwinIdentifier  = "darwin"
	linuxIdentifier   = "linux"
	windowsIdentifier = "windows"
)

func Current() Platform {
	return Platform{identifier: runtime.GOOS}
}

func (p Platform) IsDarwin() bool {
	return p.identifier == darwinIdentifier
}

func (p Platform) IsLinux() bool {
	return p.identifier == linuxIdentifier
}

func (p Platform) IsWindows() bool {
	return p.identifier == windowsIdentifier
}

func (p Platform) String() string {
	return p.identifier
}

func ExecutableExtension() string {
	if Current().IsWindows() {
		return ".exe"
	}
	return ""
}

// Hostname returns the operating system host name, or an empty string when it
// cannot be determined.
func Hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return ""
	}
	return name
}

// DeviceName returns a human-friendly name for the current machine, the name a
// user recognizes in system settings, or an empty string when it cannot be
// determined. It falls back to Hostname when no friendlier source is available.
func DeviceName() string {
	platform := Current()
	switch {
	case platform.IsDarwin():
		if name := commandOutput("scutil", "--get", "ComputerName"); name != "" {
			return name
		}
	case platform.IsWindows():
		if name := strings.TrimSpace(os.Getenv("COMPUTERNAME")); name != "" {
			return name
		}
	case platform.IsLinux():
		if name := commandOutput("hostnamectl", "--pretty"); name != "" {
			return name
		}
	}
	return Hostname()
}

func commandOutput(name string, args ...string) string {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
