package main

import fmtp "github.com/kr/pretty"
import "testing"
import . "launchpad.net/gocheck"
import "github.com/BurntSushi/xgb/randr"
import "github.com/BurntSushi/xgb"

func Test(t *testing.T) { TestingT(t) }

func init() {
	X, _ = xgb.NewConn()
	randr.Init(X)
	randr.QueryVersion(X, 1, 3).Reply()
	initDisplay()
	ops := []randr.Output{0x45, 0x43}
	m := NewMonitor(ops)
	Suite(m)
	Suite(DPY)
}

func (dpy *Display) TestInfo(c *C) {
	/*fmtp.Println("DPY:", dpy)*/
}

func (m *Monitor) TestInfo(c *C) {
	fmtp.Println("Monitor:", m)
}

func (dpy *Display) TestTryJoin(c *C) {
	dpy.tryJoin("a_b")
}
