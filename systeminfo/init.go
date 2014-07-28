package systeminfo

import "dde-daemon"

func init() {
	loader.Register(&loader.Module{"systeminfo", Start, Stop, true})
}
