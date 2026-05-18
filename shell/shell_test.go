package shell

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"
)

func TestResolve_FailsWhenNothingMatches(t *testing.T) {
	candidates := NewResolver().
		Lookup("definitely-not-a-real-binary-xyz").
		Fallback("/definitely/not/a/path/binary")

	_, err := candidates.Resolve()
	if !errors.Is(err, ErrBinaryNotFound) {
		t.Fatalf("expected ErrBinaryNotFound, got %v", err)
	}
}

func TestResolve_FindsExplicitCandidateFile(t *testing.T) {
	dir := t.TempDir()
	binary := filepath.Join(dir, "fakebin")
	if err := os.WriteFile(binary, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("write fake binary: %v", err)
	}

	candidates := NewResolver().
		Lookup("definitely-not-a-real-binary-xyz").
		Fallback(binary)

	resolved, err := candidates.Resolve()
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if resolved.AbsolutePath() != binary {
		t.Fatalf("expected %q, got %q", binary, resolved.AbsolutePath())
	}
	if resolved.Source() != SourceFallback {
		t.Fatalf("expected SourceFallback, got %s", resolved.Source())
	}
}

func TestResolve_PrefersPATHOverCandidates(t *testing.T) {
	candidates := NewResolver().
		Lookup("sh").
		Fallback("/definitely/not/a/path/binary")

	resolved, err := candidates.Resolve()
	if err != nil {
		t.Skipf("sh not available on this system: %v", err)
	}
	if resolved.Source() != SourcePath {
		t.Fatalf("expected SourcePath, got %s", resolved.Source())
	}
}

func TestResolver_IgnoresEmptyInputs(t *testing.T) {
	candidates := NewResolver().
		Lookup("").
		Fallback("")

	_, err := candidates.Resolve()
	if !errors.Is(err, ErrBinaryNotFound) {
		t.Fatalf("expected ErrBinaryNotFound, got %v", err)
	}
}

func TestResolver_FallbacksAddsAll(t *testing.T) {
	dir := t.TempDir()
	binary := filepath.Join(dir, "fakebin")
	if err := os.WriteFile(binary, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("write fake binary: %v", err)
	}

	candidates := NewResolver().
		Lookup("definitely-not-a-real-binary-xyz").
		Fallbacks([]string{"", "/definitely/not/here", binary})

	resolved, err := candidates.Resolve()
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if resolved.AbsolutePath() != binary {
		t.Fatalf("expected %q, got %q", binary, resolved.AbsolutePath())
	}
}

func TestListUserLocalBinDirs_IncludesHomeLocal(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		t.Skipf("no home directory available: %v", err)
	}
	dirs := ListUserLocalBinDirs()
	if !slices.Contains(dirs, filepath.Join(home, ".local", "bin")) {
		t.Fatalf("expected ~/.local/bin in %v", dirs)
	}
}

func TestListSystemBinDirs_MatchesPlatform(t *testing.T) {
	dirs := ListSystemBinDirs()
	switch runtime.GOOS {
	case "windows":
		if len(dirs) != 0 {
			t.Fatalf("expected no system bin dirs on Windows, got %v", dirs)
		}
	case "darwin":
		if !slices.Contains(dirs, "/opt/homebrew/bin") {
			t.Fatalf("expected /opt/homebrew/bin on darwin, got %v", dirs)
		}
	case "linux":
		if !slices.Contains(dirs, "/usr/local/bin") {
			t.Fatalf("expected /usr/local/bin on linux, got %v", dirs)
		}
	}
}

func TestListWindowsApplicationDirs_EmptyOnNonWindows(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("test targets non-Windows platforms")
	}
	dirs := ListWindowsApplicationDirs("anything")
	if len(dirs) != 0 {
		t.Fatalf("expected empty slice on non-Windows, got %v", dirs)
	}
}

func TestListWindowsApplicationDirs_EmptyWhenNameMissing(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skipf("test targets Windows platform")
	}
	dirs := ListWindowsApplicationDirs("")
	if len(dirs) != 0 {
		t.Fatalf("expected empty slice when name empty, got %v", dirs)
	}
}

func TestListNpmGlobalBinDirs_ReturnsPlatformSpecificPaths(t *testing.T) {
	dirs := ListNpmGlobalBinDirs()
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			t.Skipf("APPDATA not set on this Windows host")
		}
		want := filepath.Join(appData, "npm")
		if !slices.Contains(dirs, want) {
			t.Fatalf("expected %q in %v", want, dirs)
		}
		return
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		t.Skipf("no home directory available: %v", err)
	}
	want := filepath.Join(home, ".npm-global", "bin")
	if !slices.Contains(dirs, want) {
		t.Fatalf("expected %q in %v", want, dirs)
	}
}
