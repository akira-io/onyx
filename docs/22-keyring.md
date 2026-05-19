# keyring

Store, retrieve, and delete secrets in the system credential store via per-platform backends.

```go
import "github.com/akira-io/onyx/keyring"
```

## API

| Symbol | Kind | Summary |
|--------|------|---------|
| `Set(service, account, secret string) error` | func | Store or overwrite a secret. |
| `Get(service, account string) (string, error)` | func | Read the stored secret. |
| `Delete(service, account string) error` | func | Remove the entry. |
| `ErrKeyringUnavailable` | sentinel | No supported backend reachable. |
| `ErrNotFound` | sentinel | Entry is not in the credential store. |
| `ErrEmptyService` | sentinel | Service string empty after trim. |
| `ErrEmptyAccount` | sentinel | Account string empty after trim. |

Backend failures wrap `*exec.ExitError` via `fmt.Errorf("keyring %s via %s: %w", action, backend, err)`.

## Platform backends

| Platform | Backend |
|----------|---------|
| macOS | `security add-generic-password / find-generic-password / delete-generic-password` |
| Linux | `secret-tool store / lookup / clear` (libsecret over D-Bus) |
| Windows | `cmdkey /generic` for write/delete, PowerShell + `CredentialManager` for read |

The Linux backend requires a running Secret Service implementation (`gnome-keyring`, `kwallet5`, `KeePassXC` with the Secret Service plugin). On headless CI, install `gnome-keyring` and start a session daemon — or accept that `Get`/`Set` will return `ErrKeyringUnavailable`.

The Windows read path uses PowerShell because `cmdkey` does not print stored passwords. The script imports `CredentialManager` (a PowerShell Gallery module) — when the module is absent, the call returns `ErrKeyringUnavailable` with PowerShell exit code `2`.

## Examples

```go
import (
    "errors"
    "github.com/akira-io/onyx/keyring"
)

keyring.Set("io.akira.unified-dev", "kid@example.com", "hunter2")

secret, err := keyring.Get("io.akira.unified-dev", "kid@example.com")
switch {
case err == nil:
    useToken(secret)
case errors.Is(err, keyring.ErrNotFound):
    promptLogin()
default:
    return err
}

keyring.Delete("io.akira.unified-dev", "kid@example.com")
```

Service / account convention:

- `service` — application identifier, typically the bundle id (`io.akira.unified-dev`).
- `account` — customer identifier within that app (email, user id, or `"default"`).

This matches the macOS `kSecAttrService` / `kSecAttrAccount` pair and the Linux `secret-tool` `service` / `account` attribute schema.

## Behaviour

- Validation runs **before** any IPC. Empty service or account returns `ErrEmptyService` / `ErrEmptyAccount` without touching the OS keychain.
- The macOS backend uses `security add-generic-password -U` so a second `Set` overwrites the previous secret without prompting.
- Linux backend writes the secret to `secret-tool`'s stdin (not as a CLI argument) so it does not appear in process listings.
- Windows backend joins `service:account` into the credential target name — `Delete` and `Get` reconstruct the same join.
- Trailing newlines from the backend output are stripped on `Get`.
- macOS exit code 44, Linux exit code 1, and PowerShell exit code 3 are translated to `ErrNotFound`.

## Errors

- `ErrEmptyService` / `ErrEmptyAccount` — caller passed empty (or whitespace-only) values.
- `ErrNotFound` — backend reported the entry is missing.
- `ErrKeyringUnavailable` — no backend reachable (unsupported OS, or Windows without the `CredentialManager` PowerShell module).
- Wrapped error — backend was located but the call failed. `errors.Unwrap` reveals the underlying `*exec.ExitError`.

## Security notes

- `onyx` does not zero the returned `string`. The Go runtime makes proper zeroing difficult — assume the secret may linger in memory until GC.
- Secrets are passed as CLI arguments on macOS and Windows (`security -w` and `cmdkey /pass`). On a system you do not control, this may be visible in `ps`. The Linux backend writes via stdin to avoid this.
- The Windows backend invokes PowerShell with `-NoProfile` to skip user-customised profile scripts.

## Dependencies

- `osinfo` — `Current()` selects the per-OS code path.

## Related packages

- [`paths`](25-paths.md) — store non-secret config alongside the keychain entry.

## Cross-module parity

Mirrors the Rust crate's `keyring` module one-to-one: same backend choice, same exit-code parsing, same stdin handling.

---

Navigation: [← Files](21-files.md) · **Keyring** · [Notify →](23-notify.md)
