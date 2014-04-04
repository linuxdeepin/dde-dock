package main

// #cgo CFLAGS: -DLIBEXECDIR=""
// #cgo amd64 386 CFLAGS: -g
// #cgo pkg-config:glib-2.0 gtk+-3.0 x11 xext xtst xi upower-glib libnotify libcanberra-gtk3 gudev-1.0 xrandr
// #cgo LDFLAGS: -lm
// #define GNOME_DESKTOP_USE_UNSTABLE_API
// #include "libgnome-desktop/gnome-idle-monitor.h"
// #include "gsd-power-manager.h"
// #include "power-force-idle.h"
// GsdPowerManager *manager;
// int deepin_power_manager_start()
// {
//      manager = gsd_power_manager_new();
//      GError *error = NULL;
//      gtk_init(0,NULL);
//      g_setenv("G_MESSAGES_DEBUG","all",FALSE);
//      notify_init("deepin power manager");
//      XInitThreads();
//      gsd_power_manager_start(manager,&error);
//      return 0;
// }
import "C"

import (
    "dbus/org/freedesktop/upower"
    "dlib"
    "dlib/dbus"
    "dlib/dbus/property"
    "dlib/gio-2.0"
    //"dlib/logger"
    "fmt"
    "os"
    "os/user"
    //"reflect"
    "reflect"
    "regexp"
    "time"
    //"unsafe"
    "github.com/BurntSushi/xgb"
    "github.com/BurntSushi/xgb/randr"
    "github.com/BurntSushi/xgb/xproto"
)

type dbusBattery struct {
    bus_name    string
    object_path string

    //device upower.device
}

const (
    power_bus_name    = "com.deepin.daemon.Power"
    power_object_path = "/com/deepin/daemon/Power"
    power_interface   = "com.deepin.daemon.Power"

    schema_gsettings_power                        = "com.deepin.daemon.power"
    schema_gsettings_power_settings_id_specific   = "com.deepin.daemon.power.settings.specific"
    schema_gsettings_power_settings_id_common     = "com.deepin.daemon.power.settings.common"
    schema_gsettings_power_settings_common_path   = "/com/deepin/daemon/power/"
    schema_gsettings_power_settings_specific_path = "/com/deepin/daemon/power/profiles/"
    //schema_gsettings_screensaver         = "org.gnome.desktop.screensaver"
)

const (
    POWER_DEST = "com.deepin.daemon.Power"
    POWER_PATH = "/com/deepin/daemon/Power"
    POWER_IFC  = "com.deepin.daemon.Power"

    ORG_FREEDESKTOP_SS_DEST = "org.freedesktop.ScreenSaver"
    ORG_FREEDESKTOP_SS_PATH = "/org/freedesktop/ScreenSaver"
    ORG_FREEDESKTOP_SS_IFC  = "org.freedesktop.ScreenSaver"

    LOGIND_DEST = "org.freedesktop.login1"
    LOGIND_PATH = "/org/freedesktop/login1"
    LOGIND_IFC  = "org.freedesktop.login1.Manager"

    DM_DEST        = "org.freedesktop.DisplayManager"
    DM_PATH        = "/org/freedesktop/DisplayManager"
    DM_IFC         = "org.freedesktop.DisplayManager"
    DM_SESSION_IFC = "org.freedesktop.DisplayManager.Session"

    DEEPIN_SESSION_DEST = "com.deepin.SessionManager"
    DEEPIN_SESSION_PATH = "/com/deepin/SessionManager"
    DEEPIN_SESSION_IFC  = "com.deepin.SessionManager"
)

//var l = logger.NewLogger("power")

const (
    MEDIA_KEY_DEST = "com.deepin.daemon.KeyBinding"
    MEDIA_KEY_PATH = "/com/deepin/daemon/MediaKey"
    MEDIA_KEY_IFC  = "com.deepin.daemon.MediaKey"

    MEDIA_KEY_SCHEMA_ID = "com.deepin.dde.key-binding.mediakey"

    SIGNAL_POWER     = "PowerOff"
    SIGNAL_SUSPEND   = "Suspend"
    SIGNAL_SLEEP     = "Sleep"
    SIGNAL_HIBERNATE = "Hibernate"
)

