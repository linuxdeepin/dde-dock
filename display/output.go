package main

import "github.com/BurntSushi/xgb/xproto"
import "github.com/BurntSushi/xgb/randr"
import "dlib/dbus"
import "fmt"

type Mode struct {
	Width  uint16
	Height uint16
	Rate   uint16
}
type Output struct {
	bestMode  randr.Mode
	modes     []Mode
	rotations uint16
	crtc      randr.Crtc

	Identify randr.Output
	Name     string
	Type     uint8

	Mode         Mode             `access:"readwrite"`
	Allocation   xproto.Rectangle `access:"readwrite"`
	AdjustMethod uint8            `access:"readwrite"`

	Rotation   uint16  `access:"readwrite"`
	Opened     bool    `access:"readwrite"`
	Brightness float64 `access:"readwrite"`
}

func (output *Output) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Display",
		fmt.Sprintf("/com/deepin/daemon/Display/Output%d", output.Identify),
		"com.deepin.daemon.Display.Output",
	}
}

func (op *Output) ListModes() []Mode {
	return op.modes
}
func (op *Output) ListRotations() []uint8 {
	return parseRotations(op.rotations)
}

func (op *Output) updateCrtc(dpy *Display) {
	if op.crtc != 0 {
		info, err := randr.GetCrtcInfo(X, op.crtc, 0).Reply()
		if err != nil {
			panic(err)
		}
		op.rotations = info.Rotations
		op.Rotation = info.Rotation

		op.Allocation = xproto.Rectangle{info.X, info.Y, info.Width, info.Height}

		op.Mode = buildMode(dpy.modes[info.Mode])
	} else {
		op.Rotation = 0
		op.Allocation = xproto.Rectangle{0, 0, 0, 0}

		op.Mode = Mode{0, 0, 0}
	}
	dbus.NotifyChange(op, "Allocation")
	dbus.NotifyChange(op, "Rotation")
}

func (op *Output) update(dpy *Display, info *randr.GetOutputInfoReply) {
	op.crtc = info.Crtc
	op.Opened = info.Crtc != 0
	dbus.NotifyChange(op, "Opened")
	op.bestMode = info.Modes[0]
	for _, m := range info.Modes {
		info := dpy.modes[m]
		op.modes = append(op.modes, buildMode(info))
	}
	dbus.NotifyChange(op, "Mode")
}

func (op *Output) setOpened(v bool) {
	if op.Opened != v {
		if v == true {
			oinfo, err := randr.GetOutputInfo(X, op.Identify, 0).Reply()
			if err != nil {
				panic(err)
			}
			for _, crtc := range oinfo.Crtcs {
				s, err := randr.SetCrtcConfig(X, crtc, 0, 0, op.Allocation.X, op.Allocation.Y, op.bestMode, 1, []randr.Output{op.Identify}).Reply()
				if err == nil {
					break
				}
				fmt.Println("AAAA:", s, err, crtc, op.bestMode, op.Rotation, op.Identify)
			}
		} else {
			fmt.Println("OXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
			s, err := randr.SetCrtcConfig(X, op.crtc, 0, 0, op.Allocation.X, op.Allocation.Y, 0, op.Rotation, nil).Reply()
			fmt.Println(s, err)
		}
	}
}

func NewOutput(dpy *Display, core randr.Output) *Output {
	info, err := randr.GetOutputInfo(X, core, xproto.TimeCurrentTime).Reply()
	if err != nil {
		panic(err)
		return nil
	}
	if info.Connection != randr.ConnectionConnected {
		return nil
	}

	// Nvidia driver which support Randr 1.4 will show an additional connected output which I didn't know it's exactly function. So simply filter it.
	if info.MmWidth == 0 || info.MmHeight == 0 {
		return nil
	}

	edidProp, _ := randr.GetOutputProperty(X, core, atomEDID, xproto.AtomInteger, 0, 1024, false, false).Reply()

	op := &Output{
		Identify: core,
		Name:     getOutputName(edidProp.Data, string(info.Name)),
	}
	op.update(dpy, info)
	op.updateCrtc(dpy)
	return op
}
