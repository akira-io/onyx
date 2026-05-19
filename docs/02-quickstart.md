# Quickstart

One short snippet per package. Every call returns `error` — handle or wrap.

## Clipboard

```go
import "github.com/akira-io/onyx/clipboard"

if err := clipboard.Write("hello onyx"); err != nil {
    return err
}
text, err := clipboard.Read()
if err != nil {
    return err
}
fmt.Println(text)
```

## Files

```go
import "github.com/akira-io/onyx/files"

files.OpenPath("/Users/me/Downloads/report.pdf")
files.OpenURL("https://akira.foundation")
files.RevealInFileManager("/Users/me/Downloads/report.pdf")
```

## Keyring

```go
import "github.com/akira-io/onyx/keyring"

keyring.Set("io.akira.app", "kid@example.com", "hunter2")

secret, err := keyring.Get("io.akira.app", "kid@example.com")
if errors.Is(err, keyring.ErrNotFound) {
    // first run, no secret yet
}

keyring.Delete("io.akira.app", "kid@example.com")
```

## Notify

```go
import "github.com/akira-io/onyx/notify"

notify.Show("Build complete", "Spectra v0.9.0 is ready to install.")
```

## OS info

```go
import "github.com/akira-io/onyx/osinfo"

p := osinfo.Current()
fmt.Printf("%s (exe ext: %q)\n", p.String(), osinfo.ExecutableExtension())

if p.IsDarwin() { /* ... */ }
```

## Paths

```go
import "github.com/akira-io/onyx/paths"

app := paths.For("hyperion")
config, _ := app.Config() // ~/Library/Application Support/hyperion on macOS
data, _   := app.Data()
cache, _  := app.Cache()
logs, _   := app.Logs()
```

## Shell

```go
import "github.com/akira-io/onyx/shell"

claude, err := shell.NewResolver().
    Lookup("claude").
    Lookup("/opt/homebrew/bin/claude").
    Resolve()
if err != nil {
    return err
}
fmt.Println("found at", claude)
```

## Combining packages

```go
import (
    "os"
    "path/filepath"

    "github.com/akira-io/onyx/files"
    "github.com/akira-io/onyx/paths"
)

logs, _ := paths.For("hyperion").Logs()
os.MkdirAll(logs, 0o755)
log := filepath.Join(logs, "today.log")
os.WriteFile(log, []byte("boot ok\n"), 0o644)
files.RevealInFileManager(log)
```

`osinfo.Current()` is the only place the module inspects `runtime.GOOS` — every other package asks `osinfo` instead of switching directly. Reuse the same primitive in your app.

## What to read next

- [03-architecture](03-architecture.md) — package layout and design principles
- [04-conventions](04-conventions.md) — naming, error, doc rules every package follows
- Package reference under [20-clipboard](20-clipboard.md) → [26-shell](26-shell.md)

---

Navigation: [← Installation](01-installation.md) · **Quickstart** · [Architecture →](03-architecture.md)