const (
    ACTION_POWEROFF    = "shutdown"
    ACTION_SUSPEND     = "suspend"
    ACTION_INTERACTIVE = "interactive"
    ACTION_NOTHING     = "nothing"
    ACTION_LOGOUT      = "logout"
    ACTION_HIBERNATE   = "hibernate"
    ACTION_BLANK       = "blank"
)

type Inhibitor struct {
    appName string
    reason  string
}

type ScreenSaver struct {
    cookies     map[int]Inhibitor
    timers      map[int]*time.Timer
    n           uint32
    IsInhibited bool
}

func (ss *ScreenSaver) GetDBusInfo() dbus.DBusInfo {
    return dbus.DBusInfo{
        ORG_FREEDESKTOP_SS_DEST, //bus name
        ORG_FREEDESKTOP_SS_PATH, //object path
        ORG_FREEDESKTOP_SS_IFC,
    }
}

func NewScreenSaver() (*ScreenSaver, error) {
    ss := &ScreenSaver{}
    ss.cookies = make(map[int]Inhibitor)
    ss.timers = make(map[int]*time.Timer)
    ss.n = 0
    ss.IsInhibited = false

    return ss, nil
}

func (ss *ScreenSaver) Inhibit(appName, reason string) uint32 {
    newIn := Inhibitor{appName, reason}
    fmt.Println("New inhibtor: ", newIn)
    for key, v := range ss.cookies {
        if reflect.DeepEqual(newIn, v) {
            fmt.Println(newIn, ",", v, ",", ss.cookies)
            return uint32(key)
        }
    }
    res := ss.n
    cookie := ss.n
    ss.cookies[int(ss.n)] = newIn
    ss.timers[int(ss.n)] = time.NewTimer(time.Minute)
    ss.n = ss.n + 1
    fmt.Println(ss.cookies)
    if len(ss.cookies) > 0 {
        ss.IsInhibited = true
        dbus.NotifyChange(ss, "IsInhibited")
    }
    go func() {
        <-ss.timers[int(cookie)].C
        ss.UnInhibit(cookie)
    }()
    return res
}

func (ss *ScreenSaver) Tick(cookie uint32) {
    ss.timers[int(cookie)].Reset(time.Duration(time.Minute))
}

func (ss *ScreenSaver) UnInhibit(cookie uint32) {
    delete(ss.cookies, int(cookie))
    fmt.Println("After delete", ":", ss.cookies)
    if len(ss.cookies) == 0 {
        ss.IsInhibited = false
        dbus.NotifyChange(ss, "IsInhibited")
    }
}

type Power struct {
    //plugins.power keys
    powerProfile          *gio.Settings
    powerSettingsSpecific *gio.Settings
    powerSettingsCommon   *gio.Settings
    mediaKey              *gio.Settings

    //gsettings properties
    ButtonHibernate *property.GSettingsStringProperty `access:"readwrite"`
    ButtonPower     *property.GSettingsStringProperty `access:"readwrite"`
    ButtonSleep     *property.GSettingsStringProperty `access:"readwrite"`
    ButtonSuspend   *property.GSettingsStringProperty `access:"readwrite"`

    CriticalBatteryAction *property.GSettingsStringProperty `access:"read"`
    LidCloseACAction      *property.GSettingsStringProperty `access:"readwrite"`
    LidCloseBatteryAction *property.GSettingsStringProperty `access:"readwrite"`

    IdleDelay                   *property.GSettingsIntProperty `access:"readwrite"`
    SleepInactiveAcTimeout      *property.GSettingsIntProperty `access:"readwrite"`
    SleepInactiveBatteryTimeout *property.GSettingsIntProperty `access:"readwrite"`

    SleepInactiveAcType      *property.GSettingsStringProperty `access:"readwrite"`
    SleepInactiveBatteryType *property.GSettingsStringProperty `access:"readwrite"`

    LockEnabled *property.GSettingsBoolProperty `access:"readwrite"`

    CurrentProfile *property.GSettingsStringProperty `access:"readwrite"`

    //dbus
    conn          *dbus.Conn
    systemBus     *dbus.Conn
    logind        *dbus.Object
    display       *dbus.Object
    dmSession     *dbus.Object
    deepinSession *dbus.Object

    //upower interface
    upower *upower.Upower

    //upower battery interface
    upowerBattery     *upower.Device
    BatteryIsPresent  dbus.Property `access:"read` //battery present
    IsRechargable     dbus.Property `access:"read`
    BatteryPercentage dbus.Property `access:"read` //
    Model             dbus.Property `access:"read`
    Vendor            dbus.Property `access:"read`
    TimeToEmpty       dbus.Property `access:"read` //
    TimeToFull        dbus.Property `access:"read` //time to fully charged
    State             dbus.Property `access:"read` //1 for in,2 for out
    Type              dbus.Property `access:"read` //type,2

    //gnome.desktop.screensaver keys
    //screensaverSettings *gio.Settings

    //states to keep track of changes
    LidIsPresent bool
    LidIsClosed  bool
    BatteryIsLow bool
}

