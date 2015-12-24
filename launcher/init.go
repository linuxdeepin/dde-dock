package launcher

import (
	"pkg.deepin.io/dde/daemon/launcher/log"
	"pkg.deepin.io/dde/daemon/loader"
)

func init() {
	loader.Register(NewLauncherDaemon(log.Log))
}
