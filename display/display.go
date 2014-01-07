package main

import (
	"dlib/dbus"
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
)

var (
	_        = fmt.Println
	X, _     = xgb.NewConn()
	Root     = xproto.Setup(X).DefaultScreen(X).Root
	atomEDID = getAtom(X, "EDID")
)

func init() {
	randr.Init(X)
	randr.QueryVersion(X, 1, 13)
}

type Display struct {
	modes map[randr.Mode]randr.ModeInfo

	Outputs []*Output

	Width  uint16
	Height uint16

	Rotation  uint16 `access:readwrite`
	rotations uint16

	PrimaryOutput *Output `access:readwrite`
	//used by deepin-dock/launcher/desktop
	PrimaryRect    xproto.Rectangle
	PrimaryChanged func(xproto.Rectangle)
}

func (dpy *Display) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Display",
		"/com/deepin/daemon/Display",
		"com.deepin.daemon.Display",
	}
}

func (dpy *Display) updateScreenInfo() {
}

func NewDisplay() *Display {
	dpy := &Display{}

	dpy.modes = make(map[randr.Mode]randr.ModeInfo)

	sinfo, err := getScreenInfo(Root)
	dpy.rotations = uint16(sinfo.Rotations)
	dpy.Rotation = sinfo.Rotation

	if err != nil {
		panic(err)
	}

	resources, err := randr.GetScreenResources(X, Root).Reply()

	for _, m := range resources.Modes {
		dpy.modes[randr.Mode(m.Id)] = m
	}

	for _, output := range resources.Outputs {
		dpy.updateOutputList(output)
	}

	size := sinfo.Sizes[sinfo.SizeID]
	dpy.Width = size.Width
	dpy.Height = size.Height

	dpy.updatePrimary()

	randr.SelectInput(X, Root, randr.NotifyMaskOutputChange|randr.NotifyMaskOutputProperty|randr.NotifyMaskScreenChange)
	go dpy.listener()
	return dpy
}

func (dpy *Display) ShowInfoOnScreen() {
}
func (dpy *Display) ListRotations() []uint8 {
	return parseRotations(dpy.rotations)
}

func (dpy *Display) updatePrimary() {
	r, _ := randr.GetOutputPrimary(X, Root).Reply()
	if r.Output == 0 {
		dpy.PrimaryOutput = nil
		dpy.PrimaryRect = xproto.Rectangle{0, 0, dpy.Width, dpy.Height}
	} else {
		dpy.PrimaryOutput = queryOutput(dpy, r.Output)
		dpy.PrimaryRect = dpy.PrimaryOutput.Allocation
	}

	if dpy.PrimaryChanged != nil {
		dpy.PrimaryChanged(dpy.PrimaryRect)
	}
}

func (dpy *Display) updateOutputList(output randr.Output) {
	op := queryOutput(dpy, output)
	if op == nil {
		if op = NewOutput(dpy, output); op != nil {
			dpy.Outputs = append(dpy.Outputs, op)
		}
	} else {
		info, err := randr.GetOutputInfo(X, output, xproto.TimeCurrentTime).Reply()
		if err != nil {
			panic(err)
		}
		op.update(dpy, info)
	}
}

func (dpy *Display) SetPrimary(output uint32) {
	randr.SetOutputPrimary(X, Root, randr.Output(output))
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
			case randr.NotifyCrtcChange:
				info := ee.U.Cc
				if op := queryOutputByCrtc(dpy, info.Crtc); op != nil {
					op.updateCrtc(dpy)
				}
			case randr.NotifyOutputChange:
				info := ee.U.Oc
				if op := queryOutput(dpy, info.Output); op != nil {
					outputInfo, err := randr.GetOutputInfo(X, info.Output, xproto.TimeCurrentTime).Reply()
					if err != nil {
						fmt.Println(err)
						continue
					}
					op.update(dpy, outputInfo)
				}
			}
		case randr.ScreenChangeNotifyEvent:
			ee := e.(randr.ScreenChangeNotifyEvent)
			dpy.updatePrimary()
			dpy.updateScreenSize(ee.Width, ee.Height)
			dpy.updateRotation(uint16(ee.Rotation))
		}
	}
}

func (dpy *Display) updateScreenSize(width uint16, height uint16) {
	if dpy.Width != width {
		dpy.Width = width
		dbus.NotifyChange(dpy, "Width")
	}
	if dpy.Height != height {
		dpy.Height = height
		dbus.NotifyChange(dpy, "Height")
	}
}
func (dpy *Display) updateRotation(rotation uint16) {
	if dpy.Rotation != rotation {
		dpy.Rotation = rotation
		dbus.NotifyChange(dpy, "Rotation")
	}
}

func main() {
	dpy := NewDisplay()
	dbus.InstallOnSession(dpy)
	dbus.DealWithUnhandledMessage()

	select {}
}
