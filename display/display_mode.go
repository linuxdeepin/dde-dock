package main

import (
	"fmt"
	"github.com/BurntSushi/xgb/randr"
)

const (
	DisplayModeUnknow  = -100
	DisplayModeMirrors = -1
	DisplayModeCustom  = 0
	DisplayModeOnlyOne = 1
)

func (dpy *Display) SaveConfiguration() {
	dpy.configuration = GenerateCurrentConfig(dpy)
	dpy.configuration.save()
}

func (dpy *Display) SetDisplayMode(mode int16) {
	if mode == dpy.DisplayMode {
		return
	}

	dpy.setPropDisplayMode(mode)

	if dpy.DisplayMode == DisplayModeMirrors {
		for _, op := range dpy.Outputs {
			op.setOpened(true)
			op.SetPos(0, 0)
			op.pendingConfig.SetScale(1, 1)
		}
		w, h := getMirrorSize(DPY.Outputs)
		for _, op := range dpy.Outputs {
			op.ensureSize(w, h, EnsureSizeHintAuto)
		}
		fmt.Println("GetMirrorSize:", w, h)
		fmt.Println("Mirrors mode...")
	} else if dpy.DisplayMode == DisplayModeCustom {
		for _, config := range dpy.configuration.Outputs {
			for _, op := range dpy.Outputs {
				if op.Name == config.Name {
					fmt.Println("OUTPUT:", op.Name, "Enabled:", config.Enabled)
					op.setOpened(config.Enabled)
					if config.Enabled {
						op.SetPos(config.X, config.Y)
						op.pendingConfig.SetScale(1, 1)
						op.setReflect(config.Reflect)
						op.setRotation(config.Rotation)
						op.SetMode(uint32(guestMode(op, config.Width, config.Height, config.RefreshRate)))
					}
					if config.Primary {
						dpy.PrimaryOutput = op
					}
				}
			}
		}
		fmt.Println("Cusstom mode...")
	} else if dpy.DisplayMode >= DisplayModeOnlyOne && int(dpy.DisplayMode) <= len(dpy.Outputs) {
		reserveed := dpy.Outputs[dpy.DisplayMode-1]
		fmt.Println("Reserverd Output:", reserveed.Name)
		for _, op := range dpy.Outputs {
			if op != reserveed {
				op.setOpened(false)
			}
		}
		reserveed.setOpened(true)

		reserveed.SetPos(0, 0)
		reserveed.pendingConfig.SetScale(1, 1)
	}

	dpy.ApplyChanged()
}

func (dpy *Display) ApplyChanged() {
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
		tmpClosedOutput := dpy.adjustScreenSize()
		for _, op := range tmpClosedOutput {
			op.setOpened(true)
		}
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
