package main

import (
	"dlib/dbus"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"fmt"
	"math"
)

var (
	_        = fmt.Println
	X, _     = xgb.NewConn()
	atomEDID = getAtom(X, "EDID")
)

func getAtom(c *xgb.Conn, name string) xproto.Atom {
	r, err := xproto.InternAtom(c, false, uint16(len(name)), name).Reply()
	if err != nil {
		return xproto.AtomNone
	}
	return r.Atom
}

func getOutputName(edid []byte) string {
	if len(edid) == 128 {
		timingDescriptor := edid[36:]
		for i := 0; i < 4; i++ {
			block := timingDescriptor[i*18 : (i+1)*18]
			if block[3] == 0xfc { //descriptor type == Monitor Name
				data := block[5:]
				for i := 0; i < 13; i++ {
					if data[i] == 0xa && data[i+1] == 0x20 && data[i+2] == 0x20 {
						return string(data[:i])
					}
				}
			}
		}
	}
	return "Unknow"
}

type Mode struct {
	Width  uint16
	Height uint16
	Rate   uint16
}
type Output struct {
	core  randr.Output
	Name  string
	Modes []Mode
}

func NewOutput(display *Display, core randr.Output) *Output {
	info, err := randr.GetOutputInfo(X, core, xproto.TimeCurrentTime).Reply()
	if err != nil {
		panic(err)
	}
	if len(info.Modes) == 0 {
		return nil
	}

	r, _ := randr.GetOutputProperty(X, core, atomEDID, xproto.AtomInteger, 0, 1024, false, false).Reply()

	ret := &Output{core, getOutputName(r.Data), nil}
	for _, m := range info.Modes {
		info := display.modes[m]
		vTotal := info.Vtotal

		if info.ModeFlags&randr.ModeFlagDoubleScan != 0 {
			vTotal *= 2
		}
		if info.ModeFlags&randr.ModeFlagInterlace != 0 {
			vTotal /= 2
		}

		rate := float64(info.DotClock) / float64(uint32(info.Htotal)*uint32(vTotal))
		rate = math.Floor(rate + 0.5)
		ret.Modes = append(ret.Modes, Mode{info.Width, info.Height, uint16(rate)})
	}
	fmt.Println("OutputMode:", info)
	return ret
}

func (output *Output) MonitorName() string {
	r, _ := randr.GetOutputProperty(X, output.core, atomEDID, xproto.AtomInteger, 0, 1024, false, false).Reply()
	return getOutputName(r.Data)
}

type Display struct {
	outputs          []*Output
	modes            map[randr.Mode]randr.ModeInfo
	SupportRotations []uint32
	CurrentRotation  uint32
	CurrentMode      Mode
}

func (display *Display) GetOutpus() []*Output {
	return display.outputs
}
func (output *Output) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Display",
		"/com/deepin/daemon/Display/Output" + output.Name,
		"com.deepin.daemon.Display.Output",
	}
}

func (display *Display) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Display",
		"/com/deepin/daemon/Display",
		"com.deepin.daemon.Display",
	}
}

func (display *Display) updateOutput() {
	root := xproto.Setup(X).DefaultScreen(X).Root
	resources, err := randr.GetScreenResources(X, root).Reply()
	if err != nil {
		panic(err)
	}
	// Iterate through all of the outputs and show some of their info.
	for _, output := range resources.Outputs {
		o := NewOutput(display, output)
		if o != nil {
			display.outputs = append(display.outputs, o)
		}
	}
}

// xgb/randr.GetScreenInfo will panic at this moment.( https://github.com/BurntSushi/xgb/issues/20)
func getScreenInfo(root xproto.Window) (*randr.GetScreenInfoReply, error) {
	cook := randr.GetScreenInfo(X, root)

	buf, err := cook.Cookie.Reply()
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}
	v := new(randr.GetScreenInfoReply)
	b := 1 // skip reply determinant

	v.Rotations = buf[b]
	b += 1

	v.Sequence = xgb.Get16(buf[b:])
	b += 2

	v.Length = xgb.Get32(buf[b:]) // 4-byte units
	b += 4

	v.Root = xproto.Window(xgb.Get32(buf[b:]))
	b += 4

	v.Timestamp = xproto.Timestamp(xgb.Get32(buf[b:]))
	b += 4

	v.ConfigTimestamp = xproto.Timestamp(xgb.Get32(buf[b:]))
	b += 4

	v.NSizes = xgb.Get16(buf[b:])
	b += 2

	v.SizeID = xgb.Get16(buf[b:])
	b += 2

	v.Rotation = xgb.Get16(buf[b:])
	b += 2

	v.Rate = xgb.Get16(buf[b:])
	b += 2

	v.NInfo = xgb.Get16(buf[b:])
	b += 2

	b += 2 // padding

	v.Sizes = make([]randr.ScreenSize, v.NSizes)
	b += randr.ScreenSizeReadList(buf[b:], v.Sizes)

	return v, nil
}

func NewDisplay() *Display {
	randr.Init(X)
	randr.QueryVersion(X, 1, 13)
	r := &Display{}
	r.modes = make(map[randr.Mode]randr.ModeInfo)

	root := xproto.Setup(X).DefaultScreen(X).Root
	resources, err := randr.GetScreenResources(X, root).Reply()
	sinfo, err := getScreenInfo(root)
	size := sinfo.Sizes[sinfo.SizeID]
	r.CurrentMode = Mode{size.Width, size.Height, sinfo.Rate}
	fmt.Println("CurrentRate:", sinfo.Rate, sinfo.SizeID, sinfo.Sizes)
	rotations := sinfo.Rotations
	if rotations&randr.RotationRotate0 == randr.RotationRotate0 {
		r.SupportRotations = append(r.SupportRotations, randr.RotationRotate0)
	}
	if rotations&randr.RotationRotate90 == randr.RotationRotate90 {
		r.SupportRotations = append(r.SupportRotations, randr.RotationRotate90)
	}
	if rotations&randr.RotationRotate180 == randr.RotationRotate180 {
		r.SupportRotations = append(r.SupportRotations, randr.RotationRotate180)
	}
	if rotations&randr.RotationRotate270 == randr.RotationRotate270 {
		r.SupportRotations = append(r.SupportRotations, randr.RotationRotate270)
	}
	if rotations&randr.RotationReflectX == randr.RotationReflectX {
		r.SupportRotations = append(r.SupportRotations, randr.RotationReflectX)
	}
	if rotations&randr.RotationReflectY == randr.RotationReflectY {
		r.SupportRotations = append(r.SupportRotations, randr.RotationReflectY)
	}
	r.CurrentRotation = uint32(sinfo.Rotation)
	if err != nil {
		panic(err)
	}
	for _, m := range resources.Modes {
		r.modes[randr.Mode(m.Id)] = m
	}
	r.updateOutput()
	return r
}

func main() {
	dpy := NewDisplay()
	dbus.InstallOnSession(dpy)
	dbus.DealWithUnhandledMessage()
	select {}
}
