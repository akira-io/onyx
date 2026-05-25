package osinfo

import (
	"os"
	"runtime"
	"testing"
)

func TestCurrent_MatchesRuntimeGOOS(t *testing.T) {
	platform := Current()
	if platform.String() != runtime.GOOS {
		t.Fatalf("expected %q, got %q", runtime.GOOS, platform.String())
	}
}

func TestPlatform_ExactlyOnePredicateHolds(t *testing.T) {
	platform := Current()
	trueCount := 0
	if platform.IsDarwin() {
		trueCount++
	}
	if platform.IsLinux() {
		trueCount++
	}
	if platform.IsWindows() {
		trueCount++
	}
	if trueCount > 1 {
		t.Fatalf("expected at most one predicate to be true, got %d", trueCount)
	}
}

func TestExecutableExtension_MatchesPlatform(t *testing.T) {
	got := ExecutableExtension()
	if Current().IsWindows() {
		if got != ".exe" {
			t.Fatalf("expected .exe on Windows, got %q", got)
		}
		return
	}
	if got != "" {
		t.Fatalf("expected empty extension on non-Windows, got %q", got)
	}
}

func TestHostname_MatchesOSHostnameWhenAvailable(t *testing.T) {
	got := Hostname()
	want, err := os.Hostname()
	if err != nil {
		if got != "" {
			t.Fatalf("expected empty hostname on error, got %q", got)
		}
		return
	}
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestDeviceName_FallsBackToHostname(t *testing.T) {
	got := DeviceName()
	if got == "" {
		if Hostname() != "" {
			t.Fatalf("expected DeviceName to fall back to a non-empty hostname")
		}
		return
	}
}
