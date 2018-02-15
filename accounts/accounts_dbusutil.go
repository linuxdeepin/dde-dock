package accounts

func (v *Manager) setPropAllowGuest(value bool) (changed bool) {
	v.PropsMaster.Lock()
	if v.AllowGuest != value {
		v.AllowGuest = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "AllowGuest", value)
	}
	return
}

func (v *Manager) getPropAllowGuest() bool {
	v.PropsMaster.RLock()
	value := v.AllowGuest
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropUserName(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.UserName != value {
		v.UserName = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "UserName", value)
	}
	return
}

func (v *User) getPropUserName() string {
	v.PropsMaster.RLock()
	value := v.UserName
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropFullName(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.FullName != value {
		v.FullName = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "FullName", value)
	}
	return
}

func (v *User) getPropFullName() string {
	v.PropsMaster.RLock()
	value := v.FullName
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropUid(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.Uid != value {
		v.Uid = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Uid", value)
	}
	return
}

func (v *User) getPropUid() string {
	v.PropsMaster.RLock()
	value := v.Uid
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropGid(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.Gid != value {
		v.Gid = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Gid", value)
	}
	return
}

func (v *User) getPropGid() string {
	v.PropsMaster.RLock()
	value := v.Gid
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropHomeDir(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.HomeDir != value {
		v.HomeDir = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "HomeDir", value)
	}
	return
}

func (v *User) getPropHomeDir() string {
	v.PropsMaster.RLock()
	value := v.HomeDir
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropShell(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.Shell != value {
		v.Shell = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Shell", value)
	}
	return
}

func (v *User) getPropShell() string {
	v.PropsMaster.RLock()
	value := v.Shell
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropLocale(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.Locale != value {
		v.Locale = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Locale", value)
	}
	return
}

func (v *User) getPropLocale() string {
	v.PropsMaster.RLock()
	value := v.Locale
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropLayout(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.Layout != value {
		v.Layout = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Layout", value)
	}
	return
}

func (v *User) getPropLayout() string {
	v.PropsMaster.RLock()
	value := v.Layout
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropIconFile(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.IconFile != value {
		v.IconFile = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "IconFile", value)
	}
	return
}

func (v *User) getPropIconFile() string {
	v.PropsMaster.RLock()
	value := v.IconFile
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropDesktopBackgrounds(value []string) {
	v.PropsMaster.Lock()
	v.DesktopBackgrounds = value
	v.PropsMaster.Unlock()
	if v.service != nil {
		v.PropsMaster.NotifyChanged(v, v.service, "DesktopBackgrounds", value)
	}
}

func (v *User) getPropDesktopBackgrounds() []string {
	v.PropsMaster.RLock()
	value := v.DesktopBackgrounds
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropGreeterBackground(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.GreeterBackground != value {
		v.GreeterBackground = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "GreeterBackground", value)
	}
	return
}

func (v *User) getPropGreeterBackground() string {
	v.PropsMaster.RLock()
	value := v.GreeterBackground
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropXSession(value string) (changed bool) {
	v.PropsMaster.Lock()
	if v.XSession != value {
		v.XSession = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "XSession", value)
	}
	return
}

func (v *User) getPropXSession() string {
	v.PropsMaster.RLock()
	value := v.XSession
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropLocked(value bool) (changed bool) {
	v.PropsMaster.Lock()
	if v.Locked != value {
		v.Locked = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "Locked", value)
	}
	return
}

func (v *User) getPropLocked() bool {
	v.PropsMaster.RLock()
	value := v.Locked
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropAutomaticLogin(value bool) (changed bool) {
	v.PropsMaster.Lock()
	if v.AutomaticLogin != value {
		v.AutomaticLogin = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "AutomaticLogin", value)
	}
	return
}

func (v *User) getPropAutomaticLogin() bool {
	v.PropsMaster.RLock()
	value := v.AutomaticLogin
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropSystemAccount(value bool) (changed bool) {
	v.PropsMaster.Lock()
	if v.SystemAccount != value {
		v.SystemAccount = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "SystemAccount", value)
	}
	return
}

func (v *User) getPropSystemAccount() bool {
	v.PropsMaster.RLock()
	value := v.SystemAccount
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropNoPasswdLogin(value bool) (changed bool) {
	v.PropsMaster.Lock()
	if v.NoPasswdLogin != value {
		v.NoPasswdLogin = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "NoPasswdLogin", value)
	}
	return
}

func (v *User) getPropNoPasswdLogin() bool {
	v.PropsMaster.RLock()
	value := v.NoPasswdLogin
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropAccountType(value int32) (changed bool) {
	v.PropsMaster.Lock()
	if v.AccountType != value {
		v.AccountType = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "AccountType", value)
	}
	return
}

func (v *User) getPropAccountType() int32 {
	v.PropsMaster.RLock()
	value := v.AccountType
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropLoginTime(value uint64) (changed bool) {
	v.PropsMaster.Lock()
	if v.LoginTime != value {
		v.LoginTime = value
		changed = true
	}
	v.PropsMaster.Unlock()
	if v.service != nil && changed {
		v.PropsMaster.NotifyChanged(v, v.service, "LoginTime", value)
	}
	return
}

func (v *User) getPropLoginTime() uint64 {
	v.PropsMaster.RLock()
	value := v.LoginTime
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropIconList(value []string) {
	v.PropsMaster.Lock()
	v.IconList = value
	v.PropsMaster.Unlock()
	if v.service != nil {
		v.PropsMaster.NotifyChanged(v, v.service, "IconList", value)
	}
}

func (v *User) getPropIconList() []string {
	v.PropsMaster.RLock()
	value := v.IconList
	v.PropsMaster.RUnlock()
	return value
}

func (v *User) setPropHistoryLayout(value []string) {
	v.PropsMaster.Lock()
	v.HistoryLayout = value
	v.PropsMaster.Unlock()
	if v.service != nil {
		v.PropsMaster.NotifyChanged(v, v.service, "HistoryLayout", value)
	}
}

func (v *User) getPropHistoryLayout() []string {
	v.PropsMaster.RLock()
	value := v.HistoryLayout
	v.PropsMaster.RUnlock()
	return value
}
