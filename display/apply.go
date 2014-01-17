package main

import "github.com/BurntSushi/xgb/randr"
import "github.com/BurntSushi/xgb"
import "github.com/BurntSushi/xgb/xproto"
import "dlib/logger"
import "sync"
import "math"
import "fmt"

import "github.com/BurntSushi/xgb/render"

const (
	_PendingMaskMode      = 1 << 0
	_PendingMaskPos       = 1 << 1
	_PendingMaskBorder    = 1 << 2
	_PendingMaskTransform = 1 << 3
	_PendingMaskRotation  = 1 << 4
	_PendingMaskGramma    = 1 << 5
)

const (
	AdjustModeNone uint8 = iota
	AdjustModBorder
	AdjustModePanning
	AdjustModeScale
	AdjustModeAuto
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
func (c *pendingConfig) SetBorder(b Border) {
	c.mask = c.mask | _PendingMaskBorder

	c.border = b
}
func (c *pendingConfig) SetTransform(matrix render.Transform, filterName string, params []render.Fixed) {
	c.mask = c.mask | _PendingMaskTransform

	c.transform = matrix
	c.filterName = filterName
	c.filterParams = params
}

func (c *pendingConfig) SetGamma(red, green, blue []uint16) *pendingConfig {
	c.mask = c.mask | _PendingMaskGramma

	c.gammaRed = red
	c.gammaGreen = green
	c.gammaBlue = blue
	return c
}

func (c *pendingConfig) SetScale(xScale, yScale float32) {
	c.mask = c.mask | _PendingMaskTransform

	c.transform.Matrix11 = double2fixed(xScale)
	c.transform.Matrix22 = double2fixed(yScale)
	c.transform.Matrix33 = double2fixed(1)
	if xScale != 1 || yScale != 1 {
		c.filterName = "bilinear"
	} else {
		c.filterName = "nearest"
	}
}

func (c *pendingConfig) Apply() error {
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
func (c *pendingConfig) appliedAllocation() (r xproto.Rectangle) {
	minfo := DPY.modes[c.mode]
	x1, y1, x2, y2 := calcBound(c.transform, c.rotation, minfo.Width, minfo.Height)
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

	if dpy.mirrorMode {
		mainOP := getMirrorOutput(dpy)
		w, h := mainOP.Allocation.Width, mainOP.Allocation.Height
		for _, op := range dpy.Outputs {
			if op.Opened && op != mainOP {
				op.pendingConfig = nil
				op.SetAllocation(0, 0, w, h, AdjustModeAuto)
			}
		}
	}

	dpy.adjustScreenSize()

	for _, op := range dpy.Outputs {
		if op.pendingConfig != nil {
			if err := op.pendingConfig.Apply(); err != nil {
				fmt.Println("Apply", op.Name, "failed", err)
			}
			op.pendingConfig = nil
		}
	}
}

func (dpy *Display) adjustScreenSize() {
	var ops []*Output
	var w, h uint16
	for _, op := range dpy.Outputs {
		if op.Opened {
			if op.pendingConfig != nil {
				ops = append(ops, op)
				w, h = boundAggregate(w, h, op.pendingConfig.appliedAllocation())
			} else {
				w, h = boundAggregate(w, h, NewPendingConfig(op).appliedAllocation())
			}
		}
	}
	dpy.setScreenSize(w, h)
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
