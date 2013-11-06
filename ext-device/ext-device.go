package main

/*import "fmt"*/
/*import "time"*/
import "dlib/dbus"

type ExtDevice struct {
    UseHabitChanged        func (string, string)
    MoveSpeedChanged       func (string, float64)
    MoveAccuracyChanged    func (string, float64)
    ClickFrequencyChanged  func (string, float64)
    DragDelayChanged       func (string, float64)
}

type ExtProp struct {
    UseHabit        string
    MoveSpeed       float64
    MoveAccuracy    float64
    ClickFrequency  float64
}

type MouseSettings struct {
    Prop ExtProp
}

type TPadSettings struct {
    Prop ExtProp
    DragDelay float64
}

func (ext *ExtDevice) GetDBusInfo () dbus.DBusInfo {
    return dbus.DBusInfo {
        "com.deepin.dss.ExtDevice",
        "/com/deepin/dss/ExtDevice",
        "com.deepin.dss.ExtDevice",
    }
}

func (ext *ExtDevice) SetUseHabit (device, value string) bool {
    if device == "mouse" {
        //set mouse UseHabit
    } else if device == "touchpad" {
        //set touchpad UseHabit
    }
    return true
}

func (ext *ExtDevice) SetMoveSpeed (device string, spd float64) bool {
    if device == "mouse" {
        //set mouse MoveSpeed
    } else if device == "touchpad" {
        //set touchpad MoveSpeed
    }
    return true
}

func (ext *ExtDevice) SetMoveAccuracy (device string, acy float64) bool {
    if device == "mouse" {
        //set mouse MoveAccuracy
    } else if device == "touchpad" {
        //set touchpad MoveAccuracy
    }
    return true
}

func (ext *ExtDevice) SetClickFrequency (device string, f float64) bool {
    if device == "mouse" {
        //set mouse ClickFrequency
    } else if device == "touchpad" {
        //set touchpad ClickFrequency
    }
    return true
}

func (ext *ExtDevice) SetDragDelay (device string, delay float64) bool {
    if device == "touchpad" {
        //set touchpad DragDelay
    }

    return true
}

func (ext *ExtDevice) GetUseHabit (device string) string {

    return ""
}

func (ext *ExtDevice) GetMoveSpeed (device string) float64 {

    return 0
}

func (ext *ExtDevice) GetMoveAccuracy (device string) float64 {

    return 0
}

func (ext *ExtDevice) GetClickFrequency (device string) float64 {

    return 0
}

func (ext *ExtDevice) GetDragDelay (device string) float64 {
    if device != "touchpad" {
        return 0
    }

    return 0
}

func main () {
    m := ExtDevice {}
    dbus.InstallOnSession (&m)
    select {}
    /*timer := time.NewTimer (time.Second * 20)*/
    /*<-timer.C*/
}
