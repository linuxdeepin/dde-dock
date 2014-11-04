package network

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "network",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
