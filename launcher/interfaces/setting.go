package interfaces

type SettingCoreInterface interface {
	GetEnum(string) int32
	SetEnum(string, int32) bool
	Connect(string, interface{})
	Unref()
}
