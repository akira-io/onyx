# Installation

## Add the module

```bash
go get github.com/akira-io/onyx@latest
```

Or pin to a specific tag:

```bash
go get github.com/akira-io/onyx@v1.0.2
```

Go 1.23+. `go.mod` declares `go 1.23`.

## Imports

```go
import (
    "github.com/akira-io/onyx/clipboard"
    "github.com/akira-io/onyx/files"
    "github.com/akira-io/onyx/keyring"
    "github.com/akira-io/onyx/notify"
    "github.com/akira-io/onyx/osinfo"
    "github.com/akira-io/onyx/paths"
    "github.com/akira-io/onyx/shell"
)
```

Each package is independently importable — pull in only what you use.

## Runtime dependencies

Zero non-stdlib runtime deps. Every backend call shells out via `os/exec`. No CGo, no FFI bindings, no transitive C libraries.

## OS support

| Package | macOS | Linux | Windows |
|---------|-------|-------|---------|
| `clipboard` | pbcopy / pbpaste | wl-clipboard / xclip / xsel | PowerShell + WinForms |
| `files` | open | xdg-open | cmd `start` |
| `keyring` | `security` | `secret-tool` | `cmdkey` + PowerShell + `CredentialManager` |
| `notify` | osascript | notify-send | BurntToast / `msg` |
| `osinfo` | always | always | always |
| `paths` | `~/Library/...` | XDG Base Directory | Known Folders |
| `shell` | `PATH` lookup | `PATH` lookup | `PATH` lookup + Program Files heuristics |

Packages that shell out report a sentinel error when no backend is reachable (`clipboard.ErrClipboardUnavailable`, `keyring.ErrKeyringUnavailable`, `notify.ErrNotifyUnavailable`) — callers branch via `errors.Is` and surface install instructions.

## Verify

```bash
go build ./...
go test ./...
go vet ./...
```

The unit suite skips cleanly when a backend is missing (returns early on `Unavailable`), so the suite stays green on minimal CI runners. `gofmt -l` should report no diff.

---

Navigation: [← Index](00-index.md) · **Installation** · [Quickstart →](02-quickstart.md)
