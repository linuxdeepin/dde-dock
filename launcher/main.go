package launcher

import (
	. "dlib/gettext"
	"dlib/gio-2.0"
	l "dlib/logger"
)

var logger *l.Logger = l.NewLogger("dde-daemon/launcher-daemon")

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
