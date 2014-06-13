package mpris

import "dde-daemon"

func init() {
	loader.Register(&loader.Module{"mpris", Start, Stop, true})
}
