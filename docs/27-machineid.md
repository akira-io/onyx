# machineid

A stable, per-application identifier for the current machine, persisted in the system keyring so it survives restarts and reinstalls of the application data.

```go
import "github.com/akira-io/onyx/machineid"
```

## API

| Symbol | Kind | Summary |
|--------|------|---------|
| `GetOrCreate(application string) (string, error)` | func | Returns the machine identifier for the application, creating one on first use. |
| `ErrEmptyApplication` | error | Returned when `application` is empty. |

## Why

Desktop apps need a durable device identifier for licensing, telemetry, and device management. Generating one inline in every app leads to subtly different schemes and storage locations. `machineid` keeps the identifier in the OS keyring (Keychain, Secret Service, Credential Manager) under the application's namespace, so it is both stable and outside the app's own data directory.

## Behaviour

- First call generates a 256-bit random identifier (hex-encoded) and stores it.
- Subsequent calls return the stored value unchanged.
- The identifier is scoped to `application`; two applications get independent identifiers on the same machine.

## Errors

- `ErrEmptyApplication` when no application name is supplied.
- Wrapped keyring errors when the credential store is unavailable.

## Dependencies

- [keyring](./22-keyring.md) for secure persistence.

## Related modules

- [keyring](./22-keyring.md) backs storage.

---

Navigation: [← Shell](26-shell.md) · **Machine ID** · [Appearance →](28-appearance.md)
