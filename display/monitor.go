package main

import "github.com/BurntSushi/xgb/randr"
import "dlib/dbus"
import "fmt"
import "strings"
import "math"
import "runtime"

const joinSeparator = "="

type Monitor struct {
	cfg     *ConfigMonitor
	Outputs []string

	BestMode Mode

	IsComposited bool
	Name         string
	FullName     string

	X int16
	Y int16

	Opened   bool
	Rotation uint16
	Reflect  uint16

	CurrentMode Mode

	Width  uint16
	Height uint16
}

func (m *Monitor) ListRotations() []uint16 {
	set := newSetUint16()
	for _, oname := range m.Outputs {
		if op, ok := GetDisplayInfo().outputNames[oname]; ok {
			oinfo, err := randr.GetOutputInfo(xcon, op, LastConfigTimeStamp).Reply()
			if err == nil && oinfo.Connection == randr.ConnectionConnected && oinfo.Crtc != 0 {
				cinfo, err := randr.GetCrtcInfo(xcon, oinfo.Crtc, LastConfigTimeStamp).Reply()
				if err != nil {
					continue
				}
				set.Add(parseRotations(cinfo.Rotations)...)
			}
		}
	}
	return set.Set()
}
func (m *Monitor) ListReflect() []uint16 {
	set := newSetUint16()
	for _, oname := range m.Outputs {
		if op, ok := GetDisplayInfo().outputNames[oname]; ok {
			oinfo, err := randr.GetOutputInfo(xcon, op, LastConfigTimeStamp).Reply()
			if err == nil && oinfo.Connection == randr.ConnectionConnected && oinfo.Crtc != 0 {
				cinfo, err := randr.GetCrtcInfo(xcon, oinfo.Crtc, LastConfigTimeStamp).Reply()
				if err != nil {
					continue
				}
				set.Add(parseReflects(cinfo.Rotations)...)
			}
		}
	}
	return set.Set()
}
func (m *Monitor) ListModes() []Mode {
	set := make(map[Mode]int)
	for _, oname := range m.Outputs {
		if op, ok := GetDisplayInfo().outputNames[oname]; ok {
			oinfo, err := randr.GetOutputInfo(xcon, op, LastConfigTimeStamp).Reply()
			if err != nil {
				continue
			}
			for _, m := range oinfo.Modes {
				mode := GetDisplayInfo().modes[m]
				set[mode] += 1
			}
		}
	}

	var r []Mode
	for k, n := range set {
		if n == len(m.Outputs) {
			r = append(r, k)
		}
	}
	return r
}

func (m *Monitor) SetRotation(v uint16) error {
	switch v {
	case 1, 2, 4, 8:
		break
	default:
		err := fmt.Errorf("SetRotation with invalid value ", v)
		Logger.Error(err)
		return err
	}
	m.cfg.Rotation = v
	m.setPropRotation(v)
	return nil
}
func (m *Monitor) SetReflect(v uint16) error {
	switch v {
	case 0, 16, 32, 48:
		break
	default:
		err := fmt.Errorf("SetReflect with invalid value ", v)
		Logger.Error(err)
		return err
	}
	m.cfg.Reflect = v
	m.setPropReflect(v)
	return nil
}

func (m *Monitor) SetPos(x, y int16) {
	m.cfg.X, m.cfg.Y = x, y
	m.setPropXY(x, y)
}

func (m *Monitor) SwitchOn(v bool) {
	m.cfg.Enabled = v
	m.setPropOpened(v)
}

func (m *Monitor) SetMode(id uint32) {
	for _, _m := range m.ListModes() {
		if _m.ID == id {
			m.setPropCurrentMode(_m)
			m.cfg.Width, m.cfg.Height, m.cfg.RefreshRate = _m.Width, _m.Height, _m.Rate
			return
		}
	}
}

