package grub2

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "grub2",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
