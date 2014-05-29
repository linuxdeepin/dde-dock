package dock

import (
	"bytes"
	"crypto/md5"
	"dlib/dbus"
	"github.com/BurntSushi/xgb/composite"
	"github.com/BurntSushi/xgb/damage"
	"github.com/BurntSushi/xgb/xfixes"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
)

var (
	TrayXU, _ = xgbutil.NewConn()
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

func (m *TrayManager) addTrayIcon(xid xproto.Window) {
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
			LOGGER.Debug("DamageCreate Failed:", err)
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
		LOGGER.Debug("WmNameGet failed:", err, xid)
	}
	m.nameInfo[xid] = name
	m.notifyInfo[xid] = true
	if m.Added != nil {
		m.Added(uint32(xid))
	}
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
			LOGGER.Infof("handleTrayDamage: %s(%d) changed (%v)", m.nameInfo[xid], xid, md5)
		}
	}
}

var TRAYMANAGER *TrayManager

var _NET_SYSTEM_TRAY_S0, _ = xprop.Atm(TrayXU, "_NET_SYSTEM_TRAY_S0")
var _NET_SYSTEM_TRAY_OPCODE, _ = xprop.Atm(TrayXU, "_NET_SYSTEM_TRAY_OPCODE")

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

func initTrayManager() {
	composite.Init(TrayXU.Conn())
	composite.QueryVersion(TrayXU.Conn(), 0, 4)
	damage.Init(TrayXU.Conn())
	damage.QueryVersion(TrayXU.Conn(), 1, 1)
	xfixes.Init(TrayXU.Conn())
	xfixes.QueryVersion(TrayXU.Conn(), 5, 0)

	visualId := findRGBAVisualID()
	win, _ := xwindow.Generate(TrayXU)
	xproto.CreateWindowChecked(TrayXU.Conn(), 0, win.Id, TrayXU.RootWin(), 0, 0, 1, 1, 0, xproto.WindowClassInputOnly, visualId, 0, nil)
	TrayXU.Sync()
	win.Listen(xproto.EventMaskStructureNotify)

	TRAYMANAGER = &TrayManager{
		owner:      win.Id,
		visual:     visualId,
		nameInfo:   make(map[xproto.Window]string),
		notifyInfo: make(map[xproto.Window]bool),
		md5Info:    make(map[xproto.Window][]byte),
		dmageInfo:  make(map[xproto.Window]damage.Damage),
	}

	dbus.InstallOnSession(TRAYMANAGER)
	TRAYMANAGER.tryOwner()
}

func (m *TrayManager) RequireManageTrayIcons() {
	mstype, err := xprop.Atm(TrayXU, "MANAGER")
	if err != nil {
		LOGGER.Error("Get MANAGER Failed")
		return
	}

	timeStamp, _ := ewmh.WmUserTimeGet(TrayXU, TRAYMANAGER.owner)
	cm, err := xevent.NewClientMessage(
		32,
		TrayXU.RootWin(),
		mstype,
		int(timeStamp),
		int(_NET_SYSTEM_TRAY_S0),
		int(TRAYMANAGER.owner),
	)

	if err != nil {
		LOGGER.Error("Send MANAGER Request failed:", err)
		return
	}

	// !!! ewmh.ClientEvent not use EventMaskStructureNotify.
	xevent.SendRootEvent(TrayXU, cm,
		uint32(xproto.EventMaskStructureNotify))
}

func (m *TrayManager) tryOwner() bool {
	// Make a check, the tray application MUST be 1.
	_trayInstance := xproto.GetSelectionOwner(TrayXU.Conn(), _NET_SYSTEM_TRAY_S0)
	reply, err := _trayInstance.Reply()
	if err != nil {
		LOGGER.Fatal(err)
	}
	if reply.Owner == 0 {
		timeStamp, _ := ewmh.WmUserTimeGet(TrayXU, TRAYMANAGER.owner)
		err := xproto.SetSelectionOwnerChecked(
			TrayXU.Conn(),
			m.owner,
			_NET_SYSTEM_TRAY_S0,
			xproto.Timestamp(timeStamp),
		).Check()
		if err != nil {
			LOGGER.Info("Set Selection Owner failed: ", err)
			return false
		}
		go TRAYMANAGER.startListener()

		//owner the _NET_SYSTEM_TRAY_Sn
		LOGGER.Info("Required _NET_SYSTEM_TRAY_S0")

		m.RequireManageTrayIcons()

		xprop.ChangeProp32(
			TrayXU,
			TRAYMANAGER.owner,
			"_NET_SYSTEM_TRAY_VISUAL",
			"VISUALID",
			uint(TRAYMANAGER.visual),
		)
		xprop.ChangeProp32(
			TrayXU,
			TRAYMANAGER.owner,
			"_NET_SYSTEM_TRAY_ORIENTAION",
			"CARDINAL",
			0,
		)
		xfixes.SelectSelectionInput(TrayXU.Conn(), TrayXU.RootWin(), _NET_SYSTEM_TRAY_S0, xfixes.SelectionEventMaskSelectionClientClose)
		return true
	} else {
		LOGGER.Info("Another System tray application is running")
		return false
	}
}

func (m *TrayManager) unmanage() {
	timeStamp, _ := ewmh.WmUserTimeGet(TrayXU, TRAYMANAGER.owner)
	// FIXME:
	// is 0 right?
	xproto.SetSelectionOwner(TrayXU.Conn(), 0,
		_NET_SYSTEM_TRAY_S0, xproto.Timestamp(timeStamp)).Check()
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
				// LOGGER.Info("ClientMessageEvent")
				if ev.Type == _NET_SYSTEM_TRAY_OPCODE {
					// timeStamp = ev.Data.Data32[0]
					opCode := ev.Data.Data32[1]
					// LOGGER.Info("TRAY_OPCODE")

					switch opCode {
					case OpCodeSystemTrayRequestDock:
						xid := xproto.Window(ev.Data.Data32[2])
						LOGGER.Info("Tray Get Request Dock: ", xid)
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
				//clean up
				m.unmanage()
			case xfixes.SelectionNotifyEvent:
				m.tryOwner()
			}
		}
	}
}

func icon2md5(xid xproto.Window) []byte {
	pixmap, _ := xproto.NewPixmapId(TrayXU.Conn())
	defer xproto.FreePixmap(TrayXU.Conn(), pixmap)
	if err := composite.NameWindowPixmapChecked(TrayXU.Conn(), xid, pixmap).Check(); err != nil {
		LOGGER.Warning("NameWindowPixmap failed:", err, xid)
		return nil
	}
	im, err := xgraphics.NewDrawable(TrayXU, xproto.Drawable(pixmap))
	if err != nil {
		LOGGER.Warning("Create xgraphics.Image failed:", err, pixmap)
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
