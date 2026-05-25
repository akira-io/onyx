# process

Launches and relaunches the host application across platforms, hiding the per-OS command needed to open a fresh instance.

```go
import "github.com/akira-io/onyx/process"
```

## API

| Symbol | Kind | Summary |
|--------|------|---------|
| `Relaunch(applicationPath string) error` | func | Starts a fresh instance of the application at the given path. |
| `ErrEmptyPath` | error | Returned when `applicationPath` is empty. |

## Platform behavior

| OS | Command |
|----|---------|
| macOS | `open -n <app>.app` (resolved via [shell](./26-shell.md)). |
| Windows | `cmd /c start "" <exe>`. |
| Linux | exec the binary directly. |

## Behaviour

`Relaunch` only starts the new instance; the caller is responsible for quitting the current process afterwards (typically right after a successful update is staged). The new process is detached — `Relaunch` returns once it has started, not when it exits.

## Errors

- `ErrEmptyPath` when no path is supplied.
- Wrapped errors when the launcher (`open`) cannot be located or the process fails to start.

## Dependencies

- [osinfo](./24-osinfo.md) to select the platform command.
- [shell](./26-shell.md) to resolve `open` on macOS.

## Related modules

- [shell](./26-shell.md) resolves system binaries.

---

Navigation: [← Appearance](28-appearance.md) · **Process** · [Index →](00-index.md)
