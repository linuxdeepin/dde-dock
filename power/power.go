package main

// #cgo CFLAGS: -DLIBEXECDIR=""
// #cgo amd64 386 CFLAGS: -g -Wall
// #cgo pkg-config:glib-2.0 gtk+-3.0 x11 xext xtst xi gnome-desktop-3.0 upower-glib libnotify libcanberra-gtk3 gudev-1.0
// #cgo LDFLAGS: -lm
// #define GNOME_DESKTOP_USE_UNSTABLE_API
// #include "gnome-idle-monitor.h"
// #include "gsd-power-manager.h"
// int deepin_power_manager_start()
// {
//      GsdPowerManager *manager = gsd_power_manager_new();
//      GError *error = NULL;
//      gtk_init(0,NULL);
//      g_setenv("G_MESSAGES_DEBUG","all",FALSE);
//      notify_init("gsd-power-manager");
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
	"regexp"
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

	schema_gsettings_power               = "com.deepin.daemon.power"
	schema_gsettings_power_settings_id   = "com.deepin.daemon.power.settings"
	schema_gsettings_power_settings_path = "/com/deepin/daemon/power/profiles/"
	schema_gsettings_screensaver         = "org.gnome.desktop.screensaver"
)

const (
	LOGIND_DEST = "org.freedesktop.login1"
	LOGIND_PATH = "/org/freedesktop/login1"
	LOGIND_IFC  = "org.freedesktop.login1.Manager"

	DM_DEST = "org.freedesktop.DisplayManager"
	DM_PATH = "/org/freedesktop/DisplayManager"
	DM_IFC  = "org.freedesktop.DisplayManager"

	DM_SESSION_IFC = "org.freedesktop.DisplayManager.Session"
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

type Power struct {
	//plugins.power keys
	powerProfile  *gio.Settings
	powerSettings *gio.Settings
	mediaKey      *gio.Settings

	//gsettings properties
	ButtonHibernate *property.GSettingsStringProperty `access:"readwrite"`
	ButtonPower     *property.GSettingsStringProperty `access:"readwrite"`
	ButtonSleep     *property.GSettingsStringProperty `access:"readwrite"`
	ButtonSuspend   *property.GSettingsStringProperty `access:"readwrite"`

	CriticalBatteryAction *property.GSettingsStringProperty `access:"read"`
	LidCloseACAction      *property.GSettingsStringProperty `access:"readwrite"`
	LidCloseBatteryAction *property.GSettingsStringProperty `access:"readwrite"`

	//ShowTray *property.GSettingsBoolProperty `access:"readwrite"`

	//SleepDisplayAc      *property.GSettingsIntProperty `access:"readwrite"`
	//SleepDisplayBattery *property.GSettingsIntProperty `access:"readwrite"`
	IdleDelay *property.GSettingsIntProperty `access:"readwrite"`

	SleepInactiveAcTimeout      *property.GSettingsIntProperty `access:"readwrite"`
	SleepInactiveBatteryTimeout *property.GSettingsIntProperty `access:"readwrite"`

	SleepInactiveAcType      *property.GSettingsStringProperty `access:"readwrite"`
	SleepInactiveBatteryType *property.GSettingsStringProperty `access:"readwrite"`

	CurrentProfile *property.GSettingsStringProperty `access:"readwrite"`

	//dbus
	conn      *dbus.Conn
	systemBus *dbus.Conn
	logind    *dbus.Object
	display   *dbus.Object
	dmSession *dbus.Object

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
	screensaverSettings *gio.Settings
	LockEnabled         *property.GSettingsBoolProperty `access:"readwrite"`

	//states to keep track of changes
	LidIsPresent bool
	LidIsClosed  bool
	BatteryIsLow bool
}

func NewPower() (*Power, error) {
	power := &Power{}

	power.powerProfile = gio.NewSettings(schema_gsettings_power)
	power.CurrentProfile = property.NewGSettingsStringProperty(
		power, "CurrentProfile", power.powerProfile,
		"current-profile")

	power.powerSettings = gio.NewSettingsWithPath(
		schema_gsettings_power_settings_id,
		string(schema_gsettings_power_settings_path)+
			power.CurrentProfile.Get()+"/")

	power.screensaverSettings = gio.NewSettings(schema_gsettings_screensaver)
	power.getGsettingsProperty()

	power.upower, _ = upower.NewUpower("/org/freedesktop/UPower")
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
			power.upowerBattery, _ = upower.NewDevice(dbus.ObjectPath(paths[0]))
			if power.upowerBattery != nil {
				power.getUPowerProperty()
			}
		} else {
			println("upower battery interface not found\n")
		}
	}
	var err error
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

	power.dmSession = power.getDMSession()

	return power, nil
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

