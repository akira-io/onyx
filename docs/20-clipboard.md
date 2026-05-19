# clipboard

Read and write the system clipboard as plain text. Per-platform backends, no CGo, only stdlib `os/exec`.

```go
import "github.com/akira-io/onyx/clipboard"
```

## API

| Symbol | Kind | Summary |
|--------|------|---------|
| `Read() (string, error)` | func | Returns current clipboard text. |
| `Write(text string) error` | func | Replaces clipboard text. Empty string clears. |
| `ErrClipboardUnavailable` | sentinel | No supported backend is reachable. |

Backend failures are wrapped via `fmt.Errorf("clipboard %s via %s: %w", action, backend, err)`. Branch on the underlying `*exec.ExitError` if you need exit codes.

## Platform backends

| Platform | Backend (priority order) |
|----------|--------------------------|
| macOS | `pbcopy` / `pbpaste` (always present) |
| Windows | PowerShell hosting `System.Windows.Forms.Clipboard` (STA) |
| Linux | `wl-copy` / `wl-paste` (Wayland) → `xclip -selection clipboard` → `xsel --clipboard`. First successful one wins. |

On Linux, install one of `wl-clipboard`, `xclip`, or `xsel` to enable. The Wayland backend is preferred because it works on both Wayland and X11 sessions when XWayland is available.

## Behaviour

- Trailing newlines and carriage returns are stripped from the clipboard contents on read — backends emit them inconsistently.
- `Write("")` clears the clipboard on all platforms.
- The Windows backend passes the payload via the `ONYX_CLIP_TEXT` environment variable, not as a command-line argument, to avoid escaping issues with quotes, backticks, and newlines.

## Examples

```go
import "github.com/akira-io/onyx/clipboard"

if err := clipboard.Write("hello onyx"); err != nil {
    return err
}

text, err := clipboard.Read()
if err != nil {
    return err
}
fmt.Println(text)
```

Best-effort branching when the user might not have a clipboard tool installed:

```go
text, err := clipboard.Read()
switch {
case err == nil:
    fmt.Println("clipboard:", text)
case errors.Is(err, clipboard.ErrClipboardUnavailable):
    fmt.Fprintln(os.Stderr, "install wl-clipboard, xclip, or xsel")
default:
    return err
}
```

## Errors

- `ErrClipboardUnavailable` — no backend reachable. Only on Linux when `wl-clipboard`/`xclip`/`xsel` are all missing, or on unsupported OS targets.
- Wrapped `*exec.ExitError` — backend located but the call itself failed (binary missing permission, display server not running, sandbox preventing IPC).

## Dependencies

- `osinfo` — `Current()` selects the per-OS code path.

## Related packages

- [`files`](21-files.md) — opening files with the default app.
- [`shell`](26-shell.md) — resolving CLI binaries when you want to detect which backend is present.

## Cross-module parity

Mirrors the Rust crate's `clipboard` module one-to-one: same backend priority, same `Unavailable` semantics, same trimming behaviour.

---

Navigation: [← Release flow](05-release-flow.md) · **Clipboard** · [Files →](21-files.md)
