package inputdevices

func (m *Mouse) Reset() {
	for _, key := range m.setting.ListKeys() {
		m.setting.Reset(key)
	}
}

func (tpad *Touchpad) Reset() {
	for _, key := range tpad.setting.ListKeys() {
		tpad.setting.Reset(key)
	}
}

func (w *Wacom) Reset() {
	for _, key := range w.setting.ListKeys() {
		w.setting.Reset(key)
	}
}

func (kbd *Keyboard) Reset() {
	for _, key := range kbd.setting.ListKeys() {
		kbd.setting.Reset(key)
	}
}

func (kbd *Keyboard) LayoutList() map[string]string {
	return kbd.layoutDescMap
}

func (kbd *Keyboard) GetLayoutDesc(layout string) string {
	if len(layout) == 0 {
		return ""
	}

	desc, ok := kbd.layoutDescMap[layout]
	if !ok {
		return ""
	}

	return desc
}

func (kbd *Keyboard) AddUserLayout(layout string) {
	kbd.addUserLayout(layout)
}

func (kbd *Keyboard) DeleteUserLayout(layout string) {
	kbd.delUserLayout(layout)
}

func (kbd *Keyboard) AddLayoutOption(option string) {
	kbd.addUserOption(option)
}

func (kbd *Keyboard) DeleteLayoutOption(option string) {
	kbd.delUserOption(option)
}

func (kbd *Keyboard) ClearLayoutOption() {
	kbd.UserOptionList.Set([]string{})
}
