package grub2

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{"grub2", Start, Stop, true})
}
