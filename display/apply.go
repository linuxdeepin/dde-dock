package main

import "github.com/BurntSushi/xgb/randr"
import "github.com/BurntSushi/xgb"
import "github.com/BurntSushi/xgb/xproto"
import "sync"
import "math"
import "fmt"

import "github.com/BurntSushi/xgb/render"

const (
	_PendingMaskMode = 1 << iota
	_PendingMaskPos
	_PendingMaskBorder
	_PendingMaskTransform
	_PendingMaskRotation
	_PendingMaskGramma
)

var UnitMatrix = render.Transform{65536, 0, 0, 0, 65536, 0, 0, 0, 65536}

const (
	EnsureSizeHintAuto uint8 = iota
	EnsureSizeHintPanning
	EnsureSizeHintBorderScale
)

var changeLock, changeUnlock = func() (func(), func()) {
	var locker sync.Mutex
	return func() {
			locker.Lock()
			xproto.GrabServer(X)
		}, func() {
			xproto.UngrabServer(X)
			locker.Unlock()
		}
}()

type pendingConfig struct {
	crtc   randr.Crtc
	output randr.Output
	mask   int

	mode     randr.Mode
	posX     int16
	posY     int16
	rotation uint16

	transform    render.Transform
	filterName   string
	filterParams []render.Fixed

	border Border

	// doesn't influence allocation
	gammaRed   []uint16
	gammaGreen []uint16
	gammaBlue  []uint16
}

func NewPendingConfig(op *Output) *pendingConfig {
	if op.pendingConfig != nil {
		return op.pendingConfig
	}
	r := &pendingConfig{}
	r.crtc = op.crtc
	r.output = op.Identify

	r.mode = randr.Mode(op.Mode.ID)
	validMode := false
	for _, m := range op.ListModes() {
		if r.mode == randr.Mode(m.ID) {
			validMode = true
			break
		}
	}
	if !validMode {
		r.mode = op.bestMode
	}

	r.posX = op.Allocation.X
	r.posY = op.Allocation.Y
	r.rotation = op.Rotation | op.Reflect

	tinfo, err := randr.GetCrtcTransform(X, op.crtc).Reply()
	if err != nil {
		panic(fmt.Sprintf("NewPendingCofing failed at GetCrtcTransform(crtc=%v)", r.crtc))
	}
	r.transform = tinfo.CurrentTransform
	r.filterName = tinfo.CurrentFilterName
	r.filterParams = tinfo.CurrentParams
	return r
}

func (c *pendingConfig) SetMode(m randr.Mode) *pendingConfig {
	c.mask = c.mask | _PendingMaskMode

	c.mode = m
	return c
}
func (c *pendingConfig) SetPos(x, y int16) *pendingConfig {
	c.mask = c.mask | _PendingMaskPos

	c.posX = x
	c.posY = y

	return c
}

func (c *pendingConfig) SetRotation(r uint16) *pendingConfig {
	c.mask = c.mask | _PendingMaskRotation

	c.rotation = r
	return c
}
func (c *pendingConfig) SetBorder(b Border) *pendingConfig {
	c.mask = c.mask | _PendingMaskBorder

	c.border = b
	return c
}
func (c *pendingConfig) SetTransform(matrix render.Transform, filterName string, params []render.Fixed) *pendingConfig {
	c.mask = c.mask | _PendingMaskTransform

	c.transform = matrix
	c.filterName = filterName
	c.filterParams = params

	return c
}

func (c *pendingConfig) SetGamma(red, green, blue []uint16) *pendingConfig {
	c.mask = c.mask | _PendingMaskGramma

	c.gammaRed = red
	c.gammaGreen = green
	c.gammaBlue = blue
	return c
}

func (c *pendingConfig) SetScale(xScale, yScale float32) *pendingConfig {
	c.mask = c.mask | _PendingMaskTransform

	c.transform.Matrix11 = double2fixed(xScale)
	c.transform.Matrix22 = double2fixed(yScale)
	c.transform.Matrix33 = double2fixed(1)
	if xScale != 1 || yScale != 1 {
		c.filterName = "bilinear"
	} else {
		c.filterName = "nearest"
	}

	return c
}

/*func (op *Output) deduceAdjustMethod(w, h uint16) uint8 {*/
/*ratio := float64(w) / float64(h)*/
/*for _, modeinfo := range op.ListModes() {*/
/*mratio, mw, mh := float64(modeinfo.Width)/float64(modeinfo.Height), modeinfo.Width, modeinfo.Height*/
/*if math.Abs(mratio-ratio) < 0.00001 {*/
/*scale := float32(mw) / float32(w)*/
/*op.pendingConfig = NewPendingConfig(op).SetScale(scale, scale)*/
/*} else {*/
/*scale := float32(mh) / float32(h)*/
/*op.pendingConfig = NewPendingConfig(op).SetScale(scale, scale)*/
/*}*/
/*}*/
/*return AdjustModeNone*/
/*}*/

func (c *pendingConfig) ensureSameRatio(dw, dh uint16) {
}

