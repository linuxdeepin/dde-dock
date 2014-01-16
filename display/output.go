package main

import "github.com/BurntSushi/xgb/xproto"
import "github.com/BurntSushi/xgb/randr"
import "dlib/logger"
import "dlib/dbus"
import "fmt"

type Mode struct {
	ID     uint32
	Width  uint16
	Height uint16
	Rate   float64
}
type Output struct {
	bestMode  randr.Mode
	modes     []Mode
	rotations uint16
	crtc      randr.Crtc

	Identify randr.Output
	Name     string
	Type     uint8

	Mode         Mode
	Allocation   xproto.Rectangle
	AdjustMethod uint8

	Rotation   uint16  `access:"readwrite"`
	Reflect    uint16  `access:"readwrite"`
	Opened     bool    `access:"readwrite"`
	Brightness float64 `access:"readwrite"`
}

type Border struct {
	Left uint16
	Top uint16
	Right uint16
	Bottom uint16
}

func (op *Output) ListModes() []Mode {
	return op.modes
}
func (op *Output) ListRotations() []uint16 {
	return parseRotations(op.rotations)
}
func (op *Output) ListReflect() []uint16 {
	return parseReflects(op.rotations)
}

func (op *Output) SetMode(id uint32) {
	_, err := randr.SetCrtcConfig(X, op.crtc, xproto.TimeCurrentTime, LastConfigTimeStamp,
		op.Allocation.X, op.Allocation.Y, randr.Mode(id), op.Rotation|op.Reflect, []randr.Output{op.Identify}).Reply()
	if err != nil {
		logger.Println(fmt.Sprintln("SetCrtcConfig failed:", err, op.crtc))
	}
}

func (op *Output) SetAllocation(x, y, width, height, adjMethod int16) {
	//TODO: handle adjMethod with `width` and `height`

	cinfo, err := randr.GetCrtcInfo(X, op.crtc, LastConfigTimeStamp).Reply()
	found := false
	for _, po := range cinfo.Possible {
		if po == op.Identify {
			found = true
			break
		}
	}
	if !found {
		panic(fmt.Sprintf("%s Crtc %d is only support %v, but op which will be SetAllocation is %d", op.Name, op.crtc, cinfo.Possible, op.Identify))
	}

	if uint32(cinfo.Mode) != op.Mode.ID {
		panic(fmt.Sprintf("%s SetAllocation check failed at mode : %v != %v(op)", op.Name, cinfo.Mode, op.Mode))
	}
	if cinfo.Rotation != op.Reflect|op.Rotation {
		panic(fmt.Sprintf("%s SetAllocation check failed at rotation: %v != %v(op)", op.Name, cinfo.Rotation, op.Reflect|op.Rotation))
	}

	_, err = randr.SetCrtcConfig(X, op.crtc, xproto.TimeCurrentTime, LastConfigTimeStamp, x, y, cinfo.Mode, cinfo.Rotation, cinfo.Outputs).Reply()
	if err != nil {
		panic(fmt.Sprintf("%s SetAllocation(%d,%d,%d,%d) screenSize(%d,%d) failed at SetCrtcConfig:%s",
			op.Name, x, y, width, height, DefaultScreen.WidthInPixels, DefaultScreen.HeightInPixels, err.Error()))
	}
}

func (op *Output) Debug() string {
	return fmt.Sprintf("Allocation:%v Rotation:%d", op.Allocation, op.Rotation)
}

//-------------- internal output methods---------------------

func (op *Output) setBrightness(brightness float64) {
	if brightness < 0.01 || brightness > 1 {
		brightness = 1
	}
	gammaSize, err := randr.GetCrtcGammaSize(X, op.crtc).Reply()
	if err != nil {
		panic(fmt.Sprintf("GetCrtcGrammSize(crtc:%d) failed: %s", op.crtc, err.Error()))
	}
	red, green, blue := genGammaRamp(gammaSize.Size, brightness)
	randr.SetCrtcGamma(X, op.crtc, gammaSize.Size, red, green, blue)
	op.setPropBrightness(brightness)
}

