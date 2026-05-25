package osinfo

import (
	"os"
	"runtime"
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
