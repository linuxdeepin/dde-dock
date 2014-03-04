package main

import "github.com/BurntSushi/xgb/randr"
import "dlib/dbus"
import "fmt"
import "strings"
import "math"
import "runtime"

const joinSeparator = "="

type Monitor struct {
	outputs   []randr.Output
	rotations []uint16
	reflects  []uint16
	modes     []Mode
	BestMode  Mode
	scaleMode bool

	backlightRange map[string]int32
	IsComposited   bool
	Name           string
	FullName       string

	IsPrimary       bool
	X               int16
	Y               int16
	relativePosInfo [2]string

	Opened   bool
	Rotation uint16
	Reflect  uint16

	Brightness float64

	CurrentMode Mode

	Width  uint16
	Height uint16
}

func (m *Monitor) ListRotations() []uint16 {
	return m.rotations
}
func (m *Monitor) ListReflect() []uint16 {
	return m.reflects
}
func (m *Monitor) ListModes() []Mode {
	return m.modes
}
func (m *Monitor) SetRotation(v uint16) {
	m.setPropRotation(v)
}
func (m *Monitor) SetReflect(v uint16) {
	m.setPropReflect(v)
}

func (m *Monitor) ChangeBrightness(v float64) {
	if v < 0.1 {
		v = 0.1
	} else if v > 1 {
		v = 1
	}
	names := strings.Split(m.Name, joinSeparator)
	code := "xrandr "
	for _, name := range names {
		code = fmt.Sprintf("%s --output %s", code, name)
		maxBacklight := m.backlightRange[name]
		if maxBacklight == 0 {
			code = fmt.Sprintf("%s --brightness %f", code, v)
		} else if maxBacklight > 0 {
			code = fmt.Sprintf("%s --brightness 1 --set Backlight %d", code, int32(v*float64(maxBacklight)))
		} else {
			code = fmt.Sprintf("%s --brightness 1", code)
		}
	}
	runCode(code)
	m.setPropBrightness(v)
}

func (m *Monitor) SetPos(x, y int16) {
	m.relativePosInfo[0] = ""
	m.relativePosInfo[1] = ""
	m.setPropXY(x, y)
}

func (m *Monitor) SetRelativePos(reference string, pos string) {
	switch pos {
	case "above", "below", "left-of", "right-of":
		for _, name := range strings.Split(m.Name, joinSeparator) {
			if name == reference {
				return
			}
		}
		m.relativePosInfo[0], m.relativePosInfo[1] = pos, reference
	}
}

func (m *Monitor) SwitchOn(v bool) {
	m.setPropOpened(v)
}

func (m *Monitor) setPrimary(v bool) {
	m.setPropIsPrimary(v)
}

func (m *Monitor) SetMode(id uint32) {
	for _, _m := range m.ListModes() {
		if _m.ID == id {
			m.setPropCurrentMode(_m)
			m.scaleMode = false
			return
		}
	}
}