func NewPower() (*Power, error) {
    var err error
    power := &Power{}

    power.powerProfile = gio.NewSettings(schema_gsettings_power)
    power.CurrentProfile = property.NewGSettingsStringProperty(
        power, "CurrentProfile", power.powerProfile, "current-profile")
    power.CurrentProfile.ConnectChanged(power.profileChanged)

    power.getPowerSettings()

    //power.screensaverSettings = gio.NewSettings(schema_gsettings_screensaver)

    power.getPowerSettingsProperty()

    power.upower, _ = upower.NewUpower("org.freedesktop.UPower", "/org/freedesktop/UPower")
    if power.upower == nil {
        println("WARNING:UPower not provided by dbus\n")
    } else {

        power.LidIsPresent = power.upower.LidIsPresent.Get()
        power.LidIsClosed = power.upower.LidIsClosed.Get()
        power.BatteryIsLow = power.upower.OnLowBattery.Get()

        power.upower.ConnectChanged(power.upowerChanged)

        println("enumerating devices\n")
        devices, _ := power.upower.EnumerateDevices()
        paths := getUpowerDeviceObjectPath(devices)
        println(paths)
        if len(paths) >= 1 {
            power.upowerBattery, _ = upower.NewDevice("org.freedesktop.UPower", dbus.ObjectPath(paths[0]))
            if power.upowerBattery != nil {
                power.getUPowerProperty()
                power.upowerBattery.ConnectChanged(func() {
                    power.getUPowerProperty()
                    dbus.NotifyChange(power, "BatteryIsPresent")
                    dbus.NotifyChange(power, "BatteryPercentage")
                    dbus.NotifyChange(power, "TimeToEmpty")
                    dbus.NotifyChange(power, "TimeToFull")
                    dbus.NotifyChange(power, "State")
                    dbus.NotifyChange(power, "Type")
                })

            } else {
                //power.BatteryIsPresent = dbus.Property{}
                //power.IsRechargable = false
                //power.BatteryPercentage = 0
                //power.Model = ""
                //power.Vendor = ""
                //power.TimeToEmpty = 0
                //power.TimeToFull = 0
                //power.State = 0
                //power.Type = 0

            }
        } else {
            println("upower battery interface not found\n")
        }
    }

    power.conn, err = dbus.SessionBus()
    if err != nil {
        fmt.Print(os.Stderr, "Failed to connect to session bus")
    }

    power.systemBus, err = dbus.SystemBus()
    if err != nil {
        fmt.Print(os.Stderr, "Failed to connect to system bus")
    }

    power.logind = power.systemBus.Object(LOGIND_DEST, LOGIND_PATH)
    power.engineButton()
    power.listenSleep()

    power.dmSession = power.getDMSession()
    power.deepinSession = power.getDeepinSession()

    return power, nil
}

func (power *Power) OnPropertiesChanged(name string, oldv interface{}) {
    defer func() {
        if err := recover(); err != nil {
            fmt.Print(err)
        }
    }()
    switch name {
    case "CurrentProfile":
        //fmt.Println("sleep inactive ac timeout: ", power.SleepInactiveAcTimeout.Get())
        //power.powerSettings = power.getPowerSettings()
        //power.getPowerSettingsProperty()
        //dbus.InstallOnSession(power)
        break
    }
    dbus.NotifyChange(power, name)
}

func (power *Power) getPowerSettings() {
    power.powerSettingsCommon = gio.NewSettingsWithPath(
        schema_gsettings_power_settings_id_common,
        string(schema_gsettings_power_settings_common_path))
    power.powerSettingsSpecific = gio.NewSettingsWithPath(
        schema_gsettings_power_settings_id_specific,
        string(schema_gsettings_power_settings_specific_path)+
            power.CurrentProfile.Get()+"/")

}

