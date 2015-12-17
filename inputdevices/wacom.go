package inputdevices

import (
	"pkg.deepin.io/dde/api/dxinput"
	dxutils "pkg.deepin.io/dde/api/dxinput/utils"
	"pkg.deepin.io/lib/dbus/property"
	"gir/gio-2.0"
)

const (
	wacomSchema = "com.deepin.dde.wacom"

	wacomKeyLeftHanded        = "left-handed"
	wacomKeyCursorMode        = "cursor-mode"
	wacomKeyUpAction          = "keyup-action"
	wacomKeyDownAction        = "keydown-action"
	wacomKeyDoubleDelta       = "double-delta"
	wacomKeyPressureSensitive = "pressure-sensitive"
)

const (
	btnNumUpKey   int32 = 3
	btnNumDownKey       = 2
)

var actionMap = map[string]string{
	"LeftClick":   "button 1",
	"MiddleClick": "button 2",
	"RightClick":  "button 3",
	"PageUp":      "key KP_Page_Up",
	"PageDown":    "key KP_Page_Down",
}

// Soften(x1<y1 x2<y2) --> Firmer(x1>y1 x2>y2)
var pressureLevel = map[uint32][]int {
	1: []int{0, 100, 0, 100},
	2: []int{20, 80, 20, 80},
	3: []int{30, 70, 30, 70},
	4: []int{0, 0, 100, 100},
	5: []int{60, 40, 60, 40},
	6: []int{70, 30, 70, 30}, // default
	7: []int{75, 25, 75, 25},
	8: []int{80, 20, 80, 20},
	9: []int{90, 10, 90, 10},
	10: []int{100, 0, 100, 0},
}

type ActionInfo struct {
	Action string
	Desc   string
}
type ActionInfos []*ActionInfo

type Wacom struct {
	LeftHanded *property.GSettingsBoolProperty `access:"readwrite"`
	CursorMode *property.GSettingsBoolProperty `access:"readwrite"`

	KeyUpAction   *property.GSettingsStringProperty `access:"readwrite"`
	KeyDownAction *property.GSettingsStringProperty `access:"readwrite"`

	DoubleDelta       *property.GSettingsUintProperty `access:"readwrite"`
	PressureSensitive *property.GSettingsUintProperty `access:"readwrite"`

	DeviceList  dxutils.DeviceInfos
	ActionInfos ActionInfos
	Exist       bool

	dxWacoms map[int32]*dxinput.Wacom
	setting  *gio.Settings
}

var _wacom *Wacom

func getWacom() *Wacom {
	if _wacom == nil {
		_wacom = NewWacom()

		_wacom.init()
		_wacom.handleGSettings()
	}

	return _wacom
}

func NewWacom() *Wacom {
	var w = new(Wacom)

	w.setting = gio.NewSettings(wacomSchema)
	w.LeftHanded = property.NewGSettingsBoolProperty(
		w, "LeftHanded",
		w.setting, wacomKeyLeftHanded)
	w.CursorMode = property.NewGSettingsBoolProperty(
		w, "CursorMode",
		w.setting, wacomKeyCursorMode)

	w.KeyUpAction = property.NewGSettingsStringProperty(
		w, "KeyUpAction",
		w.setting, wacomKeyUpAction)
	w.KeyDownAction = property.NewGSettingsStringProperty(
		w, "KeyDownAction",
		w.setting, wacomKeyDownAction)

	w.DoubleDelta = property.NewGSettingsUintProperty(
		w, "DoubleDelta",
		w.setting, wacomKeyDoubleDelta)
	w.PressureSensitive = property.NewGSettingsUintProperty(
		w, "PressureSensitive",
		w.setting, wacomKeyPressureSensitive)

	w.updateDeviceList()
	w.dxWacoms = make(map[int32]*dxinput.Wacom)
	w.updateDXWacoms()

	return w
}

func (w *Wacom) init() {
	if !w.Exist {
		return
	}

	w.enableCursorMode()
	w.enableLeftHanded()
	w.setKeyAction(btnNumUpKey, w.KeyUpAction.Get())
	w.setKeyAction(btnNumDownKey, w.KeyDownAction.Get())
	w.setPressure()
	w.setClickDelta()
}

func (w *Wacom) handleDeviceChanged() {
	w.updateDeviceList()
	w.updateDXWacoms()
	w.init()
}

func (w *Wacom) updateDeviceList() {
	w.DeviceList = getWacomInfos(false)
	if len(w.DeviceList) == 0 {
		w.setPropExist(false)
	} else {
		w.setPropExist(true)
	}
}

func (w *Wacom) updateDXWacoms() {
	for _, info := range w.DeviceList {
		_, ok := w.dxWacoms[info.Id]
		if ok {
			continue
		}

		dxw, err := dxinput.NewWacom(info.Id)
		if err != nil {
			logger.Warning(err)
			continue
		}
		w.dxWacoms[info.Id] = dxw
	}
}

func (w *Wacom) setKeyAction(btnNum int32, action string) {
	value, ok := actionMap[action]
	if !ok {
		return
	}

	for _, v := range w.dxWacoms {
		err := v.SetButton(int(btnNum), value)
		if err != nil {
			logger.Debugf("Set btn mapping for '%v - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (w *Wacom) enableLeftHanded() {
	var rotate string = "none"
	if w.LeftHanded.Get() {
		rotate = "half"
	}

	for _, v := range w.dxWacoms {
		err := v.SetRotate(rotate)
		if err != nil {
			logger.Debugf("Enable left handed for '%v - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (w *Wacom) enableCursorMode() {
	var mode string = "Absolute"
	if w.CursorMode.Get() {
		mode = "Relative"
	}

	for _, v := range w.dxWacoms {
		err := v.SetMode(mode)
		if err != nil {
			logger.Debugf("Enable cursor mode for '%v - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (w *Wacom) setPressure() {
	array, ok := pressureLevel[w.PressureSensitive.Get()]
	if !ok {
		return
	}

	for _, v := range w.dxWacoms {
		err := v.SetPressureCurve(array[0], array[1], array[2], array[3])
		if err != nil {
			logger.Debugf("Set pressure curve for '%v - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (w *Wacom) setClickDelta() {
	for _, v := range w.dxWacoms {
		err := v.SetSuppress(int(w.DoubleDelta.Get()))
		if err != nil {
			logger.Debugf("Set double click delta for '%v - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}
