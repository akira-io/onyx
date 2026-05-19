# osinfo

Single source of truth for platform identity. Every other package asks `osinfo` instead of branching on `runtime.GOOS` directly.

```go
import "github.com/akira-io/onyx/osinfo"
```

## API

| Symbol | Kind | Summary |
|--------|------|---------|
| `Platform` | struct | Value wrapping the OS identifier. |
| `Current() Platform` | func | Returns the host platform. |
| `(Platform).IsDarwin() bool` | method | macOS predicate. |
| `(Platform).IsLinux() bool` | method | Linux predicate. |
| `(Platform).IsWindows() bool` | method | Windows predicate. |
| `(Platform).String() string` | method | `"darwin" \| "linux" \| "windows" \| <other runtime.GOOS>`. |
| `ExecutableExtension() string` | func | `".exe"` on Windows, `""` otherwise. |

`Platform` is a struct with one unexported field. Compare with `==`. Pass by value — it is cheap to copy.

## Why

Spreading `if runtime.GOOS == "darwin"` throughout the codebase makes adding a new target a search-and-replace job. Concentrating that logic here means a new platform target adds one constant in `osinfo` and propagates everywhere.

Predicates over string comparisons also kill the typo-class of bugs: `osinfo.Current().IsDawin()` would not compile, whereas `runtime.GOOS == "dawin"` silently compiles and never matches.

## Examples

```go
import "github.com/akira-io/onyx/osinfo"

p := osinfo.Current()
fmt.Printf("%s (exe ext: %q)\n", p, osinfo.ExecutableExtension())

switch {
case p.IsDarwin():
    useMacKeychain()
case p.IsLinux():
    useSecretService()
case p.IsWindows():
    useCredentialManager()
}
```

Cache the value once when you need many predicates in a tight loop — `Current()` is cheap (one field read) but explicit threading is clearer:

```go
func pickBackend(p osinfo.Platform) string {
    switch {
    case p.IsDarwin():  return "security"
    case p.IsLinux():   return "secret-tool"
    case p.IsWindows(): return "cmdkey"
    }
    return "unsupported"
}
```

## ExecutableExtension

Helper for building binary file names that work across platforms:

```go
target := "hyperion" + osinfo.ExecutableExtension()
// "hyperion" on macOS/Linux, "hyperion.exe" on Windows
```

Use this when constructing paths to bundled binaries — e.g. embedding a sidecar via `embed` and writing it to disk before launch.

## Behaviour

- `Current()` always succeeds — `runtime.GOOS` is a compile-time constant; there is no runtime call.
- Predicates return `false` on any unsupported OS rather than panicking. `String()` exposes whatever `runtime.GOOS` reports (e.g. `"freebsd"`, `"haiku"`), so consumers can branch:
  ```go
  if !p.IsDarwin() && !p.IsLinux() && !p.IsWindows() {
      return fmt.Errorf("unsupported os: %s", p)
  }
  ```

## Errors

`osinfo` has no error type. Every function returns a plain value.

## Dependencies

None (uses `runtime` only).

## Related packages

Every other package depends on `osinfo`. When adding a new package, ask `Current()` instead of touching `runtime.GOOS` outside the test gates.

## Cross-module parity

Mirrors the Rust crate's `osinfo` module: same predicate names (`IsDarwin` ↔ `is_darwin`), same `ExecutableExtension` helper, same one-source-of-truth principle.

---

Navigation: [← Notify](23-notify.md) · **OS info** · [Paths →](25-paths.md)
