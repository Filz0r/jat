package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

var (
	Version   = "undefined"
	Commit    = "undefined"
	BuildDate = "undefined"
)

func GoInstallVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}
	return ""
}

func EffectiveVersion() string {
	if v := GoInstallVersion(); v != "" {
		return v
	}
	return Version
}

func Info() string {
	goVersion := runtime.Version()
	os := runtime.GOOS
	arch := runtime.GOARCH
	return fmt.Sprintf("%sVERSION LOG%s\nBuild Version: %s\nEffective Version: %s\nCommit: %s\nBuild Date: %s\nGo: %s OS: %s Arch: %s",
		"\x1b[35m", "\x1b[0m", Version,
		EffectiveVersion(), Commit, BuildDate, goVersion, os, arch)
}
