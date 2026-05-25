package machineid

import (
	"testing"

	"github.com/akira-io/onyx/keyring"
)

func TestGetOrCreate_RequiresApplication(t *testing.T) {
	if _, err := GetOrCreate(""); err == nil {
		t.Fatal("expected ErrEmptyApplication for empty application")
	}
}

func TestGetOrCreate_StableAcrossCalls(t *testing.T) {
	const app = "onyx-machineid-test"
	first, err := GetOrCreate(app)
	if err != nil {
		t.Skipf("keyring unavailable: %v", err)
	}
	t.Cleanup(func() { _ = keyring.Delete(app, idAccount) })

	if _, err := keyring.Get(app, idAccount); err != nil {
		t.Skipf("keyring read backend unavailable, identity cannot persist: %v", err)
	}

	second, err := GetOrCreate(app)
	if err != nil {
		t.Fatalf("second GetOrCreate: %v", err)
	}
	if first != second {
		t.Fatalf("identity not stable: %q != %q", first, second)
	}
}

func TestNewID_IsLongHex(t *testing.T) {
	id, err := newID()
	if err != nil {
		t.Fatalf("newID: %v", err)
	}
	if len(id) != 64 {
		t.Fatalf("expected 64 hex chars, got %d (%q)", len(id), id)
	}
}
