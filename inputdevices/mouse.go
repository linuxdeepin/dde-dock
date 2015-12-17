package inputdevices

import (
	"pkg.deepin.io/dde/api/dxinput"
	dxutils "pkg.deepin.io/dde/api/dxinput/utils"
	"pkg.deepin.io/lib/dbus/property"
	"gir/gio-2.0"
)

const (
	mouseSchema = "com.deepin.dde.mouse"

	mouseKeyLeftHanded      = "left-handed"
	mouseKeyDisableTouchpad = "disable-touchpad"
	mouseKeyMiddleButton    = "middle-button-enabled"
	mouseKeyNaturalScroll   = "natural-scroll"
	mouseKeyAcceleration    = "motion-acceleration"
	mouseKeyThreshold       = "motion-threshold"
	mouseKeyDoubleClick     = "double-click"
	mouseKeyDragThreshold   = "drag-threshold"
)

type Mouse struct {
	LeftHanded            *property.GSettingsBoolProperty `access:"readwrite"`
	DisableTpad           *property.GSettingsBoolProperty `access:"readwrite"`
	NaturalScroll         *property.GSettingsBoolProperty `access:"readwrite"`
	MiddleButtonEmulation *property.GSettingsBoolProperty `access:"readwrite"`

	MotionAcceleration *property.GSettingsFloatProperty `access:"readwrite"`
	MotionThreshold    *property.GSettingsFloatProperty `access:"readwrite"`

	DoubleClick   *property.GSettingsIntProperty `access:"readwrite"`
	DragThreshold *property.GSettingsIntProperty `access:"readwrite"`

	DeviceList dxutils.DeviceInfos
	Exist      bool

	dxMouses map[int32]*dxinput.Mouse
	setting  *gio.Settings
}

var _mouse *Mouse

func getMouse() *Mouse {
	if _mouse == nil {
		_mouse = NewMouse()

		_mouse.init()
		_mouse.handleGSettings()
	}

	return _mouse
}

func NewMouse() *Mouse {
	var m = new(Mouse)

	m.setting = gio.NewSettings(mouseSchema)
	m.LeftHanded = property.NewGSettingsBoolProperty(
		m, "LeftHanded",
		m.setting, mouseKeyLeftHanded)
	m.DisableTpad = property.NewGSettingsBoolProperty(
		m, "DisableTpad",
		m.setting, mouseKeyDisableTouchpad)
	m.NaturalScroll = property.NewGSettingsBoolProperty(
		m, "NaturalScroll",
		m.setting, mouseKeyNaturalScroll)
	m.MiddleButtonEmulation = property.NewGSettingsBoolProperty(
		m, "MiddleButtonEmulation",
		m.setting, mouseKeyMiddleButton)

	m.MotionAcceleration = property.NewGSettingsFloatProperty(
		m, "MotionAcceleration",
		m.setting, mouseKeyAcceleration)
	m.MotionThreshold = property.NewGSettingsFloatProperty(
		m, "MotionThreshold",
		m.setting, mouseKeyThreshold)

	m.DoubleClick = property.NewGSettingsIntProperty(
		m, "DoubleClick",
		m.setting, mouseKeyDoubleClick)
	m.DragThreshold = property.NewGSettingsIntProperty(
		m, "DragThreshold",
		m.setting, mouseKeyDragThreshold)

	m.updateDeviceList()
	m.dxMouses = make(map[int32]*dxinput.Mouse)
	m.updateDXMouses()

	return m
}

func (m *Mouse) init() {
	if !m.Exist {
		if getTouchpad().Exist {
			getTouchpad().enable(true)
		}
		return
	}

	m.enableLeftHanded()
	m.enableMidBtnEmu()
	m.enableNaturalScroll()
	m.motionAcceleration()
	m.motionThreshold()
	if m.DisableTpad.Get() {
		m.disableTouchpad()
	}
}

func (m *Mouse) handleDeviceChanged() {
	m.updateDeviceList()
	m.updateDXMouses()
	m.init()
}

func (m *Mouse) updateDeviceList() {
	m.DeviceList = getMouseInfos(false)
	if len(m.DeviceList) == 0 {
		m.setPropExist(false)
	} else {
		m.setPropExist(true)
	}
}

func (m *Mouse) updateDXMouses() {
	for _, info := range m.DeviceList {
		_, ok := m.dxMouses[info.Id]
		if ok {
			continue
		}

		dxm, err := dxinput.NewMouse(info.Id)
		if err != nil {
			logger.Warning(err)
			continue
		}
		m.dxMouses[info.Id] = dxm
	}
}

func (m *Mouse) disableTouchpad() {
	if !m.Exist {
		return
	}

	tpad := getTouchpad()
	if !tpad.Exist {
		return
	}

	if !m.DisableTpad.Get() && tpad.TPadEnable.Get() {
		tpad.enable(true)
		return
	}

	tpad.enable(false)
}

func (m *Mouse) enableLeftHanded() {
	for _, v := range m.dxMouses {
		err := v.EnableLeftHanded(m.LeftHanded.Get())
		if err != nil {
			logger.Debugf("Enable left handed for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (m *Mouse) enableNaturalScroll() {
	for _, v := range m.dxMouses {
		err := v.EnableNaturalScroll(m.NaturalScroll.Get())
		if err != nil {
			logger.Debugf("Enable natural scroll for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (m *Mouse) enableMidBtnEmu() {
	for _, v := range m.dxMouses {
		err := v.EnableMiddleButtonEmulation(
			m.MiddleButtonEmulation.Get())
		if err != nil {
			logger.Debugf("Enable mid btn emulation for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (m *Mouse) motionAcceleration() {
	for _, v := range m.dxMouses {
		err := v.SetMotionAcceleration(
			float32(m.MotionAcceleration.Get()))
		if err != nil {
			logger.Debugf("Set acceleration for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (m *Mouse) motionThreshold() {
	for _, v := range m.dxMouses {
		err := v.SetMotionThreshold(float32(m.MotionThreshold.Get()))
		if err != nil {
			logger.Debugf("Set threshold for '%d - %v' failed: %v",
				v.Id, v.Name, err)
		}
	}
}

func (m *Mouse) doubleClick() {
	xsSetInt32(xsPropDoubleClick, m.DoubleClick.Get())
}

func (m *Mouse) dragThreshold() {
	xsSetInt32(xsPropDragThres, m.DragThreshold.Get())
}
