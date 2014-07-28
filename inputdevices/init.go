package inputdevices

import "dde-daemon"

func init() {
	loader.Register(&loader.Module{"inputdevices", Start, Stop, true})
}
