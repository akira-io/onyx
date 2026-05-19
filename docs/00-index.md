# onyx — Reference

Cross-platform desktop toolkit for Go. Sister module to [`akira-io/onyx-rs`](https://github.com/akira-io/onyx-rs) (Rust) — same module surface, same conventions, idiomatic in each language.

The module packages thin, opinionated wrappers around the standard library (and direct OS calls when no library exists), behind a single, consistent, intention-revealing API.

## Meta

| File | Topic |
|------|-------|
| [01-installation](01-installation.md) | Add the module, Go version requirements |
| [02-quickstart](02-quickstart.md) | 60-second snippet for each package |
| [03-architecture](03-architecture.md) | Package layout, SOLID/DRY/KISS principles |
| [04-conventions](04-conventions.md) | Naming, function design, comments, errors |
| [05-release-flow](05-release-flow.md) | Versioning, branching, tagging, publishing |

## Packages

| File | Topic |
|------|-------|
| [20-clipboard](20-clipboard.md) | Read/write the system clipboard as plain text |
| [21-files](21-files.md) | Open files / URLs, reveal in file manager |
| [22-keyring](22-keyring.md) | Store, retrieve, delete secrets in the OS credential store |
| [23-notify](23-notify.md) | Desktop notifications |
| [24-osinfo](24-osinfo.md) | Platform identity — single source of truth for OS branching |
| [25-paths](25-paths.md) | Per-application config / data / cache / log directories |
| [26-shell](26-shell.md) | Locate command-line executables |

## What this module is

- **Cross-platform** — single import works on `darwin`, `linux`, `windows`.
- **Idiomatic Go** — packages, functions, and errors follow standard Go conventions.
- **Low-dependency** — only standard library imports at runtime.
- **Mirror of the Rust crate** — same module names, same surface, idiomatic in each language.
- **Open source** — MIT-licensed, lives at `github.com/akira-io/onyx`.

## What this module is not

- It is not a UI framework. No window, menu, or rendering surface.
- It is not a desktop runtime. Use it from Wails, Fyne, Gio, or a CLI.
- It is not a kitchen sink. Packages are small, focused, and independently usable.

## Status

Pre-1.0. API stable within a minor version; breaking changes only on `v0.X` bumps.

---

Navigation: [README](../README.md) · **Index** · [Installation →](01-installation.md)
