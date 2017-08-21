package trayicon

import (
	"fmt"
	"pkg.deepin.io/lib/dbus"

	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/damage"
)

type TrayIcon struct {
	win    x.Window
	notify bool
	md5    []byte
	damage damage.Damage
}

func NewTrayIcon(win x.Window) *TrayIcon {
	return &TrayIcon{
		win:    win,
		notify: true,
	}
}

func (icon *TrayIcon) getName() string {
	wmName, _ := ewmhConn.GetWMName(icon.win).Reply(ewmhConn)
	if wmName != "" {
		return wmName
	}

	wmNameTextProp, err := icccmConn.GetWMName(icon.win).Reply(icccmConn)
	if err == nil {
		wmName, _ := wmNameTextProp.GetStr()
		if wmName != "" {
			return wmName
		}
	}

	wmClass, err := icccmConn.GetWMClass(icon.win).Reply(icccmConn)
	if err == nil {
		return fmt.Sprintf("[%s|%s]", wmClass.Class, wmClass.Instance)
	}

	return ""
}

func (m *TrayManager) addIcon(win x.Window) {
	m.checkValid()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.icons[win]
	if ok {
		logger.Debugf("addIcon failed: %v existed", win)
		return
	}
	damageId, err := XConn.GenerateID()
	if err != nil {
		logger.Debug("addIcon failed, new damage id failed:", err)
		return
	}
	d := damage.Damage(damageId)

	icon := NewTrayIcon(win)
	icon.damage = d

	err = damage.CreateChecked(XConn, d, x.Drawable(win), damage.ReportLevelRawRectangles).Check(XConn)
	if err != nil {
		logger.Debug("addIcon failed, damage create failed:", err)
		return
	}

	const valueMask = x.CWBackPixel | x.CWEventMask
	valueList := &x.ChangeWindowAttributesValueList{
		BackgroundPixel: 0,
		EventMask:       x.EventMaskVisibilityChange | x.EventMaskStructureNotify,
	}

	x.ChangeWindowAttributes(XConn, win, valueMask, valueList)

	dbus.Emit(m, "Added", uint32(win))
	logger.Infof("Add tray icon %v name: %q", win, icon.getName())
	m.icons[win] = icon
	m.updateTrayIcons()
}

func (m *TrayManager) removeIcon(win x.Window) {
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
