package launcher

import (
	. "pkg.linuxdeepin.com/lib/gettext"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/log"
)

var logger = log.NewLogger("dde-daemon/launcher-daemon")

func Stop() {
	logger.EndTracing()
}

func Start() {
	logger.BeginTracing()

	InitI18n()
	// DesktopAppInfo.ShouldShow does not know deepin.
	gio.DesktopAppInfoSetDesktopEnv("Deepin")

	initCategory()
	logger.Info("init category done")

	initItems()
	logger.Info("init items done")

	initDBus()
	logger.Info("init dbus done")
}
