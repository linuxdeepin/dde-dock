package main

import "github.com/BurntSushi/xgb/xproto"
import "github.com/BurntSushi/xgb/randr"
import "github.com/BurntSushi/xgb/render"
import "github.com/BurntSushi/xgb"
import "dlib/logger"
import "math"

var backlightAtom = getAtom(X, "Backlight")

func getAtom(c *xgb.Conn, name string) xproto.Atom {
	r, err := xproto.InternAtom(c, false, uint16(len(name)), name).Reply()
	if err != nil {
		return xproto.AtomNone
	}
	return r.Atom
}
func queryAtomName(c *xgb.Conn, atom xproto.Atom) string {
	r, err := xproto.GetAtomName(c, atom).Reply()
	if err != nil {
		return ""
	}
	return r.Name

}

func queryOutput(dpy *Display, o randr.Output) *Output {
	for _, op := range dpy.Outputs {
		if op.Identify == o {
			return op
		}
	}
	return nil
}
func queryOutputByCrtc(dpy *Display, crtc randr.Crtc) *Output {
	for _, op := range dpy.Outputs {
		if op.crtc == crtc {
			return op
		}
	}
	return nil
}

var (
	edidAtom    = getAtom(X, "EDID")
	borderAtom  = getAtom(X, "Border")
	unknownAtom = getAtom(X, "unknown")
)

func getOutputName(data [128]byte, defaultName string) string {
	timingDescriptor := data[36:]
	for i := 0; i < 4; i++ {
		block := timingDescriptor[i*18 : (i+1)*18]
		if block[3] == 0xfc { //descriptor type == Monitor Name
			data := block[5:]
			for i := 0; i < 13; i++ {
				if data[i] == 0xa {
					return string(data[:i])
				}
			}
		}
	}
	return defaultName
}

func buildMode(info randr.ModeInfo) Mode {
	vTotal := info.Vtotal

	if info.ModeFlags&randr.ModeFlagDoubleScan != 0 {
		vTotal *= 2
	}
	if info.ModeFlags&randr.ModeFlagInterlace != 0 {
		vTotal /= 2
	}

	rate := float64(info.DotClock) / float64(uint32(info.Htotal)*uint32(vTotal))
	rate = math.Floor(rate*10+0.5) / 10
	return Mode{info.Id, info.Width, info.Height, rate}
}

// xgb/randr.GetScreenInfo will panic at this moment.( https://github.com/BurntSushi/xgb/issues/20)
func getScreenInfo() (*randr.GetScreenInfoReply, error) {
	cook := randr.GetScreenInfo(X, Root)

	buf, err := cook.Cookie.Reply()
	if err != nil {
		return nil, err
	}
	if buf == nil {
		return nil, nil
	}
	v := new(randr.GetScreenInfoReply)
	b := 1 // skip reply determinant

	v.Rotations = buf[b]
	b += 1

	v.Sequence = xgb.Get16(buf[b:])
	b += 2

	v.Length = xgb.Get32(buf[b:]) // 4-byte units
	b += 4

	v.Root = xproto.Window(xgb.Get32(buf[b:]))
	b += 4

	v.Timestamp = xproto.Timestamp(xgb.Get32(buf[b:]))
	b += 4

	v.ConfigTimestamp = xproto.Timestamp(xgb.Get32(buf[b:]))
	b += 4

	v.NSizes = xgb.Get16(buf[b:])
	b += 2

	v.SizeID = xgb.Get16(buf[b:])
	b += 2

	v.Rotation = xgb.Get16(buf[b:])
	b += 2

	v.Rate = xgb.Get16(buf[b:])
	b += 2

	v.NInfo = xgb.Get16(buf[b:])
	b += 2

	b += 2 // padding

	v.Sizes = make([]randr.ScreenSize, v.NSizes)
	b += randr.ScreenSizeReadList(buf[b:], v.Sizes)

	return v, nil
}

func parseRandR(randr uint16) (uint16, uint16) {
	rotation := randr & 0xf
	reflect := randr & 0xf0
	switch rotation {
	case 1, 2, 4, 8:
		break
	default:
		logger.Println("invalid rotation value", rotation, randr)
		rotation = 1
	}
	switch reflect {
	case 0, 16, 32, 48:
		break
	default:
		logger.Println("invalid reflect value", reflect, randr)
		reflect = 0
	}
	return rotation, reflect
}

