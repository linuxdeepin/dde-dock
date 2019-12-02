/*
 * Copyright (C) 2015 ~ 2017 Deepin Technology Co., Ltd.
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

package lastore

/*
#include <stdlib.h>
#include <sys/statvfs.h>
*/
import "C"

import (
	"os"
	"path"
	"sort"
	"strings"
	"unsafe"
)

// QueryLang return user lang.
// the rule is document at man gettext(3)
func QueryLang() string {
	return QueryLangs()[0]
}

// QueryLangs return array of user lang, split by ":".
// the rule is document at man gettext(3)
func QueryLangs() []string {
	LC_ALL := os.Getenv("LC_ALL")
	LC_MESSAGE := os.Getenv("LC_MESSAGE")
	LANGUAGE := os.Getenv("LANGUAGE")
	LANG := os.Getenv("LANG")

	cutoff := func(s string) string {
		for i, c := range s {
			if c == '.' {
				return s[:i]
			}
		}
		return s
	}

	if LC_ALL != "C" && LANGUAGE != "" {
		var r []string
		for _, lang := range strings.Split(LANGUAGE, ":") {
			r = append(r, cutoff(lang))
		}
		return r
	}

	if LC_ALL != "" {
		return []string{cutoff(LC_ALL)}
	}
	if LC_MESSAGE != "" {
		return []string{cutoff(LC_MESSAGE)}
	}
	if LANG != "" {
		return []string{cutoff(LANG)}
	}
	return []string{""}
}

func PackageName(pkg string, lang string) string {
	names := make(map[string]struct {
		Id         string            `json:"id"`
		Name       string            `json:"name"`
		LocaleName map[string]string `json:"locale_name"`
	})

	err := decodeJson(path.Join(varLibDir, "applications.json"), &names)
	if err != nil {
		logger.Warning(err)
	}

	info, ok := names[pkg]
	if !ok {
		return pkg
	}

	name := info.LocaleName[lang]
	if name == "" {
		if info.Name != "" {
			name = info.Name
		} else {
			name = pkg
		}
	}
	return name
}

func strSliceSetEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	for i, va := range a {
		if va != b[i] {
			return false
		}
	}
	return true
}

func queryVFSAvailable(path string) (uint64, error) {
	var vfsStat C.struct_statvfs
	path0 := C.CString(path)
	_, err := C.statvfs(path0, &vfsStat)
	C.free(unsafe.Pointer(path0))
	if err != nil {
		return 0, err
	}
	avail := uint64(vfsStat.f_bavail) * uint64(vfsStat.f_bsize)
	return avail, nil
}
