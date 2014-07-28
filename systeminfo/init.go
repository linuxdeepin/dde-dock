package systeminfo

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{"systeminfo", Start, Stop, true})
}
