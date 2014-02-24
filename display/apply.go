package main

import "github.com/BurntSushi/xgb/randr"
import "github.com/BurntSushi/xgb/xproto"
import "math"
import "sync"
import "fmt"

import "github.com/BurntSushi/xgb/render"

const (
	_PendingMaskMode = 1 << iota
	_PendingMaskPos
	_PendingMaskTransform
	_PendingMaskRotation
	_PendingMaskGramma
	_PendingMaskBacklight
)

const (
	EnsureSizeHintAuto uint8 = iota
	EnsureSizeHintPanning
	EnsureSizeHintScale
)

func min(a, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}
func max(a, b uint16) uint16 {
	if a > b {
		return a
	}
	return b
}

var changeLock, changeUnlock = func() (func(), func()) {
	var locker sync.Mutex
	return func() {
			locker.Lock()
			xproto.GrabServer(X)
			fmt.Println("-------------------------GRABSERVER-----------------------")
		}, func() {
			xproto.UngrabServer(X)
			locker.Unlock()
			fmt.Println("-------------------------UNGRABSERVER-----------------------")
		}
}()

type pendingConfig struct {
	crtc             randr.Crtc
	output           randr.Output
	mask             int
	supportBacklight bool

	mode     randr.Mode
	posX     int16
	posY     int16
	rotation uint16

	transform    render.Transform
	filterName   string
	filterParams []render.Fixed

	borderCompensationX uint16
	borderCompensationY uint16

	// doesn't influence allocation
	gammaRed   []uint16
	gammaGreen []uint16
	gammaBlue  []uint16

	backlight float64
}

func ClonePendingConfig(c *pendingConfig) *pendingConfig {
	r := &pendingConfig{}
	r.crtc = c.crtc
	r.output = c.output
	r.mask = c.mask
	r.mode = c.mode
	r.rotation = c.rotation
	r.posX, r.posY = c.posX, c.posY
	r.transform, r.filterName, r.filterParams = c.transform, c.filterName, c.filterParams
	r.borderCompensationX, r.borderCompensationY = c.borderCompensationX, c.borderCompensationY
	r.supportBacklight, r.backlight = c.supportBacklight, c.backlight
	copy(r.gammaRed, c.gammaRed)
	copy(r.gammaGreen, c.gammaGreen)
	copy(r.gammaBlue, c.gammaBlue)

	return r
}

func (c *pendingConfig) String() string {
	return fmt.Sprintf(`
	Crtc:%v, Output:%v, mode:%v, pos:(%v,%v), rotation:%v transform:(%v,%v,%v), borderComp:(%v,%v), backlight:%v
	`, c.crtc, c.output, c.mode, c.posX, c.posY, c.rotation, c.transform.Matrix11, c.transform.Matrix22, c.transform.Matrix33,
		c.borderCompensationX, c.borderCompensationY, c.backlight)
}
func NewPendingConfigWithoutCache(op *Output) *pendingConfig {
	//don't call any pendingConfig.SetXX  otherwise ApplyChanged will apply this changed.
	c := &pendingConfig{}

	if op.crtc != 0 {
		c.crtc = op.crtc
	} else if op.savedConfig != nil {
		c.crtc = op.savedConfig.crtc
	}
	if op.Mode.ID != 0 {
		c.mode = randr.Mode(op.Mode.ID)
	} else if op.savedConfig != nil {
		c.mode = op.savedConfig.mode
	}
	c.output = op.Identify

	validMode := false
	for _, m := range op.ListModes() {
		if c.mode == randr.Mode(m.ID) {
			validMode = true
			break
		}
	}
	if !validMode {
		c.mode = op.bestMode
	}

	cinfo, err := randr.GetCrtcInfo(X, c.crtc, LastConfigTimeStamp).Reply()
	if err != nil {
		panic(fmt.Sprintf("NewPendingCofing failed at GetCrtcInfo(crtc=%v,op=%v,opened=%v):%v", c.crtc, c.output, op.Opened, err))
	}
	if cinfo.Rotation&0xf == 0 {
		panic("Rotation err")
	}
	c.rotation = cinfo.Rotation

	tinfo, err := randr.GetCrtcTransform(X, c.crtc).Reply()
	if err != nil {
		panic(fmt.Sprintf("NewPendingCofing failed at GetCrtcTransform(crtc=%v,op=%v,opened=%v):%v", c.crtc, c.output, op.Opened, err))
	}
	c.transform = tinfo.CurrentTransform
	c.filterName = tinfo.CurrentFilterName
	c.filterParams = tinfo.CurrentParams

	if DPY.DisplayMode == DisplayModeMirrors {
		c.posX, c.posY = 0, 0
		if cinfo.X < 0 {
			c.borderCompensationX = uint16(-cinfo.X)
		}
		if cinfo.X < 0 {
			c.borderCompensationY = uint16(-cinfo.Y)
		}
	} else {
		c.posX, c.posY = cinfo.X, cinfo.Y
		c.borderCompensationX, c.borderCompensationY = 0, 0
	}

	c.supportBacklight, c.backlight = supportedBacklight(X, c.output)
	return c
}

