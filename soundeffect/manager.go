package soundeffect

import (
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus/property"
	"pkg.deepin.io/lib/gio-2.0"
	"dbus/com/deepin/api/sound"
)

const (
	soundEffectSchema = "com.deepin.dde.sound-effect"

	keyLogin         = "login"
	keyShutdown      = "shutdown"
	keyLogout        = "logout"
	keyWakeup        = "wakeup"
	keyNotification  = "notification"
	keyUnableOperate = "unable-operate"
	keyEmptyTrash    = "empty-trash"
	keyVolumeChange  = "volume-change"
	keyBatteryLow    = "battery-low"
	keyPowerPlug     = "power-plug"
	keyPowerUnplug   = "power-unplug"
	keyDevicePlug    = "device-plug"
	keyDeviceUnplug  = "device-unplug"
	keyIconToDesktop = "icon-to-desktop"
	keyScreenshot    = "screenshot"
)

const (
	dbusDest = "com.deepin.daemon.SoundEffect"
	dbusPath = "/com/deepin/daemon/SoundEffect"
	dbusIFC  = dbusDest
)

type Manager struct {
	Login         *property.GSettingsBoolProperty `access:"readwrite"`
	Shutdown      *property.GSettingsBoolProperty `access:"readwrite"`
	Logout        *property.GSettingsBoolProperty `access:"readwrite"`
	Wakeup        *property.GSettingsBoolProperty `access:"readwrite"`
	Notification  *property.GSettingsBoolProperty `access:"readwrite"`
	UnableOperate *property.GSettingsBoolProperty `access:"readwrite"`
	EmptyTrash    *property.GSettingsBoolProperty `access:"readwrite"`
	VolumeChange  *property.GSettingsBoolProperty `access:"readwrite"`
	BatteryLow    *property.GSettingsBoolProperty `access:"readwrite"`
	PowerPlug     *property.GSettingsBoolProperty `access:"readwrite"`
	PowerUnplug   *property.GSettingsBoolProperty `access:"readwrite"`
	DevicePlug    *property.GSettingsBoolProperty `access:"readwrite"`
	DeviceUnplug  *property.GSettingsBoolProperty `access:"readwrite"`
	IconToDesktop *property.GSettingsBoolProperty `access:"readwrite"`
	Screenshot    *property.GSettingsBoolProperty `access:"readwrite"`

	player *sound.Sound
	setting *gio.Settings
}

func NewManager() (*Manager, error) {
	var m = new(Manager)

	var err error
	m.player, err = sound.NewSound("com.deepin.api.Sound",
		"/com/deepin/api/Sound")
	if err != nil {
		return nil, err
	}

	m.setting = gio.NewSettings(soundEffectSchema)
	m.Login = property.NewGSettingsBoolProperty(
		m, "Login",
		m.setting, keyLogin)
	m.Shutdown = property.NewGSettingsBoolProperty(
		m, "Shutdown",
		m.setting, keyShutdown)
	m.Logout = property.NewGSettingsBoolProperty(
		m, "Logout",
		m.setting, keyLogout)
	m.Wakeup = property.NewGSettingsBoolProperty(
		m, "Wakeup",
		m.setting, keyWakeup)
	m.Notification = property.NewGSettingsBoolProperty(
		m, "Notification",
		m.setting, keyNotification)
	m.UnableOperate = property.NewGSettingsBoolProperty(
		m, "UnableOperate",
		m.setting, keyUnableOperate)
	m.EmptyTrash = property.NewGSettingsBoolProperty(
		m, "EmptyTrash",
		m.setting, keyEmptyTrash)
	m.VolumeChange = property.NewGSettingsBoolProperty(
		m, "VolumeChange",
		m.setting, keyVolumeChange)
	m.BatteryLow = property.NewGSettingsBoolProperty(
		m, "BatteryLow",
		m.setting, keyBatteryLow)
	m.PowerPlug = property.NewGSettingsBoolProperty(
		m, "PowerPlug",
		m.setting, keyPowerPlug)
	m.PowerUnplug = property.NewGSettingsBoolProperty(
		m, "PowerUnplug",
		m.setting, keyPowerUnplug)
	m.DevicePlug = property.NewGSettingsBoolProperty(
		m, "DevicePlug",
		m.setting, keyDevicePlug)
	m.DeviceUnplug = property.NewGSettingsBoolProperty(
		m, "DeviceUnplug",
		m.setting, keyDeviceUnplug)
	m.IconToDesktop = property.NewGSettingsBoolProperty(
		m, "IconToDesktop",
		m.setting, keyIconToDesktop)
	m.Screenshot = property.NewGSettingsBoolProperty(
		m, "Screenshot",
		m.setting, keyScreenshot)

	return m, nil
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) handleGSetting() {
	m.setting.Connect("changed::shutdown", func(s *gio.Settings, key string){
		// TODO: write to users config file
	})
}
