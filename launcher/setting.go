package launcher

import (
	"errors"
	"fmt"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	. "pkg.deepin.io/dde/daemon/launcher/setting"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gio-2.0"
	"sync"
)

// Setting存储launcher相关的设置。
type Setting struct {
	core SettingCoreInterface
	lock sync.Mutex

	categoryDisplayMode CategoryDisplayMode
	// CategoryDisplayModeChanged当分类的显示模式改变后触发。
	CategoryDisplayModeChanged func(int64)

	sortMethod SortMethod
	// SortMethodChanged在排序方式改变后触发。
	SortMethodChanged func(int64)
}

func NewSetting(core SettingCoreInterface) (*Setting, error) {
	if core == nil {
		return nil, errors.New("get failed")
	}
	s := &Setting{
		core: core,
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

	// at least one read operation must be called after signal connected, otherwise,
	// the signal connection won't work from glib 2.43.
	// NB: https://github.com/GNOME/glib/commit/8ff5668a458344da22d30491e3ce726d861b3619
	s.categoryDisplayMode = CategoryDisplayMode(core.GetEnum(CategoryDisplayModeKey))
	s.sortMethod = SortMethod(core.GetEnum(SortMethodkey))

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

// GetCategoryDisplayMode获取launcher当前的分类显示模式。
func (s *Setting) GetCategoryDisplayMode() int64 {
	s.lock.Lock()
	defer s.lock.Unlock()
	return int64(s.categoryDisplayMode)
}

// SetCategoryDisplayMode设置launcher的分类显示模式。
func (s *Setting) SetCategoryDisplayMode(newMode int64) {
	if CategoryDisplayMode(newMode) != s.categoryDisplayMode {
		s.core.SetEnum(CategoryDisplayModeKey, int32(newMode))
	}
}

// GetSortMethod获取launcher当前的排序模式。
func (s *Setting) GetSortMethod() int64 {
	s.lock.Lock()
	defer s.lock.Unlock()

	return int64(s.sortMethod)
}

// SetSortMethod设置launcher的排序方法。
func (s *Setting) SetSortMethod(newMethod int64) {
	if SortMethod(newMethod) != s.sortMethod {
		s.core.SetEnum(SortMethodkey, int32(newMethod))
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