func (c *pendingConfig) appliedAllocation() (r xproto.Rectangle) {
	minfo := DPY.modes[c.mode]
	width := minfo.Width - c.border.Left - c.border.Right
	height := minfo.Height - c.border.Top - c.border.Bottom
	x1, y1, x2, y2 := calcBound(c.transform, c.rotation, width, height)
	r.X = int16(int(c.posX) + x1)
	r.Y = int16(int(c.posY) + y1)
	r.Width = uint16(x2 - x1)
	r.Height = uint16(y2 - y1)

	//remove border space
	/*r.X = r.X - int16(c.border.Left)*/
	/*r.Y = r.Y - int16(c.border.Top)*/
	/*r.Width = r.Width - c.border.Right*/
	/*r.Height = r.Height - c.border.Bottom*/
	/*fmt.Println(c.border)*/
	return
}
func (c *pendingConfig) EnsureSize(width, height uint16, methodHint uint8) *pendingConfig {
	minfo := DPY.modes[c.mode]
	if minfo.Width == width && minfo.Height == height {
		fmt.Println("SameSize....")
		return c
	}
	ow := int16(minfo.Width - width)
	oh := int16(minfo.Height - height)
	ratio := minfo.Width / minfo.Height
	fmt.Println("WTF", ow, oh)
	switch {
	case ow > 0 && oh > 0:
		fmt.Println("Fucn..")
		c.SetBorder(Border{uint16(ow / 2), uint16(oh / 2), uint16(ow / 2), uint16(oh / 2)})
		c.SetPos(-int16(ow/2), -int16(oh/2))

	case ow < 0 && oh < 0:
		if ratio == width/height {
			scale := 1 + float32(-ow)/float32(minfo.Width)
			c.SetScale(scale, scale)
			fmt.Printf("Here!%v/%v=%v'\n", float32(-ow), float32(width), scale)
		} else {
			panic("XX")
		}

	case ow > 0 && oh < 0:
		margin := width - ratio*height
		scale := float32(-oh) / float32(height)
		c.SetBorder(Border{margin / 2, 0, margin / 2, 0})
		c.SetScale(scale, scale)
	case ow < 0 && oh > 0:
		fmt.Println("GoHere....", ow, oh)
		margin := height - ratio*width
		scale := float32(-ow) / float32(width)
		c.SetBorder(Border{0, margin / 2, 0, margin / 2})
		c.SetScale(scale, scale)

		/*case ow < 0 && oh < 0:*/
		/*if ratio == width/height {*/
		/*scale := float32(-ow) / float32(width)*/
		/*if methodHint == EnsureSizeHintPanning {*/
		/*//setPanning*/
		/*} else {*/
		/*c.SetScale(scale, scale)*/
		/*}*/
		/*} else {*/
		/*if methodHint == EnsureSizeHintBorderScale {*/
		/*c.SetScale(scale, scale)*/
		/*} else {*/
		/*//setPanning*/
		/*}*/
		/*}*/
	}

	return c
}

func (c *pendingConfig) apply() error {
	//setCrtcConfig: pos, mode, rotation
	//setCrtcGamma: gamma
	//setCrtcTransform: transform, filter
	//setOutputProperty: border

	var err error
	if c.mask&_PendingMaskGramma == _PendingMaskGramma {
		err = randr.SetCrtcGammaChecked(X, c.crtc, uint16(len(c.gammaRed)), c.gammaRed, c.gammaGreen, c.gammaBlue).Check()
		if err != nil {
			return fmt.Errorf("PendingConfig apply failed when SetCrtcGammaCheched: %v %v", err, c)
		}
	}
	if c.mask&_PendingMaskBorder == _PendingMaskBorder {
		var buf [2 * 4]byte
		xgb.Put16(buf[0:2], c.border.Left)
		xgb.Put16(buf[2:4], c.border.Top)
		xgb.Put16(buf[4:6], c.border.Right)
		xgb.Put16(buf[6:8], c.border.Bottom)
		err = randr.ChangeOutputPropertyChecked(X, c.output, borderAtom, xproto.AtomInteger, 16, xproto.PropModeReplace, 4, buf[:]).Check()
		if err != nil {
			return fmt.Errorf("PendingConfig apply failed when ChangeOutputProperty: %v %v", err, c)
		}
	}
	if c.mask&_PendingMaskTransform == _PendingMaskTransform {
		err = randr.SetCrtcTransformChecked(X, c.crtc, c.transform, uint16(len(c.filterName)), c.filterName, c.filterParams).Check()
		if err != nil {
			return fmt.Errorf("PendingConfig apply failed when SetCrtcTransform: %v %v", err, c)
		}
	}

	if c.mask&_PendingMaskPos|_PendingMaskMode|_PendingMaskRotation != 0 {
		_, err = randr.SetCrtcConfig(X, c.crtc, xproto.TimeCurrentTime, LastConfigTimeStamp, c.posX, c.posY, c.mode, c.rotation, []randr.Output{c.output}).Reply()
		if err != nil {
			return fmt.Errorf("PendingConfig apply failed when SetCrtcConfig: %v %v", err, c)
		}
	}
	return nil
}

