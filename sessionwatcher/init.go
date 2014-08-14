package sessionwatcher

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{"sessionwatcher", Start, Stop, true})
}