func (m *Monitor) generateShell() string {
	code := ""
	names := strings.Split(m.Name, joinSeparator)
	for _, name := range names {
		code = fmt.Sprintf("%s --output %s", code, name)

		if m.IsPrimary {
			code = fmt.Sprintf(" %s --primary", code)
		}
		if m.Opened {
			if m.CurrentMode.ID == 0 {
				m.setPropCurrentMode(m.BestMode)
				code = fmt.Sprintf("%s --auto", code)
			} else {
				fmt.Println("CurrentModeID:", m.Name, m.CurrentMode.ID)
				code = fmt.Sprintf("%s --mode %dx%d --rate %f", code, m.CurrentMode.Width, m.CurrentMode.Height, m.CurrentMode.Rate)
			}
			if len(m.relativePosInfo[0]) != 0 && len(m.relativePosInfo[1]) != 0 {
				code = fmt.Sprintf(" %s --%s %s", code, m.relativePosInfo[0], m.relativePosInfo[1])
			} else {
				code = fmt.Sprintf(" %s --pos %dx%d", code, m.X, m.Y)
			}

			if m.scaleMode {
				code = fmt.Sprintf("%s --scale-from %dx%d", code, m.Width, m.Height)
			} else {
				code = fmt.Sprintf("%s --scale 1x1", code)
			}

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
	op := m.outputs[0]
	oinfo, err := randr.GetOutputInfo(X, op, LastConfigTimeStamp).Reply()
	if err != nil {
	}
	if oinfo.Crtc == 0 {
		m.SwitchOn(false)
		m.setPropXY(0, 0)
		m.setPropWidth(0)
		m.setPropHeight(0)
		m.setPropRotation(1)
		m.setPropReflect(0)
		m.setPropCurrentMode(Mode{})
	} else {
		m.SwitchOn(true)
		cinfo, err := randr.GetCrtcInfo(X, oinfo.Crtc, LastConfigTimeStamp).Reply()
		if err != nil {
		}
		m.setPropXY(cinfo.X, cinfo.Y)
		m.setPropWidth(cinfo.Width)
		m.setPropHeight(cinfo.Height)
		rotation, reflect := parseRandR(cinfo.Rotation)
		m.setPropRotation(rotation)
		m.setPropReflect(reflect)
		m.setPropCurrentMode(DPY.modes[cinfo.Mode])
	}
	if ok, backlight := supportedBacklight(X, op); ok {
		m.setPropBrightness(backlight)
	}
}

func (m *Monitor) WorkaroundBacklight() {
	if ok, backlight := supportedBacklight(X, m.outputs[0]); ok {
		m.setPropBrightness(backlight)
	}
}

func NewMonitor(outputs []randr.Output) *Monitor {
	m := &Monitor{}
	runtime.SetFinalizer(m, func(o interface{}) { dbus.UnInstallObject(m) })
	m.outputs = make([]randr.Output, len(outputs))
	m.backlightRange = make(map[string]int32)
	m.setPropBrightness(1)
	m.IsComposited = len(outputs) > 1
	copy(m.outputs, outputs)

	modeMap := make(map[randr.Mode]int)
	rotationMap := make(map[uint16]int)
	reflectMap := make(map[uint16]int)
	names := make([]string, 0)

	for _, op := range m.outputs {
		oinfo, err := randr.GetOutputInfo(X, op, LastConfigTimeStamp).Reply()
		if err != nil {
			continue
		}
		m.backlightRange[string(oinfo.Name)] = queryBacklightRange(X, op)
		names = append(names, string(oinfo.Name))
		for i := uint16(0); i < oinfo.NumModes; i++ {
			modeMap[oinfo.Modes[i]] += 1
			if i == 0 {
				m.BestMode = DPY.modes[oinfo.Modes[0]]
				fmt.Println("BestID:", m.Name, m.BestMode.ID)
			}
		}
		if oinfo.Crtc != 0 {
			cinfo, err := randr.GetCrtcInfo(X, oinfo.Crtc, LastConfigTimeStamp).Reply()
			if err != nil {
				continue
			}
			for _, rotation := range parseRotations(cinfo.Rotations) {
				rotationMap[rotation] += 1
			}
			for _, reflect := range parseReflects(cinfo.Rotations) {
				reflectMap[reflect] += 1
			}
		}

	}
	m.modes = make([]Mode, 0)
	m.rotations = make([]uint16, 0)
	m.reflects = make([]uint16, 0)
	for mode, value := range modeMap {
		if value == len(m.outputs) {
			m.modes = append(m.modes, DPY.modes[mode])
		}
	}
	for r, value := range rotationMap {
		if value == len(m.outputs) {
			m.rotations = append(m.rotations, r)
		}
	}
	for r, value := range reflectMap {
		if value == len(m.outputs) {
			m.reflects = append(m.reflects, r)
		}
	}

	m.Name = strings.Join(names, joinSeparator)
	m.FullName = m.Name

	m.updateInfo()

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
			m.scaleMode = true
		}
	}
}

func (m *Monitor) isContain(op randr.Output) bool {
	for _, o := range m.outputs {
		if o == op {
			return true
		}
	}
	return false
}
