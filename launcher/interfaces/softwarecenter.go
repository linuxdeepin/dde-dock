package interfaces

type SoftwareCenterInterface interface {
	GetPkgNameFromPath(string) (string, error)
	UninstallPkg(string, bool) error
	Connectupdate_signal(func(message [][]interface{})) func()
}
