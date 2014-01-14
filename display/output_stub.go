package main

import "dlib/dbus"
import "fmt"

import "github.com/BurntSushi/xgb/randr"
import "github.com/BurntSushi/xgb/xproto"

func (output *Output) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Display",
		fmt.Sprintf("/com/deepin/daemon/Display/Output%d", output.Identify),
		"com.deepin.daemon.Display.Output",
	}
}

func (op *Output) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	switch name {
	case "Rotation":
		op.setRotation(op.Rotation)
	case "Reflect":
		op.setReflect(op.Reflect)
	case "Opened":
		op.setOpened(op.Opened)
	case "Brightness":
		op.setBrightness(op.Brightness)
	}
}

func (op *Output) setPropIdentify(v randr.Output) {
	if op.Identify != v {
		op.Identify = v
		dbus.NotifyChange(op, "Identify")
	}
}

func (op *Output) setPropName(v string) {
	if op.Name != v {
		op.Name = v
		dbus.NotifyChange(op, "Name")
	}
}

func (op *Output) setPropType(v uint8) {
	if op.Type != v {
		op.Type = v
		dbus.NotifyChange(op, "Type")
	}
}

func (op *Output) setPropMode(v Mode) {
	if op.Mode != v {
		op.Mode = v
		dbus.NotifyChange(op, "Mode")
	}
}

func (op *Output) setPropAllocation(v xproto.Rectangle) {
	if op.Allocation != v {
		op.Allocation = v
		dbus.NotifyChange(op, "Allocation")
	}
}

func (op *Output) setPropAdjustMethod(v uint8) {
	if op.AdjustMethod != v {
		op.AdjustMethod = v
		dbus.NotifyChange(op, "AdjustMethod")
	}
}

func (op *Output) setPropRotation(v uint16) {
	if op.Rotation != v {
		op.Rotation = v
		dbus.NotifyChange(op, "Rotation")
	}
}
func (op *Output) setPropReflect(v uint16) {
	if op.Reflect != v {
		op.Reflect = v
		dbus.NotifyChange(op, "Reflect")
	}
}

func (op *Output) setPropOpened(v bool) {
	if op.Opened != v {
		op.Opened = v
		dbus.NotifyChange(op, "Opened")
	}
}
func (op *Output) setPropBrightness(v float64) {
	if op.Brightness != v {
		op.Brightness = v
		dbus.NotifyChange(op, "Brightness")
	}
}
