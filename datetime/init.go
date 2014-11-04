package datetime

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "datetime",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
