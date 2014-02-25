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
	_              = fmt.Println
	DPY            *Display
	X, _           = xgb.NewConn()
	Root           = xproto.Setup(X).DefaultScreen(X).Root
	ScreenWidthMm  = xproto.Setup(X).DefaultScreen(X).WidthInMillimeters
	ScreenHeightMm = xproto.Setup(X).DefaultScreen(X).HeightInMillimeters

	LastConfigTimeStamp = xproto.Timestamp(0)

	MinWidth, MinHeight, MaxWidth, MaxHeight uint16
)

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

	DisplayMode   int16
	BuiltinOutput *Output

	listening     bool
	configuration DisplayConfiguration
}

func initDisplay() *Display {
	dpy := &Display{}
	dbus.InstallOnSession(dpy)
	DPY = dpy

	dpy.DisplayMode = DisplayModeUnknow
	dpy.update()

	dpy.configuration = LoadDisplayConfiguration(dpy)
	dpy.SetDisplayMode(dpy.configuration.DisplayMode)

	dpy.setPropBuiltinOutput(guestBuiltIn(dpy.Outputs))

	randr.SelectInput(X, Root, randr.NotifyMaskOutputChange|randr.NotifyMaskOutputProperty|randr.NotifyMaskCrtcChange|randr.NotifyMaskScreenChange)
	dpy.startListen()

	return dpy
}

func (dpy *Display) update() {
	sinfo, err := getScreenInfo()
	if err != nil {
		fmt.Println("GetScreenInfo Failed:" + err.Error())
		return
	}

	{
		dpy.setPropRotation(uint16(sinfo.Rotations))
		rotation, reflect := parseRandR(sinfo.Rotation)
		dpy.setPropRotation(rotation)
		dpy.setPropReflect(reflect)
	}

	{
		sizeinfo := xproto.Setup(X).DefaultScreen(X)
		dpy.setPropWidth(sizeinfo.WidthInPixels)
		dpy.setPropHeight(sizeinfo.HeightInPixels)
	}

	{
		// update output list
		resources, err := randr.GetScreenResources(X, Root).Reply()
		LastConfigTimeStamp = resources.ConfigTimestamp

		if err != nil {
			panic("GetScreenResources failed:" + err.Error())
		}

		dpy.modes = make(map[randr.Mode]randr.ModeInfo)
		for _, m := range resources.Modes {
			dpy.modes[randr.Mode(m.Id)] = m
		}

		// clean old outputs
		for _, op := range dpy.Outputs {
			op.destroy()
		}
		// build new outputs
		ops := make([]*Output, 0)
		for _, output := range resources.Outputs {
			op := NewOutput(dpy, output)
			if op != nil {
				ops = append(ops, op)
			}
		}
		dpy.setPropOutputs(ops)
	}

	dpy.adjustScreenSize()

	{
		// update primary output and primary rectangle
		pinfo, err := randr.GetOutputPrimary(X, Root).Reply()
		var op *Output
		if err == nil {
			op = queryOutput(dpy, randr.Output(pinfo.Output))
		}
		if op != nil {
			dpy.setPropPrimaryOutput(op)
			dpy.setPropPrimaryRect(op.pendingAllocation())
		} else {
			dpy.setPropPrimaryOutput(nil)
			dpy.setPropPrimaryRect(xproto.Rectangle{0, 0, dpy.Width, dpy.Height})
		}
		fmt.Println("PrimaryOutput:", op, dpy.PrimaryRect)
	}
	{
	}

	if dpy.DisplayMode == DisplayModeCustom {
		dpy.configuration = GenerateDefaultConfig(dpy)
		dpy.configuration.save()
	}
}

func (dpy *Display) Reset() {
	for _, op := range dpy.Outputs {
		op.setBrightness(1)
	}
	dpy.ApplyChanged()
}

func (dpy *Display) ShowInfoOnScreen() {
}
func (dpy *Display) ListRotations() []uint16 {
	return parseRotations(dpy.rotations)
}
func (dpy *Display) ListReflect() []uint16 {
	return parseReflects(dpy.rotations)
}

func (dpy *Display) stopListen() {
	dpy.listening = false
}
func (dpy *Display) startListen() {
	dpy.listening = true
	go dpy.listener()
}
func (dpy *Display) listener() {
	for {
		if !dpy.listening {
			return
		}
		e, err := X.WaitForEvent()
		if err != nil {
			continue
		}
		switch ee := e.(type) {
		case randr.NotifyEvent:
			switch ee.SubCode {
			case randr.NotifyCrtcChange:
			case randr.NotifyOutputChange:
			case randr.NotifyOutputProperty:
				fmt.Println("OutputPropertyChange...")
				info := ee.U.Op
				//some driver use "BACKLIGHT" instead "Backlight" value as info.Atom, so we didn't check it.
				if support, value := supportedBacklight(X, info.Output); support {
					if op := queryOutput(dpy, info.Output); op != nil {
						fmt.Println("TrySetBrightness:", value)
						op.setPropBrightness(value)
					}
				}
			}
		case randr.ScreenChangeNotifyEvent:
			dpy.update()
			if ee.ConfigTimestamp > LastConfigTimeStamp {
				//output connection changed
				LastConfigTimeStamp = ee.ConfigTimestamp
				dpy.setPropBuiltinOutput(guestBuiltIn(dpy.Outputs))
			}
			fmt.Println("BeginUpdate...")
			fmt.Println("endUpdate...")
		}
	}
}

func (dpy *Display) setScreenSize(width uint16, height uint16) {
	if width < MinWidth || width > MaxWidth || height < MinHeight || height > MaxHeight {
		logger.Println("update1ScreenSize with invalid value:", width, height)
		return
	}

	if (width != dpy.Width) || (height != dpy.Height) {
		err := randr.SetScreenSizeChecked(X, Root, width, height, uint32(ScreenWidthMm), uint32(ScreenHeightMm)).Check()

		if err != nil {
			(fmt.Println("randr.SetScreenSize to :", width, height, dpy.Width, dpy.Height, err))
			/*panic(fmt.Sprintln("randr.SetScreenSize to :", width, height, dpy.Width, dpy.Height, err))*/
		} else {
			dpy.setPropWidth(width)
			dpy.setPropHeight(height)
		}
	}
}

func main() {
	randr.Init(X)
	ver, err := randr.QueryVersion(X, 1, 3).Reply()
	fmt.Println("VER:", ver)
	if err != nil {
		panic(fmt.Sprintln("randr.QueryVersion error:", err))
	}
	if ver.MajorVersion != 1 || ver.MinorVersion != 3 {
		panic(fmt.Sprintln("randr version is too low:", ver.MajorVersion, ver.MinorVersion, "this program require at least randr 1.3"))
	}

	rng, err := randr.GetScreenSizeRange(X, Root).Reply()
	MinWidth, MinHeight, MaxWidth, MaxHeight = rng.MinWidth, rng.MinHeight, rng.MaxWidth, rng.MaxHeight

	if err != nil {
		panic(fmt.Sprintln("randr.GetSceenSizeRange failed :", err))
	}

	initDisplay()

	dbus.DealWithUnhandledMessage()
	select {}
}
