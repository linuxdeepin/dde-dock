/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package appearance

// #cgo pkg-config:  gtk+-3.0
// #include <stdlib.h>
// #include "cursor.h"
import "C"

func (*Manager) listenCursorChanged() {
	C.handle_gtk_cursor_changed()
}

func (*Manager) endCursorChangedHandler() {
	C.end_cursor_changed_handler()
}
