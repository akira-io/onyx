package process

import (
	"errors"
	"testing"
)

func TestRelaunch_RequiresPath(t *testing.T) {
	if err := Relaunch(""); !errors.Is(err, ErrEmptyPath) {
		t.Fatalf("expected ErrEmptyPath, got %v", err)
	}
}
