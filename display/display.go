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

const (
	DisplayModeOnlyPrimary uint8 = iota
	DisplayModeOnlySecondary
	DisplayModeMirrors
	DisplayModeExtend
)

func init() {
	randr.Init(X)
	ver, err := randr.QueryVersion(X, 1, 4).Reply()
	if err != nil {
		panic(fmt.Sprintln("randr.QueryVersion error:", err))
	}
	if ver.MajorVersion < 1 && ver.MinorVersion < 4 {
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
	DisplayMode  uint8
	MirrorOutput *Output `access:readwrite`

	listening bool
}

func initDisplay() *Display {
	dpy := &Display{}
	dbus.InstallOnSession(dpy)
	DPY = dpy

	dpy.modes = make(map[randr.Mode]randr.ModeInfo)
	DPY.MirrorMode = true

	sinfo, err := getScreenInfo()
	dpy.setPropRotation(uint16(sinfo.Rotations))
	dpy.updateRotationAndRelfect(sinfo.Rotation)

	if err != nil {
		panic("GetScreenInfo Failed:" + err.Error())
	}
	dpy.updateResources()
	dpy.updateScreenSize(xproto.Setup(X).DefaultScreen(X).WidthInPixels, xproto.Setup(X).DefaultScreen(X).HeightInPixels)

	dpy.setPropPrimaryRect(xproto.Rectangle{0, 0, dpy.Width, dpy.Height})

	randr.SelectInput(X, Root, randr.NotifyMaskOutputChange|randr.NotifyMaskOutputProperty|randr.NotifyMaskCrtcChange|randr.NotifyMaskScreenChange)
	dpy.startListen()
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

func (dpy *Display) updateOutputList(output randr.Output) {
	op := queryOutput(dpy, output)
	if op == nil {
		if op = NewOutput(dpy, output); op != nil {
			dpy.setPropOutputs(append(dpy.Outputs, op))
		}
	} else {
		op.update()
	}
}
func (dpy *Display) removeOutput(output randr.Output) {
	var newOutputs []*Output
	var removingOp *Output
	for _, op := range dpy.Outputs {
		if op.Identify != output {
			newOutputs = append(newOutputs, op)
		} else {
			removingOp = op
		}
	}
	if removingOp != nil {
		dpy.setPropOutputs(newOutputs)

		dbus.UnInstallObject(removingOp)

		if dpy.MirrorOutput.Identify == output {
			dpy.SetMirrorOutput(deduceMirrorOutput(dpy.Outputs))
		}
		randr.SetCrtcConfig(X, removingOp.crtc, LastConfigTimeStamp, xproto.TimeCurrentTime, 0, 0, 0, randr.RotationRotate0, nil).Reply()
	}
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
				info := ee.U.Cc
				if op := queryOutputByCrtc(dpy, info.Crtc); op != nil {
					op.update()
					w, h := parseScreenSize(dpy.Outputs)
					fmt.Println("NotifyCrtcChange....:", op.Name, w, h)

					/*dpy.SetMirrorOutput(deduceMirrorOutput(dpy.Outputs))*/
				}

			case randr.NotifyOutputChange:
				info := ee.U.Oc
				switch info.Connection {
				case randr.ConnectionConnected:
					dpy.updateOutputList(info.Output)

					dpy.updatePrimaryOutput()
					dpy.SetMirrorOutput(deduceMirrorOutput(dpy.Outputs))

					fmt.Println("OutputChanged....", info.Output, info.Mode, queryOutput(DPY, info.Output))
				case randr.ConnectionDisconnected, randr.ConnectionUnknown:
					dpy.removeOutput(info.Output)

					fmt.Println("OutputChanged lost....", info.Output)
				}
			case randr.NotifyOutputProperty:
				info := ee.U.Op
				//some driver use "BACKLIGHT" instead "Backlight" value as info.Atom, so we didn't check it.
				if support, value := supportedBacklight(X, info.Output); support {
					if op := queryOutput(dpy, info.Output); op != nil {
						op.setPropBrightness(value)
					}
				}
			}
		case randr.ScreenChangeNotifyEvent:
			if LastConfigTimeStamp < ee.ConfigTimestamp {
				dpy.updateResources()
				LastConfigTimeStamp = ee.ConfigTimestamp
				//TODO: monitor changed.
				dpy.updateRotationAndRelfect(uint16(ee.Rotation))
				if dpy.MirrorMode {
					if ee.Width != dpy.MirrorOutput.pendingAllocation().Width || ee.Height != dpy.MirrorOutput.pendingAllocation().Width {
						dpy.ApplyChanged()
					}
				} else {
					dpy.updateScreenSize(ee.Width, ee.Height)
				}

			}
		}
	}
}

func (dpy *Display) updateRotationAndRelfect(randr uint16) {
	rotation, reflect := parseRandR(randr)

	dpy.setPropRotation(rotation)
	dpy.setPropReflect(reflect)
}

func (dpy *Display) updateScreenSize(width uint16, height uint16) {
	dpy.setPropWidth(width)
	dpy.setPropHeight(height)
	fmt.Println("UpdateScreenSize:", width, height, dpy.PrimaryOutput)
	dpy.updatePrimaryOutput()
}

func (dpy *Display) SetPrimaryOutput(output uint32) {
	op := queryOutput(dpy, randr.Output(output))
	dpy.setPropPrimaryOutput(op)
	if op != nil && op.Opened {
		dpy.setPropPrimaryRect(op.pendingAllocation())
	} else {
		dpy.setPropPrimaryRect(xproto.Rectangle{0, 0, dpy.Width, dpy.Height})
	}
}

func (dpy *Display) updatePrimaryOutput() {
	pinfo, err := randr.GetOutputPrimary(X, Root).Reply()
	if err == nil {
		dpy.SetPrimaryOutput(uint32(pinfo.Output))
	} else {
		dpy.SetPrimaryOutput(0)
	}
}
func (dpy *Display) setScreenSize(width uint16, height uint16) {
	if width < MinWidth || width > MaxWidth || height < MinHeight || height > MaxHeight {
		logger.Println("updateScreenSize with invalid value:", width, height)
		return
	}

	if (width != DPY.Width) || (height != DPY.Height) {
		fmt.Println("SetScreenSize...................................................", width, height, DPY.Width, DPY.Height)
		err := randr.SetScreenSizeChecked(X, Root, width, height, uint32(ScreenWidthMm), uint32(ScreenHeightMm)).Check()

		if err != nil {
			logger.Println("randr.SetScreenSize to :", width, height, DPY.Width, DPY.Height, err)
			/*panic(fmt.Sprintln("randr.SetScreenSize to :", width, height, err))*/
		} else {
			dpy.updateScreenSize(width, height)
		}
	}
}

func TT() {
	fmt.Println("sssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss")
	for _, op := range DPY.Outputs {
		if op.Name == "LVDS1" && op.Opened {
			op.EnsureSize(1153, 864, EnsureSizeHintAuto)
			op.EnsureSize(1024, 768, EnsureSizeHintAuto)
			/*op.EnsureSize(1280, 800, EnsureSizeHintAuto)*/
			/*op.EnsureSize(800, 600, EnsureSizeHintAuto)*/
			c := op.pendingConfig
			DPY.ApplyChanged()
			fmt.Println("OPSSSSSS>>>", &(op.pendingConfig), op.Name, c, fixed2double(c.transform.Matrix22))
		}
	}
}

func main() {
	dbus.DealWithUnhandledMessage()
	DPY.SetMirrorMode(true)
	/*TT()*/
	select {}
}
