package main

import "github.com/BurntSushi/xgb/xproto"
import "math"
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
	savedConfig   *pendingConfig

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

func (op *Output) EnsureSize2(width, height uint16, hint uint8) {
	op.EnsureSize(width, height, hint)
	DPY.ApplyChanged()
}
func (op *Output) EnsureSize(width, height uint16, hint uint8) {
	if !op.Opened {
		return
	}
	matchedRate := op.Mode.Rate
	matchedMode := randr.Mode(op.Mode.ID)
	acc := float64(width + height)
	for _, minfo := range op.ListModes() {
		newAcc := math.Abs(float64(int16(minfo.Width)-int16(width))) + math.Abs(float64(int16(minfo.Height)-int16(height)))
		if newAcc < acc {
			matchedMode = randr.Mode(minfo.ID)
			matchedRate = minfo.Rate
			acc = newAcc
		} else if newAcc == acc {
			if matchedRate < minfo.Rate {
				matchedMode = randr.Mode(minfo.ID)
				matchedRate = minfo.Rate
				acc = newAcc
			}
		}
	}
	fmt.Println("OPEnsureSize:", matchedMode, buildMode(DPY.modes[matchedMode]), width, height)
	op.pendingConfig = NewPendingConfig(op).SetMode(matchedMode).EnsureSize(width, height, hint)
}

func (op *Output) SetMode(id uint32) {
	op.pendingConfig = NewPendingConfig(op).SetMode(randr.Mode(id))
	DPY.ApplyChanged()
}

func (op *Output) Debug() string {
	return fmt.Sprintf("Allocation:%v Rotation:%d", op.Allocation, op.Rotation)
}

//-------------- internal output methods---------------------

func (op *Output) setRotation(rotation uint16) {
	op.pendingConfig = NewPendingConfig(op).SetRotation(rotation | op.Reflect)
}
func (op *Output) setBrightness(brightness float64) {
	op.pendingConfig = NewPendingConfig(op).SetBrightness(brightness)
}

func (op *Output) setReflect(reflect uint16) {
	switch reflect {
	case 0, 16, 32, 48:
		break
	default:
		panic(fmt.Sprintf("setReflect Value%d Error", reflect))
	}

	op.pendingConfig = NewPendingConfig(op).SetRotation(op.Rotation | reflect)
}

func (op *Output) update() {
	info, err := randr.GetOutputInfo(X, op.Identify, xproto.TimeCurrentTime).Reply()
	if err != nil {
		fmt.Println("Output.update failed at GetOutputInfo", err)
	}
	op.crtc = info.Crtc
	op.setPropOpened(info.Crtc != 0)
	op.bestMode = info.Modes[0]

	op.modes = nil
	for _, m := range info.Modes {
		info := DPY.modes[m]
		op.modes = append(op.modes, buildMode(info))
	}

	if op.crtc != 0 {
		cinfo, err := randr.GetCrtcInfo(X, op.crtc, LastConfigTimeStamp).Reply()
		op.setPropMode(buildMode(DPY.modes[cinfo.Mode]))
		if err != nil {
			panic(fmt.Sprintf("Op.crtc(%d) != 0 && can't GetCrtcInfo (%s)", op.crtc, err.Error()))
		}
	}

	if op.crtc != 0 {
		info, err := randr.GetCrtcInfo(X, op.crtc, LastConfigTimeStamp).Reply()
		if err != nil {
			panic("Opps:" + err.Error())
		}
		op.rotations = info.Rotations

		rotation, reflect := parseRandR(info.Rotation)
		op.setPropRotation(rotation)
		op.setPropReflect(reflect)

		op.setPropMode(buildMode(DPY.modes[info.Mode]))
		op.setPropAllocation(op.pendingAllocation())
	} else {
		op.setPropRotation(1)
		op.setPropReflect(0)
		op.setPropAllocation(xproto.Rectangle{0, 0, 0, 0})
		op.setPropMode(Mode{0, 0, 0, 0})
	}
}
func (op *Output) tryOpen() {
	oinfo, _ := randr.GetOutputInfo(X, op.Identify, LastConfigTimeStamp).Reply()
	for _, crtc := range oinfo.Crtcs {
		if isCrtcConnected(X, crtc) == false {
			_, err := randr.SetCrtcConfig(X, crtc, LastConfigTimeStamp, xproto.TimeCurrentTime, 0, 0, oinfo.Modes[0], randr.RotationRotate0, []randr.Output{op.Identify}).Reply()
			if err != nil {
				fmt.Println("TryOpenOutput", op.Identify, "Failed", err)
			}
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
	if info.Crtc == 0 {
		//workaround: nvidia driver has an VGA-1 output which connection status equal connected but hasn't corresponding crtc, it will panic Xserver when trying to connect the output.
		return nil
	}

	op := &Output{
		Identify:   core,
		Name:       getOutputName(core, string(info.Name)),
		Brightness: 1, //TODO: init this value
	}
	if info.Crtc == 0 {
		op.tryOpen()
	}
	op.update()
	dbus.InstallOnSession(op)
	return op
}
