package launcher

type SettingInterface interface {
	GetCategoryDisplayMode() int64
	SetCategoryDisplayMode(newMode int64)
	GetSortMethod() int64
	SetSortMethod(newMethod int64)
	destroy()
}
