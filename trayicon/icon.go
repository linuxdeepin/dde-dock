package trayicon

import (
	"fmt"
	"github.com/BurntSushi/xgb/composite"
	"github.com/BurntSushi/xgb/damage"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xwindow"
	"pkg.deepin.io/lib/dbus"
)

type TrayIcon struct {
	win    xproto.Window
	notify bool
	md5    []byte
	damage damage.Damage
}

func NewTrayIcon(win xproto.Window) *TrayIcon {
	return &TrayIcon{
		win:    win,
		notify: true,
	}
}

func (icon *TrayIcon) getName() string {
	name, err := ewmh.WmNameGet(XU, icon.win)
	if err != nil || name == "" {
		name, err = icccm.WmNameGet(XU, icon.win)

		if err != nil || name == "" {
			wmclass, _ := icccm.WmClassGet(XU, icon.win)
			if wmclass != nil {
				name = fmt.Sprintf("[%s|%s]", wmclass.Class, wmclass.Instance)
			}
		}
	}
	return name
}

func (m *TrayManager) addIcon(win xproto.Window) {
	m.checkValid()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.icons[win]
	if ok {
		logger.Debugf("addIcon failed: %v existed", win)
		return
	}
	xConn := XU.Conn()
	d, err := damage.NewDamageId(xConn)
	if err != nil {
		logger.Debug("addIcon failed, new damage id failed:", err)
		return
	}
	icon := NewTrayIcon(win)
	icon.damage = d

	err = damage.CreateChecked(xConn, d, xproto.Drawable(win), damage.ReportLevelRawRectangles).Check()
	if err != nil {
		logger.Debug("addIcon failed, damage create failed:", err)
		return
	}

	composite.RedirectWindow(xConn, win, composite.RedirectAutomatic)

	iconWin := xwindow.New(XU, win)
	iconWin.Listen(xproto.EventMaskVisibilityChange | damage.Notify | xproto.EventMaskStructureNotify)
	iconWin.Change(xproto.CwBackPixel, 0)

	dbus.Emit(m, "Added", uint32(win))
	logger.Infof("Add tray icon %v name: %q", win, icon.getName())
	m.icons[win] = icon
	m.updateTrayIcons()
}

func (m *TrayManager) removeIcon(win xproto.Window) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.icons[win]
	if !ok {
		logger.Debugf("removeIcon failed: %v not exist", win)
		return
	}
	delete(m.icons, win)
	dbus.Emit(m, "Removed", uint32(win))
	logger.Debugf("remove tray icon %v", win)
	m.updateTrayIcons()
}

func (m *TrayManager) updateTrayIcons() {
	var icons []uint32
	for _, icon := range m.icons {
		icons = append(icons, uint32(icon.win))
	}
	m.TrayIcons = icons
	dbus.NotifyChange(m, "TrayIcons")
}
