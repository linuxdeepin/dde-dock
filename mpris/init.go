package mpris

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{"mpris", Start, Stop, true})
}
