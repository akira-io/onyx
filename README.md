# onyx

Cross-platform Go toolkit for building desktop applications without rewriting OS-specific glue every time.

`onyx` packages thin, opinionated wrappers around the best individual community libraries (and direct OS calls when no library exists), behind a single, consistent, intention-revealing API.

```go
import (
    "github.com/akira-io/onyx/paths"
    "github.com/akira-io/onyx/files"
    "github.com/akira-io/onyx/shell"
)

app := paths.For("Hyperion")
config, _ := app.Config()
_ = files.RevealInFileManager(config)

claude, err := shell.NewCandidates().
    WithName("claude").
    WithCandidate("/opt/homebrew/bin/claude").
    Resolve()
```

## Status

Pre-1.0. Public API stable within a minor version.

## Modules

| Module | Purpose |
| --- | --- |
| [paths](./docs/modules/paths.md) | Configuration, data, cache, and log directories per platform. |
| [files](./docs/modules/files.md) | Open paths and URLs, reveal in file manager. |
| [shell](./docs/modules/shell.md) | Resolve CLI binaries via PATH and explicit candidates. |
| [osinfo](./docs/modules/osinfo.md) | Typed runtime detection helpers shared by every module. |

Planned: `clipboard`, `notify`, `keyring`.

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