func (power *Power) getDMSession() *dbus.Object {
    dm := power.systemBus.Object(DM_DEST, DM_PATH)
    sessions, err := dm.GetProperty(DM_IFC + ".Sessions")
    if err != nil {
        panic(err)
    }
    //_sessions := reflect.ValueOf(sessions.Value()).
    var _sessions []dbus.ObjectPath = sessions.Value().([]dbus.ObjectPath)
    for _, value := range _sessions {
        obj := power.systemBus.Object(DM_DEST, value)
        username, err := obj.GetProperty(DM_SESSION_IFC + ".UserName")
        if err != nil {
            panic(err)
        }
        user, err := user.Current()
        if err != nil {
            panic(err)
        }
        myusername := user.Username
        fmt.Println("username: ", username.String(), ",", myusername)
        if "\""+myusername+"\"" == username.String() {
            fmt.Println("username: equals", username.String(), ",", myusername)
            return obj
        }
    }

    return nil
}

func (power *Power) getDeepinSession() *dbus.Object {
    obj := power.conn.Object(DEEPIN_SESSION_DEST, DEEPIN_SESSION_PATH)

    return obj
}

func (power *Power) profileChanged() {
    defer func() {
        if err := recover(); err != nil {
            fmt.Print(err)
        }
    }()
    //name := power.CurrentProfile.Get()
    //switch name {
    //case "CurrentProfile":
    fmt.Println("sleep inactive ac timeout: ", power.SleepInactiveAcTimeout.Get())
    power.getPowerSettings()
    power.getPowerSettingsProperty()
    //dbus.InstallOnSession(power)
    dbus.NotifyChange(power, "CurrentProfile")
    //break
    //}

}

func (power *Power) upowerChanged() {

    isPresent := power.upower.LidIsPresent.Get()
    if isPresent {
        closed := power.upower.LidIsClosed.Get()
        fmt.Println("lid is closed ? ", closed)
        if closed == power.LidIsClosed {
            return
        } else {
            power.LidIsClosed = closed
        }
        if closed {
            power.doLidCloseAction()
        } else {
            power.doLidOpenAction()
        }
    }
}

func (power *Power) listenSleep() {
    power.systemBus.BusObject().Call("org.freedesktop.DBus.AddMatch",
        0, "type='signal',path='/org/freedesktop/login1', interface='org.freedesktop.login1.Manager',member='PrepareForSleep'")

    go func() {
        c := make(chan *dbus.Signal, 10)
        power.systemBus.Signal(c)
        for v := range c {
            fmt.Println(v)
            switch v.Body[0].(bool) {
            case true:
                fmt.Println("Preparing to sleep")
                if power.LockEnabled.Get() {
                    //power.actionLock()
                }
                break
            case false:
                fmt.Println("Preparing to wake from  sleep")
                power.actionLock()
                break
            }
        }
    }()
}

func (power *Power) externalMonitorIsOn() bool {
    X, _ := xgb.NewConn()

    err := randr.Init(X)
    if err != nil {
        panic(err)
    }

    //root window on the default screen
    root := xproto.Setup(X).DefaultScreen(X).Root

    resources, err := randr.GetScreenResources(X, root).Reply()
    if err != nil {
        panic(err)
    }

    var on int = 0
    for _, output := range resources.Outputs {
        info, err := randr.GetOutputInfo(X, output, 0).Reply()
        if err != nil {
            panic(err)
        }

        if info.Connection == randr.ConnectionConnected {
            on += 1
        }

        if on >= 2 {
            return true
        }
    }

    return false
}

func (power *Power) doLidCloseAction() {
    battery := power.upower.OnBattery.Get()
    var action string
    if battery {
        action = power.LidCloseBatteryAction.Get()
    } else {
        action = power.LidCloseACAction.Get()
    }

    fmt.Println("lid is closed: ", action)
    switch action {
    case ACTION_NOTHING:
        power.actionNothing()
        break
    case ACTION_BLANK:
        power.actionBlank()
        break
    case ACTION_LOGOUT:
        power.actionLogout()
        break
    case ACTION_SUSPEND:
        if power.externalMonitorIsOn() {
            break
        } else {
            if power.LockEnabled.Get() {
                go power.actionLock()
            }
            power.actionSuspend()
        }
        break
    case ACTION_POWEROFF:
        power.actionPowerOff()
        break
    case ACTION_HIBERNATE:
        power.actionHibernate()
        break
    }
}

