package dock

import (
	"bytes"
	"crypto/md5"
	"github.com/BurntSushi/xgb/composite"
	"github.com/BurntSushi/xgb/damage"
	"github.com/BurntSushi/xgb/xfixes"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
	"pkg.linuxdeepin.com/lib/dbus"
)

var (
	TRAYMANAGER *TrayManager
)

const (
	OpCodeSystemTrayRequestDock   uint32 = 0
	OpCodeSystemTrayBeginMessage  uint32 = 1
	OpCodeSystemTrayCancelMessage uint32 = 2
)

type TrayManager struct {
	owner  xproto.Window
	visual xproto.Visualid

	TrayIcons []uint32

	Removed func(id uint32)
	Added   func(id uint32)
	Changed func(id uint32)

	nameInfo   map[xproto.Window]string
	notifyInfo map[xproto.Window]bool
	md5Info    map[xproto.Window][]byte
	dmageInfo  map[xproto.Window]damage.Damage
}

func (m *TrayManager) isValidWindow(xid xproto.Window) bool {
	r, err := xproto.GetWindowAttributes(TrayXU.Conn(), xid).Reply()
	return r != nil && err == nil
}

func (m *TrayManager) checkValid() {
	for _, id := range m.TrayIcons {
		xid := xproto.Window(id)
		if m.isValidWindow(xid) {
			continue
		}

		m.removeTrayIcon(xid)
	}
}

func (m *TrayManager) addTrayIcon(xid xproto.Window) {
	m.checkValid()
	for _, id := range m.TrayIcons {
		if xproto.Window(id) == xid {
			return
		}
	}

	if d, err := damage.NewDamageId(TrayXU.Conn()); err != nil {
		return
	} else {
		m.dmageInfo[xid] = d
		if err := damage.CreateChecked(TrayXU.Conn(), d, xproto.Drawable(xid), damage.ReportLevelRawRectangles).Check(); err != nil {
			logger.Debug("DamageCreate Failed:", err)
			return
		}
	}
	composite.RedirectWindow(TrayXU.Conn(), xid, composite.RedirectAutomatic)

	m.TrayIcons = append(m.TrayIcons, uint32(xid))
	icon := xwindow.New(TrayXU, xid)
	icon.Listen(xproto.EventMaskVisibilityChange | damage.Notify | xproto.EventMaskStructureNotify)
	icon.Change(xproto.CwBackPixel, 0)

	name, err := ewmh.WmNameGet(TrayXU, xid)
	if err != nil {
		logger.Debug("WmNameGet failed:", err, xid)
	}
	m.nameInfo[xid] = name
	m.notifyInfo[xid] = true
	if m.Added != nil {
		m.Added(uint32(xid))
	}
	logger.Infof("Added try icon: \"%s\"(%d)", name, uint32(xid))
}
func (m *TrayManager) removeTrayIcon(xid xproto.Window) {
	delete(m.dmageInfo, xid)
	delete(m.nameInfo, xid)
	delete(m.notifyInfo, xid)
	delete(m.md5Info, xid)
	if m.Removed != nil {
		m.Removed(uint32(xid))
	}
	var newIcons []uint32
	for _, id := range m.TrayIcons {
		if id != uint32(xid) {
			newIcons = append(newIcons, id)
		}
	}
	m.TrayIcons = newIcons
}

func (m *TrayManager) GetName(xid uint32) string {
	return m.nameInfo[xproto.Window(xid)]
}

func (m *TrayManager) EnableNotification(xid uint32, enable bool) {
	m.notifyInfo[xproto.Window(xid)] = enable
}

func (m *TrayManager) handleTrayDamage(xid xproto.Window) {
	if m.notifyInfo[xid] && m.Changed != nil {
		if md5 := icon2md5(xid); !md5Equal(m.md5Info[xid], md5) {
			m.md5Info[xid] = md5
			if m.Changed != nil {
				m.Changed(uint32(xid))
			}
			logger.Infof("handleTrayDamage: %s(%d) changed (%v)", m.nameInfo[xid], xid, md5)
		}
	}
}

