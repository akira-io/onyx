# files

Open files / URLs with the user's default application and reveal a path inside the platform file manager.

```go
import "github.com/akira-io/onyx/files"
```

## API

| Symbol | Kind | Summary |
|--------|------|---------|
| `OpenPath(path string) error` | func | Open `path` with the OS default app. |
| `OpenURL(url string) error` | func | Open `url` in the default browser. |
| `RevealInFileManager(path string) error` | func | Highlight `path` inside Finder / Explorer / Files. |
| `ErrPathRequired` | sentinel | Empty path or URL. |
| `ErrUnsupportedPlatform` | sentinel | Running on a platform `onyx` does not target. |

Backend failures are wrapped via `fmt.Errorf("%s: %w", commandName, err)`.

## Platform backends

| Platform | `OpenPath` / `OpenURL` | `RevealInFileManager` |
|----------|------------------------|------------------------|
| macOS | `open <target>` | `open -R <path>` |
| Linux | `xdg-open <target>` | `xdg-open <parent-dir>` |
| Windows | `cmd /c start "" <target>` | `explorer /select,<path>` |

Linux does not have a portable "highlight this file" command, so `RevealInFileManager` falls back to opening the parent directory via `filepath.Dir(path)`. On macOS the `-R` flag tells Finder to select the file in its enclosing folder. On Windows, `/select` ensures Explorer opens with the file highlighted.

## Examples

```go
import "github.com/akira-io/onyx/files"

if err := files.OpenPath("/Users/me/Downloads/report.pdf"); err != nil {
    return err
}
files.OpenURL("https://akira.foundation")
files.RevealInFileManager("/Users/me/Downloads/report.pdf")
```

Wails integration:

```go
//go:wails:command
func Reveal(path string) error {
    return files.RevealInFileManager(path)
}
```

## Behaviour

- The spawned process is detached — `onyx` calls `Run` but the launched helper detaches on its own (system `open` / `xdg-open` / `cmd start`).
- An empty path or URL returns `ErrPathRequired` immediately (no spawn).
- On unsupported platforms (BSD, Solaris, …), `ErrUnsupportedPlatform` is returned.

## Errors

- `ErrPathRequired` — caller passed an empty path or empty URL.
- `ErrUnsupportedPlatform` — `osinfo.Current()` returned a value other than `darwin`/`linux`/`windows`.
- Wrapped `*exec.ExitError` — the spawned helper failed before detaching. Rare; usually means the helper binary is missing.

## Dependencies

- `osinfo` — `Current()` selects the per-OS code path.

## Related packages

- [`paths`](25-paths.md) — choose where to write the file you are about to reveal.
- [`shell`](26-shell.md) — resolve a custom binary if you need to override the default app.

## Cross-module parity

Mirrors the Rust crate's `files` module one-to-one: same commands, same arguments, same parent-directory fallback on Linux.

---

Navigation: [← Clipboard](20-clipboard.md) · **Files** · [Keyring →](22-keyring.md)
