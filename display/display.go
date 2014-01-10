package main

import (
	"dlib/dbus"
	"dlib/logger"
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
)

var (
	_                   = fmt.Println
	X, _                = xgb.NewConn()
	DefaultScreen       = xproto.Setup(X).DefaultScreen(X)
	Root                = DefaultScreen.Root
	atomEDID            = getAtom(X, "EDID")
	LastConfigTimeStamp = xproto.Timestamp(0)
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
	Reflect   uint16 `access:readwrite`
	rotations uint16

	PrimaryOutput *Output `access:readwrite`
	//used by deepin-dock/launcher/desktop
	PrimaryRect    xproto.Rectangle
	PrimaryChanged func(xproto.Rectangle)
}

func NewDisplay() *Display {
	dpy := &Display{}

	dpy.modes = make(map[randr.Mode]randr.ModeInfo)

	sinfo, err := getScreenInfo(Root)
	dpy.setPropRotation(uint16(sinfo.Rotations))
	dpy.updateRotationAndRelfect(sinfo.Rotation)

	if err != nil {
		panic("GetScreenInfo Failed:" + err.Error())
	}

	resources, err := randr.GetScreenResources(X, Root).Reply()
	LastConfigTimeStamp = resources.ConfigTimestamp

	if err != nil {
		panic("GetScreenResources failed:" + err.Error())
	}

	for _, m := range resources.Modes {
		dpy.modes[randr.Mode(m.Id)] = m
	}

	for _, output := range resources.Outputs {
		dpy.updateOutputList(output)
	}

	dpy.updateScreenSize(DefaultScreen.WidthInPixels, DefaultScreen.HeightInPixels)

	dpy.updatePrimary()

	randr.SelectInput(X, Root, randr.NotifyMaskOutputChange|randr.NotifyMaskOutputProperty|randr.NotifyMaskScreenChange)
	go dpy.listener()
	return dpy
}

func (dpy *Display) ShowInfoOnScreen() {
}
func (dpy *Display) ListRotations() []uint16 {
	return parseRotations(dpy.rotations)
}
func (dpy *Display) ListReflect() []uint16 {
	return parseReflects(dpy.rotations)
}

func (dpy *Display) updatePrimary() {
	r, _ := randr.GetOutputPrimary(X, Root).Reply()
	if r.Output == 0 {
		dpy.setPropPrimaryOutput(nil)
		dpy.setPropPrimaryRect(xproto.Rectangle{0, 0, dpy.Width, dpy.Height})
	} else if dpy.setPropPrimaryOutput(queryOutput(dpy, r.Output)); dpy.PrimaryOutput == nil {
		//to avoid repeatedly trigger ScreenChangeNotifyEvent
		if len(dpy.Outputs) != 0 {
			//this output is invalid or disconnected, so set OutputPrimary to None
			randr.SetOutputPrimary(X, Root, 0)
		}
		return
	} else {
		dpy.setPropPrimaryRect(dpy.PrimaryOutput.Allocation)
	}
}

func (dpy *Display) updateOutputList(output randr.Output) {
	op := queryOutput(dpy, output)
	if op == nil {
		if op = NewOutput(dpy, output); op != nil {
			dpy.setPropOutputs(append(dpy.Outputs, op))
		}
	} else {
		info, err := randr.GetOutputInfo(X, output, xproto.TimeCurrentTime).Reply()
		if err != nil {
			panic("GetOutputInfo failed:" + err.Error())
		}
		op.update(dpy, info)
	}
}
func (dpy *Display) removeOutput(output randr.Output) {
	var newOutput []*Output
	for _, op := range dpy.Outputs {
		if op.Identify != output {
			newOutput = append(newOutput, op)
		} else {
			dbus.UnInstallObject(op)
		}
	}
	if len(newOutput) != len(dpy.Outputs) {
		dpy.setPropOutputs(newOutput)
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
				switch info.Connection {
				case randr.ConnectionConnected:
					dpy.updateOutputList(info.Output)
				case randr.ConnectionDisconnected, randr.ConnectionUnknown:
					dpy.removeOutput(info.Output)
				}
			}
		case randr.ScreenChangeNotifyEvent:
			ee := e.(randr.ScreenChangeNotifyEvent)
			if ee.ConfigTimestamp <= LastConfigTimeStamp {
				logger.Println("Recived an invalid ScreenChangeNotifyEvent", ee)
				continue
			}
			LastConfigTimeStamp = ee.ConfigTimestamp

			DefaultScreen = xproto.Setup(X).DefaultScreen(X)
			dpy.updateRotationAndRelfect(uint16(ee.Rotation))

			for _, op := range dpy.Outputs {
				op.updateCrtc(dpy)
			}
			width, height := parseScreenSize(dpy.Outputs)
			dpy.updateScreenSize(width, height)
			fmt.Println("UpdateScreenSize:", width, height, DefaultScreen.WidthInPixels, DefaultScreen.HeightInPixels)
			dpy.updatePrimary() //depend on updateScreenSize when there hasn't an primary output
		}
	}
}

func (dpy *Display) updateScreenSize(width uint16, height uint16) {
	if DefaultScreen.WidthInPixels != width || DefaultScreen.HeightInPixels != height {
		//SetScreenSize will cause emit ScreenChangeNotifyEvent, so next time we will jump in the real "else" branch to set dpy.Width/Height
		randr.SetScreenSize(X, Root, width, height, uint32(DefaultScreen.WidthInMillimeters), uint32(DefaultScreen.HeightInMillimeters)).Reply()
	}
	dpy.setPropWidth(width)
	dpy.setPropHeight(height)
}
func (dpy *Display) updateRotationAndRelfect(randr uint16) {
	rotation, reflect := parseRandR(randr)

	dpy.setPropRotation(rotation)
	dpy.setPropReflect(reflect)
}

func main() {
	dpy := NewDisplay()
	dbus.InstallOnSession(dpy)
	dbus.DealWithUnhandledMessage()

	select {}
}
