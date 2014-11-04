package inputdevices

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "inputdevices",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
