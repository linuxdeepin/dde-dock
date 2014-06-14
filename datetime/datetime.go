package datetime

import (
	"dbus/com/deepin/api/setdatetime"
	"dlib/dbus"
	"dlib/dbus/property"
	. "dlib/gettext"
	"dlib/gio-2.0"
	"dlib/logger"
	libutils "dlib/utils"
	"github.com/howeyc/fsnotify"
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

	objUtils         = libutils.NewUtils()
	setDate          *setdatetime.SetDateTime
	zoneWatcher      *fsnotify.Watcher
	Logger           = logger.NewLogger("dde-daemon/datetime")
	changeLocaleFlag = false
)

type Manager struct {
	AutoSetTime      *property.GSettingsBoolProperty `access:"readwrite"`
	Use24HourDisplay *property.GSettingsBoolProperty `access:"readwrite"`
	CurrentTimezone  string
	UserTimezoneList []string
	LocaleListMap    map[string]string
	CurrentLocale    string

	LocaleStatus func(bool, string)

	ntpRunning bool
	quitChan   chan bool
}

func (op *Manager) SetDate(d string) (bool, error) {
	ret, err := setDate.SetCurrentDate(d)
	if err != nil {
		Logger.Warning("Set Date - '%s' Failed: %s\n",
			d, err)
		return false, err
	}
	return ret, nil
}

func (op *Manager) SetTime(t string) (bool, error) {
	ret, err := setDate.SetCurrentTime(t)
	if err != nil {
		Logger.Warning("Set Time - '%s' Failed: %s\n",
			t, err)
		return false, err
	}
	return ret, nil
}

func (op *Manager) TimezoneCityList() map[string]string {
	//return getZoneCityList()
	return zoneCityMap
}

func (op *Manager) SetTimeZone(zone string) bool {
	_, err := setDate.SetTimezone(zone)
	if err != nil {
		Logger.Warning("Set TimeZone - '%s' Failed: %s\n",
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
	if objUtils.IsElementExist(tz, list) {
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
	if !objUtils.IsElementExist(tz, list) {
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

	sendNotify("", "", Tr("Changing system language, please wait"))
	setDate.GenLocale(locale)
	changeLocaleFlag = true
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
	m.LocaleListMap = make(map[string]string)
	m.LocaleListMap = localDescMap

	m.ntpRunning = false
	m.quitChan = make(chan bool)

	return m
}

func Init() {
	var err error

	setDate, err = setdatetime.NewSetDateTime("com.deepin.api.SetDateTime", "/com/deepin/api/SetDateTime")
	if err != nil {
		Logger.Error("New SetDateTime Failed:", err)
		panic(err)
	}

	zoneWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		Logger.Error("New FS Watcher Failed:", err)
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
	Logger.BeginTracing()

	var err error

	Init()

	date := GetManager()
	err = dbus.InstallOnSession(date)
	if err != nil {
		Logger.Fatal("Install Session DBus Failed:", err)
	}

	if date.AutoSetTime.Get() {
		date.setAutoSetTime(true)
	}
	date.listenLocaleChange()
}

func Stop() {
	zoneWatcher.Close()
	dbus.UnInstallObject(GetManager())

	Logger.EndTracing()
}