func (power *Power) doLidOpenAction() {
    fmt.Println("lid is open")
}

func (power *Power) engineButton() {

    power.conn.BusObject().Call("org.freedesktop.DBus.AddMatch",
        0, "type='signal',path='/com/deepin/daemon/MediaKey',interface='com.deepin.daemon.MediaKey',member='PowerOff',sender='com.deepin.daemon.KeyBinding'")
    power.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
        "type='signal',sender='com.deepin.daemon.KeyBinding',path='/com/deepin/daemon/MediaKey',interface='com.deepin.daemon.MediaKey',member='PowerOff'")
    go func() {
        c := make(chan *dbus.Signal, 16)
        power.conn.Signal(c)
        for v := range c {
            var action string
            switch v.Name {
            case MEDIA_KEY_IFC + "." + SIGNAL_POWER:
                //power.logind.Call(""
                action = power.ButtonPower.Get()
                break
            case MEDIA_KEY_IFC + "." + SIGNAL_SLEEP:
                action = power.ButtonSleep.Get()
                break
            case MEDIA_KEY_IFC + "." + SIGNAL_HIBERNATE:
                action = power.ButtonHibernate.Get()
                break
            case MEDIA_KEY_IFC + "." + SIGNAL_SUSPEND:
                action = power.ButtonSuspend.Get()
                break
            }

            fmt.Println(v.Name, v.Body, " action: ", action)
            switch action {
            case ACTION_NOTHING:
                power.actionNothing()
                break
            case ACTION_INTERACTIVE:
                power.actionInteractive()
                break
            case ACTION_BLANK:
                break
            case ACTION_LOGOUT:
                power.actionLogout()
                break
            case ACTION_SUSPEND:
                if power.LockEnabled.Get() {
                    go power.actionLock()
                }
                power.actionSuspend()
                break
            case ACTION_POWEROFF:
                power.actionPowerOff()
                break
            case ACTION_HIBERNATE:
                power.actionHibernate()
                break
            }
        }
    }()
    fmt.Println("listening to power events")
}

func (power *Power) actionNothing() {
    if power.deepinSession != nil {
        power.actionInteractive()
    } else {
        fmt.Print("actionNothing() do nothing")
    }
}

func (power *Power) actionInteractive() {
    if power.deepinSession == nil {
        fmt.Print("deepin session doesn't exist,can't choose")
        return
    }
    fmt.Println("Interactive interface now")
    power.deepinSession.Call(DEEPIN_SESSION_IFC+".PowerOffChoose", 0)
}

func (power *Power) actionBlank() {

}

func (power *Power) actionLock() {
    fmt.Print("actionLock(): Locking\n")
    power.dmSession.Call(DM_SESSION_IFC+".Lock", 0).Store()
}

func (power *Power) actionLogout() {
    if power.deepinSession == nil {
        fmt.Print("deepin session doesn't exist,can't logout")
        return
    }

    var can bool
    err := power.deepinSession.Call(DEEPIN_SESSION_IFC+".CanLogout", 0).Store(&can)
    if err != nil {
        fmt.Print(err)
        return
    } else {
        if can {
            power.deepinSession.Call(DEEPIN_SESSION_IFC+".Logout", 0)
        }
    }
}

func (power *Power) actionPowerOff() {
    var can string
    err := power.logind.Call(LOGIND_IFC+".CanPowerOff", 0).Store(&can)
    if err != nil {
        panic(err)
    }
    if can == "yes" {
        call := power.logind.Call(LOGIND_IFC+".PowerOff", 0,
            true)
        if call.Err != nil {
            fmt.Println(call.Err)
        }
    }
}

func (power *Power) actionSuspend() {
    var can string
    err := power.logind.Call(LOGIND_IFC+".CanSuspend", 0).Store(&can)
    if err != nil {
        panic(err)
    }
    if can == "yes" {
        fmt.Print("actionSuspend():suspending\n")
        //power.systemBus.Emit(power.logind., "org.freedesktop.login1.Manager.PrepareForSleep", true)
        call := power.logind.Call(LOGIND_IFC+".Suspend", 0,
            true)
        err = call.Store()
        if err != nil {
            fmt.Println(err)
        }
    }
}

