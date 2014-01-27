package main

/*import fmtp "github.com/kr/pretty"*/
import "fmt"
import "testing"
import "time"
import "github.com/BurntSushi/xgb/randr"
import "github.com/BurntSushi/xgb/xproto"
import . "launchpad.net/gocheck"

func delay() {
	<-time.After(time.Millisecond * 500)
}
func Test(t *testing.T) { TestingT(t) }

func init() {
	/*Suite(DPY)*/
	for _, op := range DPY.Outputs {
		if op.Name == "LVDS1" {
			Suite(op)
		}
		/*Suite(op)*/
	}
}

func (op *Output) TestEnsureSize(c *C) {
	op.setOpened(true)
	sizes := []struct {
		w uint16
		h uint16
	}{
		{1152, 864},
		/*{1153, 864},*/
		/*{800, 600},*/
		/*{1280, 800},*/
		/*{640, 400},*/
		/*{1280, 800},*/
		/*{1153, 864},*/
	}
	fmt.Println("Begin.....:", DPY.Width, DPY.Height)
	for _, size := range sizes {
		op.EnsureSize(size.w, size.h, EnsureSizeHintAuto)
		DPY.ApplyChanged()

		info, err := randr.GetCrtcInfo(X, op.crtc, LastConfigTimeStamp).Reply()
		if err != nil {
			c.Skip(fmt.Sprintln("GetCrtcInfo EFailed...", op.crtc, err))
		}
		c.Assert(err, Not(NotNil))
		c.Check(size.w, Equals, uint16(int16(info.Width)+info.X+info.X))
		c.Check(size.h, Equals, uint16(int16(info.Height)+info.Y+info.Y))
		delay()
		delay()
		delay()
		delay()
		delay()
		delay()
		c.Check(size.w, Equals, DPY.Width)
		c.Check(size.h, Equals, DPY.Height)
		fmt.Println("::::", info.X, info.Y, info.Width, info.Height, DPY.Width, DPY.Height)
		fmt.Println("\n")
	}

}

func (dpy *Display) TestScreenInfo(c *C) {
	/*return*/
	delay()

	for _, r := range dpy.ListRotations() {
		if r == dpy.Rotation {
			return
		}
	}
	for _, r := range dpy.ListReflect() {
		if r == dpy.Reflect {
			return
		}
	}
	c.Fail()
}

func (dpy *Display) TestMirrorMode(c *C) {

	for _, op := range dpy.Outputs {
		if op.Name != "LVDS1" {
			fmt.Println("Set MirrorOutput to ", op.Name)
			op.SetMode(0x9d)
			dpy.ApplyChanged()
			/*dpy.SetMirrorMode(true)*/
			/*dpy.SetMirrorOutput(op)*/
		}
	}
	delay()
	delay()
	delay()
}

func (dpy *Display) TestOutputList(c *C) {
	c.Check(len(dpy.Outputs) >= 1, Equals, true)
}
func (dpy *Display) TestApply(c *C) {
	/*return*/
	for _, op := range dpy.Outputs {
		if op.Name == "LVDS1" {
			op.pendingConfig = NewPendingConfig(op)
			op.pendingConfig.SetScale(1.125, 1.125)
			fmt.Println("RECT:", op.pendingConfig.AppliedAllocation())
		} else {
			op.pendingConfig = NewPendingConfig(op)
			op.pendingConfig.SetRotation(randr.RotationRotate0 | randr.RotationReflectY | randr.RotationReflectX)

		}
	}
	dpy.ApplyChanged()
}

func (dpy *Display) TestPrimaryOutput(c *C) {
	/*return*/
	po := dpy.PrimaryOutput
	savedIdentify := uint32(0)
	if po != nil {
		savedIdentify = uint32(po.Identify)
	}
	defer func() {
		dpy.SetPrimaryOutput(savedIdentify)
	}()

	dpy.SetPrimaryOutput(0)
	delay()
	delay()
	delay()

	for _, op := range dpy.Outputs {
		if op.Opened {
			dpy.SetPrimaryOutput(uint32(op.Identify))
			delay()
			delay()
			delay()
			c.Assert(dpy.PrimaryOutput, Equals, op)
			c.Check(dpy.PrimaryRect.Width, Equals, op.Allocation.Width)
			c.Check(dpy.PrimaryRect.Height, Equals, op.Allocation.Height)
			fmt.Println("PrimARY:", dpy.PrimaryRect, "DPY:", op.Allocation)
		}
	}

	if po != nil {
		dpy.SetPrimaryOutput(uint32(po.Identify))
		delay()
		c.Check(dpy.PrimaryOutput, Equals, po)
	}

	dpy.SetPrimaryOutput(0)
	delay()
	delay()
	delay()
	c.Assert(dpy.PrimaryOutput, Not(NotNil))
	c.Check(dpy.PrimaryRect.Width, Equals, dpy.Width)
	c.Check(dpy.PrimaryRect.Height, Equals, dpy.Height)
}

