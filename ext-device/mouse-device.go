package main

/*import "fmt"*/
import "time"
import "dlib/dbus"

type MouseDevice struct {
	MouseUseHabit       string
	MouseMoveSpeed      float64
	MouseMoveAccuracy   float64
	MouseClickFrequency float64

	MouseUseHabitChanged       func(string)
	MouseMoveSpeedChanged      func(float64)
	MouseMoveAccuracyChanged   func(float64)
	MouseClickFrequencyChanged func(float64)
}

func (mouse *MouseDevice) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.dss.ExtDevice",
		"/com/deepin/dss/ExtDevice",
		"com.deepin.dss.ExtDevice",
	}
}

func (mouse *MouseDevice) MouseSetUseHabit(value string) bool {
	if mouse == nil {
		return false
	}

	mouse.MouseUseHabit = value
	return true
}

func (mouse *MouseDevice) MouseSetMoveSpeed(spd float64) bool {
	if mouse == nil {
		return false
	}

	mouse.MouseMoveSpeed = spd
	return true
}

func (mouse *MouseDevice) MouseSetMoveAccuracy(acy float64) bool {
	if mouse == nil {
		return false
	}

	mouse.MouseMoveAccuracy = acy
	return true
}

func (mouse *MouseDevice) MouseSetClickFrequency(f float64) bool {
	if mouse == nil {
		return false
	}

	mouse.MouseClickFrequency = f
	return true
}

func main() {
	m := MouseDevice{}
	dbus.InstallOnSession(&m)
	/*select {}*/
	timer := time.NewTimer(time.Second * 20)
	<-timer.C
}