func (power *Power) upowerChanged() {

	isPresent := power.upower.LidIsPresent.Get()
	if isPresent {
		closed := power.upower.LidIsClosed.Get()
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

	power.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"type='signal',path='/com/deepin/daemon/MediaKey',interface='com.deepin.daemon.MediaKey',member='PowerOff',sender='com.deepin.daemon.KeyBinding'")
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
				break
			case ACTION_INTERACTIVE:
				break
			case ACTION_BLANK:
				break
			case ACTION_LOGOUT:
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

func (power *Power) actionBlank() {

}

func (power *Power) actionLock() {
	fmt.Print("actionLock(): Locking\n")
	power.dmSession.Call(DM_SESSION_IFC+".Lock", 0)
}

func (power *Power) actionLogout() {

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
		call := power.logind.Call(LOGIND_IFC+".Suspend", 0,
			true)
		if call.Err != nil {
			fmt.Println(call.Err)
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
		"com.deepin.daemon.Power",  //bus name
		"/com/deepin/daemon/Power", //object path
		"com.deepin.daemon.Power",
	}
}

func (power *Power) getGsettingsProperty() int32 {
	power.CurrentProfile = property.NewGSettingsStringProperty(
		power, "CurrentProfile", power.powerProfile, "current-profile")
	power.ButtonHibernate = property.NewGSettingsStringProperty(
		power, "ButtonHibernate", power.powerSettings, "button-hibernate")
	power.ButtonPower = property.NewGSettingsStringProperty(
		power, "ButtonPower", power.powerSettings, "button-power")
	power.ButtonSleep = property.NewGSettingsStringProperty(
		power, "ButtonSleep", power.powerSettings, "button-sleep")
	power.ButtonSuspend = property.NewGSettingsStringProperty(
		power, "ButtonSuspend", power.powerSettings, "button-suspend")

	power.CriticalBatteryAction = property.NewGSettingsStringProperty(
		power, "CriticalBatteryAction", power.powerSettings, "critical-battery-action")
	power.LidCloseACAction = property.NewGSettingsStringProperty(
		power, "LidCloseACAction", power.powerSettings, "lid-close-ac-action")
	power.LidCloseBatteryAction = property.NewGSettingsStringProperty(
		power, "LidCloseBatteryAction", power.powerSettings, "lid-close-battery-action")
	//power.ShowTray = property.NewGSettingsBoolProperty(
	//power, "ShowTray", power.powerSettings, "show-tray")
	power.SleepInactiveAcTimeout = property.NewGSettingsIntProperty(
		power, "SleepInactiveAcTimeout", power.powerSettings, "sleep-inactive-ac-timeout")
	power.SleepInactiveBatteryTimeout = property.NewGSettingsIntProperty(
		power, "SleepInactiveBatteryTimeout", power.powerSettings, "sleep-inactive-battery-timeout")
	power.IdleDelay = property.NewGSettingsIntProperty(
		power, "IdleDelay", power.powerSettings, "idle-delay")
	//power.SleepDisplayAc = property.NewGSettingsIntProperty(
	//power, "SleepDisplayAc", power.powerSettings, "sleep-display-ac")
	//power.SleepDisplayBattery = property.NewGSettingsIntProperty(
	//power, "SleepDisplayBattery", power.powerSettings, "sleep-display-battery")

	power.SleepInactiveAcType = property.NewGSettingsStringProperty(
		power, "SleepInactiveAcType", power.powerSettings,
		"sleep-inactive-ac-type")
	power.SleepInactiveBatteryType = property.NewGSettingsStringProperty(
		power, "SleepInactiveBatteryType", power.powerSettings, "sleep-inactive-battery-type")

	power.LockEnabled = property.NewGSettingsBoolProperty(
		power, "LockEnabled", power.screensaverSettings, "lock-enabled")

	return 0
}

func (p *Power) getUPowerProperty() int32 {
	if p.upowerBattery == nil {
		return -1
	}
	p.BatteryIsPresent = property.NewWrapProperty(p, "IsPresent", p.upowerBattery.IsPresent)
	p.IsRechargable = property.NewWrapProperty(p, "IsRechargable", p.upowerBattery.IsRechargeable)
	p.BatteryPercentage = property.NewWrapProperty(p, "BatteryPercentage", p.upowerBattery.Percentage)
	p.TimeToEmpty = property.NewWrapProperty(p, "TimeToEmpty", p.upowerBattery.TimeToEmpty)
	p.TimeToFull = property.NewWrapProperty(p, "TimeToFull", p.upowerBattery.TimeToFull)
	p.Model = property.NewWrapProperty(p, "Model", p.upowerBattery.Model)
	p.Vendor = property.NewWrapProperty(p, "Vendor", p.upowerBattery.Vendor)
	p.State = property.NewWrapProperty(p, "State", p.upowerBattery.State)
	p.Type = property.NewWrapProperty(p, "Type", p.upowerBattery.Type)
	return 1
}

func (power *Power) EnumerateDevices() []dbus.ObjectPath {
	if power.upower == nil {
		println("WARNING:Upower object it nil\n")
	}
	devices, _ := power.upower.EnumerateDevices()
	for _, v := range devices {
		println(v)
	}
	//devices := []dbus.ObjectPath{"testing"}
	return devices
}

func (p *Power) Test() uint32 {
	return 937
}

func getUpowerDeviceObjectPath(devices []dbus.ObjectPath) []dbus.ObjectPath {
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
	dbus.InstallOnSession(power)
	dbus.DealWithUnhandledMessage()
	fmt.Print("power module started,looping")
	go dlib.StartLoop()
	if err := dbus.Wait(); err != nil {
		//l.Error("lost dbus session:", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
