package main

import (
	"dbus/com/deepin/api/graphic"
	"dlib/gio-2.0"
	"fmt"
)

const (
	PersonalizationID string = "com.deepin.dde.personalization"
	CurrentBgKey      string = "current-picture"
)

type Background struct {
	settings   *gio.Settings
	imgHandler *graphic.Graphic
	changed    chan bool
}

func (b *Background) init() error {
	var err error
	b.imgHandler, err = graphic.NewGraphic("com.deepin.api.Graphic", "/com/deepin/api/Graphic")
	if err != nil {
		return err
	}

	b.changed = make(chan bool)
	b.settings = gio.NewSettings(PersonalizationID)
	detailSignal := fmt.Sprintf("changed::%s", CurrentBgKey)
	b.settings.Connect(detailSignal, func(s *gio.Settings, key string, d interface{}) {
		logger.Info(key)
		uri := s.GetString(key)
		logger.Info(uri)
		b.changed <- true
	})

	return nil
}

func (b *Background) currentBg() string {
	pict := b.settings.GetString(CurrentBgKey)

	// status:
	// -1: invalid pict passed, return default pict
	//  0: blur pic
	//  1: original pic
	status, blurPath, err := b.imgHandler.BackgroundBlurPictPath(pict, "", 30, 1)
	if err != nil {
		logger.Info("BackgroundBlurPictPath:", err)
		return DefaultBackgroundImage
	}

	fmt.Printf("status:%d, pict: %s\n", status, blurPath)
	return blurPath
}
