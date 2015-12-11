package appearance

// #cgo pkg-config:  gtk+-3.0
// #include <stdlib.h>
// #include "cursor.h"
import "C"

func (*Manager) listenCursorChanged() {
	C.handle_gtk_cursor_changed()
}