func findRGBAVisualID() xproto.Visualid {
	for _, dinfo := range TrayXU.Screen().AllowedDepths {
		for _, vinfo := range dinfo.Visuals {
			if dinfo.Depth == 32 {
				return vinfo.VisualId
			}
		}
	}
	return TrayXU.Screen().RootVisual
}

func (m *TrayManager) destroyOwnerWindow() {
	if m.owner != 0 {
		xproto.DestroyWindow(TrayXU.Conn(), m.owner)
	}
	m.owner = 0
}
func (m *TrayManager) Manage() bool {
	m.destroyOwnerWindow()

	win, _ := xwindow.Generate(TrayXU)
	m.owner = win.Id

	xproto.CreateWindowChecked(TrayXU.Conn(), 0, m.owner, TrayXU.RootWin(), 0, 0, 1, 1, 0, xproto.WindowClassInputOnly, m.visual, 0, nil)
	TrayXU.Sync()
	win.Listen(xproto.EventMaskStructureNotify)
	return m.tryOwner()
}

func (m *TrayManager) RetryManager() {
	m.Unmanage()
	m.Manage()

	if m.Added != nil {
		for _, icon := range m.TrayIcons {
			m.Added(icon)
		}
	}
}

func initTrayManager() {
	composite.Init(TrayXU.Conn())
	composite.QueryVersion(TrayXU.Conn(), 0, 4)
	damage.Init(TrayXU.Conn())
	damage.QueryVersion(TrayXU.Conn(), 1, 1)
	xfixes.Init(TrayXU.Conn())
	xfixes.QueryVersion(TrayXU.Conn(), 5, 0)

	visualId := findRGBAVisualID()

	TRAYMANAGER = &TrayManager{
		owner:      0,
		visual:     visualId,
		nameInfo:   make(map[xproto.Window]string),
		notifyInfo: make(map[xproto.Window]bool),
		md5Info:    make(map[xproto.Window][]byte),
		dmageInfo:  make(map[xproto.Window]damage.Damage),
	}
	TRAYMANAGER.Manage()

	dbus.InstallOnSession(TRAYMANAGER)

	xfixes.SelectSelectionInput(
		TrayXU.Conn(),
		TrayXU.RootWin(),
		_NET_SYSTEM_TRAY_S0,
		xfixes.SelectionEventMaskSelectionClientClose,
	)
	go TRAYMANAGER.startListener()
}

func (m *TrayManager) RequireManageTrayIcons() {
	mstype, err := xprop.Atm(TrayXU, "MANAGER")
	if err != nil {
		logger.Warning("Get MANAGER Failed")
		return
	}

	timeStamp, _ := ewmh.WmUserTimeGet(TrayXU, m.owner)
	cm, err := xevent.NewClientMessage(
		32,
		TrayXU.RootWin(),
		mstype,
		int(timeStamp),
		int(_NET_SYSTEM_TRAY_S0),
		int(m.owner),
	)

	if err != nil {
		logger.Warning("Send MANAGER Request failed:", err)
		return
	}

	// !!! ewmh.ClientEvent not use EventMaskStructureNotify.
	xevent.SendRootEvent(TrayXU, cm,
		uint32(xproto.EventMaskStructureNotify))
}

func (m *TrayManager) getSelectionOwner() (*xproto.GetSelectionOwnerReply, error) {
	_trayInstance := xproto.GetSelectionOwner(TrayXU.Conn(), _NET_SYSTEM_TRAY_S0)
	return _trayInstance.Reply()
}