func NewPendingConfig(op *Output) *pendingConfig {
	if op.pendingConfig != nil {
		return op.pendingConfig
	}
	return NewPendingConfigWithoutCache(op)
}

func (c *pendingConfig) apply() error {
	//setCrtcConfig: pos, mode, rotation
	//setCrtcGamma: gamma
	//setCrtcTransform: transform, filter
	// allocation of the output maybe changed when rotation/transform changed without change mode
	/*defer func() { c.mask = 0 }()*/

	if c.mode == 0 {
		_, err := randr.SetCrtcConfig(X, c.crtc, xproto.TimeCurrentTime, LastConfigTimeStamp,
			0, 0, 0, c.rotation, nil).Reply()
		return err
	}

	var err error
	if c.mask&_PendingMaskGramma == _PendingMaskGramma {
		err = randr.SetCrtcGammaChecked(X, c.crtc, uint16(len(c.gammaRed)), c.gammaRed, c.gammaGreen, c.gammaBlue).Check()
		if err != nil {
			return fmt.Errorf("PendingConfig apply failed when SetCrtcGammaCheched: %v %v", err, c)
		}
	}
	if c.mask&_PendingMaskTransform == _PendingMaskTransform {
		err = randr.SetCrtcTransformChecked(X, c.crtc, c.transform, uint16(len(c.filterName)), c.filterName, c.filterParams).Check()
		if err != nil {
			return fmt.Errorf("PendingConfig apply failed when SetCrtcTransform: %v %v", err, c)
		}
	}
	if c.mask&_PendingMaskBacklight == _PendingMaskBacklight {
		setOutputBacklight(c.output, c.backlight)
	}

	if c.mask&_PendingMaskPos|_PendingMaskMode|_PendingMaskRotation != 0 {
		var outputs []randr.Output = nil
		if c.mode != 0 {
			outputs = []randr.Output{c.output}
		}

		if DPY.DisplayMode == DisplayModeMirrors {
			_, err = randr.SetCrtcConfig(X, c.crtc, xproto.TimeCurrentTime, LastConfigTimeStamp,
				int16(-c.borderCompensationX), int16(-c.borderCompensationY), c.mode, c.rotation, outputs).Reply()
		} else {
			_, err = randr.SetCrtcConfig(X, c.crtc, xproto.TimeCurrentTime, LastConfigTimeStamp,
				c.posX, c.posY, c.mode, c.rotation, outputs).Reply()
		}
		if err != nil {
			panic(fmt.Errorf("PendingConfig apply failed when SetCrtcConfig(2): %v %v", err, c).Error())
			return fmt.Errorf("PendingConfig apply failed when SetCrtcConfig(2): %v %v", err, c)
		}
		fmt.Println("|||||||||d...", queryOutput(DPY, c.output).Name, c.appliedAllocation(), c.posX, c.posY, c.borderCompensationX)
	}

	{
		/*outputs := []randr.Output{c.output}*/
		/*_, err := randr.SetCrtcConfig(X, c.crtc, xproto.TimeCurrentTime, LastConfigTimeStamp,*/
		cinfo, err := randr.GetCrtcInfo(X, c.crtc, xproto.TimeCurrentTime).Reply()
		rect := xproto.Rectangle{cinfo.X, cinfo.Y, cinfo.Width - uint16(c.borderCompensationX), cinfo.Height - uint16(c.borderCompensationY)}
		if err != nil || rect != c.appliedAllocation() {
			/*panic(fmt.Sprintln("Apply failed...", cinfo, c))*/
			fmt.Println("Apply failed...", queryOutput(DPY, c.output).Name, rect, c.appliedAllocation(), c.posX, c.posY, c.borderCompensationX)
		}
	}
	return nil
}

