# shell

Locate command-line executables. Each `Lookup` target is treated as a `PATH` name when it has no separators, or as a file path otherwise.

```go
import "github.com/akira-io/onyx/shell"
```

## API

| Symbol | Kind | Summary |
|--------|------|---------|
| `Resolver` | struct | Immutable builder collecting lookup targets. |
| `NewResolver() Resolver` | func | Empty resolver. |
| `(Resolver).Lookup(target string) Resolver` | method | Append one target. |
| `(Resolver).Lookups(targets []string) Resolver` | method | Append many. |
| `(Resolver).Resolve() (string, error)` | method | Try each target in order; first matching executable wins. |
| `ErrBinaryNotFound` | sentinel | Resolver could not find any of the targets. |
| `ListNpmGlobalBinDirs() []string` | func | Candidate npm global bin directories. |
| `ListUserLocalBinDirs() []string` | func | `~/.local/bin`, `~/bin`. |
| `ListSystemBinDirs() []string` | func | Platform-specific system bin directories. |
| `ListWindowsApplicationDirs(applicationName string) []string` | func | `LOCALAPPDATA\Programs\<app>`, `ProgramFiles\<app>`, `ProgramFiles(x86)\<app>`. |
| `LoginPath() string` | func | PATH as seen by the user's login shell; recovers tool dirs GUI launchers omit. |
| `EnrichedEnviron() []string` | func | `os.Environ()` with PATH merged from the process and login-shell PATH. |

`Resolver` is a value type. `Lookup` / `Lookups` return a **new** `Resolver` rather than mutating the receiver — chain calls freely.

## Lookup semantics

`Resolver` accepts two kinds of targets:

- **Bare name** (`"claude"`, `"node"`) — searched on `PATH` via `exec.LookPath`. Same algorithm as `which`.
- **Path-like** (`"./bin/foo"`, `"/opt/homebrew/bin/claude"`, `"C:\Program Files\app.exe"`) — checked directly. The substring rules:
  - Contains `/` or `\` → path-like.
  - Has a Windows drive letter (`X:`) → path-like.

Targets are tried in the order they were appended. First one that resolves to an existing file wins. Empty inputs are skipped silently.

`isExecutableFile` checks `os.Stat` then `!info.IsDir()`. The Unix executable bit is not consulted — most desktop apps embed a sidecar with a known name; the file existence check is sufficient.

## Examples

```go
import "github.com/akira-io/onyx/shell"

claude, err := shell.NewResolver().
    Lookup("claude").                                       // PATH
    Lookup("/opt/homebrew/bin/claude").                     // explicit
    Lookup("/Applications/Claude.app/Contents/MacOS/claude").
    Resolve()
if err != nil {
    return err
}
```

Cross-platform binary lookup using `ExecutableExtension`:

```go
import (
    "github.com/akira-io/onyx/osinfo"
    "github.com/akira-io/onyx/shell"
)

name := "hyperion" + osinfo.ExecutableExtension()
bin, err := shell.NewResolver().
    Lookup(name).
    Lookups([]string{
        "/usr/local/bin/hyperion",
        "/opt/homebrew/bin/hyperion",
    }).
    Resolve()
```

Building a candidate list dynamically:

```go
candidates := append(shell.ListUserLocalBinDirs(), shell.ListSystemBinDirs()...)
targets := make([]string, 0, len(candidates))
for _, d := range candidates {
    targets = append(targets, filepath.Join(d, "hyperion"))
}
bin, err := shell.NewResolver().Lookups(targets).Resolve()
```

## Bin-dir helpers

| Helper | Returns |
|--------|---------|
| `ListNpmGlobalBinDirs()` | `~/.npm-global/bin`, `~/.local/share/npm/bin` on Unix; `%APPDATA%\npm` on Windows. |
| `ListUserLocalBinDirs()` | `~/.local/bin`, `~/bin` on Unix. Empty on Windows (no convention). |
| `ListSystemBinDirs()` | `/usr/local/bin`, `/opt/homebrew/bin`, `/usr/bin` on macOS; `/usr/local/bin`, `/usr/bin` on Linux. Empty on Windows. |
| `ListWindowsApplicationDirs(app)` | `LOCALAPPDATA\Programs\<app>`, `ProgramFiles\<app>`, `ProgramFiles(x86)\<app>`. Empty on non-Windows or when `app` is empty. |

The helpers consult environment variables (`HOME`, `APPDATA`, `LOCALAPPDATA`, `ProgramFiles`, `ProgramFiles(x86)`). Missing vars yield an empty slice — the helper never errors.

## Behaviour

- `Resolver{}.Lookup("")` is a no-op (empty strings are filtered out before storage).
- `Resolver.Resolve()` returns `ErrBinaryNotFound` when no target matches. Always check the error.
- The package does not spawn the resolved binary — that is the caller's job. Use `exec.Command(resolved)`.
- Path-like targets that exist but are not regular files (e.g. directories) are skipped.

## Errors

- `ErrBinaryNotFound` — no target resolved. The sentinel carries no payload — log the targets you tried at the call site if you need diagnostics.

## Dependencies

- `osinfo` — `Current()` drives the OS-specific bin-dir lists.
- `os.Getenv`, `os.UserHomeDir`, `os/exec.LookPath` from stdlib.

## Related packages

- [`files`](21-files.md) — open the binary's installer or repository page if missing.
- [`osinfo`](24-osinfo.md) — `ExecutableExtension()` for `<name>.exe` on Windows.

## Cross-module parity

Mirrors the Rust crate's `shell` module one-to-one: same `Resolver` builder shape, same path-like heuristic, same priority order for bin-dir helpers.

---

Navigation: [← Paths](25-paths.md) · **Shell** · [Machine ID →](27-machineid.md)