func (power *Power) actionHibernate() {
    var can string
    err := power.logind.Call(LOGIND_IFC+".CanHibernate", 0).Store(&can)
    if err != nil {
        panic(err)
    }
    if can == "yes" {
        call := power.logind.Call(LOGIND_IFC+".Hibernate", 0,
            true)
        if call.Err != nil {
            fmt.Println(call.Err)
        }
    }
}

func (power *Power) actionHybridSleep() {
    var can string
    err := power.logind.Call(LOGIND_IFC+".CanHybridSleep", 0).Store(&can)
    if err != nil {
        panic(err)
    }
    if can == "yes" {
        call := power.logind.Call(LOGIND_IFC+".HybridSleep", 0,
            true)
        if call.Err != nil {
            fmt.Println(call.Err)
        }
    }
}

func (p *Power) GetDBusInfo() dbus.DBusInfo {
    return dbus.DBusInfo{
        POWER_DEST, //bus name
        POWER_PATH, //object path
        POWER_IFC,
    }
}

func (power *Power) getPowerSettingsProperty() int32 {
    power.CurrentProfile = property.NewGSettingsStringProperty(
        power, "CurrentProfile", power.powerProfile, "current-profile")
    power.ButtonHibernate = property.NewGSettingsStringProperty(
        power, "ButtonHibernate", power.powerSettingsCommon, "button-hibernate")
    power.ButtonPower = property.NewGSettingsStringProperty(
        power, "ButtonPower", power.powerSettingsCommon, "button-power")
    power.ButtonSleep = property.NewGSettingsStringProperty(
        power, "ButtonSleep", power.powerSettingsCommon, "button-sleep")
    power.ButtonSuspend = property.NewGSettingsStringProperty(
        power, "ButtonSuspend", power.powerSettingsCommon, "button-suspend")

    power.CriticalBatteryAction = property.NewGSettingsStringProperty(
        power, "CriticalBatteryAction", power.powerSettingsCommon, "critical-battery-action")
    power.LidCloseACAction = property.NewGSettingsStringProperty(
        power, "LidCloseACAction", power.powerSettingsCommon, "lid-close-ac-action")
    power.LidCloseBatteryAction = property.NewGSettingsStringProperty(
        power, "LidCloseBatteryAction", power.powerSettingsCommon, "lid-close-battery-action")
    power.SleepInactiveAcTimeout = property.NewGSettingsIntProperty(
        power, "SleepInactiveAcTimeout", power.powerSettingsSpecific, "sleep-inactive-ac-timeout")

    fmt.Println("settings profile:", power.CurrentProfile.Get(), ",", power.SleepInactiveAcTimeout.Get())

    power.SleepInactiveBatteryTimeout = property.NewGSettingsIntProperty(
        power, "SleepInactiveBatteryTimeout", power.powerSettingsSpecific, "sleep-inactive-battery-timeout")
    power.IdleDelay = property.NewGSettingsIntProperty(
        power, "IdleDelay", power.powerSettingsSpecific, "idle-delay")
    //power.SleepDisplayAc = property.NewGSettingsIntProperty(
    //power, "SleepDisplayAc", power.powerSettings, "sleep-display-ac")
    //power.SleepDisplayBattery = property.NewGSettingsIntProperty(
    //power, "SleepDisplayBattery", power.powerSettings, "sleep-display-battery")

    power.SleepInactiveAcType = property.NewGSettingsStringProperty(
        power, "SleepInactiveAcType", power.powerSettingsSpecific,
        "sleep-inactive-ac-type")
    power.SleepInactiveBatteryType = property.NewGSettingsStringProperty(
        power, "SleepInactiveBatteryType", power.powerSettingsSpecific, "sleep-inactive-battery-type")

    power.LockEnabled = property.NewGSettingsBoolProperty(
        power, "LockEnabled", power.powerSettingsCommon, "lock-enabled")

    power.signalPowerSettingsChange()

    return 0
}