func (c *pendingConfig) SetMode(m randr.Mode) *pendingConfig {
	if c.mode != m {
		c.mask = c.mask | _PendingMaskMode

		c.mode = m
	}
	return c
}
func (c *pendingConfig) SetPos(x, y int16) *pendingConfig {
	if c.posX != x || c.posY != y {
		c.mask = c.mask | _PendingMaskPos

		c.posX = x
		c.posY = y
	}
	return c
}

func (c *pendingConfig) setCompensation(x, y uint16) *pendingConfig {
	c.borderCompensationX = x
	c.borderCompensationY = y
	/*c.SetPos(-c.borderCompensationX, -c.borderCompensationY)*/
	return c
}

func (c *pendingConfig) SetRotation(r uint16) *pendingConfig {
	if c.rotation != r {
		c.mask = c.mask | _PendingMaskRotation

		if r&0xf == 0 {
			panic("SetRotation Error..")
		}
		c.rotation = r
	}
	return c
}

func (c *pendingConfig) SetTransform(matrix render.Transform, filterName string, params []render.Fixed) *pendingConfig {
	//ignore params at this moment
	if c.transform != matrix || c.filterName != filterName {
		c.mask = c.mask | _PendingMaskTransform

		c.transform = matrix
		c.filterName = filterName
		c.filterParams = params
	}
	return c
}

func (c *pendingConfig) setGamma(red, green, blue []uint16) *pendingConfig {
	c.mask = c.mask | _PendingMaskGramma

	c.gammaRed = red
	c.gammaGreen = green
	c.gammaBlue = blue
	return c
}

func (c *pendingConfig) SetBrightness(brightness float64) *pendingConfig {
	fmt.Println("SetBrightness..........................................", brightness)
	if brightness < 0.1 {
		brightness = 0.1
	}
	if brightness > 1 {
		brightness = 1
	}

	if c.supportBacklight {
		c.backlight = brightness
		c.mask = c.mask | _PendingMaskBacklight
		return c
	} else {
		gammaSize, err := randr.GetCrtcGammaSize(X, c.crtc).Reply()
		if err != nil {
			panic(fmt.Sprintf("GetCrtcGrammSize(crtc:%d) failed: %s", c.crtc, err.Error()))
		}
		red, green, blue := genGammaRamp(gammaSize.Size, brightness)
		return c.setGamma(red, green, blue)
	}
}

func (c *pendingConfig) SetScale(xScale, yScale float32) *pendingConfig {
	c.transform.Matrix11 = double2fixed(xScale)
	c.transform.Matrix22 = double2fixed(yScale)
	c.transform.Matrix33 = double2fixed(1)
	if xScale != 1 || yScale != 1 {
		c.SetTransform(c.transform, "bilinear", nil)
	} else {
		c.SetTransform(c.transform, "nearest", nil)
	}

	return c
}

func (c *pendingConfig) ensureSameRatio(dw, dh uint16) {
}

