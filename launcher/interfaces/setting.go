package interfaces

type SettingCoreInterface interface {
	GetEnum(string) int
	SetEnum(string, int) bool
	Connect(string, interface{})
	Unref()
}
