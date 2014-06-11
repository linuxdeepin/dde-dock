package mime

import "dde-daemon"

func init() {
	loader.Register(&loader.Module{"mime", Start, nil, true})
}
