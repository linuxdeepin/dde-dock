/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package bluetooth

import (
	"errors"
)

var (
	bluezErrorInvalidKey = errors.New("org.bluez.Error.Failed:Resource temporarily unavailable")
	errBluezRejected     = errors.New("org.bluez.Error.Rejected")
	errBluezCanceled     = errors.New("org.bluez.Error.Canceled")
)
