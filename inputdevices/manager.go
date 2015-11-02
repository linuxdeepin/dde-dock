package inputdevices


type devicePathInfo struct {
	Path string
	Type string
}
type devicePathInfos []*devicePathInfo

type Manager struct {
	Infos devicePathInfos

	kbd *Keyboard
	mouse *Mouse
	tpad *Touchpad
	wacom *Wacom
}

func NewManager() *Manager {
	var m = new(Manager)

	m.Infos = devicePathInfos{
		&devicePathInfo{
			Path:"com.deepin.daemon.InputDevice.Keyboard",
			Type:"keyboard",
		},
		&devicePathInfo{
			Path:"com.deepin.daemon.InputDevice.Mouse",
			Type:"mouse",
		},
		&devicePathInfo{
			Path:"com.deepin.daemon.InputDevice.TouchPad",
			Type:"touchpad",
		},
	}

	m.kbd = getKeyboard()
	m.wacom = getWacom()
	m.tpad = getTouchpad()
	m.mouse = getMouse()

	return m
}
