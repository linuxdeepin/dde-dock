package dock

import (
	"dbus/com/deepin/daemon/display"
	"github.com/BurntSushi/xgb/xproto"
)

var displayRect xproto.Rectangle = xproto.Rectangle{0, 0, 0, 0}

func init() {
	display, err := display.NewDisplay("com.deepin.daemon.Display",
		"/com/deepin/daemon/Display")
	if err != nil {
		logger.Error("connect failed:", err)
		return
	}
	setDisplayRect(display.PrimaryRect.Get())
	display.ConnectPrimaryChanged(setDisplayRect)
}

func setDisplayRect(rect []interface{}) {
	displayRect.X = rect[0].(int16)
	displayRect.Y = rect[1].(int16)
	displayRect.Width = rect[2].(uint16)
	displayRect.Height = rect[3].(uint16)
}
