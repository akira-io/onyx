# shell

Locates command-line executables on the user's machine. Each `Lookup` accepts either a name (resolved via `PATH`) or a path (checked as a file). Try in order; first match wins.

Used by Go desktop applications that wrap a third-party CLI (`claude`, `gh`, `git`, `ffmpeg`) when `PATH` alone may not be enough.

## Public API

| Symbol | Purpose |
| --- | --- |
| `Resolver` | Builder that collects lookup targets (names or paths). |
| `NewResolver()` | Returns an empty `Resolver`. |
| `Resolver.Lookup(target string) Resolver` | Adds one target. If `target` contains a path separator it is checked as a file; otherwise resolved via `PATH`. |
| `Resolver.Lookups(targets []string) Resolver` | Adds many targets. Empty entries are ignored. |
| `Resolver.Resolve() (string, error)` | Returns the absolute path of the first target that exists and is executable. |
| `ListNpmGlobalBinDirs() []string` | Conventional directories where npm global packages install binaries. |
| `ListUserLocalBinDirs() []string` | Conventional per-user bin directories (`~/.local/bin`, `~/bin`). |
| `ListSystemBinDirs() []string` | Conventional system-wide bin directories per platform. |
| `ListWindowsApplicationDirs(applicationName string) []string` | Conventional Windows install directories for a named application. |
| `ErrBinaryNotFound` | Returned when no target resolves. |

## Example

Resolving the `claude` CLI installed via npm, Homebrew, the official installer, or PATH:

```go
import (
    "path/filepath"

    "github.com/akira-io/onyx/osinfo"
    "github.com/akira-io/onyx/shell"
)

binary := "claude" + osinfo.ExecutableExtension()

dirs := append(shell.ListNpmGlobalBinDirs(), shell.ListUserLocalBinDirs()...)
dirs = append(dirs, shell.ListSystemBinDirs()...)
dirs = append(dirs, shell.ListWindowsApplicationDirs("claude")...)

paths := make([]string, 0, len(dirs))
for _, dir := range dirs {
    paths = append(paths, filepath.Join(dir, binary))
}

resolved, err := shell.NewResolver().
    Lookup(binary).
    Lookups(paths).
    Resolve()
if err != nil {
    return err
}

cmd := exec.Command(resolved, "-p", prompt)
```

## Errors

| Sentinel | When |
| --- | --- |
| `ErrBinaryNotFound` | None of the supplied targets resolved. |

## Dependencies

None beyond the standard library. `os/exec.LookPath` is used internally for `PATH`-based lookup.

## Related

- [files](./files.md): launch a file with the default application instead of a specific CLI.
- [osinfo](./osinfo.md): pick the right executable extension to pass to `Lookup`.
