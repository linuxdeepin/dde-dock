/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package appinfo

import "strings"

func NormalizeAppID(candidateID string) string {
	normalizedAppID := strings.ToLower(NormalizeAppIDWithCaseSensitive(candidateID))
	return normalizedAppID
}

func NormalizeAppIDWithCaseSensitive(candidateID string) string {
	normalizedAppID := strings.Replace(candidateID, "_", "-", -1)
	return normalizedAppID
}