func (p *Power) getUPowerProperty() int32 {
    if p.upowerBattery == nil {
        return -1
    }
    p.BatteryIsPresent = property.NewWrapProperty(p, "IsPresent", p.upowerBattery.IsPresent)
    dbus.NotifyChange(p, "BatteryIsPresent")
    p.IsRechargable = property.NewWrapProperty(p, "IsRechargable", p.upowerBattery.IsRechargeable)
    dbus.NotifyChange(p, "IsRechargable")
    p.BatteryPercentage = property.NewWrapProperty(p, "BatteryPercentage", p.upowerBattery.Percentage)
    dbus.NotifyChange(p, "BatteryPercentage")
    p.TimeToEmpty = property.NewWrapProperty(p, "TimeToEmpty", p.upowerBattery.TimeToEmpty)
    dbus.NotifyChange(p, "TimeToEmpty")
    p.TimeToFull = property.NewWrapProperty(p, "TimeToFull", p.upowerBattery.TimeToFull)
    dbus.NotifyChange(p, "TimeToFull")
    p.Model = property.NewWrapProperty(p, "Model", p.upowerBattery.Model)
    dbus.NotifyChange(p, "Model")
    p.Vendor = property.NewWrapProperty(p, "Vendor", p.upowerBattery.Vendor)
    dbus.NotifyChange(p, "Vendor")
    p.State = property.NewWrapProperty(p, "State", p.upowerBattery.State)
    dbus.NotifyChange(p, "State")
    p.Type = property.NewWrapProperty(p, "Type", p.upowerBattery.Type)
    dbus.NotifyChange(p, "Type")

    return 1
}

func (power *Power) signalPowerSettingsChange() int32 {
    dbus.NotifyChange(power, "CurrentProfile")
    dbus.NotifyChange(power, "ButtonHibernate")
    dbus.NotifyChange(power, "ButtonPower")
    dbus.NotifyChange(power, "ButtonSleep")
    dbus.NotifyChange(power, "ButtonSuspend")
    dbus.NotifyChange(power, "CriticalBatteryAction")
    dbus.NotifyChange(power, "LidCloseACAction")
    dbus.NotifyChange(power, "LidCloseBatteryAction")
    dbus.NotifyChange(power, "SleepInactiveAcTimeout")
    dbus.NotifyChange(power, "SleepInactiveBatteryTimeout")
    dbus.NotifyChange(power, "IdleDelay")
    dbus.NotifyChange(power, "SleepInactiveAcType")
    dbus.NotifyChange(power, "SleepInactiveBatteryType")
    dbus.NotifyChange(power, "LockEnabled")

    return 0
}

//exported functions

func (power *Power) EnumerateDevices() []dbus.ObjectPath {
    if power.upower == nil {
        println("WARNING:Upower object is nil\n")
    }
    devices, _ := power.upower.EnumerateDevices()
    for _, v := range devices {
        println(v)
    }
    //devices := []dbus.ObjectPath{"testing"}
    return devices
}

func (power *Power) StartDim() int32 {
    fmt.Println("Starting dim")
    C.start_dim(C.manager)
    return 0
}

func (power *Power) StopDim() int32 {
    fmt.Println("Stoping dim")
    C.stop_dim(C.manager)
    return 0
}

func getUpowerDeviceObjectPath(devices []dbus.ObjectPath) []dbus.ObjectPath {
    if len(devices) == 0 {
        return nil
    }
    paths := make([]dbus.ObjectPath, len(devices))
    batPattern, err := regexp.Compile(
        "/org/freedesktop/UPower/devices/battery_BAT[[:digit:]]+")
    if err != nil {
        panic(err)
    }
    linePattern, err := regexp.Compile(
        "/org/freedesktop/UPower/devices/line_power_ADP[[:digit:]]+")
    if err != nil {
        panic(err)
    }

    i := 0
    for _, path := range devices {
        ret := batPattern.FindString(string(path))
        println("findString " + ret)
        if ret == "" {
            ret = linePattern.FindString(string(path))
            if ret == "" {
                continue
            } else {
                println("findString " + ret)
                paths[1] = path
                i = i + 1
            }
        } else {
            paths[0] = path
            i = i + 1
        }
    }
    return paths[0:i]
}

func main() {

    C.deepin_power_manager_start()
    power, err := NewPower()
    if err != nil {
        return
    }
    screenSaver, err := NewScreenSaver()
    if err != nil {
        return
    }
    dbus.InstallOnSession(screenSaver)
    dbus.InstallOnSession(power)

    dbus.DealWithUnhandledMessage()
    fmt.Print("power module started,looping")
    go func() {
        if err := dbus.Wait(); err != nil {
            //l.Error("lost dbus session:", err)
            os.Exit(1)
        } else {
            os.Exit(0)
        }
    }()

    dlib.StartLoop()
}

//export
func doPowerActionType(action int) {
}