func (c *pendingConfig) appliedAllocation() (r xproto.Rectangle) {
	minfo := DPY.modes[c.mode]
	if minfo.Width == 0 || minfo.Height == 0 {
		panic("No modeinfo")
	}
	x1, y1, x2, y2 := calcBound(c.transform, c.rotation, minfo.Width, minfo.Height)
	r.X = int16(int(c.posX)+x1) - int16(c.borderCompensationX)
	r.Y = int16(int(c.posY)+y1) - int16(c.borderCompensationY)
	r.Width = uint16(x2-x1) - 2*uint16(c.borderCompensationX)
	r.Height = uint16(y2-y1) - 2*uint16(c.borderCompensationY)
	if r.Width > 1440 {
		fmt.Println("AppliedAllocation ppp:", r, c)
	}

	return
}
func (c *pendingConfig) EnsureSize(width, height uint16, methodHint uint8) *pendingConfig {
	minfo := DPY.modes[c.mode]
	if minfo.Width == width && minfo.Height == height {
		c.SetScale(1, 1).setCompensation(0, 0)
		if DPY.DisplayMode == DisplayModeMirrors {
			c.SetPos(0, 0)
			c.mask = c.mask | _PendingMaskPos
		}
		fmt.Println("Find best mode:", c.appliedAllocation())
		return c
	}
	ow := int16(minfo.Width - width)
	oh := int16(minfo.Height - height)
	ratio := float32(minfo.Width) / float32(minfo.Height)
	switch {
	case ow >= 0 && oh >= 0:
		c.SetScale(1, 1)
		c.setCompensation(uint16(ow/2), uint16(oh/2))
	case ow < 0 && oh <= 0:
		if ratio == float32(width)/float32(height) {
			scale := 1 + float32(-ow)/float32(minfo.Width)
			c.SetScale(scale, scale)
			c.setCompensation(0, 0)
			fmt.Printf("Here!%v/%v=%v'\n", float32(-ow), float32(width), scale)
		} else {
			if -ow < -oh {
				panic("YAHOOO1")
				//width offset smaller
			} else {
				ratio := float32(width) / float32(minfo.Width) /// float32(height)
				margin := int16(float32(minfo.Height)*ratio - float32(height))
				c.setCompensation(0, uint16(margin/2))
				c.SetScale(ratio, ratio)
				fmt.Println("YAHOOO2", c.posX, c.posY, ratio)
				//height offset smaller
			}
		}

	case ow >= 0 && oh <= 0:
		scale := float32(height) / float32(minfo.Height)

		/*margin := minfo.Width - width*/
		margin := int16(math.Ceil(float64(scale*float32(minfo.Width)-float32(width)) / 2))

		c.setCompensation(uint16(margin), 0)
		c.SetScale(scale, scale)

		fmt.Println("XX:", &c, c.appliedAllocation(), margin)

	case ow <= 0 && oh >= 0:
		margin := int16(ratio*float32(width) - float32(height))
		fmt.Printf("Ratio:%v, height:%v, width:%v\n", ratio, height, width)
		c.setCompensation(0, uint16(margin/2))
		scale := float32(width) / float32(minfo.Width)
		c.SetScale(scale, scale)
		fmt.Println("XX2:", c.appliedAllocation())

	}
	{
		alloc := c.appliedAllocation()
		if width != alloc.Width || height != alloc.Height {
			cinfo, _ := randr.GetCrtcInfo(X, c.crtc, 0).Reply()
			rect := xproto.Rectangle{cinfo.X, cinfo.Y, cinfo.Width, cinfo.Height}
			fmt.Println("Ensure to Size failed:", ow, oh, width, height, "DesginAllocation:", c.appliedAllocation(), rect, "------------------")
		}
	}

	return c
}

