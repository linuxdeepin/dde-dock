package utils

// #cgo pkg-config: glib-2.0
// #include <glib.h>
import "C"

// GReloadUserSpecialDirsCache reloads user special dirs cache.
func GReloadUserSpecialDirsCache() {
	C.g_reload_user_special_dirs_cache()
}

func uniqueStringList(l []string) []string {
	m := make(map[string]bool, 0)
	for _, v := range l {
		m[v] = true
	}
	var n []string
	for k := range m {
		n = append(n, k)
	}
	return n
}
