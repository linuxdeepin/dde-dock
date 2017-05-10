package grub2

import (
	"pkg.deepin.io/lib/dbus"
)

var _g *Grub2

func Start() error {
	initPolkit()
	_g = New()
	err := dbus.InstallOnSystem(_g)
	if err != nil {
		return err
	}

	return dbus.InstallOnSystem(_g.theme)
}

func CanSafelyExit() bool {
	return _g.canSafelyExit()
}

// write default config
// write default /etc/default/grub
// generate theme background image file
// call from deepin-installer hooks/in_chroot/50_setup_bootloader_x86.job
func Setup(resolution string) error {
	config := NewConfig()
	config.UseDefault()

	w, h, err := parseResolution(resolution)
	if err != nil {
		return err
	}

	config.Resolution = resolution
	err = config.Save()
	if err != nil {
		return err
	}

	err = writeGrubParams(config)
	if err != nil {
		return err
	}

	return generateThemeBackground(w, h)
	// no run update-grub
}

// call from grub-themes-deepin debian/postinst
func SetupTheme() error {
	config, _ := loadConfig()
	w, h, err := parseResolution(config.Resolution)
	if err != nil {
		// keep background image size
		return nil
	}

	return generateThemeBackground(w, h)
}
