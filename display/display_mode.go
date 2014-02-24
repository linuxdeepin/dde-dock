package main

import (
	"fmt"
	"github.com/BurntSushi/xgb/randr"
)

const (
	DisplayModeMirrors = -1
	DisplayModeCustom  = 0
	DisplayModeOnlyOne = 1
)

func (dpy *Display) SetDisplayMode(mode int16) {
	// TODO: rewrite
	dpy.setPropDisplayMode(mode)
	dpy.ApplyChanged()
}

func (dpy *Display) ApplyChanged() {
	if dpy.DisplayMode == DisplayModeMirrors {
		w, h := getMirrorSize(DPY.Outputs)
		builtIn := guestBuiltIn(dpy.Outputs)
		dpy.ApplyChanged2()
		for _, op := range dpy.Outputs {
			op.setOpened(true)
			op.SetPos(0, 0)
			if op != builtIn {
				op.EnsureSize(w, h, EnsureSizeHintAuto)
			}
		}
		dpy.ApplyChanged2()
		fmt.Println("GetMirrorSize:", w, h)
		fmt.Println("Mirrors mode...")
	} else if dpy.DisplayMode == DisplayModeCustom {
		x := int16(0)
		for _, op := range dpy.Outputs {
			op.setOpened(true)
			op.SetPos(x, 0)
			fmt.Println("Set:", op.Identify, x, 0)
			x += int16(op.pendingAllocation().Width)
		}
		dpy.ApplyChanged2()
		fmt.Println("Cusstom mode...")
	} else if dpy.DisplayMode >= DisplayModeOnlyOne && int(dpy.DisplayMode) <= len(dpy.Outputs) {
		reserveed := dpy.Outputs[dpy.DisplayMode-1]
		reserveed.setOpened(true)
		reserveed.SetPos(0, 0)
		for _, op := range dpy.Outputs {
			if op != reserveed {
				op.setOpened(false)
			}
		}
		dpy.ApplyChanged2()
	}
}

func (dpy *Display) ApplyChanged2() {
	changeLock()
	defer func() {
		changeUnlock()
		if err := recover(); err != nil {
			var buf []byte
			/*runtime.Stack(buf, true)*/
			fmt.Println("***************************************************ApplyChanged Panic:", err, buf)
		}
	}()
	nothingWillChange := true
	for _, op := range dpy.Outputs {
		if op.pendingConfig != nil && op.pendingConfig.mask != 0 {
			nothingWillChange = false
			break
		}
	}
	if nothingWillChange {
		return
	}

	dpy.stopListen()
	defer dpy.startListen()

	tmpClosedOutput := dpy.adjustScreenSize()

	for _, op := range dpy.Outputs {
		if op.pendingConfig != nil {
			if err := op.pendingConfig.apply(); err != nil {
				panic(fmt.Sprintln("Apply", op.Name, "failed", err))
				fmt.Println("Apply", op.Name, "failed", err)
			}
			op.pendingConfig = nil
			fmt.Println("Clearn config...", op.Name)
		}
	}

	for _, op := range tmpClosedOutput {
		op.setOpened(true)
	}

	if dpy.PrimaryOutput != nil {
		randr.SetOutputPrimary(X, Root, dpy.PrimaryOutput.Identify)
	} else {
		randr.SetOutputPrimary(X, Root, 0)
	}

}
