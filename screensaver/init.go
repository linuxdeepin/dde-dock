package screensaver

import "dde-daemon"

func init() {
	loader.Register(&loader.Module{"screensaver", Start, nil, true})
}
