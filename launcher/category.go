/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package launcher

import (
	"fmt"
	"strings"
)

type CategoryID int64

const (
	CategoryInternet CategoryID = iota
	CategoryChat
	CategoryMusic
	CategoryVideo
	CategoryGraphics
	CategoryGame
	CategoryOffice
	CategoryReading
	CategoryDevelopment
	CategorySystem
	CategoryOthers
)

func (cid CategoryID) String() string {
	var prefix string
	switch cid {
	case CategoryInternet:
		prefix = "Internet"
	case CategoryChat:
		prefix = "Chat"
	case CategoryMusic:
		prefix = "Music"
	case CategoryVideo:
		prefix = "Video"
	case CategoryGraphics:
		prefix = "Graphics"
	case CategoryOffice:
		prefix = "Office"
	case CategoryGame:
		prefix = "Game"
	case CategoryReading:
		prefix = "Reading"
	case CategoryDevelopment:
		prefix = "Development"
	case CategorySystem:
		prefix = "System"
	case CategoryOthers:
		prefix = "Others"
	default:
		prefix = "Unknown"
	}
	return fmt.Sprintf("%s(%d)", prefix, int(cid))
}

func (cid CategoryID) Pinyin() string {
	switch cid {
	case CategoryInternet:
		return "wangluo"
	case CategoryChat:
		return "shejiaogoutong"
	case CategoryMusic:
		return "yinyuexinshang"
	case CategoryVideo:
		return "shipinbofang"
	case CategoryGraphics:
		return "tuxintuxiang"
	case CategoryOffice:
		return "bangongxuexi"
	case CategoryGame:
		return "youxiyule"
	case CategoryReading:
		return "yuedufanyi"
	case CategoryDevelopment:
		return "bianchengkaifai"
	case CategorySystem:
		return "xitongguanli"
	case CategoryOthers:
		return "qita"
	default:
		return "qita"
	}
}

var categoryNameTable = map[string]CategoryID{
	"internet":    CategoryInternet,
	"chat":        CategoryChat,
	"music":       CategoryMusic,
	"video":       CategoryVideo,
	"graphics":    CategoryGraphics,
	"office":      CategoryOffice,
	"game":        CategoryGame,
	"reading":     CategoryReading,
	"development": CategoryDevelopment,
	"system":      CategorySystem,
	"others":      CategoryOthers,
}

func parseCategoryString(str string) (CategoryID, bool) {
	if str == "" {
		return CategoryOthers, false
	}

	cid, ok := categoryNameTable[str]
	if !ok {
		return CategoryOthers, false
	}
	return cid, true
}