func (m *Monitor) generateShell() string {
	code := ""
	names := strings.Split(m.Name, joinSeparator)
	for _, name := range names {
		code = fmt.Sprintf("%s --output %s", code, name)

		if m.Opened {
			if m.CurrentMode.ID == 0 {
				code = fmt.Sprintf("%s --auto", code)
			} else {
				code = fmt.Sprintf("%s --mode %dx%d --rate %f", code, m.CurrentMode.Width, m.CurrentMode.Height, m.CurrentMode.Rate)
			}
			code = fmt.Sprintf(" %s --pos %dx%d", code, m.X, m.Y)

			code = fmt.Sprintf("%s --scale 1x1", code)

			switch m.Rotation {
			case randr.RotationRotate90:
				code = fmt.Sprintf("%s --rotate right", code)
			case randr.RotationRotate180:
				code = fmt.Sprintf("%s --rotate inverted", code)
			case randr.RotationRotate270:
				code = fmt.Sprintf("%s --rotate left", code)
			default:
				code = fmt.Sprintf("%s --rotate normal", code)
			}
			switch m.Reflect {
			case randr.RotationReflectX:
				code = fmt.Sprintf("%s --reflect x", code)
			case randr.RotationReflectY:
				code = fmt.Sprintf("%s --reflect y", code)
			case randr.RotationReflectX | randr.RotationReflectY:
				code = fmt.Sprintf("%s --reflect xy", code)
			default:
				code = fmt.Sprintf("%s --reflect normal", code)
			}
		} else {
			code = fmt.Sprintf(" %s --off", code)
		}
	}
	return code + " "
}

func (m *Monitor) updateInfo() {
	op := GetDisplayInfo().outputNames[m.Outputs[0]]
	oinfo, err := randr.GetOutputInfo(xcon, op, LastConfigTimeStamp).Reply()
	if err != nil {
		Logger.Warning("updateInfo error:", err)
		return
	}
	if oinfo.Crtc == 0 {
		m.SwitchOn(false)
		m.setPropXY(0, 0)
		m.setPropWidth(0)
		m.setPropHeight(0)
		m.setPropRotation(1)
		m.setPropReflect(0)
		m.SetMode(0)
	} else {
		m.SwitchOn(true)
		cinfo, err := randr.GetCrtcInfo(xcon, oinfo.Crtc, LastConfigTimeStamp).Reply()
		if err != nil {
		}
		m.SetMode(uint32(cinfo.Mode))
		m.setPropXY(cinfo.X, cinfo.Y)
		m.setPropWidth(cinfo.Width)
		m.setPropHeight(cinfo.Height)
		rotation, reflect := parseRandR(cinfo.Rotation)
		m.setPropRotation(rotation)
		m.setPropReflect(reflect)
		m.setPropCurrentMode(GetDisplayInfo().modes[cinfo.Mode])
	}
}

func NewMonitor(dpy *Display, info *ConfigMonitor) *Monitor {
	m := &Monitor{}
	m.cfg = info
	m.Name = info.Name
	m.X, m.Y, m.Width, m.Height = info.X, info.Y, info.Width, info.Height
	m.Rotation, m.Reflect = info.Rotation, info.Reflect
	m.Opened = info.Enabled

	m.Outputs = info.Outputs
	runtime.SetFinalizer(m, func(o interface{}) { dbus.UnInstallObject(m) })
	m.IsComposited = len(m.Outputs) > 1

	m.FullName = m.Name

	if m.IsComposited {
		best := Mode{}
		for _, mode := range m.ListModes() {
			if mode.Width+mode.Height > best.Width+best.Height {
				best = mode
			}
		}
		m.BestMode = best
	}

	if info.currentMode != 0 {
		m.CurrentMode = GetDisplayInfo().modes[info.currentMode]
	} else {
		m.CurrentMode = m.BestMode
	}

	return m
}

func (m *Monitor) ensureSize(w, h uint16) {
	//find the nearest mode
	delta := float64(w + h)
	modeID := uint32(0)
	for _, mInfo := range m.ListModes() {
		t := math.Abs(float64((mInfo.Width + mInfo.Height) - (w + h)))
		if t <= delta {
			delta = t
			modeID = mInfo.ID
			if modeID == m.BestMode.ID {
				break
			}
		}
	}
	if modeID != 0 {
		m.SetMode(modeID)
		m.setPropWidth(w)
		m.setPropHeight(h)
		if delta != 0 {
			Logger.Warningf("Can't ensureSize(%s) to %d %d", m.Name, w, h)
		}
	}
}
