package interfaces

// SettingCore is interface for setting.
type SettingCore interface {
	GetEnum(string) int32
	SetEnum(string, int32) bool
	Connect(string, interface{})
	Unref()
}

// Setting is the interface for setting.
type Setting interface {
	GetCategoryDisplayMode() int64
	SetCategoryDisplayMode(newMode int64)
	GetSortMethod() int64
	SetSortMethod(newMethod int64)
	Destroy()
}
