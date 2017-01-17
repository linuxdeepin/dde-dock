/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package appearance

import (
	"pkg.deepin.io/lib/locale"
)

// key is locale code
type DefaultFontConfig map[string]FontConfigItem

type FontConfigItem struct {
	Standard  string
	Monospace string `json:"Mono"`
}

func (cfg DefaultFontConfig) Get() (standard, monospace string) {
	languages := locale.GetLanguageNames()
	for _, lang := range languages {
		if item, ok := cfg[lang]; ok {
			return item.Standard, item.Monospace
		}
	}

	defaultItem := cfg["en_US"]
	return defaultItem.Standard, defaultItem.Monospace
}