func (op *Output) TestInfo(c *C) {
	/*return*/
	c.Check(op.Brightness >= 0 && op.Brightness <= 1, Equals, true)

	delay()
	find := false
	for _, r := range op.ListModes() {
		if r == op.Mode {
			find = true
		}
	}
	if op.Opened {
		c.Check(find, Equals, true)
	} else {
		fmt.Println("OP:", op.Name, "Mode:", op.Mode)
		c.Check(find, Equals, false)
	}

	crtcInfo, err := randr.GetCrtcInfo(X, op.crtc, 0).Reply()
	if op.Opened {
		c.Assert(err, Equals, nil)
		c.Check(op.Mode, Equals, buildMode(DPY.modes[crtcInfo.Mode]))
		c.Check(op.Rotation, Equals, uint16(crtcInfo.Rotation))
	} else {
		c.Assert(err, NotNil)
	}

	c.Check(op.Opened, Equals, op.crtc != 0)

	_, err = randr.GetOutputInfo(X, op.Identify, 0).Reply()
	c.Check(err, Equals, nil)

	op.ListModes()
	op.ListRotations()
	delay()
}

func (op *Output) TestClose(c *C) {
	/*return*/
	delay()

	v := op.Opened
	delay()
	op.setOpened(true)

	delay()
	op.setOpened(v)
	delay()
	c.Check(op.Opened, Equals, v)

	delay()
	delay()
}

func (op *Output) TestRandr(c *C) {
	/*return*/
	rv := op.Rotation
	fv := op.Reflect

	for _, reflect := range op.ListReflect() {
		delay()
		op.setReflect(reflect)
		delay()
		for _, r := range op.ListRotations() {
			fmt.Println("op.setReflect>", reflect, r)
			op.setRotation(r)
			delay()
			cinfo, err := randr.GetCrtcInfo(X, op.crtc, 0).Reply()
			if err != nil {
				panic(err)
			}
			c.Check(op.Allocation.X, Equals, cinfo.X)
			c.Check(op.Allocation.Y, Equals, cinfo.Y)
			c.Check(op.Allocation.Width, Equals, cinfo.Width)
			c.Check(op.Allocation.Height, Equals, cinfo.Height)
		}
	}
	delay()

	op.setRotation(rv)
	delay()
	op.setReflect(fv)
	delay()
	c.Check(rv, Equals, op.Rotation)
	c.Check(fv, Equals, op.Reflect)
}

func (op *Output) TestMode(c *C) {
	/*return*/
	vm := op.Mode
	vw, vh := op.Allocation.Width, op.Allocation.Height
	for _, m := range op.ListModes() {
		delay()
		op.SetMode(m.ID)
		rect := op.pendingConfig.AppliedAllocation()
		DPY.ApplyChanged()
		delay()
		delay()
		c.Check(op.Allocation.Width, Equals, rect.Width)
		c.Check(op.Allocation.Height, Equals, rect.Height)
	}
	op.SetMode(vm.ID)
	DPY.ApplyChanged()
	delay()
	delay()
	c.Check(vw, Equals, op.Allocation.Width)
	c.Check(vh, Equals, op.Allocation.Height)
}
func (op *Output) TestAllocation(c *C) {
	/*return*/
	delay()
	delay()
	if op.Name == "LVDS1" {
		op.SetPos(0, 100)
		c.Check(op.pendingConfig.AppliedAllocation(), Equals, xproto.Rectangle{0, 100, op.Allocation.Width, op.Allocation.Height})
		DPY.ApplyChanged()
		c.Assert(op.pendingConfig, Equals, (*pendingConfig)(nil))
		delay()
		delay()
		c.Check(op.Allocation, Equals, xproto.Rectangle{0, 100, op.Allocation.Width, op.Allocation.Height})
	} else {
		op.SetPos(1280, 0)
	}
}

func (op *Output) TestPos(c *C) {
	/*return*/
	DPY.SetPrimaryOutput(uint32(op.Identify))
	delay()
	op.SetPos(0, 100)
	rect := op.pendingConfig.AppliedAllocation()
	DPY.ApplyChanged()
	delay()
	delay()
	c.Check(DPY.PrimaryOutput, Equals, op)
	c.Check(DPY.PrimaryRect, Equals, rect)
	c.Check(DPY.PrimaryRect, Equals, op.Allocation)
	c.Check(rect.Y, Equals, int16(100))
}

func (op *Output) TestGramm(c *C) {
	/*return*/
	vb := op.Brightness

	for i := 0.1; i < 1; i = i + 0.1 {
		op.setBrightness(i)
		delay()
		delay()
	}
	for i := 0.1; i < 1; i = i + 0.001 {
		<-time.After(time.Millisecond * 10)
		op.setBrightness(i)
	}
	for i := 1.0; i > 0; i = i - 0.001 {
		<-time.After(time.Millisecond * 10)
		op.setBrightness(i)
	}

	delay()
	delay()
	op.setBrightness(vb)
	delay()
	delay()
	c.Check(op.Brightness, Equals, vb)

}

func (op *Output) TestChangeMode(c *C) {
	/*return*/
	op.setRotation(1)
	delay()
	delay()
	op.setRotation(2)
	delay()
	delay()
	delay()
}

func TestEnsure(t *testing.T) {
	return
	if DPY.PrimaryOutput == nil {
		rect := xproto.Rectangle{0, 0, DPY.Width, DPY.Height}
		if DPY.PrimaryRect != rect {
			t.Fatal("PriamryRect not mathced when no primary output", DPY.PrimaryRect, rect)
		}
	} else {
		if DPY.PrimaryRect != DPY.PrimaryOutput.Allocation {
			t.Fatal("PriamryRect not mathced when primary output with allocation:", DPY.PrimaryRect, DPY.PrimaryOutput.Allocation)
		}
	}
}
