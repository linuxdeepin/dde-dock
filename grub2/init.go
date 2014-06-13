package grub2

import "dde-daemon"

func init() {
	loader.Register(&loader.Module{"grub2", Start, Stop, true})
}
