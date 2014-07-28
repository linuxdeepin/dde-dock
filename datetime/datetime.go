package datetime

import (
	"dbus/com/deepin/api/setdatetime"
	"github.com/howeyc/fsnotify"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/dbus/property"
	. "pkg.linuxdeepin.com/lib/gettext"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/log"
)

const (
	_DATE_TIME_DEST = "com.deepin.daemon.DateAndTime"
	_DATE_TIME_PATH = "/com/deepin/daemon/DateAndTime"
	_DATA_TIME_IFC  = "com.deepin.daemon.DateAndTime"

	_DATE_TIME_SCHEMA = "com.deepin.dde.datetime"
	_TIME_ZONE_FILE   = "/etc/timezone"
)

var (
	busConn      *dbus.Conn
	dateSettings = gio.NewSettings(_DATE_TIME_SCHEMA)

	setDate          *setdatetime.SetDateTime
	zoneWatcher      *fsnotify.Watcher
	logger           = log.NewLogger(_DATE_TIME_DEST)
	changeLocaleFlag = false
)

type Manager struct {
	AutoSetTime      *property.GSettingsBoolProperty `access:"readwrite"`
	Use24HourDisplay *property.GSettingsBoolProperty `access:"readwrite"`
	CurrentTimezone  string
	UserTimezoneList []string
	CurrentLocale    string

	LocaleStatus func(bool, string)

	ntpRunning bool
	quitChan   chan bool
}

func (op *Manager) SetDate(d string) (bool, error) {
	ret, err := setDate.SetCurrentDate(d)
	if err != nil {
		logger.Warning("Set Date - '%s' Failed: %s\n",
			d, err)
		return false, err
	}
	return ret, nil
}

func (op *Manager) SetTime(t string) (bool, error) {
	ret, err := setDate.SetCurrentTime(t)
	if err != nil {
		logger.Warning("Set Time - '%s' Failed: %s\n",
			t, err)
		return false, err
	}
	return ret, nil
}

func (op *Manager) TimezoneCityList() []zoneCityInfo {
	return zoneInfos
}

func (op *Manager) SetTimeZone(zone string) bool {
	_, err := setDate.SetTimezone(zone)
	if err != nil {
		logger.Warning("Set TimeZone - '%s' Failed: %s\n",
			zone, err)
		return false
	}
	op.setPropName("CurrentTimezone")
	return true
}

func (op *Manager) SyncNtpTime() bool {
	return op.syncNtpTime()
}

func (op *Manager) AddUserTimezoneList(tz string) {
	if !timezoneIsValid(tz) {
		return
	}

	list := dateSettings.GetStrv("user-timezone-list")
	if isElementExist(tz, list) {
		return
	}

	list = append(list, tz)
	dateSettings.SetStrv("user-timezone-list", list)
}

func (op *Manager) DeleteTimezoneList(tz string) {
	if !timezoneIsValid(tz) {
		return
	}

	list := dateSettings.GetStrv("user-timezone-list")
	if !isElementExist(tz, list) {
		return
	}

	tmp := []string{}
	for _, v := range list {
		if v == tz {
			continue
		}
		tmp = append(tmp, v)
	}
	dateSettings.SetStrv("user-timezone-list", tmp)
}

func (op *Manager) SetLocale(locale string) {
	if len(locale) < 1 {
		return
	}

	if op.CurrentLocale == locale {
		return
	}

	sendNotify("", "", Tr("Language changing, please wait"))
	setDate.GenLocale(locale)
	changeLocaleFlag = true
	op.CurrentLocale = locale
	dbus.NotifyChange(op, "CurrentLocale")
}

func (m *Manager) GetLocaleList() []localeInfo {
	return getLocaleInfoList()
}

func NewDateAndTime() *Manager {
	m := &Manager{}

	m.AutoSetTime = property.NewGSettingsBoolProperty(
		m, "AutoSetTime",
		dateSettings, "is-auto-set")
	m.Use24HourDisplay = property.NewGSettingsBoolProperty(
		m, "Use24HourDisplay",
		dateSettings, "is-24hour")

	m.setPropName("CurrentTimezone")
	m.setPropName("UserTimezoneList")
	m.setPropName("CurrentLocale")
	m.listenSettings()
	m.listenZone()
	m.AddUserTimezoneList(m.CurrentTimezone)

	m.ntpRunning = false
	m.quitChan = make(chan bool)

	return m
}

func Init() {
	var err error

	setDate, err = setdatetime.NewSetDateTime("com.deepin.api.SetDateTime", "/com/deepin/api/SetDateTime")
	if err != nil {
		logger.Error("New SetDateTime Failed:", err)
		panic(err)
	}

	zoneWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		logger.Error("New FS Watcher Failed:", err)
		panic(err)
	}
}

var _manager *Manager

func GetManager() *Manager {
	if _manager == nil {
		_manager = NewDateAndTime()
	}

	return _manager
}

func Start() {
	logger.BeginTracing()

	var err error

	Init()

	initZoneInfos()

	date := GetManager()
	err = dbus.InstallOnSession(date)
	if err != nil {
		logger.Fatal("Install Session DBus Failed:", err)
	}

	if date.AutoSetTime.Get() {
		date.setAutoSetTime(true)
	}
	date.listenLocaleChange()
}

func Stop() {
	zoneWatcher.Close()
	dbus.UnInstallObject(GetManager())

	logger.EndTracing()
}
