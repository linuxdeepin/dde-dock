package interfaces

// DStore is interface for deepin store.
type DStore interface {
	GetPkgNameFromPath(string) (string, error)
	UninstallPkg(string, bool) error
	Connectupdate_signal(func(message [][]interface{})) func()
}
