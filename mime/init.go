package mime

import "pkg.linuxdeepin.com/dde-daemon/loader"

func init() {
	loader.Register(NewDaemon())
}
