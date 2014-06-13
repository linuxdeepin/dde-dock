package dsc

import "dde-daemon"

func init() {
	loader.Register(&loader.Module{"dsc", Start, Stop, true})
}