func parseRotations(rotations uint16) (ret []uint16) {
	if rotations&randr.RotationRotate0 == randr.RotationRotate0 {
		ret = append(ret, randr.RotationRotate0)
	}
	if rotations&randr.RotationRotate90 == randr.RotationRotate90 {
		ret = append(ret, randr.RotationRotate90)
	}
	if rotations&randr.RotationRotate180 == randr.RotationRotate180 {
		ret = append(ret, randr.RotationRotate180)
	}
	if rotations&randr.RotationRotate270 == randr.RotationRotate270 {
		ret = append(ret, randr.RotationRotate270)
	}
	return
}
func parseReflects(rotations uint16) (ret []uint16) {
	ret = append(ret, 0) //the normal reflect

	if rotations&randr.RotationReflectX == randr.RotationReflectX {
		ret = append(ret, randr.RotationReflectX)
	}
	if rotations&randr.RotationReflectY == randr.RotationReflectY {
		ret = append(ret, randr.RotationReflectY)
	}

	if (rotations&randr.RotationReflectX == randr.RotationReflectX) && (rotations&randr.RotationReflectY == randr.RotationReflectY) {
		ret = append(ret, randr.RotationReflectX|randr.RotationReflectY)
	}

	return
}

func isCrtcConnected(c *xgb.Conn, crtc randr.Crtc) bool {
	cinfo, err := randr.GetCrtcInfo(c, crtc, 0).Reply()
	if err != nil {
		panic(err)
	}
	if cinfo.Mode == 0 {
		return false
	} else if cinfo.NumOutputs == 0 {
		return false
	} else {
		oinfo, _ := randr.GetOutputInfo(X, cinfo.Outputs[0], 0).Reply()
		if oinfo.Crtc != crtc {
			return false
		}
	}
	return true
}

func supportedBacklight(c *xgb.Conn, output randr.Output) (bool, float64) {
	prop, err := randr.GetOutputProperty(c, output, backlightAtom, xproto.AtomAny, 0, 1, false, false).Reply()
	if err != nil || prop.NumItems != 1 {
		return false, 100
	}
	return true, float64(xgb.Get32(prop.Data)) / 100
}

func parseRotationSize(rotation, width, height uint16) (uint16, uint16) {
	if rotation == randr.RotationRotate90 || rotation == randr.RotationRotate270 {
		return height, width
	} else {
		return width, height
	}
}

func parseScreenSize(ops []*Output) (width, height uint16) {
	for _, op := range ops {
		if op.Opened {
			alloc := op.pendingAllocation()
			width = max(width, alloc.Width)
			height = max(height, alloc.Height)
		}
	}
	return
}

func fixed2double(v render.Fixed) float32 {
	return float32(v) / 65536
}
func double2fixed(v float32) render.Fixed {
	return render.Fixed(v * 65536)
}

func genGammaRamp(size uint16, brightness float64) (red []uint16, green []uint16, blue []uint16) {
	red = make([]uint16, size)
	green = make([]uint16, size)
	blue = make([]uint16, size)

	step := uint16(65536 / uint32(size))
	for i := uint16(0); i < size; i++ {
		red[i] = uint16(float64(step*i) * brightness)
		green[i] = uint16(float64(step*i) * brightness)
		blue[i] = uint16(float64(step*i) * brightness)
	}
	return
}

func genTransformByScale(xScale float32, yScale float32) render.Transform {
	m := render.Transform{}
	m.Matrix11 = double2fixed(xScale)
	m.Matrix22 = double2fixed(yScale)
	m.Matrix33 = double2fixed(1)
	return m
}

func setOutputBorder(op randr.Output, border Border) {
	var buf [2 * 4]byte
	xgb.Put16(buf[0:2], border.Left)
	xgb.Put16(buf[2:4], border.Top)
	xgb.Put16(buf[4:6], border.Right)
	xgb.Put16(buf[6:8], border.Bottom)

	err := randr.ChangeOutputPropertyChecked(X, op, borderAtom,
		xproto.AtomInteger, 16, xproto.PropModeReplace, 4,
		buf[:]).Check()
	if err != nil {
		logger.Println(err)
	}
}
func setOutputBacklight(op randr.Output, light float64) {
	var buf [4]byte
	xgb.Put32(buf[0:4], uint32(light*100))

	err := randr.ChangeOutputPropertyChecked(X, op, backlightAtom,
		xproto.AtomInteger, 32, xproto.PropModeReplace, 1,
		buf[:]).Check()
	if err != nil {
		logger.Println(err)
	}
}