func (m *TrayManager) tryOwner() bool {
	// Make a check, the tray application MUST be 1.
	reply, err := m.getSelectionOwner()
	if err != nil {
		logger.Error(err)
		return false
	}
	if reply.Owner != 0 {
		logger.Warning("Another System tray application is running")
		return false
	}

	timeStamp, _ := ewmh.WmUserTimeGet(TrayXU, m.owner)
	err = xproto.SetSelectionOwnerChecked(
		TrayXU.Conn(),
		m.owner,
		_NET_SYSTEM_TRAY_S0,
		xproto.Timestamp(timeStamp),
	).Check()
	if err != nil {
		logger.Warning("Set Selection Owner failed: ", err)
		return false
	}

	//owner the _NET_SYSTEM_TRAY_Sn
	logger.Info("Required _NET_SYSTEM_TRAY_S0 successful")

	m.RequireManageTrayIcons()

	xprop.ChangeProp32(
		TrayXU,
		m.owner,
		"_NET_SYSTEM_TRAY_VISUAL",
		"VISUALID",
		uint(TRAYMANAGER.visual),
	)
	xprop.ChangeProp32(
		TrayXU,
		m.owner,
		"_NET_SYSTEM_TRAY_ORIENTAION",
		"CARDINAL",
		0,
	)
	reply, err = m.getSelectionOwner()
	if err != nil {
		logger.Warning(err)
		return false
	}
	return reply.Owner != 0
}

func (m *TrayManager) Unmanage() bool {
	reply, err := m.getSelectionOwner()
	if err != nil {
		logger.Info("get selection owner failed:", err)
		return false
	}
	if reply.Owner != m.owner {
		logger.Info("not selection owner")
		return false
	}

	m.destroyOwnerWindow()
	timeStamp, _ := ewmh.WmUserTimeGet(TrayXU, m.owner)
	return xproto.SetSelectionOwnerChecked(
		TrayXU.Conn(),
		0,
		_NET_SYSTEM_TRAY_S0,
		xproto.Timestamp(timeStamp),
	).Check() == nil
}

var isListened bool = false

func (m *TrayManager) startListener() {
	// to avoid creating too much listener when SelectionNotifyEvent occurs.
	if isListened {
		return
	}
	isListened = true

	for {
		if e, err := TrayXU.Conn().WaitForEvent(); err == nil {
			switch ev := e.(type) {
			case xproto.ClientMessageEvent:
				// logger.Info("ClientMessageEvent")
				if ev.Type == _NET_SYSTEM_TRAY_OPCODE {
					// timeStamp = ev.Data.Data32[0]
					opCode := ev.Data.Data32[1]
					// logger.Info("TRAY_OPCODE")

					switch opCode {
					case OpCodeSystemTrayRequestDock:
						xid := xproto.Window(ev.Data.Data32[2])
						m.addTrayIcon(xid)
					case OpCodeSystemTrayBeginMessage:
					case OpCodeSystemTrayCancelMessage:
					}
				}
			case damage.NotifyEvent:
				m.handleTrayDamage(xproto.Window(ev.Drawable))
			case xproto.DestroyNotifyEvent:
				m.removeTrayIcon(ev.Window)
			case xproto.SelectionClearEvent:
				m.Unmanage()
			case xfixes.SelectionNotifyEvent:
				m.Manage()
			}
		}
	}
}

func icon2md5(xid xproto.Window) []byte {
	pixmap, _ := xproto.NewPixmapId(TrayXU.Conn())
	defer xproto.FreePixmap(TrayXU.Conn(), pixmap)
	if err := composite.NameWindowPixmapChecked(TrayXU.Conn(), xid, pixmap).Check(); err != nil {
		logger.Warning("NameWindowPixmap failed:", err, xid)
		return nil
	}
	im, err := xgraphics.NewDrawable(TrayXU, xproto.Drawable(pixmap))
	if err != nil {
		logger.Warning("Create xgraphics.Image failed:", err, pixmap)
		return nil
	}
	buf := bytes.NewBuffer(nil)
	im.WritePng(buf)
	hasher := md5.New()
	hasher.Write(buf.Bytes())
	return hasher.Sum(nil)
}
func md5Equal(a []byte, b []byte) bool {
	if len(a) != 16 || len(b) != 16 {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
func (*TrayManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.dde.TrayManager",
		"/com/deepin/dde/TrayManager",
		"com.deepin.dde.TrayManager",
	}
}
