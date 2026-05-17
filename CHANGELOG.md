# Changelog

All notable changes to `desktopkit` are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- Swap EOL git-cliff-action for taiki-e/install-action.

## [Unreleased]

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
