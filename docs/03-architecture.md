# Architecture

## Package layout

```
onyx/
├── clipboard/clipboard.go    pbcopy/pbpaste, wl-copy/xclip/xsel, PowerShell
├── files/files.go            open, xdg-open, cmd start
├── keyring/keyring.go        security, secret-tool, cmdkey + CredentialManager
├── notify/notify.go          osascript, notify-send, BurntToast / msg
├── osinfo/osinfo.go          Platform + ExecutableExtension
├── paths/paths.go            AppPaths { Config, Data, Cache, Logs }
└── shell/shell.go            Resolver + bin-dir lookup helpers
```

One file per package keeps the surface narrow. Tests live next to source as `<package>_test.go`. No `util.go`, `helpers.go`, or `common.go` — every helper belongs to a domain package.

## Dependency direction

```
   osinfo      ◄────── single source of truth for OS branching
    ▲
    │
   ┌┴──────────────────────────────────────┐
   │                                       │
   clipboard  files  keyring  notify       paths   shell
```

`osinfo.Current()` is the only path that touches `runtime.GOOS`. Every other package asks `osinfo` instead of switching directly. When you add a new platform target, the change set converges on `osinfo` first.

`paths` and `shell` also lean on `os.Getenv` / `os.UserHomeDir` for environment lookups (`HOME`, `XDG_*`, `APPDATA`, `LOCALAPPDATA`, `PATH`). These are scoped per-package — no helper layer in between.

## Trust model

`onyx` runs in the same process as the application. It does not sandbox itself, does not validate inputs beyond shape (empty path, empty service, etc.), and does not redact secrets. Wrap secrets with your own zeroing logic at the call site if you need stronger guarantees — `keyring.Get` returns a plain `string`.

The module is library code: no global mutex, no spawned background goroutines, no environment mutation. It is safe to call from any goroutine. The shell-out calls block; wrap in your own concurrency control if needed.

## Backend selection

Packages that shell out follow the same pattern:

1. Ask `osinfo.Current()`.
2. Branch into a per-OS implementation.
3. Linux paths fall through a priority-ordered backend list (`wl-clipboard` → `xclip` → `xsel`, BurntToast → `msg`, etc.). First successful one wins; the `Unavailable` sentinel is returned only when every option fails.

The backend list lives inline in the package that owns it — no global registry, no plugin mechanism. Adding a new backend means editing one function and one constant.

## Error model

One sentinel error per failure mode per package, declared at the package top:

```go
var (
    ErrClipboardUnavailable = errors.New("clipboard: no supported backend available")
    ErrEmptyService         = errors.New("keyring: service must not be empty")
    ErrNotFound             = errors.New("keyring: secret not found")
)
```

Backend failures are wrapped via `fmt.Errorf("%s: %w", action, err)` so callers can unwrap to the underlying `*exec.ExitError` when needed.

Callers branch via `errors.Is` (not string comparison):

```go
if errors.Is(err, keyring.ErrNotFound) {
    promptLogin()
}
```

## Concurrency

All public functions are synchronous. The shell-out calls are blocking by design — desktop primitives are not hot-path code. Wrap with a `go` keyword + `chan error` for fire-and-forget, or use your runtime's worker pool.

Packages are stateless. `osinfo.Platform`, `paths.AppPaths`, and `shell.Resolver` are immutable value types; the builders (`Lookup`, `Lookups`) return new values rather than mutating the receiver.

## SOLID / DRY / KISS

- **Single responsibility** — one package, one concern.
- **Open/closed** — extend by adding a new package; do not edit unrelated packages.
- **Liskov** — `Platform` is one shape; `Resolver` is one shape. No interface hierarchies.
- **Interface segregation** — the module exports free functions where the contract is small (`clipboard.Read`), and typed values where state is required (`Resolver`, `AppPaths`).
- **Dependency inversion** — packages depend on the `osinfo` contract, not on `runtime.GOOS`.
- **DRY** — every OS switch lives in `osinfo` or the package that owns the backend list.
- **KISS** — the simplest correct API wins. Reflection, init-time side effects, and channels are last resorts.

## Cross-module parity

The Rust crate at [`akira-io/onyx-rs`](https://github.com/akira-io/onyx-rs) ships the same package names with idiomatic Rust shapes (typed `Result` enums, `&Path` for paths). Behavioural parity is the goal — same backends, same priority order, same error categories. The two crates diverge only where the language idiom requires.

---

Navigation: [← Quickstart](02-quickstart.md) · **Architecture** · [Conventions →](04-conventions.md)
