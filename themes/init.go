package themes

import "dde-daemon"

func init() {
	loader.Register(&loader.Module{"themes", Start, nil, true})
}
