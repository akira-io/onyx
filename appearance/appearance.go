package appearance

import (
	"os/exec"
	"strings"

	"github.com/akira-io/onyx/osinfo"
)

// IsDark reports whether the operating system is currently using a dark color
// scheme. It is best-effort: when the preference cannot be determined it
// returns false (light).
func IsDark() bool {
	platform := osinfo.Current()
	switch {
	case platform.IsDarwin():
		return darwinIsDark()
	case platform.IsWindows():
		return windowsIsDark()
	case platform.IsLinux():
		return linuxIsDark()
	}
	return false
}

func darwinIsDark() bool {
	out, err := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(out)), "dark")
}

func windowsIsDark() bool {
	out, err := exec.Command(
		"reg", "query",
		`HKCU\Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`,
		"/v", "AppsUseLightTheme",
	).Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "0x0")
}

func linuxIsDark() bool {
	if isDarkSetting("color-scheme") {
		return true
	}
	return isDarkSetting("gtk-theme")
}

func isDarkSetting(key string) bool {
	out, err := exec.Command("gsettings", "get", "org.gnome.desktop.interface", key).Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(out)), "dark")
}
