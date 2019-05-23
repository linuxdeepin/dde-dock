package audio

import (
	"encoding/json"

	"github.com/linuxdeepin/go-dbus-factory/com.deepin.api.soundthemeplayer"
	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/dbus1"
)

const (
	gsKeyAudioVolumeChange       = "audio-volume-change"
	gsKeyCameraShutter           = "camera-shutter"
	gsKeyCompleteCopy            = "complete-copy"
	gsKeyCompletePrint           = "complete-print"
	gsKeyDesktopLogin            = "desktop-login"
	gsKeyDesktopLogout           = "desktop-logout"
	gsKeyDeviceAdded             = "device-added"
	gsKeyDeviceRemoved           = "device-removed"
	gsKeyDialogErrorCritical     = "dialog-error-critical"
	gsKeyDialogError             = "dialog-error"
	gsKeyDialogErrorSerious      = "dialog-error-serious"
	gsKeyMessage                 = "message"
	gsKeyPowerPlug               = "power-plug"
	gsKeyPowerUnplug             = "power-unplug"
	gsKeyPowerUnplugBatteryLow   = "power-unplug-battery-low"
	gsKeyScreenCaptureComplete   = "screen-capture-complete"
	gsKeyScreenCapture           = "screen-capture"
	gsKeySuspendResume           = "suspend-resume"
	gsKeySystemShutdown          = "system-shutdown"
	gsKeyTrashEmpty              = "trash-empty"
	gsKeyXDeepinAppSentToDesktop = "x-deepin-app-sent-to-desktop"
)

type syncSoundEffect struct {
	Enabled                 bool `json:"enabled"`
	AudioVolumeChange       bool `json:"audio_volume_change"`
	CameraShutter           bool `json:"camera_shutter"`
	CompleteCopy            bool `json:"complete_copy"`
	CompletePrint           bool `json:"complete_print"`
	DesktopLogin            bool `json:"desktop_login"`
	DesktopLogout           bool `json:"desktop_logout"`
	DeviceAdded             bool `json:"device_added"`
	DeviceRemoved           bool `json:"device_removed"`
	DialogErrorCritical     bool `json:"dialog_error_critical"`
	DialogError             bool `json:"dialog_error"`
	DialogErrorSerious      bool `json:"dialog_error_serious"`
	Message                 bool `json:"message"`
	PowerPlug               bool `json:"power_plug"`
	PowerUnplug             bool `json:"power_unplug"`
	PowerUnplugBatteryLow   bool `json:"power_unplug_battery_low"`
	ScreenCaptureComplete   bool `json:"screen_capture_complete"`
	ScreenCapture           bool `json:"screen_capture"`
	SuspendResume           bool `json:"suspend_resume"`
	SystemShutdown          bool `json:"system_shutdown"`
	TrashEmpty              bool `json:"trash_empty"`
	XDeepinAppSentToDesktop bool `json:"x_deepin_app_sent_to_desktop"`
}

type syncData struct {
	Version     string           `json:"version"`
	SoundEffect *syncSoundEffect `json:"soundeffect"`
}

type syncConfig struct {
	a *Audio
}

const (
	syncVersion = "1.0"
)

