# notify

Show desktop notifications via per-platform backends.

```go
import "github.com/akira-io/onyx/notify"
```

## API

| Symbol | Kind | Summary |
|--------|------|---------|
| `Show(title, body string) error` | func | Display a notification with title and body. |
| `ErrEmptyTitle` | sentinel | Title was empty (after trim). |
| `ErrNotifyUnavailable` | sentinel | No supported backend is reachable. |

Backend failures wrap `*exec.ExitError` via `fmt.Errorf("notify via %s: %w", backend, err)`.

## Platform backends

| Platform | Backend (priority) |
|----------|--------------------|
| macOS | `osascript -e 'display notification "<body>" with title "<title>"'` |
| Linux | `notify-send <title> <body>` |
| Windows | `BurntToast` PowerShell module first; falls back to `msg *` |

`BurntToast` is a community PowerShell module that renders proper Windows 10/11 toasts. When it is not installed, `notify` falls back to `msg *`, which displays a blocking message box — workable but visually different from a toast.

Linux requires `notify-send` (`libnotify-bin` package on Debian/Ubuntu). The Wayland and X11 notification daemons listen on the same `org.freedesktop.Notifications` D-Bus service.

## Examples

```go
import "github.com/akira-io/onyx/notify"

if err := notify.Show("Build complete", "Spectra v0.9.0 is ready to install."); err != nil {
    log.Printf("notify failed: %v", err)
}
```

Best-effort branching when the user might not have a notifier:

```go
err := notify.Show("Heads up", "Update available")
if errors.Is(err, notify.ErrNotifyUnavailable) {
    fallBackToLogging()
}
```

## Behaviour

- Title is trimmed and rejected if empty — empty notifications are useless and most backends silently drop them.
- Body may be empty. The notifier displays a title-only toast.
- macOS: arguments are passed through an AppleScript-safe quoter that escapes `\` and `"` — pasting user input as a title is safe.
- Windows: arguments are passed through a PowerShell-safe quoter that doubles `'` — `it's fine` becomes `'it''s fine'`.
- `notify-send` accepts the title and body as plain arguments; no extra escaping is needed because `exec.Command` does not invoke a shell.

## Errors

- `ErrEmptyTitle` — caller passed an empty (or whitespace-only) title.
- `ErrNotifyUnavailable` — Windows with neither BurntToast nor `msg`; or running on an unsupported OS.
- Wrapped error — backend was reached but failed.

## Dependencies

- `osinfo` — `Current()` selects the per-OS code path.

## Related packages

- [`files`](21-files.md) — open a follow-up action target from the notification handler.
- [`shell`](26-shell.md) — detect whether `notify-send` is on `PATH` before calling.

## Cross-module parity

Mirrors the Rust crate's `notify` module one-to-one: same backends, same priority order, same quoting rules.

---

Navigation: [← Keyring](22-keyring.md) · **Notify** · [OS info →](24-osinfo.md)
