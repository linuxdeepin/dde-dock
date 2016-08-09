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
	xid    xproto.Window
	notify bool
	md5    []byte
	damage damage.Damage
}

func NewTrayIcon(xid xproto.Window) *TrayIcon {
	return &TrayIcon{
		xid:    xid,
		notify: true,
	}
}

func (icon *TrayIcon) getName() string {
	name, err := ewmh.WmNameGet(TrayXU, icon.xid)
	if err != nil || name == "" {
		name, err = icccm.WmNameGet(TrayXU, icon.xid)

		if err != nil || name == "" {
			wmclass, _ := icccm.WmClassGet(TrayXU, icon.xid)
			if wmclass != nil {
				name = fmt.Sprintf("[%s|%s]", wmclass.Class, wmclass.Instance)
			}
		}
	}
	return name
}

func (m *TrayManager) addIcon(xid xproto.Window) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.icons[xid]
	if ok {
		logger.Debugf("addIcon failed: %v existed", xid)
		return
	}
	xConn := TrayXU.Conn()
	d, err := damage.NewDamageId(xConn)
	if err != nil {
		logger.Debug("addIcon failed, new damage id failed:", err)
		return
	}
	icon := NewTrayIcon(xid)
	icon.damage = d

	err = damage.CreateChecked(xConn, d, xproto.Drawable(xid), damage.ReportLevelRawRectangles).Check()
	if err != nil {
		logger.Debug("addIcon failed, damage create failed:", err)
		return
	}

	composite.RedirectWindow(xConn, xid, composite.RedirectAutomatic)

	iconWin := xwindow.New(TrayXU, xid)
	iconWin.Listen(xproto.EventMaskVisibilityChange | damage.Notify | xproto.EventMaskStructureNotify)
	iconWin.Change(xproto.CwBackPixel, 0)

	dbus.Emit(m, "Added", uint32(xid))
	logger.Infof("Add tray icon %v name: %q", xid, icon.getName())
	m.icons[xid] = icon
	m.updateTrayIcons()
}

func (m *TrayManager) removeIcon(xid xproto.Window) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.icons[xid]
	if !ok {
		logger.Debugf("removeIcon failed: %v not exist", xid)
		return
	}
	delete(m.icons, xid)
	dbus.Emit(m, "Removed", uint32(xid))
	logger.Debugf("remove tray icon %v", xid)
	m.updateTrayIcons()
}

func (m *TrayManager) updateTrayIcons() {
	var icons []uint32
	for _, icon := range m.icons {
		icons = append(icons, uint32(icon.xid))
	}
	m.TrayIcons = icons
	dbus.NotifyChange(m, "TrayIcons")
}
