package paths

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/akira-io/onyx/osinfo"
)

var ErrMissingApplicationName = errors.New("paths: application name is required")

type AppPaths struct {
	applicationName string
	platform        osinfo.Platform
}

func For(applicationName string) *AppPaths {
	return &AppPaths{
		applicationName: applicationName,
		platform:        osinfo.Current(),
	}
}

func (a *AppPaths) Name() string {
	return a.applicationName
}

func (a *AppPaths) Config() (string, error) {
	if err := a.requireName(); err != nil {
		return "", err
	}
	switch {
	case a.platform.IsDarwin():
		return joinUnderHome("Library/Application Support", a.applicationName)
	case a.platform.IsLinux():
		return resolveXDG("XDG_CONFIG_HOME", ".config", a.applicationName)
	case a.platform.IsWindows():
		return joinUnderEnv("APPDATA", a.applicationName)
	}
	return fallbackUserConfigDir(a.applicationName)
}

func (a *AppPaths) Data() (string, error) {
	if err := a.requireName(); err != nil {
		return "", err
	}
	switch {
	case a.platform.IsDarwin():
		return joinUnderHome("Library/Application Support", a.applicationName)
	case a.platform.IsLinux():
		return resolveXDG("XDG_DATA_HOME", ".local/share", a.applicationName)
	case a.platform.IsWindows():
		return joinUnderEnv("APPDATA", a.applicationName)
	}
	return fallbackUserConfigDir(a.applicationName)
}

func (a *AppPaths) Cache() (string, error) {
	if err := a.requireName(); err != nil {
		return "", err
	}
	switch {
	case a.platform.IsDarwin():
		return joinUnderHome("Library/Caches", a.applicationName)
	case a.platform.IsLinux():
		return resolveXDG("XDG_CACHE_HOME", ".cache", a.applicationName)
	case a.platform.IsWindows():
		return joinUnderEnv("LOCALAPPDATA", filepath.Join(a.applicationName, "Cache"))
	}
	return fallbackUserCacheDir(a.applicationName)
}

func (a *AppPaths) Logs() (string, error) {
	if err := a.requireName(); err != nil {
		return "", err
	}
	switch {
	case a.platform.IsDarwin():
		return joinUnderHome("Library/Logs", a.applicationName)
	case a.platform.IsLinux():
		return resolveXDG("XDG_STATE_HOME", ".local/state", filepath.Join(a.applicationName, "logs"))
	case a.platform.IsWindows():
		return joinUnderEnv("LOCALAPPDATA", filepath.Join(a.applicationName, "Logs"))
	}
	return fallbackUserCacheDir(filepath.Join(a.applicationName, "logs"))
}

func (a *AppPaths) requireName() error {
	if a.applicationName == "" {
		return ErrMissingApplicationName
	}
	return nil
}

func joinUnderHome(segments ...string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	return filepath.Join(append([]string{home}, segments...)...), nil
}

func joinUnderEnv(envVar string, suffix string) (string, error) {
	base := os.Getenv(envVar)
	if base == "" {
		return "", fmt.Errorf("resolve %s: environment variable not set", envVar)
	}
	return filepath.Join(base, suffix), nil
}

func resolveXDG(envVar string, fallbackSubdir string, suffix string) (string, error) {
	if explicit := os.Getenv(envVar); explicit != "" {
		return filepath.Join(explicit, suffix), nil
	}
	return joinUnderHome(fallbackSubdir, suffix)
}

func fallbackUserConfigDir(suffix string) (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}
	return filepath.Join(base, suffix), nil
}

func fallbackUserCacheDir(suffix string) (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("resolve user cache dir: %w", err)
	}
	return filepath.Join(base, suffix), nil
}
