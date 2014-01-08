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

var dpy *Display

func init() {
	dpy = NewDisplay()
	Suite(dpy)
	for _, op := range dpy.Outputs {
		/*if op.Name != "LVDS1" {*/
		/*Suite(op)*/
		/*}*/
		Suite(op)
	}
	for _, op := range dpy.Outputs {
		xx, _ := randr.GetOutputInfo(X, op.Identify, 0).Reply()
		fmt.Println(op.Name, " : ", op.crtc, xx.Crtcs)
	}
}

func (dpy *Display) TestScreenInfo(c *C) {
	delay()
	c.Check(dpy.Width, Equals, DefaultScreen.WidthInPixels)
	c.Check(dpy.Height, Equals, DefaultScreen.HeightInPixels)

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

func (dpy *Display) TestOutputList(c *C) {
	c.Check(len(dpy.Outputs) >= 1, Equals, true)
}

func (dpy *Display) TestPrimaryOutput(c *C) {
	po := dpy.PrimaryOutput
	savedIdentify := uint32(0)
	if po != nil {
		savedIdentify = uint32(po.Identify)
	}
	defer func() {
		dpy.SetPrimary(savedIdentify)
	}()

	for _, op := range dpy.Outputs {
		if op.Opened {
			dpy.SetPrimary(uint32(op.Identify))
			delay()
			c.Check(dpy.PrimaryRect.Width, Equals, op.Allocation.Width)
			c.Check(dpy.PrimaryRect.Height, Equals, op.Allocation.Height)
		}
	}

	if po != nil {
		dpy.SetPrimary(uint32(po.Identify))
		delay()
		c.Check(dpy.PrimaryOutput, Equals, po)
	}

	dpy.SetPrimary(0)
	delay()
	c.Check(dpy.PrimaryRect.Width, Equals, dpy.Width)
	c.Check(dpy.PrimaryRect.Height, Equals, dpy.Height)
}

func (op *Output) TestInfo(c *C) {
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
		c.Check(op.Mode, Equals, buildMode(dpy.modes[crtcInfo.Mode]))
		c.Check(op.Rotation, Equals, uint16(crtcInfo.Rotation))
	} else {
		c.Assert(err, NotNil)
	}

	c.Check(op.Opened, Equals, op.crtc != 0)

	_, err = randr.GetOutputInfo(X, op.Identify, 0).Reply()
	c.Check(err, Equals, nil)

	op.ListModes()
	op.ListRotations()
	op.updateCrtc(dpy)
	delay()
}

func (op *Output) TestClose(c *C) {
	delay()

	v := op.Opened
	op.setOpened(false)
	delay()
	op.setOpened(true)

	delay()
	op.setOpened(v)
	delay()
	c.Check(op.Opened, Equals, v)

	return
}

func (op *Output) TestRandr(c *C) {
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
	vm := op.Mode
	for _, m := range op.ListModes() {
		delay()
		op.setMode(m.ID)
	}
	delay()
	op.setMode(vm.ID)
}
func (op *Output) TestAllocation(c *C) {
	delay()
	delay()
	if op.Name == "LVDS1" {
		op.setAllocation(0, 100, 0, 0, 0)
	} else {
		op.setAllocation(1280, 0, 0, 0, 0)
	}
	delay()
	delay()
	delay()
	delay()
}

func TestEnsure(t *testing.T) {
	if dpy.PrimaryOutput == nil {
		rect := xproto.Rectangle{0, 0, dpy.Width, dpy.Height}
		if dpy.PrimaryRect != rect {
			t.Fatal("PriamryRect not mathced when no primary output")
		}
	} else {
		if dpy.PrimaryRect != dpy.PrimaryOutput.Allocation {
			t.Fatal("PriamryRect not mathced when primary output with allocation:", dpy.PrimaryOutput.Allocation)
		}
	}
}
