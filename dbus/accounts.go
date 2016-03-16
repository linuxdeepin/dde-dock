package dbus

import (
	"dbus/com/deepin/daemon/accounts"
	"pkg.deepin.io/lib/dbus"
)

const (
	accountsDBusDest = "com.deepin.daemon.Accounts"
	accountsDBusPath = "/com/deepin/daemon/Accounts"
)

func NewAccounts() (*accounts.Accounts, error) {
	return accounts.NewAccounts(accountsDBusDest, accountsDBusPath)
}

func NewUserByName(name string) (*accounts.User, error) {
	m, err := NewAccounts()
	if err != nil {
		return nil, err
	}

	p, err := m.FindUserByName(name)
	if err != nil {
		return nil, err
	}
	return accounts.NewUser(accountsDBusDest, dbus.ObjectPath(p))
}

func NewUserByUid(uid string) (*accounts.User, error) {
	m, err := NewAccounts()
	if err != nil {
		return nil, err
	}

	p, err := m.FindUserById(uid)
	if err != nil {
		return nil, err
	}
	return accounts.NewUser(accountsDBusDest, dbus.ObjectPath(p))
}

func DestroyAccounts(act *accounts.Accounts) {
	if act == nil {
		return
	}
	accounts.DestroyAccounts(act)
}

func DestroyUser(u *accounts.User) {
	if u == nil {
		return
	}
	accounts.DestroyUser(u)
}
