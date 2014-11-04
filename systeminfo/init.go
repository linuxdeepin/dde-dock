package systeminfo

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "systeminfo",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
