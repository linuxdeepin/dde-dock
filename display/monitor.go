package main

import "github.com/BurntSushi/xgb/randr"
import "dlib/dbus"
import "strings"
import "fmt"
import "math"
import "runtime"

const joinSeparator = "|"

type Monitor struct {
	outputs   []randr.Output
	Rotations []uint16
	Reflects  []uint16
	Modes     []Mode

	BacklightRange map[string]int32

	IsComposited bool
	IsPrimary    bool

	Name     string
	FullName string

	X               int16
	Y               int16
	relativePosInfo [2]string

	Opened   bool
	Rotation uint16 `access:"readwrite"`
	Reflect  uint16 `access:"readwrite"`

	Brightness float64

	Mode     Mode
	BestMode Mode

	Width     uint16
	Height    uint16
	scaleMode bool
	Rate      float64
}

func (m *Monitor) SetBrightness(v float64) {
	m.setPropBrightness(v)
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
			if m.Mode.ID == 0 {
				m.Mode = m.BestMode
			}
			code = fmt.Sprintf("%s --mode %dx%d --rate %f", code, m.Mode.Width, m.Mode.Height, m.Mode.Rate)
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

			maxBacklight := m.BacklightRange[name]
			if maxBacklight == 0 {
				code = fmt.Sprintf("%s --brightness %f", code, m.Brightness)
			} else if maxBacklight > 0 {
				code = fmt.Sprintf("%s --brightness 1 --set Backlight %d", code, int32(m.Brightness*float64(maxBacklight)))
			} else {
				code = fmt.Sprintf("%s --brightness 1", code)
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
		m.Opened = false
		m.X = 0
		m.Y = 0
		m.Width = 0
		m.Height = 0
		m.Rotation, m.Reflect = 1, 0
		m.Mode = Mode{}
	} else {
		m.Opened = true
		cinfo, err := randr.GetCrtcInfo(X, oinfo.Crtc, LastConfigTimeStamp).Reply()
		if err != nil {
		}
		m.X = cinfo.X
		m.Y = cinfo.Y
		m.Width = cinfo.Width
		m.Height = cinfo.Height
		m.Rotation, m.Reflect = parseRandR(cinfo.Rotation)
		m.Mode = DPY.modes[cinfo.Mode]
	}
}

func NewMonitor(outputs []randr.Output) *Monitor {
	m := &Monitor{}
	runtime.SetFinalizer(m, func(o interface{}) { dbus.UnInstallObject(m) })
	m.outputs = make([]randr.Output, len(outputs))
	m.BacklightRange = make(map[string]int32)
	m.Brightness = 1
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
		m.BacklightRange[string(oinfo.Name)] = queryBacklightRange(X, op)
		names = append(names, string(oinfo.Name))
		for i := uint16(0); i < oinfo.NumModes; i++ {
			modeMap[oinfo.Modes[i]] += 1
			if i == 0 {
				m.BestMode = DPY.modes[oinfo.Modes[0]]
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
	m.Modes = make([]Mode, 0)
	m.Rotations = make([]uint16, 0)
	m.Reflects = make([]uint16, 0)
	for mode, value := range modeMap {
		if value == len(m.outputs) {
			m.Modes = append(m.Modes, DPY.modes[mode])
		}
	}
	for r, value := range rotationMap {
		if value == len(m.outputs) {
			m.Rotations = append(m.Rotations, r)
		}
	}
	for r, value := range reflectMap {
		if value == len(m.outputs) {
			m.Reflects = append(m.Reflects, r)
		}
	}

	m.Name = strings.Join(names, joinSeparator)
	m.FullName = m.Name

	m.updateInfo()

	return m
}

func (m *Monitor) SetMode(id uint32) {
	for _, _m := range m.Modes {
		if _m.ID == id {
			m.Mode = _m
			m.scaleMode = false
			return
		}
	}
}
func (m *Monitor) ensureSize(w, h uint16) {
	//find the nearest mode
	delta := float64(w + h)
	modeID := uint32(0)
	for _, mInfo := range m.Modes {
		t := math.Abs(float64((mInfo.Width + mInfo.Height) - (w + h)))
		if t < delta {
			delta = t
			modeID = mInfo.ID
		}
	}
	if modeID != 0 {
		m.SetMode(modeID)
		m.Width, m.Height = w, h
		if delta != 0 {
			m.scaleMode = true
		}
	}
	fmt.Println(m.Name, " EnsureSize", delta)
}

func (m *Monitor) SetPos(x, y int16) {
	m.relativePosInfo[0] = ""
	m.relativePosInfo[1] = ""
	m.X, m.Y = x, y
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
	m.Opened = v
}

func (m *Monitor) setPrimary(v bool) {
	m.IsPrimary = v
}

func (m *Monitor) isContain(op randr.Output) bool {
	for _, o := range m.outputs {
		if o == op {
			return true
		}
	}
	return false
}
