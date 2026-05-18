package shell

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/akira-io/onyx/osinfo"
)

var ErrBinaryNotFound = errors.New("shell: binary not found")

type ResolutionSource int

const (
	SourceUnknown ResolutionSource = iota
	SourcePath
	SourceCandidate
)

func (s ResolutionSource) String() string {
	switch s {
	case SourcePath:
		return "path"
	case SourceCandidate:
		return "candidate"
	default:
		return "unknown"
	}
}

type ResolvedExecutable struct {
	absolutePath string
	source       ResolutionSource
}

func (r ResolvedExecutable) AbsolutePath() string {
	return r.absolutePath
}

func (r ResolvedExecutable) Source() ResolutionSource {
	return r.source
}

type Candidates struct {
	names      []string
	candidates []string
}

func NewCandidates() Candidates {
	return Candidates{}
}

func (c Candidates) WithName(name string) Candidates {
	if name == "" {
		return c
	}
	c.names = append(c.names, name)
	return c
}

func (c Candidates) WithCandidate(path string) Candidates {
	if path == "" {
		return c
	}
	c.candidates = append(c.candidates, path)
	return c
}

func (c Candidates) WithCandidates(paths []string) Candidates {
	for _, path := range paths {
		c = c.WithCandidate(path)
	}
	return c
}

func (c Candidates) Resolve() (ResolvedExecutable, error) {
	for _, name := range c.names {
		if absolute, err := exec.LookPath(name); err == nil {
			return ResolvedExecutable{absolutePath: absolute, source: SourcePath}, nil
		}
	}
	for _, candidate := range c.candidates {
		if isExecutableFile(candidate) {
			return ResolvedExecutable{absolutePath: candidate, source: SourceCandidate}, nil
		}
	}
	return ResolvedExecutable{}, ErrBinaryNotFound
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
