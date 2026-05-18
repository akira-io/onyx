# onyx — Overview

`onyx` is the Akira Foundation's Go toolkit for building cross-platform desktop applications without writing the same OS-specific glue twice.

It packages thin, opinionated wrappers around the best individual community libraries (and direct OS calls when no library exists), behind a single, consistent, intention-revealing API.

## Why

Every desktop application repeats the same primitives:

- Finding the user's configuration / data / cache directory.
- Opening a file with the user's default application.
- Revealing a file inside the platform's file manager.
- Resolving the absolute path of a CLI binary the user has installed.
- Copying text or images to the system clipboard.
- Sending a system notification.
- Reading and writing secrets in the keychain.

Across macOS, Linux, and Windows each of these has subtle differences. `onyx` makes those differences disappear behind a predictable API.

## What it is

- **Cross-platform** — single import works on `darwin`, `linux`, `windows`.
- **Idiomatic Go** — packages, functions, and errors follow standard Go conventions.
- **Driver-agnostic** — `onyx` defines the contracts; you compose the modules you need.
- **Open source** — MIT-licensed, lives at `github.com/akira-io/onyx`.

## What it is not

- It is not a UI framework. It does not draw windows, menus, or render content.
- It is not a desktop runtime. Use it from inside Wails, Fyne, Gio, or a CLI.
- It is not a kitchen sink. Modules are small, focused, and independently usable.

## Reading guide

- [01-conventions.md](./01-conventions.md) — naming, function design, and documentation rules every module follows.
- [02-architecture.md](./02-architecture.md) — package layout and the SOLID/DRY/KISS principles that drive it.
- [modules/](./modules) — per-module reference.

## Status

Pre-1.0. API stable within a minor version, breaking changes only on `v0.X` bumps.
