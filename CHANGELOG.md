# Changelog

All notable changes to `onyx` are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-05-18

### Changed

- **Breaking**: project renamed from `desktopkit` to `onyx`. Module path is now `github.com/akira-io/onyx`. Update imports across all consumers.
- Repository renamed to `akira-io/onyx`. GitHub provides a redirect from the old name but new code should reference `onyx` directly.
- All documentation, internal package references, and changelog headers updated to the new identity.

### Migration

Replace `github.com/akira-io/desktopkit` with `github.com/akira-io/onyx` in `go.mod` and Go imports. Run `go mod tidy`. No API surface changed; only the module path.

## [0.2.0] - 2026-05-17

### Added

- `shell.ListNpmGlobalBinDirs`, `shell.ListUserLocalBinDirs`, `shell.ListSystemBinDirs`, and `shell.ListWindowsApplicationDirs` to remove `runtime.GOOS` branching from consumer apps that resolve third-party CLIs.
- `shell.Candidates.WithCandidates(paths []string)` for bulk candidate registration.

## [0.1.0] - 2026-05-17

### Added

- `osinfo` package with `Platform`, `Current`, and `ExecutableExtension`.
- `paths` package with `For(applicationName)` and `Config`, `Data`, `Cache`, `Logs` resolvers for macOS, Linux, and Windows.
- `files` package with `OpenPath`, `OpenURL`, and `RevealInFileManager`.
- `shell` package with the `Candidates` builder and `Resolve` returning a `ResolvedExecutable`.
- Conventions document, architecture document, and per-module documentation under `docs/`.
- MIT license.
