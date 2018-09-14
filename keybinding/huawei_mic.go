package keybinding

import (
	"os"
	"path/filepath"

	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.helper.backlight"
	"pkg.deepin.io/lib/pulse"
)

const huaweiMicLedName = "huawei::mic"

type huaweiMicLedWorkaround struct {
	backlightHelper   *backlight.Backlight
	pulseCtx          *pulse.Context
	defaultSourceName string
	defaultSourceMute bool
	quit              chan struct{}
}

func (c *AudioController) initHuaweiMicLedWorkaround(backlightHelper *backlight.Backlight) {
	fileInfo, err := os.Stat(filepath.Join("/sys/class/leds", huaweiMicLedName))
	if err != nil {
		return
	}
	if fileInfo.IsDir() {
		c.huaweiMicLedWorkaround = &huaweiMicLedWorkaround{}
		c.huaweiMicLedWorkaround.backlightHelper = backlightHelper
		c.huaweiMicLedWorkaround.quit = make(chan struct{})
		c.huaweiMicLedWorkaround.initPulseCtx()
	}
}

func (h *huaweiMicLedWorkaround) destroy() {
	close(h.quit)
}

func (h *huaweiMicLedWorkaround) initPulseCtx() {
	h.pulseCtx = pulse.GetContext()
	if h.pulseCtx != nil {
		h.defaultSourceName = h.pulseCtx.GetDefaultSource()
		for _, source := range h.pulseCtx.GetSourceList() {
			if source.Name == h.defaultSourceName {
				h.defaultSourceMute = source.Mute
				h.setHuaWeiMicLed(source.Mute)
				break
			}
		}

		eventChan := make(chan *pulse.Event, 100)
		h.pulseCtx.AddEventChan(eventChan)
		go func() {
			for {
				select {
				case ev := <-eventChan:
					switch ev.Facility {
					case pulse.FacilityServer:
						h.handlePulseServerEvent()
					case pulse.FacilitySource:
						h.handlePulseSourceEvent(ev.Type, ev.Index)
					}
				case <-h.quit:
					return
				}
			}
		}()
	}
}

func (h *huaweiMicLedWorkaround) handlePulseServerEvent() {
	logger.Debug("[Event] server")
	defaultSourceName := h.pulseCtx.GetDefaultSource()
	if h.defaultSourceName != defaultSourceName {
		// default source changed
		sources := h.pulseCtx.GetSourceList()
		for _, source := range sources {
			if source.Name == h.defaultSourceName {
				if h.defaultSourceMute != source.Mute {
					// mute changed
					h.setHuaWeiMicLed(source.Mute)
					h.defaultSourceMute = source.Mute
				}
				break
			}
		}
	}
}

func (h *huaweiMicLedWorkaround) handlePulseSourceEvent(eType int, idx uint32) {
	switch eType {
	case pulse.EventTypeChange, pulse.EventTypeNew:
		logger.Debugf("[Event] source #%d changed", idx)

		source, err := h.pulseCtx.GetSource(idx)
		if err != nil {
			logger.Warning(err)
			return
		}
		if source.Name == h.defaultSourceName {
			// is default source
			if h.defaultSourceMute != source.Mute {
				// mute changed
				h.setHuaWeiMicLed(source.Mute)
				h.defaultSourceMute = source.Mute
			}
		}
	}
}

func (h *huaweiMicLedWorkaround) setHuaWeiMicLed(mute bool) {
	logger.Debug("setHuaWeiMicLed", mute)
	var val int32
	if mute {
		val = 1
	}
	err := h.backlightHelper.SetBrightness(0, backlightTypeKeyboard, huaweiMicLedName, val)
	if err != nil {
		logger.Warning(err)
	}
}
