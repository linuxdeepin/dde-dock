package mounts

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{"mounts", Start, Stop, true})
}
