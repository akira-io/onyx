# shell

Locates command-line executables on the user's machine using `PATH` first, then a caller-supplied list of well-known install locations.

This is the package every Go desktop application reaches for when it has to wrap a third-party CLI (`claude`, `gh`, `git`, `ffmpeg`) and PATH alone is not enough.

## Public API

| Symbol | Purpose |
| --- | --- |
| `Resolver` | Builder that collects PATH names and explicit candidate file paths. |
| `Candidates.Lookup(name string) Candidates` | Adds an executable name to look up via PATH. |
| `Candidates.Fallback(path string) Candidates` | Adds an explicit absolute path to try if PATH fails. |
| `Candidates.Fallbacks(paths []string) Candidates` | Adds many explicit absolute paths. Empty entries are ignored. |
| `(Resolver) Resolve() (ResolvedExecutable, error)` | Returns the first candidate that exists and is executable. |
| `ListNpmGlobalBinDirs() []string` | Conventional directories where npm global packages install binaries. |
| `ListUserLocalBinDirs() []string` | Conventional per-user bin directories (`~/.local/bin`, `~/bin`). |
| `ListSystemBinDirs() []string` | Conventional system-wide bin directories per platform. |
| `ListWindowsApplicationDirs(applicationName string) []string` | Conventional Windows install directories for a named application. |
| `ResolvedExecutable` | The result of a successful resolution. |
| `(ResolvedExecutable) AbsolutePath() string` | Absolute path to the binary. |
| `(ResolvedExecutable) Source() ResolutionSource` | Where the binary was found (`SourcePath`, `SourceFallback`). |
| `ErrBinaryNotFound` | Returned when no candidate exists. |

## Example

Resolving the `claude` CLI installed via npm, Homebrew, the official installer, or PATH:

```go
import (
    "github.com/akira-io/onyx/osinfo"
    "github.com/akira-io/onyx/shell"
)

name := "claude" + osinfo.ExecutableExtension()
binary := "claude" + osinfo.ExecutableExtension()

dirs := append(shell.ListNpmGlobalBinDirs(), shell.ListUserLocalBinDirs()...)
dirs = append(dirs, shell.ListSystemBinDirs()...)
dirs = append(dirs, shell.ListWindowsApplicationDirs("claude")...)

candidates := []string{}
for _, dir := range dirs {
    candidates = append(candidates, filepath.Join(dir, binary))
}

resolved, err := shell.NewResolver().
    Lookup(name).
    Fallbacks(candidates).
    Resolve()
if err != nil {
    return err
}

cmd := exec.Command(resolved.AbsolutePath(), "-p", prompt)
```

## Errors

| Sentinel | When |
| --- | --- |
| `ErrBinaryNotFound` | None of the supplied names or candidate paths resolved. |

## Dependencies

None beyond the standard library. (`os/exec.LookPath` is used internally — a future minor release may switch to `github.com/cli/safeexec` for current-directory safety.)

## Related

- [files](./files.md) — when you need to launch a file with the default application instead of a specific CLI.
- [osinfo](./osinfo.md) — callers can use it to choose which executable extensions to pass to `Lookup`.
