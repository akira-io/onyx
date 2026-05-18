package clipboard

import (
	"errors"
	"testing"
)

func TestRoundTrip_PreservesText(t *testing.T) {
	const sample = "onyx-clipboard-test"
	if err := Write(sample); err != nil {
		if errors.Is(err, ErrClipboardUnavailable) {
			t.Skip("no clipboard backend available")
		}
		t.Skipf("write failed: %v", err)
	}
	got, err := Read()
	if err != nil {
		t.Skipf("read failed: %v", err)
	}
	if got != sample {
		t.Fatalf("got %q, want %q", got, sample)
	}
}

func TestWrite_EmptyStringIsAllowed(t *testing.T) {
	if err := Write(""); err != nil {
		if errors.Is(err, ErrClipboardUnavailable) {
			t.Skip("no clipboard backend available")
		}
		t.Skipf("write empty failed: %v", err)
	}
}

func TestLinuxBackends_AreDeclaredInPriorityOrder(t *testing.T) {
	readers := linuxReaders()
	writers := linuxWriters()
	if len(readers) != 3 || len(writers) != 3 {
		t.Fatalf("expected 3 readers and 3 writers, got %d/%d", len(readers), len(writers))
	}
	if readers[0].cmd != "wl-paste" || writers[0].cmd != "wl-copy" {
		t.Fatalf("Wayland should be first; got read=%q write=%q", readers[0].cmd, writers[0].cmd)
	}
}
