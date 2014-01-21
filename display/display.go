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
	_             = fmt.Println
	DPY           *Display
	X, _          = xgb.NewConn()
	DefaultScreen = xproto.Setup(X).DefaultScreen(X)
	Root          = DefaultScreen.Root

	LastConfigTimeStamp = xproto.Timestamp(0)

	MinWidth, MinHeight, MaxWidth, MaxHeight uint16
)

func init() {
	randr.Init(X)
	ver, err := randr.QueryVersion(X, 1, 13).Reply()
	if err != nil {
		panic(fmt.Sprintln("randr.QueryVersion error:", err))
	}
	if ver.MajorVersion < 1 && ver.MinorVersion < 13 {
		panic(fmt.Sprintln("randr version is too low:", ver.MajorVersion, ver.MinorVersion, "this program require at least randr 1.3"))
	}

	rng, err := randr.GetScreenSizeRange(X, Root).Reply()
	MinWidth, MinHeight, MaxWidth, MaxHeight = rng.MinWidth, rng.MinHeight, rng.MaxWidth, rng.MaxHeight
	if err != nil {
		panic(fmt.Sprintln("randr.GetSceenSizeRange failed :", err))
	}

	initDisplay()
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

	MirrorMode   bool
	MirrorOutput *Output `access:readwrite`
}

func initDisplay() *Display {
	dpy := &Display{}
	DPY = dpy

	dpy.modes = make(map[randr.Mode]randr.ModeInfo)

	sinfo, err := getScreenInfo()
	dpy.setPropRotation(uint16(sinfo.Rotations))
	dpy.updateRotationAndRelfect(sinfo.Rotation)

	if err != nil {
		panic("GetScreenInfo Failed:" + err.Error())
	}
	dpy.updateResources()
	dpy.setPropWidth(DefaultScreen.WidthInPixels)
	dpy.setPropHeight(DefaultScreen.HeightInPixels)

	dpy.updatePrimary()

	randr.SelectInput(X, Root, randr.NotifyMaskOutputChange|randr.NotifyMaskCrtcChange|randr.NotifyMaskScreenChange)
	go dpy.listener()
	return dpy
}

func (dpy *Display) updateResources() {

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
		switch ee := e.(type) {
		case randr.NotifyEvent:
			switch ee.SubCode {
			case randr.NotifyCrtcChange:
				info := ee.U.Cc
				if op := queryOutputByCrtc(dpy, info.Crtc); op != nil {
					op.updateCrtc(dpy)
				}
				w, h := parseScreenSize(dpy.Outputs)
				fmt.Println("NotifyCrtcChange....:", w, h)

			case randr.NotifyOutputChange:
				info := ee.U.Oc
				switch info.Connection {
				case randr.ConnectionConnected:
					dpy.updateOutputList(info.Output)
				case randr.ConnectionDisconnected, randr.ConnectionUnknown:
					dpy.removeOutput(info.Output)
				}
				dpy.updatePrimary()
			}
		case randr.ScreenChangeNotifyEvent:
			if LastConfigTimeStamp < ee.ConfigTimestamp {
				dpy.updateResources()
				LastConfigTimeStamp = ee.ConfigTimestamp
				//TODO: monitor changed.
				dpy.setPropWidth(DefaultScreen.WidthInPixels)
				dpy.setPropHeight(DefaultScreen.HeightInPixels)
				dpy.updateRotationAndRelfect(uint16(ee.Rotation))
				dpy.setPropWidth(DefaultScreen.WidthInPixels)
				dpy.setPropHeight(DefaultScreen.HeightInPixels)
				dpy.updatePrimary()
			}
		}
	}
}
func (dpy *Display) updateRotationAndRelfect(randr uint16) {
	rotation, reflect := parseRandR(randr)

	dpy.setPropRotation(rotation)
	dpy.setPropReflect(reflect)
}

func (dpy *Display) setScreenSize(width uint16, height uint16) {
	if width < MinWidth || width > MaxWidth || height < MinHeight || height > MaxWidth {
		logger.Println("updateScreenSize with invalid value:", width, height)
		return
	}

	err := randr.SetScreenSizeChecked(X, Root, width, height, uint32(DefaultScreen.WidthInMillimeters), uint32(DefaultScreen.HeightInMillimeters)).Check()

	if err != nil {
		logger.Println("randr.SetScreenSize to :", width, height, DefaultScreen.WidthInPixels, DefaultScreen.HeightInPixels, err)
		/*panic(fmt.Sprintln("randr.SetScreenSize to :", width, height, err))*/
	}
}

func main() {
	dbus.InstallOnSession(DPY)
	dbus.DealWithUnhandledMessage()
	select {}
}
