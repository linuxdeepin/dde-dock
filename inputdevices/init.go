package inputdevices

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{"inputdevices", Start, Stop, true})
}
