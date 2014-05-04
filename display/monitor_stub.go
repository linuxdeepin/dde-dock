package main

import "dlib/dbus"
import "fmt"
import "strings"

func (m *Monitor) GetDBusInfo() dbus.DBusInfo {
	name := strings.Replace(m.Name, "-", "_", -1)
	name = strings.Replace(name, joinSeparator, "_", -1)
	return dbus.DBusInfo{
		"com.deepin.daemon.Display",
		fmt.Sprintf("/com/deepin/daemon/Display/Monitor%s", name),
		"com.deepin.daemon.Display.Monitor",
	}
}

func (m *Monitor) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	switch name {
	case "Rotation":
	case "Reflect":
	case "Opened":
	}
}

func (m *Monitor) setPropCurrentMode(v Mode) {
	if m.CurrentMode != v {
		m.CurrentMode = v
		dbus.NotifyChange(m, "CurrentMode")
	}
}

func (m *Monitor) setPropBestMode(v Mode) {
	if m.BestMode != v {
		m.BestMode = v
		dbus.NotifyChange(m, "BestMode")
	}
}

func (m *Monitor) setPropRotation(v uint16) {
	if m.Rotation != v {
		m.Rotation = v
		dbus.NotifyChange(m, "Rotation")
	}
}
func (m *Monitor) setPropReflect(v uint16) {
	if m.Reflect != v {
		m.Reflect = v
		dbus.NotifyChange(m, "Reflect")
	}
}

func (m *Monitor) setPropOpened(v bool) {
	if m.Opened != v {
		m.Opened = v
		dbus.NotifyChange(m, "Opened")
	}
}

func (m *Monitor) setPropWidth(v uint16) {
	if m.Width != v {
		m.Width = v
		dbus.NotifyChange(m, "Width")
	}
}
func (m *Monitor) setPropHeight(v uint16) {
	if m.Height != v {
		m.Height = v
		dbus.NotifyChange(m, "Height")
	}
}
func (m *Monitor) setPropPos(x, y int16) {
	if m.X != x || m.Y != y {
		m.X, m.Y = x, y
		dbus.NotifyChange(m, "X")
		dbus.NotifyChange(m, "Y")
	}
}

func (m *Monitor) setPropIsComposited(v bool) {
	if m.IsComposited != v {
		m.IsComposited = v
		dbus.NotifyChange(m, "IsComposited")
	}
}

func (m *Monitor) setPropOutputs(v []string) {
	m.Outputs = v
	dbus.NotifyChange(m, "Outputs")
}