func smallerRectangle(a, b xproto.Rectangle) bool {
	if a.Width == b.Width && a.Height == b.Height {
		return false
	} else {
		return true
	}
	if a.Width > b.Width || a.Height > b.Height {
		return true
	} else {
		return false
	}
	if a == b {
		return false
	}
	x1, y1 := b.X, b.Y
	x2, y2 := int16(uint16(b.X)+b.Width), int16(uint16(b.Y)+b.Height)

	var inRectangleB = func(x, y int16) bool {
		return (x >= x1 && x <= x2) && (y >= y1 && y <= y2)
	}
	ret := inRectangleB(a.X, a.Y) && inRectangleB(int16(uint16(a.X)+a.Width), int16(uint16(a.Y)+a.Width))
	return ret
}

func (dpy *Display) _getoutputs() (ret []string) {
	for _, op := range dpy.Outputs {
		info, _ := randr.GetCrtcInfo(X, op.crtc, 0).Reply()
		ret = append(ret, fmt.Sprint(op.Name, op.Opened, op.crtc, info.X, info.Y, info.Width, info.Height))

	}
	return
}

func (dpy *Display) adjustScreenSize() []*Output {
	dpy.stopListen()
	defer dpy.startListen()
	var boundAggregate = func(w, h uint16, b xproto.Rectangle) (uint16, uint16) {
		return max(b.Width+uint16(b.X), w), max(b.Height+uint16(b.Y), h)
	}
	var tmpOutputs []*Output
	var w, h uint16

	if dpy.DisplayMode == DisplayModeMirrors {
		w, h = getMirrorSize(dpy.Outputs)
	} else {
		for _, op := range dpy.Outputs {
			w, h = boundAggregate(w, h, op.pendingAllocation())
			fmt.Println("HUHUHU>>>>>>>>>>>>>", op.Name, op.pendingConfig)
		}
	}

	{

		wDif := math.Abs(float64(w) - float64(dpy.Width))
		hDif := math.Abs(float64(h) - float64(dpy.Height))
		if wDif >= 4.0 || hDif >= 4.0 {
			for _, op := range dpy.Outputs {
				if op.Opened {
					info, _ := randr.GetCrtcInfo(X, op.crtc, 0).Reply()
					cw := max(op.Allocation.Width, info.Width)
					ch := max(op.Allocation.Height, info.Height)
					if cw > min(w, DPY.Width) || ch > min(h, DPY.Height) {
						op.setOpened(false)
						tmpOutputs = append(tmpOutputs, op)
					}
				}
			}
		} else {
			w, h = w+uint16(wDif), h+uint16(hDif)
		}
	}
	fmt.Println("-AdjustScreensize before:", w, h, dpy.Width, dpy.Height)
	dpy.setScreenSize(w, h)
	fmt.Println("-AdjustScreensize after:", w, h, dpy.Width, dpy.Height)

	return tmpOutputs
}

func (op *Output) pendingAllocation() (ret xproto.Rectangle) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("pendingAllocation panic:", err)
		}
	}()
	if op.Opened {
		if op.pendingConfig != nil {
			return op.pendingConfig.appliedAllocation()
		} else {
			ret := NewPendingConfig(op).appliedAllocation()
			op.pendingConfig = nil
			return ret
		}
	} else {
		return xproto.Rectangle{0, 0, 0, 0}
	}
}

func (op *Output) setOpened(v bool) {
	if op.Opened == v {
		return
	}
	op.Opened = v
	//op.Opened will be changed when we receive appropriate event
	if v == true {
		if c := op.savedConfig; c != nil {
			// there has an config we saved before
			op.pendingConfig = ClonePendingConfig(op.savedConfig)
			op.savedConfig = nil
			err := op.pendingConfig.apply()
			if err != nil {
				fmt.Println(err)
			}
			op.pendingConfig = nil
			op.Opened = false
		} else {
			op.tryOpen()
		}
	} else {
		config := NewPendingConfig(op)
		op.pendingConfig = nil
		op.savedConfig = ClonePendingConfig(config)
		/*op.Opened = false*/

		err := config.SetMode(0).apply()
		if err != nil {
			fmt.Println(err)
		}
	}
}
