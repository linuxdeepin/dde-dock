package main

import (
	"dlib/dbus"
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"math"
)

var (
	_        = fmt.Println
	X, _     = xgb.NewConn()
	atomEDID = getAtom(X, "EDID")
	Root     = xproto.Setup(X).DefaultScreen(X).Root
)

func getAtom(c *xgb.Conn, name string) xproto.Atom {
	r, err := xproto.InternAtom(c, false, uint16(len(name)), name).Reply()
	if err != nil {
		return xproto.AtomNone
	}
	return r.Atom
}

func getOutputName(edid []byte, defaultName string) string {
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
	return defaultName
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
	if info.Connection != randr.ConnectionConnected {
		return nil
	}

	// Nvidia driver which support Randr 1.4 will show an additional connected output which I didn't know it's exactly function. So simply filter it.
	if info.MmWidth == 0 || info.MmHeight == 0 {
		return nil
	}
	fmt.Println("Output:", core)

	r, _ := randr.GetOutputProperty(X, core, atomEDID, xproto.AtomInteger, 0, 1024, false, false).Reply()

	ret := &Output{core, getOutputName(r.Data, string(info.Name)), nil}
	for _, m := range info.Modes {
		info := display.modes[m]
		ret.Modes = append(ret.Modes, buildMode(info))
	}
	fmt.Println("OutputMode:", info)
	return ret
}
func buildMode(info randr.ModeInfo) Mode {
	vTotal := info.Vtotal

	if info.ModeFlags&randr.ModeFlagDoubleScan != 0 {
		vTotal *= 2
	}
	if info.ModeFlags&randr.ModeFlagInterlace != 0 {
		vTotal /= 2
	}

	rate := float64(info.DotClock) / float64(uint32(info.Htotal)*uint32(vTotal))
	rate = math.Floor(rate + 0.5)
	return Mode{info.Width, info.Height, uint16(rate)}
}

type Display struct {
	outputs          []*Output
	modes            map[randr.Mode]randr.ModeInfo
	SupportRotations []uint32
	CurrentRotation  uint32
	CurrentMode      Mode

	PrimaryRect    xproto.Rectangle
	PrimaryChanged func(xproto.Rectangle)
}

func (display *Display) updatePrimary() {
	r, _ := randr.GetOutputPrimary(X, Root).Reply()
	if r.Output == 0 {
		display.PrimaryRect = xproto.Rectangle{0, 0, display.CurrentMode.Width, display.CurrentMode.Height}
	} else {
		oinfo, err := randr.GetOutputInfo(X, r.Output, 0).Reply()
		if err == nil {
			cinfo, err := randr.GetCrtcInfo(X, oinfo.Crtc, 0).Reply()
			if err == nil {
				display.PrimaryRect = xproto.Rectangle{cinfo.X, cinfo.Y, cinfo.Width, cinfo.Height}
			}
		}
	}
	if display.PrimaryChanged != nil {
		display.PrimaryChanged(display.PrimaryRect)
	}
	fmt.Println("PrimaryMode", display.PrimaryRect)
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
	resources, err := randr.GetScreenResources(X, Root).Reply()
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
func (dispaly *Display) SetPrimary(output uint32) {
	randr.SetOutputPrimary(X, Root, randr.Output(output))
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

	resources, err := randr.GetScreenResources(X, Root).Reply()
	sinfo, err := getScreenInfo(Root)
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
	r.updatePrimary()
	return r
}

func (dpy *Display) listener() {
	for {
		e, err := X.WaitForEvent()
		if err != nil {
			continue
		}
		switch e.(type) {
		case randr.NotifyEvent:
			ee := e.(randr.NotifyEvent)
			switch ee.SubCode {
			case randr.NotifyOutputChange:
				fmt.Println("OC:", ee.U.Oc)
				dpy.updatePrimary()
				dpy.updateOutput()
			case randr.NotifyOutputProperty:
				fmt.Println("OC:", ee.U.Op)
				dpy.updatePrimary()
				dpy.updateOutput()
			}
		}
	}
}

func main() {
	dpy := NewDisplay()
	dbus.InstallOnSession(dpy)
	dbus.DealWithUnhandledMessage()

	randr.SelectInput(X, Root, randr.NotifyMaskOutputChange|randr.NotifyMaskOutputProperty|randr.NotifyMaskScreenChange)
	go dpy.listener()
	select {}
}