func getOutputBorder(op randr.Output) (ret Border) {
	prop, err := randr.GetOutputProperty(X, op, borderAtom, xproto.AtomAny, 0, 2, false, false).Reply()
	defer func() {
		if err := recover(); err != nil {
			logger.Println("getOutputBorder recvied an malformed packet", prop)
			ret = Border{}
		}
	}()
	if err != nil {
		return Border{}
	}
	switch prop.NumItems {
	case 0:
		return Border{}
	case 1:
		value := xgb.Get16(prop.Data)
		return Border{value, value, value, value}
	case 2:
		lr, tb := xgb.Get16(prop.Data), xgb.Get16(prop.Data[2:])
		return Border{lr, tb, lr, tb}
	case 4:
		l, t, r, b := xgb.Get16(prop.Data), xgb.Get16(prop.Data[2:]), xgb.Get16(prop.Data[4:]), xgb.Get16(prop.Data[6:])
		return Border{l, t, r, b}
	}
	return Border{}
}

func calcBound(m render.Transform, rotation uint16, width, height uint16) (x1, y1, x2, y2 int) {
	var applyTransform = func(m render.Transform, x float32, y float32) (int, int, int, int) {
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
	switch rotation & 0xf {
	case randr.RotationRotate90, randr.RotationRotate270:
		width, height = height, width
	}

	var min = func(a, b int) int {
		if a > b {
			return b
		}
		return a
	}
	var max = func(a, b int) int {
		if a < b {
			return b
		}
		return a
	}
	x1, y1, x2, y2 = applyTransform(m, 0, 0)

	tx1, ty1, tx2, ty2 := applyTransform(m, float32(width), 0)
	x1 = min(x1, tx1)
	y1 = min(y1, ty1)
	x2 = max(x2, tx2)
	y2 = max(y2, ty2)

	tx1, ty1, tx2, ty2 = applyTransform(m, float32(width), float32(height))
	x1 = min(x1, tx1)
	y1 = min(y1, ty1)
	x2 = max(x2, tx2)
	y2 = max(y2, ty2)

	tx1, ty1, tx2, ty2 = applyTransform(m, 0, float32(height))
	x1 = min(x1, tx1)
	y1 = min(y1, ty1)
	x2 = max(x2, tx2)
	y2 = max(y2, ty2)
	return
}
func calcBound2(m render.Transform, rotation uint16, x, y float32, width, height uint16) (x1, y1, x2, y2 int) {
	var applyTransform = func(m render.Transform, x float32, y float32) (int, int, int, int) {
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
	switch rotation & 0xf {
	case randr.RotationRotate90, randr.RotationRotate270:
		width, height = height, width
	}

	var min = func(a, b int) int {
		if a > b {
			return b
		}
		return a
	}
	var max = func(a, b int) int {
		if a < b {
			return b
		}
		return a
	}
	x1, y1, x2, y2 = applyTransform(m, x, y)

	tx1, ty1, tx2, ty2 := applyTransform(m, float32(width), y)
	x1 = min(x1, tx1)
	y1 = min(y1, ty1)
	x2 = max(x2, tx2)
	y2 = max(y2, ty2)

	tx1, ty1, tx2, ty2 = applyTransform(m, float32(width), float32(height))
	x1 = min(x1, tx1)
	y1 = min(y1, ty1)
	x2 = max(x2, tx2)
	y2 = max(y2, ty2)

	tx1, ty1, tx2, ty2 = applyTransform(m, x, float32(height))
	x1 = min(x1, tx1)
	y1 = min(y1, ty1)
	x2 = max(x2, tx2)
	y2 = max(y2, ty2)
	return
}

func guestMode(op *Output, w, h uint16, rate float64) randr.Mode {
	for _, m := range op.ListModes() {
		if m.Width == w && m.Height == h && m.Rate == rate {
			return randr.Mode(m.ID)
		}
	}
	return 0
}
