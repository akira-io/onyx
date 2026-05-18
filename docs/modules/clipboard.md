# clipboard

Read and write the system clipboard as plain text. Per-platform backends, no CGo, no third-party deps.

## API

| Symbol | Kind | Summary |
| --- | --- | --- |
| `Read() (string, error)` | func | Returns current clipboard text. |
| `Write(text string) error` | func | Replaces clipboard text. Empty string clears. |
| `ErrClipboardUnavailable` | error | Returned when no supported backend is reachable. |

## Platform behavior

| Platform | Backend (in order tried) |
| --- | --- |
| macOS | `pbcopy` / `pbpaste` (always present). |
| Windows | PowerShell `Set-Clipboard` / `Get-Clipboard`. |
| Linux | `wl-copy`/`wl-paste` (Wayland) → `xclip -selection clipboard` → `xsel --clipboard`. First successful one wins. |

On Linux, if none of the three backends are installed, `Read`/`Write` return `ErrClipboardUnavailable`. Install one of `wl-clipboard`, `xclip`, or `xsel` to enable.

## Examples

```go
import "github.com/akira-io/onyx/clipboard"

if err := clipboard.Write("hello"); err != nil {
    return err
}

text, err := clipboard.Read()
if err != nil {
    return err
}
fmt.Println(text)
```

## Errors

- `ErrClipboardUnavailable` — no backend reachable (only on Linux when wl-clipboard/xclip/xsel are all missing).
- Wrapped exec errors when a backend is found but the call itself fails (binary missing permission, display server not running, etc.).

## Dependencies

- `osinfo` for platform detection.

## Related modules

- `files` — opening files with the default app.
- `shell` — resolving CLI binaries.
