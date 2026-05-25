package machineid

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/akira-io/onyx/keyring"
)

const idAccount = "machine-id"

// ErrEmptyApplication is returned when GetOrCreate receives no application name.
var ErrEmptyApplication = errors.New("machineid: application name required")

// GetOrCreate returns a stable identifier for the current machine scoped to the
// application, generating and persisting one in the system keyring on first use.
func GetOrCreate(application string) (string, error) {
	if application == "" {
		return "", ErrEmptyApplication
	}
	if existing, err := keyring.Get(application, idAccount); err == nil && existing != "" {
		return existing, nil
	}
	id, err := newID()
	if err != nil {
		return "", err
	}
	if err := keyring.Set(application, idAccount, id); err != nil {
		return "", fmt.Errorf("machineid: persist identity: %w", err)
	}
	return id, nil
}

func newID() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("machineid: generate: %w", err)
	}
	return hex.EncodeToString(buf), nil
}
