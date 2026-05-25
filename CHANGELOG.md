# Changelog

All notable changes to `onyx` are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.0] - 2026-05-25

### Added

- `osinfo.Hostname()` returns the operating system host name, or an empty string when it cannot be determined. Best-effort: callers supply their own fallback.
- `machineid` package. `GetOrCreate(application)` returns a stable per-application identifier for the current machine, persisted in the system keyring. Sentinel `ErrEmptyApplication`.
- `appearance` package. `IsDark()` reports whether the OS uses a dark color scheme via `defaults` (macOS), the registry (Windows), and `gsettings` (Linux). Best-effort: returns `false` when undetermined.
- `shell.LoginPath()` and `shell.EnrichedEnviron()`. Recover the PATH from the user's login shell so GUI applications can locate user-installed tools (npm, php, git) that the launcher PATH omits.
- `process` package. `Relaunch(applicationPath)` starts a fresh instance of the application (`open -n` on macOS, `start` on Windows, exec on Linux). Sentinel `ErrEmptyPath`.
- `clipboard` package. `Read` and `Write` operate on the system clipboard as plain text via per-platform backends (`pbcopy`/`pbpaste` on macOS, PowerShell on Windows, Wayland/X11 tools on Linux). Sentinel `ErrClipboardUnavailable` is returned only on Linux when none of `wl-clipboard`, `xclip`, or `xsel` is installed.
- `notify` package. `Show(title, body)` displays a desktop notification via `osascript` on macOS, `notify-send` on Linux, and PowerShell `BurntToast` (with `msg.exe` fallback) on Windows. Sentinels `ErrEmptyTitle` and `ErrNotifyUnavailable` cover the error cases.
- `keyring` package. `Set`, `Get`, and `Delete` manage secrets in the system credential store via `security` (Keychain) on macOS, `secret-tool` (Secret Service) on Linux, and `cmdkey` + PowerShell `CredentialManager` on Windows. Sentinels `ErrNotFound`, `ErrKeyringUnavailable`, `ErrEmptyService`, and `ErrEmptyAccount` cover the error cases.

## [1.0.2] - 2026-05-18

### Changed

- `shell.Resolver` simplified. `Fallback` and `Fallbacks` removed; `Lookup` and `Lookups` now accept either a name (resolved via `PATH`) or a path (checked as a file). Auto-detected via path separators or Windows drive prefix.
- `Resolve` now returns `(string, error)` directly. `ResolvedExecutable` and `ResolutionSource` removed.
- Reason: one verb, one mental model. Callers no longer pick between "name" and "fallback".
- README and `docs/modules/shell.md` updated.

### Removed

- `Fallback`, `Fallbacks` methods. Replaced by `Lookup` accepting either a name (PATH lookup) or a path (file check); auto-detected via path separators or Windows drive prefix.
- `ResolvedExecutable` type and its `AbsolutePath`/`Source` getters. `Resolve` now returns `(string, error)` directly. Callers that previously called `.AbsolutePath()` on the result just use the returned string.
- `ResolutionSource` enum:
  - `SourcePath` reported the binary came from a `PATH` search via `exec.LookPath`.
  - `SourceFallback` (previously `SourceCandidate`) reported the binary came from an explicit caller-supplied path.
  - `SourceUnknown` covered uninitialized values.
  - Removed because callers almost never branched on the source. When the source matters, callers can inspect the returned path themselves (compare against `PATH`, against a known install dir, etc.). Keeping the enum in the public API forced every consumer to import and pattern-match on a value most ignored.

## [1.0.1] - 2026-05-18

### Changed

- shell API renamed for clarity. No production consumers yet, version held within v1.
  - `NewCandidates` is now `NewResolver`.
  - `Candidates` struct is now `Resolver`.
  - `WithName(s)` is now `Lookup(s)`. Semantics: PATH lookup.
  - `WithCandidate(p)` is now `Fallback(p)`. Semantics: explicit fallback path tried after PATH misses.
  - `WithCandidates(ps)` is now `Fallbacks(ps)`.
  - `SourceCandidate` is now `SourceFallback`.
- Reason: "candidate" was ambiguous. New names reveal intent.
- README before/after example uses the new API.

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