func (sc *syncConfig) Get() (interface{}, error) {
	s := gio.NewSettings(gsSchemaSoundEffect)
	defer s.Unref()
	return &syncData{
		Version: syncVersion,
		SoundEffect: &syncSoundEffect{
			Enabled:                 s.GetBoolean(gsKeyEnabled),
			AudioVolumeChange:       s.GetBoolean(gsKeyAudioVolumeChange),
			CameraShutter:           s.GetBoolean(gsKeyCameraShutter),
			CompleteCopy:            s.GetBoolean(gsKeyCompleteCopy),
			CompletePrint:           s.GetBoolean(gsKeyCompletePrint),
			DesktopLogin:            s.GetBoolean(gsKeyDesktopLogin),
			DesktopLogout:           s.GetBoolean(gsKeyDesktopLogout),
			DeviceAdded:             s.GetBoolean(gsKeyDeviceAdded),
			DeviceRemoved:           s.GetBoolean(gsKeyDeviceRemoved),
			DialogErrorCritical:     s.GetBoolean(gsKeyDialogErrorCritical),
			DialogError:             s.GetBoolean(gsKeyDialogError),
			DialogErrorSerious:      s.GetBoolean(gsKeyDialogErrorSerious),
			Message:                 s.GetBoolean(gsKeyMessage),
			PowerPlug:               s.GetBoolean(gsKeyPowerPlug),
			PowerUnplug:             s.GetBoolean(gsKeyPowerUnplug),
			PowerUnplugBatteryLow:   s.GetBoolean(gsKeyPowerUnplugBatteryLow),
			ScreenCaptureComplete:   s.GetBoolean(gsKeyScreenCaptureComplete),
			ScreenCapture:           s.GetBoolean(gsKeyScreenCapture),
			SuspendResume:           s.GetBoolean(gsKeySuspendResume),
			SystemShutdown:          s.GetBoolean(gsKeySystemShutdown),
			TrashEmpty:              s.GetBoolean(gsKeyTrashEmpty),
			XDeepinAppSentToDesktop: s.GetBoolean(gsKeyXDeepinAppSentToDesktop),
		},
	}, nil
}

func (sc *syncConfig) Set(data []byte) error {
	var info syncData
	err := json.Unmarshal(data, &info)
	if err != nil {
		return err
	}
	soundEffect := info.SoundEffect
	if soundEffect != nil {
		s := gio.NewSettings(gsSchemaSoundEffect)
		s.SetBoolean(gsKeyEnabled, soundEffect.Enabled)
		s.SetBoolean(gsKeyAudioVolumeChange, soundEffect.AudioVolumeChange)
		s.SetBoolean(gsKeyCameraShutter, soundEffect.CameraShutter)
		s.SetBoolean(gsKeyCompleteCopy, soundEffect.CompleteCopy)
		s.SetBoolean(gsKeyCompletePrint, soundEffect.CompletePrint)
		s.SetBoolean(gsKeyDesktopLogin, soundEffect.DesktopLogin)
		sc.syncConfigToSoundThemePlayer(soundEffect.DesktopLogin)
		s.SetBoolean(gsKeyDesktopLogout, soundEffect.DesktopLogout)
		s.SetBoolean(gsKeyDeviceAdded, soundEffect.DeviceAdded)
		s.SetBoolean(gsKeyDeviceRemoved, soundEffect.DeviceRemoved)
		s.SetBoolean(gsKeyDialogErrorCritical, soundEffect.DialogErrorCritical)
		s.SetBoolean(gsKeyDialogError, soundEffect.DialogError)
		s.SetBoolean(gsKeyDialogErrorSerious, soundEffect.DialogErrorSerious)
		s.SetBoolean(gsKeyMessage, soundEffect.Message)
		s.SetBoolean(gsKeyPowerPlug, soundEffect.PowerPlug)
		s.SetBoolean(gsKeyPowerUnplug, soundEffect.PowerUnplug)
		s.SetBoolean(gsKeyPowerUnplugBatteryLow, soundEffect.PowerUnplugBatteryLow)
		s.SetBoolean(gsKeyScreenCaptureComplete, soundEffect.ScreenCaptureComplete)
		s.SetBoolean(gsKeyScreenCapture, soundEffect.ScreenCapture)
		s.SetBoolean(gsKeySuspendResume, soundEffect.SuspendResume)
		s.SetBoolean(gsKeySystemShutdown, soundEffect.SystemShutdown)
		s.SetBoolean(gsKeyTrashEmpty, soundEffect.TrashEmpty)
		s.SetBoolean(gsKeyXDeepinAppSentToDesktop, soundEffect.XDeepinAppSentToDesktop)
		s.Unref()
	}
	return nil
}

func (sc *syncConfig) syncConfigToSoundThemePlayer(enabled bool) error {
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	player := soundthemeplayer.NewSoundThemePlayer(sysBus)
	return player.EnableSoundDesktopLogin(0, enabled)
}
