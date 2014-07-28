package dsc

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{"dsc", Start, Stop, true})
}
