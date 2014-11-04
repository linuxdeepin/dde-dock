package screenedge

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "screenedge",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
