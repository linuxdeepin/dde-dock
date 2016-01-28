package bluetooth

import (
	"errors"
)

var (
	bluezErrorInvalidKey = errors.New("org.bluez.Error.Failed:Resource temporarily unavailable")
	errBluezRejected     = errors.New("org.bluez.Error.Rejected")
	errBluezCanceled     = errors.New("org.bluez.Error.Canceled")
)