func (op *Output) setRotation(rotation uint16) {
	/*v := op.Opened*/
	/*defer func() { op.setOpened(v) }()*/
	/*op.setOpened(false)*/

	_, err := randr.SetCrtcConfig(X, op.crtc, xproto.TimeCurrentTime, LastConfigTimeStamp,
		op.Allocation.X, op.Allocation.Y, op.bestMode, rotation|op.Reflect, []randr.Output{op.Identify}).Reply()
	if err != nil {
		panic(fmt.Sprintln("SetRotation:", rotation, rotation|op.Reflect, err))
	}
	op.setPropRotation(rotation)

	{
		info, err := randr.GetCrtcInfo(X, op.crtc, 0).Reply()
		if err != nil {
			panic(err)
		}
		fmt.Println("SetRotation:", info.X, info.Y, info.Width, info.Height, info.Rotation)
	}
}

func (op *Output) setReflect(reflect uint16) {
	switch reflect {
	case 0, 16, 32, 48:
		break
	default:
		panic(fmt.Sprintf("setReflect Value%d Error", reflect))
	}

	v := op.Opened
	defer func() { op.setOpened(v) }()
	op.setOpened(false)

	_, err := randr.SetCrtcConfig(X, op.crtc, xproto.TimeCurrentTime, LastConfigTimeStamp,
		op.Allocation.X, op.Allocation.Y, op.bestMode, op.Rotation|reflect, []randr.Output{op.Identify}).Reply()
	op.setOpened(true)
	if err != nil {
		panic(fmt.Sprintln("SetReflect:", op.Rotation|reflect, err))
	}
	op.setPropReflect(reflect)
}

func (op *Output) setOpened(v bool) {
	fmt.Println("SetOpened.... ", v)
	//op.Opened will be changed when we receive appropriate event
	if v == true {
		oinfo, err := randr.GetOutputInfo(X, op.Identify, LastConfigTimeStamp).Reply()
		if err != nil {
			panic(err)
		}
		for _, crtc := range oinfo.Crtcs {
			if isCrtcConnected(X, crtc) {
				fmt.Println("crtc:", crtc, "Is used for ", op.Name)
				continue
			}
			s, err := randr.SetCrtcConfig(X, crtc, xproto.TimeCurrentTime, LastConfigTimeStamp,
				op.Allocation.X, op.Allocation.Y, op.bestMode, op.Rotation, []randr.Output{op.Identify}).Reply()
			if err == nil {
				fmt.Println("Crtc:", crtc, "for", op.Name, " is ok")
				break
			}
			fmt.Println("AAAA:", s, err, crtc, op.bestMode, op.Rotation, op.Identify)
		}
	} else {
		_, err := randr.SetCrtcConfig(X, op.crtc, xproto.TimeCurrentTime, LastConfigTimeStamp,
			op.Allocation.X, op.Allocation.Y, 0, op.Rotation, nil).Reply()
		if err != nil {
			panic(err)
		}
	}
}

func (op *Output) updateCrtc(dpy *Display) {
	if op.crtc != 0 {
		info, err := randr.GetCrtcInfo(X, op.crtc, LastConfigTimeStamp).Reply()
		if err != nil {
			panic("Opps:" + err.Error())
		}
		op.rotations = info.Rotations

		rotation, reflect := parseRandR(info.Rotation)
		op.setPropRotation(rotation)
		op.setPropReflect(reflect)

		op.setPropAllocation(xproto.Rectangle{info.X, info.Y, info.Width, info.Height})

		op.setPropMode(buildMode(dpy.modes[info.Mode]))
	} else {
		op.setPropRotation(1)
		op.setPropReflect(0)
		op.setPropAllocation(xproto.Rectangle{0, 0, 0, 0})
		op.setPropMode(Mode{0, 0, 0, 0})
	}
}

func (op *Output) update(dpy *Display, info *randr.GetOutputInfoReply) {
	op.crtc = info.Crtc
	op.setPropOpened(info.Crtc != 0)
	op.bestMode = info.Modes[0]

	op.modes = nil
	for _, m := range info.Modes {
		info := dpy.modes[m]
		op.modes = append(op.modes, buildMode(info))
	}

	if op.crtc != 0 {
		cinfo, err := randr.GetCrtcInfo(X, op.crtc, LastConfigTimeStamp).Reply()
		op.setPropMode(buildMode(dpy.modes[cinfo.Mode]))
		if err != nil {
			panic(fmt.Sprintf("Op.crtc(%d) != 0 && can't GetCrtcInfo (%s)", op.crtc, err.Error()))
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

	op := &Output{
		Identify:   core,
		Name:       getOutputName(core, string(info.Name)),
		Brightness: 1, //TODO: init this value
	}
	op.update(dpy, info)
	op.updateCrtc(dpy)
	dbus.InstallOnSession(op)
	return op
}
