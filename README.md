# onyx

Cross-platform Go toolkit for building desktop applications without rewriting OS-specific glue every time.

`onyx` packages thin, opinionated wrappers around the best individual community libraries (and direct OS calls when no library exists), behind a single, consistent, intention-revealing API.

### Without onyx

```go
import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "runtime"
)

func appConfigDir(app string) (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    switch runtime.GOOS {
    case "darwin":
        return filepath.Join(home, "Library", "Application Support", app), nil
    case "linux":
        if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
            return filepath.Join(xdg, app), nil
        }
        return filepath.Join(home, ".config", app), nil
    case "windows":
        if v := os.Getenv("APPDATA"); v != "" {
            return filepath.Join(v, app), nil
        }
        return filepath.Join(home, "AppData", "Roaming", app), nil
    }
    return "", fmt.Errorf("unsupported os: %s", runtime.GOOS)
}

func revealInFileManager(path string) error {
    switch runtime.GOOS {
    case "darwin":
        return exec.Command("open", "-R", path).Run()
    case "linux":
        return exec.Command("xdg-open", filepath.Dir(path)).Run()
    case "windows":
        return exec.Command("cmd", "/c", "explorer", "/select,", path).Run()
    }
    return fmt.Errorf("unsupported os: %s", runtime.GOOS)
}

func resolveClaude() (string, error) {
    if p, err := exec.LookPath("claude"); err == nil {
        return p, nil
    }
    candidates := []string{}
    switch runtime.GOOS {
    case "darwin":
        candidates = append(candidates, "/opt/homebrew/bin/claude", "/usr/local/bin/claude")
    case "linux":
        candidates = append(candidates, "/usr/local/bin/claude", "/usr/bin/claude")
    case "windows":
        if v := os.Getenv("LOCALAPPDATA"); v != "" {
            candidates = append(candidates, filepath.Join(v, "Programs", "claude", "claude.exe"))
        }
    }
    for _, c := range candidates {
        if info, err := os.Stat(c); err == nil && !info.IsDir() {
            return c, nil
        }
    }
    return "", fmt.Errorf("claude not found")
}

config, _ := appConfigDir("Hyperion")
_ = revealInFileManager(config)
claudeBin, err := resolveClaude()
```

### With onyx

```go
import (
    "github.com/akira-io/onyx/paths"
    "github.com/akira-io/onyx/files"
    "github.com/akira-io/onyx/shell"
)

app := paths.For("Hyperion")
config, _ := app.Config()
_ = files.RevealInFileManager(config)

claude, err := shell.NewResolver().
    Lookup("claude").
    Lookup("/opt/homebrew/bin/claude").
    Resolve()
```

Same behavior on macOS, Linux, and Windows. No `runtime.GOOS` switches, no hand-rolled XDG logic, no per-OS shell invocations leaked into application code.

## Status

Stable at v1.0.2. Public API stable within a major version (SemVer).

## Design notes

### `shell.Resolver`: one verb, two cases

Earlier versions of `shell` exposed two separate concepts: `WithName(s)` for `PATH` lookups and `WithCandidate(p)` (later `Fallback(p)`) for explicit file paths to try when `PATH` missed. Resolution attached a source tag (`SourcePath` vs `SourceFallback` vs `SourceUnknown`) so callers could see how the binary was found.

That split asked callers to classify each input upfront. In practice the classification is mechanical: if the string has a path separator (`/`, `\`) or a Windows drive prefix (`C:`), it is a path; otherwise it is a name. The source tag was rarely inspected.

`Resolver` collapses everything to a single ordered list of targets. `Lookup` accepts both. `Resolve` returns the absolute path of the first target that resolves, or `ErrBinaryNotFound`. The result is a plain `string`. Callers that genuinely need to know how a binary was located inspect the returned path themselves.

## Modules

| Module | Purpose |
| --- | --- |
| [paths](./docs/modules/paths.md) | Configuration, data, cache, and log directories per platform. |
| [files](./docs/modules/files.md) | Open paths and URLs, reveal in file manager. |
| [shell](./docs/modules/shell.md) | Resolve CLI binaries via PATH lookup with explicit fallback paths. |
| [clipboard](./docs/modules/clipboard.md) | Read and write the system clipboard as plain text. |
| [osinfo](./docs/modules/osinfo.md) | Typed runtime detection helpers shared by every module. |

Planned: `notify`, `keyring`.

## Reading guide

- [docs/00-overview.md](./docs/00-overview.md) — what onyx is and is not.
- [docs/01-conventions.md](./docs/01-conventions.md) — naming, function design, documentation rules every module follows.
- [docs/02-architecture.md](./docs/02-architecture.md) — package layout and the SOLID/DRY/KISS principles that drive it.
- [docs/modules/](./docs/modules/) — per-module reference.

## Platforms

- macOS — primary target, fully supported.
- Linux — supported with the XDG specification.
- Windows — supported, idiomatic `%AppData%` / `%LocalAppData%` paths.

## Installation

```sh
go get github.com/akira-io/onyx
```

Go 1.23 or later.

## Contributing

Read [docs/01-conventions.md](./docs/01-conventions.md) first. Pull requests that break the conventions are rejected.

## License

MIT. See [LICENSE](./LICENSE).
