package shell

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

func TestLoginPath_NonEmptyAndContainsProcessPath(t *testing.T) {
	got := LoginPath()
	if got == "" {
		t.Fatal("LoginPath should not be empty")
	}
}

func TestEnrichedEnviron_KeepsSinglePathEntry(t *testing.T) {
	env := EnrichedEnviron()
	pathCount := 0
	for _, kv := range env {
		if strings.HasPrefix(kv, "PATH=") {
			pathCount++
		}
	}
	if pathCount != 1 {
		t.Fatalf("expected exactly one PATH entry, got %d", pathCount)
	}
}

func TestMergePath_DropsDuplicates(t *testing.T) {
	sep := string(filepath.ListSeparator)
	got := mergePath("a"+sep+"b", "b"+sep+"c")
	want := "a" + sep + "b" + sep + "c"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestResolve_FailsWhenNothingMatches(t *testing.T) {
	_, err := NewResolver().
		Lookup("definitely-not-a-real-binary-xyz").
		Lookup("/definitely/not/a/path/binary").
		Resolve()
	if !errors.Is(err, ErrBinaryNotFound) {
		t.Fatalf("expected ErrBinaryNotFound, got %v", err)
	}
}

func TestResolve_FindsExplicitPath(t *testing.T) {
	dir := t.TempDir()
	binary := filepath.Join(dir, "fakebin")
	if err := os.WriteFile(binary, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("write fake binary: %v", err)
	}

	resolved, err := NewResolver().
		Lookup("definitely-not-a-real-binary-xyz").
		Lookup(binary).
		Resolve()
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if resolved != binary {
		t.Fatalf("expected %q, got %q", binary, resolved)
	}
}

func TestResolve_PrefersFirstMatch(t *testing.T) {
	resolved, err := NewResolver().
		Lookup("sh").
		Lookup("/definitely/not/a/path/binary").
		Resolve()
	if err != nil {
		t.Skipf("sh not available on this system: %v", err)
	}
	if !filepath.IsAbs(resolved) {
		t.Fatalf("expected absolute path, got %q", resolved)
	}
}

func TestResolver_IgnoresEmptyInputs(t *testing.T) {
	_, err := NewResolver().Lookup("").Lookup("").Resolve()
	if !errors.Is(err, ErrBinaryNotFound) {
		t.Fatalf("expected ErrBinaryNotFound, got %v", err)
	}
}

func TestResolver_LookupsBulk(t *testing.T) {
	dir := t.TempDir()
	binary := filepath.Join(dir, "fakebin")
	if err := os.WriteFile(binary, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("write fake binary: %v", err)
	}

	resolved, err := NewResolver().
		Lookup("definitely-not-a-real-binary-xyz").
		Lookups([]string{"", "/definitely/not/here", binary}).
		Resolve()
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if resolved != binary {
		t.Fatalf("expected %q, got %q", binary, resolved)
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

func TestIsPathLike(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"claude", false},
		{"/opt/homebrew/bin/claude", true},
		{"./bin/foo", true},
		{`C:\Program Files\app.exe`, true},
		{"", false},
	}
	for _, c := range cases {
		if got := isPathLike(c.in); got != c.want {
			t.Errorf("isPathLike(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}
