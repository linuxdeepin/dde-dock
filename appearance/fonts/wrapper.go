package fonts

// #cgo pkg-config: fontconfig
// #include <stdlib.h>
// #include "font_list.h"
import "C"

import (
	"os"
	"regexp"
	"strings"
	"unsafe"
)

const (
	defaultLang      = "en"
	defaultLangDelim = "|"
	defaultNameDelim = ","
	spaceTypeMono    = "100"
)

// family ex: 'sans', 'serif', 'monospace'
// cRet: `SourceCodePro-Medium.otf: "Source Code Pro" "Medium"`
func fcFontMatch(family string) string {
	cFamily := C.CString(family)
	defer C.free(unsafe.Pointer(cFamily))
	cRet := C.font_match(cFamily)
	defer C.free(unsafe.Pointer(cRet))

	ret := C.GoString(cRet)
	if len(ret) == 0 {
		return ""
	}

	tmp := strings.Split(ret, ":")
	if len(tmp) != 2 {
		return ""
	}

	// return font family id
	name := strings.Split(tmp[1], "\"")[1]
	for _, info := range ListAllFamily() {
		if info.Name == name {
			return info.Id
		}
	}
	return name
}

func fcInfosToFonts() Fonts {
	var num = C.int(0)
	list := C.list_font_info(&num)
	if num < 1 {
		return nil
	}
	defer C.free_font_info_list(list, num)

	listPtr := uintptr(unsafe.Pointer(list))
	itemLen := unsafe.Sizeof(*list)

	var infos Fonts
	for i := C.int(0); i < num; i++ {
		cItem := (*C.FcInfo)(unsafe.Pointer(
			listPtr + uintptr(i)*itemLen))

		infos = append(infos, fcInfoToFont(cItem))
	}
	return infos
}

func fcInfoToFont(cInfo *C.FcInfo) *Font {
	names := strings.Split(C.GoString(cInfo.fullname), defaultNameDelim)
	nameLang := strings.Split(C.GoString(cInfo.fullnamelang),
		defaultNameDelim)
	families := strings.Split(C.GoString(cInfo.family), defaultNameDelim)
	familyLang := strings.Split(C.GoString(cInfo.familylang),
		defaultNameDelim)

	var info = Font{
		Id: getItemByIndex(indexList(defaultLang,
			nameLang), names),
		Name: getItemByIndex(indexList(getCurLang(),
			nameLang), names),
		Family: getItemByIndex(indexList(defaultLang,
			familyLang), families),
		FamilyName: getItemByIndex(indexList(getCurLang(),
			familyLang), families),
		File: C.GoString(cInfo.filename),
		Styles: strings.Split(C.GoString(cInfo.style),
			defaultNameDelim),
		Lang: strings.Split(C.GoString(cInfo.lang),
			defaultLangDelim),
	}
	info.Monospace = isMonospace(info.Name,
		C.GoString(cInfo.spacing))
	info.Deletable = isDeletable(info.File)

	return &info
}

func isMonospace(name, spacing string) bool {
	if spacing == spaceTypeMono ||
		strings.Contains(strings.ToLower(name), "mono") {
		return true
	}

	return false
}

func isDeletable(file string) bool {
	if strings.Contains(file, os.Getenv("HOME")) {
		return true
	}
	return false
}

func getItemByIndex(idx int, list []string) string {
	if len(list) == 0 {
		return ""
	}

	if idx < 0 || len(list) <= idx {
		return list[0]
	}

	return list[idx]
}

func indexList(item string, list []string) int {
	for i, v := range list {
		if v == item {
			return i
		}
	}
	return -1
}

func getCurLang() string {
	locale := os.Getenv("LANGUAGE")
	if len(locale) == 0 {
		locale = os.Getenv("LANG")
	}

	lang := getLangFromLocale(locale)
	if len(lang) == 0 {
		return defaultLang
	}
	return lang
}

func getLangFromLocale(locale string) string {
	if len(locale) == 0 {
		return ""
	}

	locale = strings.ToLower(strings.Split(locale, ".")[0])
	var lang string
	switch locale {
	case "zh_hk":
		lang = "zh-tw"
	case "zh_cn", "zh_tw", "zh_sg", "ku_tr", "mn_mn", "pap_an", "pap_aw":
		lang = regexp.MustCompile("_").ReplaceAllString(locale, "-")
	default:
		lang = strings.Split(locale, "_")[0]
	}
	return lang
}
