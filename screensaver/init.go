package screensaver

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{"screensaver", Start, Stop, true})
}