func applyTransform(m render.Transform, x float32, y float32) (int, int, int, int) {
	rx := fixed2double(m.Matrix11)*x + fixed2double(m.Matrix12)*y + fixed2double(m.Matrix13)*1
	ry := fixed2double(m.Matrix21)*x + fixed2double(m.Matrix22)*y + fixed2double(m.Matrix23)*1
	rw := fixed2double(m.Matrix31)*x + fixed2double(m.Matrix32)*y + fixed2double(m.Matrix33)*1

	if rw == 0 {
		return 0, 0, 0, 0
	}

	rx = rx / rw
	if rx > 32767 || rx < -32767 {
		return 0, 0, 0, 0
	}

	ry = ry / rw
	if ry > 32767 || ry < -32767 {
		return 0, 0, 0, 0
	}

	rw = rw / rw
	if rw > 32767 || rw < -32767 {
		return 0, 0, 0, 0
	}
	return int(math.Floor(float64(rx))), int(math.Floor(float64(ry))), int(math.Ceil(float64(rx))), int(math.Ceil(float64(ry)))
}
func calcBound(m render.Transform, rotation uint16, width, height uint16) (x1, y1, x2, y2 int) {
	switch rotation & 0xf {
	case randr.RotationRotate90, randr.RotationRotate270:
		width, height = height, width
	}

	x1, y1, x2, y2 = applyTransform(m, 0, 0)

	tx1, ty1, tx2, ty2 := applyTransform(m, float32(width), 0)
	if tx1 < x1 {
		x1 = tx1
	}
	if ty1 < y1 {
		y1 = ty1
	}
	if tx2 > x2 {
		x2 = tx2
	}
	if ty2 > y2 {
		y2 = ty2
	}

	tx1, ty1, tx2, ty2 = applyTransform(m, float32(width), float32(height))
	if tx1 < x1 {
		x1 = tx1
	}
	if ty1 < y1 {
		y1 = ty1
	}
	if tx2 > x2 {
		x2 = tx2
	}
	if ty2 > y2 {
		y2 = ty2
	}

	tx1, ty1, tx2, ty2 = applyTransform(m, 0, float32(height))
	if tx1 < x1 {
		x1 = tx1
	}
	if ty1 < y1 {
		y1 = ty1
	}
	if tx2 > x2 {
		x2 = tx2
	}
	if ty2 > y2 {
		y2 = ty2
	}
	return
}

func boundAggregate(w, h uint16, b xproto.Rectangle) (uint16, uint16) {
	bw := uint16(b.X + int16(b.Width))
	bh := uint16(b.Y + int16(b.Height))
	if bw > w {
		w = bw
	}
	if bh > h {
		h = bh
	}
	return w, h
}

func (dpy *Display) ApplyChanged() {
	changeLock()
	defer changeUnlock()

	if mainOP := dpy.MirrorOutput; dpy.MirrorMode && mainOP != nil && mainOP.pendingConfig == nil {
		fmt.Println("MainOP:", mainOP.Name)
		w, h := mainOP.Allocation.Width, mainOP.Allocation.Height
		for _, op := range dpy.Outputs {
			if op.Opened && op != mainOP {
				op.pendingConfig = NewPendingConfig(op).SetPos(0, 0).SetBorder(Border{0, 0, 0, 0}).SetRotation(mainOP.Rotation|mainOP.Reflect).SetScale(1, 1)
				fmt.Println(op.Name, "Ensure to Size:", w, h, "DesginAllocation:", op.pendingConfig.appliedAllocation())
				op.EnsureSize(w, h, EnsureSizeHintAuto)
			}
		}
	}

	tmpClosedOutput := dpy.adjustScreenSize()

	for _, op := range dpy.Outputs {
		if op.pendingConfig != nil {
			if err := op.pendingConfig.apply(); err != nil {
				fmt.Println("Apply", op.Name, "failed", err)
			}
			op.pendingConfig = nil
		}
	}
	for _, op := range tmpClosedOutput {
		op.setOpened(true)
	}
}

func (dpy *Display) adjustScreenSize() []*Output {
	var tmpOutputs []*Output
	var w, h uint16
	for _, op := range dpy.Outputs {
		if op.Opened {
			if op.pendingConfig != nil {
				w, h = boundAggregate(w, h, op.pendingConfig.appliedAllocation())
			} else {
				w, h = boundAggregate(w, h, NewPendingConfig(op).appliedAllocation())
			}
		}
	}
	for _, op := range dpy.Outputs {
		currentWidth := uint16(op.Allocation.X + int16(op.Allocation.Width))
		currentHeight := uint16(op.Allocation.Y + int16(op.Allocation.Height))
		if currentWidth > w || currentHeight > h ||
			currentWidth > DefaultScreen.WidthInPixels || currentHeight > DefaultScreen.HeightInPixels {
			op.setOpened(false)
			tmpOutputs = append(tmpOutputs, op)
		}

	}
	dpy.setScreenSize(w, h)

	return tmpOutputs
}
