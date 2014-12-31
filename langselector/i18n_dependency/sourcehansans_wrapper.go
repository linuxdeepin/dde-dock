/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package i18n_dependency

var (
	purgePkgMap = map[string][]string{
		"zh_CN.UTF-8": []string{
			"fonts-adobe-source-han-sans-tw",
			"fonts-adobe-source-han-sans-jp",
			"fonts-adobe-source-han-sans-kr",
		},
		"zh_TW.UTF-8": []string{
			"fonts-adobe-source-han-sans-cn",
			"fonts-adobe-source-han-sans-jp",
			"fonts-adobe-source-han-sans-kr",
		},
		"zh_HK.UTF-8": []string{
			"fonts-adobe-source-han-sans-cn",
			"fonts-adobe-source-han-sans-jp",
			"fonts-adobe-source-han-sans-kr",
		},
		"ja_JP.UTF-8": []string{
			"fonts-adobe-source-han-sans-cn",
			"fonts-adobe-source-han-sans-tw",
			"fonts-adobe-source-han-sans-kr",
		},
		"ko_KR.UTF-8": []string{
			"fonts-adobe-source-han-sans-cn",
			"fonts-adobe-source-han-sans-tw",
			"fonts-adobe-source-han-sans-jp",
		},
	}
)