var xCategories = map[string]CategoryID{
	"2dgraphics":       CategoryGraphics,
	"3dgraphics":       CategoryGraphics,
	"accessibility":    CategorySystem,
	"accessories":      CategoryOthers,
	"actiongame":       CategoryGame,
	"advancedsettings": CategorySystem,
	"adventuregame":    CategoryGame,
	"amusement":        CategoryGame,
	"applet":           CategoryOthers,
	"arcadegame":       CategoryGame,
	"archiving":        CategorySystem,
	"art":              CategoryOffice,
	"artificialintelligence": CategoryOffice,
	"astronomy":              CategoryOffice,
	"audio":                  CategoryMusic,
	"audiovideo":             CategoryVideo,
	"audiovideoediting":      CategoryVideo,
	"biology":                CategoryOffice,
	"blocksgame":             CategoryGame,
	"boardgame":              CategoryGame,
	"building":               CategoryDevelopment,
	"calculator":             CategorySystem,
	"calendar":               CategorySystem,
	"cardgame":               CategoryGame,
	"cd":                     CategoryMusic,
	"chart":                  CategoryOffice,
	"chat":                   CategoryChat,
	"chemistry":              CategoryOffice,
	"clock":                  CategorySystem,
	"compiz":                 CategorySystem,
	"compression":            CategorySystem,
	"computerscience":        CategoryOffice,
	"consoleonly":            CategoryOthers,
	"contactmanagement":      CategoryChat,
	"core":                   CategoryOthers,
	"debugger":               CategoryDevelopment,
	"desktopsettings":        CategorySystem,
	"desktoputility":         CategorySystem,
	"development":            CategoryDevelopment,
	"dialup":                 CategorySystem,
	"dictionary":             CategoryOffice,
	"discburning":            CategorySystem,
	"documentation":          CategoryOffice,
	"editors":                CategoryOthers,
	"education":              CategoryOffice,
	"electricity":            CategoryOffice,
	"electronics":            CategoryOffice,
	"email":                  CategoryInternet,
	"emulator":               CategoryGame,
	"engineering":            CategorySystem,
	"favorites":              CategoryOthers,
	"filemanager":            CategorySystem,
	"filesystem":             CategorySystem,
	"filetools":              CategorySystem,
	"filetransfer":           CategoryInternet,
	"finance":                CategoryOffice,
	"game":                   CategoryGame,
	"geography":              CategoryOffice,
	"geology":                CategoryOffice,
	"geoscience":             CategoryOthers,
	"gnome":                  CategorySystem,
	"gpe":                    CategoryOthers,
	"graphics":               CategoryGraphics,
	"guidesigner":            CategoryDevelopment,
	"hamradio":               CategoryOffice,
	"hardwaresettings":       CategorySystem,
	"ide":                    CategoryDevelopment,
	"imageprocessing":        CategoryGraphics,
	"instantmessaging":       CategoryChat,
	"internet":               CategoryInternet,
	"ircclient":              CategoryChat,
	"kde":                    CategorySystem,
	"kidsgame":               CategoryGame,
	"literature":             CategoryOffice,
	"logicgame":              CategoryGame,
	"math":                   CategoryOffice,
	"medicalsoftware":        CategoryOffice,
	"meteorology":            CategoryOthers,
	"midi":                   CategoryMusic,
	"mixer":                  CategoryMusic,
	"monitor":                CategorySystem,
	"motif":                  CategoryOthers,
	"multimedia":             CategoryVideo,
	"music":                  CategoryMusic,
	"network":                CategoryInternet,
	"news":                   CategoryReading,
	"numericalanalysis":      CategoryOffice,
	"ocr":                    CategoryGraphics,
	"office":                 CategoryOffice,
	"p2p":                    CategoryInternet,
	"packagemanager":         CategorySystem,
	"panel":                  CategorySystem,
	"pda":                    CategorySystem,
	"photography":            CategoryGraphics,
	"physics":                CategoryOffice,
	"pim":                    CategoryOthers,
	"player":                 CategoryMusic,
	"playonlinux":            CategoryOthers,
	"presentation":           CategoryOffice,
	"printing":               CategoryOffice,
	"profiling":              CategoryDevelopment,
	"projectmanagement":      CategoryOffice,
	"publishing":             CategoryOffice,
	"puzzlegame":             CategoryGame,
	"rastergraphics":         CategoryGraphics,
	"recorder":               CategoryMusic,
	"remoteaccess":           CategorySystem,
	"revisioncontrol":        CategoryDevelopment,
	"robotics":               CategoryOffice,
	"roleplaying":            CategoryGame,
	"scanning":               CategoryOffice,
	"science":                CategoryOffice,
	"screensaver":            CategoryOthers,
	"sequencer":              CategoryMusic,
	"settings":               CategorySystem,
	"security":               CategorySystem,
	"simulation":             CategoryGame,
	"sportsgame":             CategoryGame,
	"spreadsheet":            CategoryOffice,
	"strategygame":           CategoryGame,
	"system":                 CategorySystem,
	"systemsettings":         CategorySystem,
	"technical":              CategoryOthers,
	"telephony":              CategorySystem,
	"telephonytools":         CategorySystem,
	"terminalemulator":       CategorySystem,
	"texteditor":             CategoryOffice,
	"texttools":              CategoryOffice,
	"transiation":            CategoryDevelopment,
	"translation":            CategoryReading,
	"trayicon":               CategorySystem,
	"tuner":                  CategoryMusic,
	"tv":                     CategoryVideo,
	"utility":                CategorySystem,
	"vectorgraphics":         CategoryGraphics,
	"video":                  CategoryVideo,
	"videoconference":        CategoryInternet,
	"viewer":                 CategoryGraphics,
	"webbrowser":             CategoryInternet,
	"webdevelopment":         CategoryDevelopment,
	"wine":                   CategoryOthers,
	"wine-programs-accessories":               CategoryOthers,
	"wordprocessor":                           CategoryOffice,
	"x-alsa":                                  CategoryMusic,
	"x-bible":                                 CategoryReading,
	"x-bluetooth":                             CategorySystem,
	"x-debian-applications-emulators":         CategoryGame,
	"x-digital_processing":                    CategorySystem,
	"x-enlightenment":                         CategorySystem,
	"x-geeqie":                                CategoryGraphics,
	"x-gnome-networksettings":                 CategorySystem,
	"x-gnome-personalsettings":                CategorySystem,
	"x-gnome-settings-panel":                  CategorySystem,
	"x-gnome-systemsettings":                  CategorySystem,
	"x-gnustep":                               CategorySystem,
	"x-islamic-software":                      CategoryReading,
	"x-jack":                                  CategoryMusic,
	"x-kde-edu-misc":                          CategoryReading,
	"x-kde-internet":                          CategorySystem,
	"x-kde-more":                              CategorySystem,
	"x-kde-utilities-desktop":                 CategorySystem,
	"x-kde-utilities-file":                    CategorySystem,
	"x-kde-utilities-peripherals":             CategorySystem,
	"x-kde-utilities-pim":                     CategorySystem,
	"x-lxde-settings":                         CategorySystem,
	"x-mandriva-office-publishing":            CategoryOthers,
	"x-mandrivalinux-internet-other":          CategorySystem,
	"x-mandrivalinux-office-other":            CategoryOffice,
	"x-mandrivalinux-system-archiving-backup": CategorySystem,
	"x-midi":                           CategoryMusic,
	"x-misc":                           CategorySystem,
	"x-multitrack":                     CategoryMusic,
	"x-novell-main":                    CategorySystem,
	"x-quran":                          CategoryReading,
	"x-red-hat-base":                   CategorySystem,
	"x-red-hat-base-only":              CategorySystem,
	"x-red-hat-extra":                  CategorySystem,
	"x-red-hat-serverconfig":           CategorySystem,
	"x-religion":                       CategoryReading,
	"x-sequencers":                     CategoryMusic,
	"x-sound":                          CategoryMusic,
	"x-sun-supported":                  CategorySystem,
	"x-suse-backup":                    CategorySystem,
	"x-suse-controlcenter-lookandfeel": CategorySystem,
	"x-suse-controlcenter-system":      CategorySystem,
	"x-suse-core":                      CategorySystem,
	"x-suse-core-game":                 CategoryGame,
	"x-suse-core-office":               CategoryOffice,
	"x-suse-sequencer":                 CategoryMusic,
	"x-suse-yast":                      CategorySystem,
	"x-suse-yast-high_availability":    CategorySystem,
	"x-synthesis":                      CategorySystem,
	"x-turbolinux-office":              CategoryOffice,
	"x-xfce":                           CategorySystem,
	"x-xfce-toplevel":                  CategorySystem,
	"x-xfcesettingsdialog":             CategorySystem,
	"x-ximian-main":                    CategorySystem,
}

func parseXCategoryString(name string) CategoryID {
	name = strings.ToLower(name)
	if id, ok := xCategories[name]; ok {
		return id
	}
	logger.Debugf("parseXCategoryString unknown category %q", name)
	return CategoryOthers
}
