package main

import "fmt"
import "testing"
import "time"
import "github.com/BurntSushi/xgb/randr"
import . "launchpad.net/gocheck"

func Test(t *testing.T) { TestingT(t) }

var dpy *Display

func init() {
	dpy = NewDisplay()
	Suite(dpy)
	for _, op := range dpy.Outputs {
		Suite(op)
	}
}

func (dpy *Display) TestScreenInfo(c *C) {
	c.Check(dpy.Width, Equals, DefaultScreen.WhitePixel)
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
	defer func() {
		if po == nil {
			dpy.SetPrimary(0)
		} else {
			dpy.SetPrimary(uint32(po.Identify))
		}
	}()

	for _, op := range dpy.Outputs {
		dpy.SetPrimary(uint32(op.Identify))
		<-time.After(time.Millisecond * 50)
		c.Check(dpy.PrimaryRect.Width, Equals, op.Allocation.Width)
		c.Check(dpy.PrimaryRect.Height, Equals, op.Allocation.Height)
	}

	if po != nil {
		dpy.SetPrimary(uint32(po.Identify))
		<-time.After(time.Millisecond * 100)
		c.Check(dpy.PrimaryOutput, Equals, po)
	}

	dpy.SetPrimary(0)
	<-time.After(time.Millisecond * 50)
	c.Check(dpy.PrimaryRect.Width, Equals, dpy.Width)
	c.Check(dpy.PrimaryRect.Height, Equals, dpy.Height)
}

func (op *Output) TestInfo(c *C) {
	c.Check(op.Brightness >= 0 && op.Brightness <= 1, Equals, true)

	find := false
	for _, r := range op.ListModes() {
		if r == op.Mode {
			find = true
		}
	}
	c.Check(find, Equals, true)

	crtcInfo, err := randr.GetCrtcInfo(X, op.crtc, 0).Reply()

	c.Assert(err, Equals, nil)
	c.Check(op.Mode, Equals, buildMode(dpy.modes[crtcInfo.Mode]))

	c.Check(op.Rotation, Equals, uint16(crtcInfo.Rotation))
	c.Check(op.Opened, Equals, op.crtc != 0)
	_, err = randr.GetOutputInfo(X, op.Identify, 0).Reply()
	c.Check(err, Equals, nil)

	op.ListModes()
	op.ListRotations()
	op.updateCrtc(dpy)
}

func (op *Output) TestRotation(c *C) {
	oi, _ := randr.GetOutputInfo(X, op.Identify, 0).Reply()
	/*ci, _ := randr.GetCrtcInfo(X, op.crtc, 0).Reply()*/
	if op.Name == "LVDS1" {
		return
	}
	return

	s, err := randr.SetCrtcConfig(X, op.crtc, 0, 0, op.Allocation.X, op.Allocation.Y, op.bestMode, 2, []randr.Output{op.Identify}).Reply()
	fmt.Println("Rotation....:", s, err, op.crtc, op.Allocation, op.bestMode)
	fmt.Println("SupportCrtcs:", oi.Crtcs, op.crtc)
	fmt.Println("SupportModes:", oi.Modes, op.bestMode)
	fmt.Println("DPY:", dpy.Width, dpy.Height)

	return
	for _, reflect := range op.ListReflect() {
		op.setReflect(reflect)
		for _, rotation := range op.ListRotations() {
			op.setRotation(uint16(rotation))
			<-time.After(time.Millisecond * 2000)
		}
	}
	op.setRotation(1)
	op.setReflect(0)
}

func (op *Output) TestClose(c *C) {
	return
	/*op.setOpened(false)*/
	/*<-time.After(time.Millisecond * 800)*/
	/*op.setOpened(true)*/
	oinfo, err := randr.GetOutputInfo(X, op.Identify, 0).Reply()
	s, err := randr.SetCrtcConfig(X, oinfo.Crtcs[1], 0, 0, op.Allocation.X, op.Allocation.Y, op.bestMode, 1, []randr.Output{op.Identify}).Reply()
	fmt.Println(s, err, oinfo.Crtcs[0])
	<-time.After(time.Millisecond * 100)
}
