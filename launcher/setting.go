package launcher

import (
	"errors"
	"fmt"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/interfaces"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/setting"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"sync"
)

type Setting struct {
	core SettingCoreInterface
	lock sync.Mutex

	categoryDisplayMode        CategoryDisplayMode
	CategoryDisplayModeChanged func(int64)

	sortMethod        SortMethod
	SortMethodChanged func(int64)
}

func NewSetting(core SettingCoreInterface) (*Setting, error) {
	if core == nil {
		return nil, errors.New("get failed")
	}
	s := &Setting{
		core:                core,
		categoryDisplayMode: CategoryDisplayMode(core.GetEnum(CategoryDisplayModeKey)),
		sortMethod:          SortMethod(core.GetEnum(SortMethodkey)),
	}

	s.listenSettingChange(CategoryDisplayModeKey, func(setting *gio.Settings, key string) {
		_newValue := int64(setting.GetEnum(key))
		newValue := CategoryDisplayMode(_newValue)
		s.lock.Lock()
		defer s.lock.Unlock()
		if newValue != s.categoryDisplayMode {
			s.categoryDisplayMode = newValue
			dbus.Emit(s, "CategoryDisplayModeChanged", _newValue)
		}
	})
	s.listenSettingChange(SortMethodkey, func(setting *gio.Settings, key string) {
		_newValue := int64(setting.GetEnum(key))
		newValue := SortMethod(_newValue)
		s.lock.Lock()
		defer s.lock.Unlock()
		if newValue != s.sortMethod {
			s.sortMethod = newValue
			dbus.Emit(s, "SortMethodChanged", _newValue)
		}
	})

	return s, nil
}

func (d *Setting) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		launcherObject,
		launcherPath,
		"com.deepin.dde.daemon.launcher.Setting",
	}
}

func (s *Setting) listenSettingChange(signalName string, handler func(*gio.Settings, string)) {
	detailSignal := fmt.Sprintf("changed::%s", signalName)
	s.core.Connect(detailSignal, handler)
}

func (s *Setting) GetCategoryDisplayMode() int64 {
	s.lock.Lock()
	defer s.lock.Unlock()
	return int64(s.categoryDisplayMode)
}

func (s *Setting) SetCategoryDisplayMode(newMode int64) {
	if CategoryDisplayMode(newMode) != s.categoryDisplayMode {
		s.core.SetEnum(CategoryDisplayModeKey, int(newMode))
	}
}

func (s *Setting) GetSortMethod() int64 {
	s.lock.Lock()
	defer s.lock.Unlock()

	return int64(s.sortMethod)
}

func (s *Setting) SetSortMethod(newMethod int64) {
	if SortMethod(newMethod) != s.sortMethod {
		s.core.SetEnum(SortMethodkey, int(newMethod))
	}
}

func (s *Setting) destroy() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.core != nil {
		s.core.Unref()
		s.core = nil
	}
	dbus.UnInstallObject(s)
}
