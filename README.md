# onyx

Cross-platform Go toolkit for building desktop applications without rewriting OS-specific glue every time. Thin, opinionated wrappers around the best community libraries (and direct OS calls when no library exists), behind a single, intention-revealing API.

> Full reference: [`docs/00-index.md`](docs/00-index.md) - one file per package, mirrored with the Rust crate's documentation tree.

## Why

```go
// Without onyx
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
```

```go
// With onyx
config, err := paths.For("Hyperion").Config()
```

Same behaviour on macOS, Linux, and Windows. No `runtime.GOOS` switches, no hand-rolled XDG logic, no per-OS shell invocations leaked into application code.

## Install

```sh
go get github.com/akira-io/onyx
```

Requires Go 1.23 or later.

## Quickstart

```go
package main

import (
    "log"

    "github.com/akira-io/onyx/files"
    "github.com/akira-io/onyx/paths"
    "github.com/akira-io/onyx/shell"
)

func main() {
    app := paths.For("Hyperion")

    config, err := app.Config()
    if err != nil {
        log.Fatal(err)
    }
    if err := files.RevealInFileManager(config); err != nil {
        log.Fatal(err)
    }

    claude, err := shell.NewResolver().
        Lookup("claude").
        Lookup("/opt/homebrew/bin/claude").
        Resolve()
    if err != nil {
        log.Fatal(err)
    }
    log.Println("claude at:", claude)
}
```

## Modules

| Module | Purpose |
| --- | --- |
| [paths](./docs/modules/paths.md) | Configuration, data, cache, and log directories per platform. |
| [files](./docs/modules/files.md) | Open paths and URLs, reveal in file manager. |
| [shell](./docs/modules/shell.md) | Resolve CLI binaries via PATH lookup with explicit fallback paths. |
| [clipboard](./docs/modules/clipboard.md) | Read and write the system clipboard as plain text. |
| [notify](./docs/modules/notify.md) | Show desktop notifications. |
| [keyring](./docs/modules/keyring.md) | Store, retrieve, and delete secrets in the system credential store. |
| [osinfo](./docs/modules/osinfo.md) | Typed runtime platform detection. |

## Documentation

The full reference lives in [`docs/`](./docs/):

- [docs/00-overview.md](./docs/00-overview.md) — what onyx is and is not.
- [docs/01-conventions.md](./docs/01-conventions.md) — naming, function design, documentation rules every module follows.
- [docs/02-architecture.md](./docs/02-architecture.md) — package layout and the SOLID/DRY/KISS principles that drive it.
- [docs/modules/](./docs/modules/) — per-module reference.

## Platforms

| OS | Status |
| --- | --- |
| macOS | Fully tested by the maintainer. Primary development target. |
| Linux | Compiled and unit-tested in CI. Not exercised against real desktop environments by the maintainer. |
| Windows | Compiled and unit-tested in CI. Not exercised against real desktop environments by the maintainer. |

CI runs on all three platforms, but a green pipeline only proves the code compiles and the platform-agnostic tests pass. The Linux and Windows backends (notification daemons, clipboard helpers, credential managers, file managers) are not verified end-to-end against live systems by the maintainer. Report issues with reproduction steps and we will work through them.

## Testing

```sh
go test ./...
```

The full suite runs on macOS, Linux, and Windows in CI on every push to `main` and on every pull request. Locally, run the suite before opening a PR. Tests that exercise OS facilities (notifications, clipboard, keychain) skip cleanly when no backend is reachable, so the suite stays green even on minimal CI images.

## Contributing

Pull requests welcome. Before opening one:

1. Read [docs/01-conventions.md](./docs/01-conventions.md). PRs that break the conventions get rejected without further review.
2. Read [docs/02-architecture.md](./docs/02-architecture.md) for the rules that govern where new code goes.
3. Add tests for every public change. Touch [CHANGELOG.md](./CHANGELOG.md) under `## [Unreleased]`.
4. Use conventional commits (`feat:`, `fix:`, `refactor:`, `docs:`, `chore:`). The changelog workflow groups bullets by prefix.

For Rust consumers, see the sister crate [`akira-io/onyx-rs`](https://github.com/akira-io/onyx-rs).

## Prior art

`onyx` is a study project. These libraries solve overlapping problems and have more battle-test. If you want a dependency you can lean on, reach for them first:

- Paths: [`github.com/adrg/xdg`](https://github.com/adrg/xdg), stdlib `os.UserConfigDir` / `os.UserCacheDir`.
- Clipboard: [`github.com/atotto/clipboard`](https://github.com/atotto/clipboard).
- Notifications: [`github.com/gen2brain/beeep`](https://github.com/gen2brain/beeep).
- Keyring: [`github.com/zalando/go-keyring`](https://github.com/zalando/go-keyring).
- Binary resolution: stdlib `os/exec.LookPath`.

## License

MIT. See [LICENSE](./LICENSE).
