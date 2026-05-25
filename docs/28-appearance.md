# appearance

Reports the operating system's current color scheme so applications can match the native light or dark theme.

```go
import "github.com/akira-io/onyx/appearance"
```

## API

| Symbol | Kind | Summary |
|--------|------|---------|
| `IsDark() bool` | func | True when the OS is using a dark color scheme. |

## Platform behavior

| OS | Source |
|----|--------|
| macOS | `defaults read -g AppleInterfaceStyle` contains `Dark`. |
| Windows | Registry `AppsUseLightTheme` under `Themes\Personalize` is `0x0`. |
| Linux | `gsettings` `color-scheme` (then `gtk-theme`) contains `dark`. |

## Behaviour

Best-effort. When the preference cannot be read — the key is absent, the tool is missing, or the OS is unsupported — `IsDark()` returns `false` (light). Read it on demand; the value can change at runtime when the user switches themes, so cache it only for the duration of a single render.

## Errors

`appearance` has no error type. `IsDark()` returns a plain `bool`.

## Dependencies

- [osinfo](./24-osinfo.md) to select the platform backend.

## Related modules

- [osinfo](./24-osinfo.md) for static platform identity.

---

Navigation: [← Machine ID](27-machineid.md) · **Appearance** · [Process →](29-process.md)
