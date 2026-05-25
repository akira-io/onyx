package shell

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/akira-io/onyx/osinfo"
)

var ErrBinaryNotFound = errors.New("shell: binary not found")

type Resolver struct {
	lookups []string
}

func NewResolver() Resolver {
	return Resolver{}
}

func (r Resolver) Lookup(target string) Resolver {
	if target == "" {
		return r
	}
	r.lookups = append(r.lookups, target)
	return r
}

func (r Resolver) Lookups(targets []string) Resolver {
	for _, t := range targets {
		r = r.Lookup(t)
	}
	return r
}

func (r Resolver) Resolve() (string, error) {
	for _, t := range r.lookups {
		if isPathLike(t) {
			if isExecutableFile(t) {
				return t, nil
			}
			continue
		}
		if absolute, err := exec.LookPath(t); err == nil {
			return absolute, nil
		}
	}
	return "", ErrBinaryNotFound
}

func isPathLike(s string) bool {
	if strings.ContainsAny(s, `/\`) {
		return true
	}
	if len(s) >= 2 && s[1] == ':' {
		return true
	}
	return false
}

func ListNpmGlobalBinDirs() []string {
	platform := osinfo.Current()
	home, _ := os.UserHomeDir()
	if platform.IsWindows() {
		out := []string{}
		if appData := os.Getenv("APPDATA"); appData != "" {
			out = append(out, filepath.Join(appData, "npm"))
		}
		return out
	}
	out := []string{}
	if home != "" {
		out = append(out,
			filepath.Join(home, ".npm-global", "bin"),
			filepath.Join(home, ".local", "share", "npm", "bin"),
		)
	}
	return out
}

func ListUserLocalBinDirs() []string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return []string{}
	}
	return []string{
		filepath.Join(home, ".local", "bin"),
		filepath.Join(home, "bin"),
	}
}

func ListSystemBinDirs() []string {
	platform := osinfo.Current()
	if platform.IsWindows() {
		return []string{}
	}
	if platform.IsDarwin() {
		return []string{
			"/usr/local/bin",
			"/opt/homebrew/bin",
			"/usr/bin",
		}
	}
	return []string{
		"/usr/local/bin",
		"/usr/bin",
	}
}

func ListWindowsApplicationDirs(applicationName string) []string {
	if !osinfo.Current().IsWindows() || applicationName == "" {
		return []string{}
	}
	out := []string{}
	if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
		out = append(out, filepath.Join(localAppData, "Programs", applicationName))
	}
	if programFiles := os.Getenv("ProgramFiles"); programFiles != "" {
		out = append(out, filepath.Join(programFiles, applicationName))
	}
	if programFilesX86 := os.Getenv("ProgramFiles(x86)"); programFilesX86 != "" {
		out = append(out, filepath.Join(programFilesX86, applicationName))
	}
	return out
}

func isExecutableFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}

// LoginPath returns the PATH as seen by the user's interactive login shell.
// GUI applications launched outside a terminal on macOS and Linux do not
// inherit the shell's PATH additions; this recovers them by asking the login
// shell. On Windows, and when the shell cannot be queried, it falls back to the
// current process PATH.
func LoginPath() string {
	if osinfo.Current().IsWindows() {
		return os.Getenv("PATH")
	}
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}
	out, err := exec.Command(shell, "-l", "-c", "echo $PATH").Output()
	if err != nil {
		return os.Getenv("PATH")
	}
	if trimmed := strings.TrimSpace(string(out)); trimmed != "" {
		return trimmed
	}
	return os.Getenv("PATH")
}

// EnrichedEnviron returns the current environment with PATH replaced by the
// union of the process PATH and the login shell PATH, preserving order and
// dropping duplicates. Use it when spawning child processes that must find
// user-installed tools.
func EnrichedEnviron() []string {
	merged := mergePath(os.Getenv("PATH"), LoginPath())
	env := os.Environ()
	for i, kv := range env {
		if strings.HasPrefix(kv, "PATH=") {
			env[i] = "PATH=" + merged
			return env
		}
	}
	return append(env, "PATH="+merged)
}

func mergePath(first, second string) string {
	seen := make(map[string]struct{})
	var ordered []string
	for _, group := range []string{first, second} {
		for _, segment := range filepath.SplitList(group) {
			if segment == "" {
				continue
			}
			if _, ok := seen[segment]; ok {
				continue
			}
			seen[segment] = struct{}{}
			ordered = append(ordered, segment)
		}
	}
	return strings.Join(ordered, string(filepath.ListSeparator))
}
