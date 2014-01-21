package main

import "github.com/BurntSushi/xgb/xproto"
import "github.com/BurntSushi/xgb/randr"
import "dlib/dbus"
import "fmt"

type Mode struct {
	ID     uint32
	Width  uint16
	Height uint16
	Rate   float64
}
type Output struct {
	bestMode      randr.Mode
	modes         []Mode
	rotations     uint16
	crtc          randr.Crtc
	pendingConfig *pendingConfig

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
	Left   uint16
	Top    uint16
	Right  uint16
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

func (op *Output) SetPos(x, y int16) {
	op.pendingConfig = NewPendingConfig(op).SetPos(x, y)
}

func (op *Output) EnsureSize(width, height uint16, hint uint8) {
	switch hint {
	case EnsureSizeHintAuto, EnsureSizeHintPanning, EnsureSizeHintBorderScale:
		op.pendingConfig = NewPendingConfig(op).EnsureSize(width, height, hint)
	}
}

func (op *Output) SetMode(id uint32) {
	op.pendingConfig = NewPendingConfig(op).SetMode(randr.Mode(id))
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
	NewPendingConfig(op).SetGamma(red, green, blue)
}

func (op *Output) setRotation(rotation uint16) {
	op.pendingConfig = NewPendingConfig(op).SetRotation(rotation | op.Reflect)
}

func (op *Output) setReflect(reflect uint16) {
	switch reflect {
	case 0, 16, 32, 48:
		break
	default:
		panic(fmt.Sprintf("setReflect Value%d Error", reflect))
	}

	NewPendingConfig(op).SetRotation(op.Rotation | reflect)
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
	} else if op.crtc != 0 {
		_, err := randr.SetCrtcConfig(X, op.crtc, xproto.TimeCurrentTime, LastConfigTimeStamp,
			op.Allocation.X, op.Allocation.Y, 0, op.Rotation, nil).Reply()
		if err != nil {
			panic(fmt.Sprintln("Close Output failed when SetCrtcConfig", op.crtc, err))
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
		op.setPropMode(buildMode(dpy.modes[info.Mode]))
		op.setPropAllocation(NewPendingConfig(op).appliedAllocation())

		op.setPropAllocation(NewPendingConfig(op).appliedAllocation())
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
