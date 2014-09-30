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

package fonts

// #cgo pkg-config: fontconfig
// #cgo CFLAGS: -Wall -g
// #include <stdlib.h>
// #include "font_list.h"
import "C"

import (
	"os"
	"regexp"
	"strings"
	"unsafe"
)

const monoSpacing = "100"

type FcInfo struct {
	Family     string
	FamilyLang string
	Style      string
	Lang       string
	// monospace: spacing=100
	Spacing  string
	Filename string
}

type StyleInfo struct {
	Id          string
	Families    []string
	FamilyLangs []string
	StyleList   []string
}

func getFcInfoList() []FcInfo {
	var fcList []FcInfo
	num := C.int(0)

	list := C.get_font_info_list(&num)
	if num < 1 {
		return fcList
	}

	tmpList := uintptr(unsafe.Pointer(list))
	length := unsafe.Sizeof(*list)
	for i := C.int(0); i < num; i++ {
		cInfo := (*C.FontInfo)(unsafe.Pointer(tmpList + uintptr(i)*length))
		info := FcInfo{
			C.GoString(cInfo.family),
			C.GoString(cInfo.familylang),
			C.GoString(cInfo.style),
			C.GoString(cInfo.lang),
			C.GoString(cInfo.spacing),
			C.GoString(cInfo.filename),
		}

		fcList = append(fcList, info)
	}
	C.font_info_list_free(list, num)

	return fcList
}

func getStyleInfoList() ([]StyleInfo, []StyleInfo) {
	var (
		standList []StyleInfo
		monoList  []StyleInfo
	)

	infoList := getFcInfoList()
	for _, info := range infoList {
		var sInfo StyleInfo

		sInfo.Families = strings.Split(info.Family, ",")
		sInfo.FamilyLangs = strings.Split(info.FamilyLang, ",")
		sInfo.StyleList = strings.Split(info.Style, ",")
		sInfo.Id = getFontId(sInfo.Families, sInfo.FamilyLangs)

		if isMonospacedFont(sInfo.Id, info.Spacing) {
			monoList = addStyleInfo(sInfo, monoList)
			continue
		}

		langList := strings.Split(info.Lang, "|")
		if isLangSupported(getCurrentLang(), langList) {
			standList = addStyleInfo(sInfo, standList)
		}
	}

	return standList, monoList
}

func isStyleInfoSame(info1, info2 StyleInfo) bool {
	if info1.Id == info2.Id {
		return true
	}

	return false
}

func isStyleInList(style string, list []string) bool {
	for _, s := range list {
		if style == s {
			return true
		}
	}

	return false
}

func isStyleInfoInList(info StyleInfo, infos []StyleInfo) (bool, int) {
	for i, sInfo := range infos {
		if isStyleInfoSame(info, sInfo) {
			return true, i
		}
	}

	return false, -1
}

func addStyleInfo(info StyleInfo, list []StyleInfo) []StyleInfo {
	found, index := isStyleInfoInList(info, list)

	if !found {
		list = append(list, info)
		return list
	}

	for _, s := range info.StyleList {
		if isStyleInList(s, list[index].StyleList) {
			continue
		}

		list[index].StyleList = append(list[index].StyleList, s)
	}

	return list
}

/**
 * fontconfig language list example:
 * lang: aa|ab|af|av|ay|ba|be|bg|bi|bin|br|bs|bua|ca|ce|
 *       da|de|el|en|eo|es|et|eu|fi|fj|fo|fr|fur|fy|gd|zh-cn|zh-tw|
 *       mn-mn|ms|na|ng|pap-an|pap-aw|rn|rw
 **/
func isLangSupported(lang string, langList []string) bool {
	if len(langList) == 0 {
		return false
	}

	match, err := regexp.Compile("^" + lang)
	if err != nil {
		return false
	}
	for _, lang := range langList {
		if match.MatchString(lang) {
			return true
		}
	}

	return false
}

func getLangIndex(key string, langList []string) int {
	for i, lang := range langList {
		if lang == key {
			return i
		}
	}

	return -1
}

func getFontId(families, langs []string) string {
	idx := getLangIndex("en", langs)
	if idx == -1 || len(families) < idx {
		return families[0]
	}

	return families[idx]
}

func isMonospacedFont(id, spacing string) bool {
	str := strings.ToLower(id)
	if spacing == monoSpacing || strings.Contains(str, "mono") {
		return true
	}

	return false
}

func getCurrentLang() string {
	lang := os.Getenv("LANG")
	if len(lang) == 0 {
		return "en"
	}

	lang = strings.Split(lang, ".")[0]
	lang = strings.Split(lang, "_")[0]
	lang = strings.ToLower(lang)

	return lang
}
