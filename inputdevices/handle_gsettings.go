package inputdevices

import (
	"gir/gio-2.0"
)

func (kbd *Keyboard) handleGSettings() {
	kbd.setting.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case kbdKeyRepeatEnable, kbdKeyRepeatDelay,
			kbdKeyRepeatInterval:
			kbd.setRepeat()
		case kbdKeyLayout:
			kbd.setLayout()
		case kbdKeyCursorBlink:
			kbd.setCursorBlink()
		case kbdKeyLayoutOptions:
			kbd.setOptions()
		case kbdKeyUserLayoutList:
			kbd.setGreeterLayoutList()
		}
	})
}

func (m *Mouse) handleGSettings() {
	m.setting.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case mouseKeyLeftHanded:
			m.enableLeftHanded()
		case mouseKeyDisableTouchpad:
			m.disableTouchpad()
		case mouseKeyNaturalScroll:
			m.enableNaturalScroll()
		case mouseKeyMiddleButton:
			m.enableMidBtnEmu()
		case mouseKeyAcceleration:
			m.motionAcceleration()
		case mouseKeyThreshold:
			m.motionThreshold()
		case mouseKeyDoubleClick:
			m.doubleClick()

			// Sync tpad and mouse double clicck time
			var tpad = getTouchpad()
			if tpad.DoubleClick.Get() == m.DoubleClick.Get() {
				return
			}
			tpad.DoubleClick.Set(m.DoubleClick.Get())
		case mouseKeyDragThreshold:
			m.dragThreshold()
		}
	})
}

func (tpad *Touchpad) handleGSettings() {
	tpad.setting.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case tpadKeyEnabled:
			tpad.enable(tpad.TPadEnable.Get())
		case tpadKeyLeftHanded:
			tpad.enableLeftHanded()
			tpad.enableTapToClick()
		case tpadKeyTapClick:
			tpad.enableTapToClick()
		case tpadKeyNaturalScroll:
			tpad.enableNaturalScroll()
		case tpadKeyScrollDelta:
			tpad.setScrollDistance()
		case tpadKeyEdgeScroll:
			tpad.enableEdgeScroll()
		case tpadKeyVertScroll, tpadKeyHorizScroll:
			tpad.enableTwoFingerScroll()
		case tpadKeyWhileTyping:
			tpad.disableWhileTyping()
		case tpadKeyDoubleClick:
			tpad.doubleClick()

			var m = getMouse()
			if tpad.DoubleClick.Get() == m.DoubleClick.Get() {
				return
			}
			m.DoubleClick.Set(tpad.DoubleClick.Get())
		case tpadKeyDragThreshold:
			tpad.dragThreshold()
		case tpadKeyAcceleration:
			tpad.motionAcceleration()
		case tpadKeyThreshold:
			tpad.motionThreshold()
		}
	})
}

func (w *Wacom) handleGSettings() {
	w.setting.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case wacomKeyLeftHanded:
			w.enableLeftHanded()
		case wacomKeyCursorMode:
			w.enableCursorMode()
		case wacomKeyUpAction:
			w.setKeyAction(btnNumUpKey, w.KeyUpAction.Get())
		case wacomKeyDownAction:
			w.setKeyAction(btnNumDownKey, w.KeyDownAction.Get())
		case wacomKeyDoubleDelta:
			w.setClickDelta()
		case wacomKeyPressureSensitive:
			w.setPressure()
		}
	})
}
