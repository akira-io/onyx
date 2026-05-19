# paths

Per-application platform paths — config, data, cache, logs. Follows macOS conventions (`~/Library/...`), XDG Base Directory on Linux, and Known Folders on Windows.

```go
import "github.com/akira-io/onyx/paths"
```

## API

| Symbol | Kind | Summary |
|--------|------|---------|
| `For(applicationName string) *AppPaths` | func | Build a path resolver for an application. |
| `(*AppPaths).Name() string` | method | The configured application name. |
| `(*AppPaths).Config() (string, error)` | method | Per-user configuration directory. |
| `(*AppPaths).Data() (string, error)` | method | Per-user data directory (state shared across machines). |
| `(*AppPaths).Cache() (string, error)` | method | Per-user cache directory (regenerable). |
| `(*AppPaths).Logs() (string, error)` | method | Per-user log directory. |
| `ErrMissingApplicationName` | sentinel | `For("")` was called and a path method was invoked. |

`AppPaths` is a struct holding the application name and the cached `osinfo.Platform`. It is safe to share across goroutines (immutable after construction).

## Path table

| Method | macOS | Linux | Windows |
|--------|-------|-------|---------|
| `Config()` | `~/Library/Application Support/<name>` | `$XDG_CONFIG_HOME/<name>` or `~/.config/<name>` | `%APPDATA%\<name>` |
| `Data()` | `~/Library/Application Support/<name>` | `$XDG_DATA_HOME/<name>` or `~/.local/share/<name>` | `%APPDATA%\<name>` |
| `Cache()` | `~/Library/Caches/<name>` | `$XDG_CACHE_HOME/<name>` or `~/.cache/<name>` | `%LOCALAPPDATA%\<name>\Cache` |
| `Logs()` | `~/Library/Logs/<name>` | `$XDG_STATE_HOME/<name>/logs` or `~/.local/state/<name>/logs` | `%LOCALAPPDATA%\<name>\Logs` |

`Config` and `Data` collapse to the same path on macOS and Windows — those platforms do not distinguish the two. On Linux they are separate (XDG defines distinct base directories).

## Examples

```go
import (
    "os"
    "path/filepath"

    "github.com/akira-io/onyx/paths"
)

app := paths.For("hyperion")

config, err := app.Config()
if err != nil {
    return err
}
data, _   := app.Data()
cache, _  := app.Cache()
logs, _   := app.Logs()

os.MkdirAll(config, 0o755)
os.WriteFile(filepath.Join(config, "settings.json"), []byte("{}"), 0o644)
```

For a long-lived app, build `AppPaths` once and share it:

```go
var (
    app    = paths.For("hyperion")
    config = mustPath(app.Config())
    cache  = mustPath(app.Cache())
)
```

## Behaviour

- Paths are **constructed**, not created. The caller calls `os.MkdirAll(path, 0o755)` before writing.
- `For("")` does not fail — the empty name is captured. The error comes from the first `Config()`/`Data()`/etc. call, which checks `requireName()`.
- XDG resolution: when `$XDG_CONFIG_HOME` (etc.) is set and non-empty, it wins. Otherwise the fallback is rooted at `$HOME`.
- Linux logs use `$XDG_STATE_HOME` (or `~/.local/state`) per the 2021 XDG update.
- Unknown platforms fall back to `os.UserConfigDir()` / `os.UserCacheDir()` from the stdlib.

## Errors

- `ErrMissingApplicationName` — empty `name` reached a path-returning method.
- Wrapped error from `os.UserHomeDir()`, `os.UserConfigDir()`, `os.UserCacheDir()`, or `os.Getenv` lookups. Format: `fmt.Errorf("resolve <thing>: %w", err)`.

## Dependencies

- `osinfo` — `Current()` selects the per-OS code path.
- `os.UserHomeDir`, `os.UserConfigDir`, `os.UserCacheDir`, `os.Getenv` from stdlib.

## Related packages

- [`files`](21-files.md) — open the directory you just resolved.
- [`keyring`](22-keyring.md) — keep secrets in the OS keychain, non-secrets in `paths.For(...)`.

## Cross-module parity

Mirrors the Rust crate's `paths` module one-to-one: same XDG semantics, same macOS `~/Library/...` paths, same Windows `Known Folder` mapping, same `Logs` subdirectory convention.

---

Navigation: [← OS info](24-osinfo.md) · **Paths** · [Shell →](26-shell.md)
