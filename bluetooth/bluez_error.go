package bluetooth

import (
	"errors"
)

var (
	bluezErrorInvalidKey = errors.New("org.bluez.Error.Failed:Resource temporarily unavailable")
)
