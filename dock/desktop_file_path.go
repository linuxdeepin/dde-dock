package dock

import (
	"path/filepath"
	"strings"
)

var pathDirCodeMap map[string]string
var pathCodeDirMap map[string]string

const desktopExt = ".desktop"

func initPathDirCodeMap() {
	pathDirCodeMap = map[string]string{
		"/usr/share/applications/":       "/S@",
		"/usr/local/share/applications/": "/L@",
	}

	dir := filepath.Join(homeDir, ".local/share/applications")
	dir = addDirTrailingSlash(dir)
	pathDirCodeMap[dir] = "/H@"

	dir = addDirTrailingSlash(scratchDir)
	pathDirCodeMap[dir] = "/D@"

	logger.Debugf("pathDirCodeMap: %#v", pathDirCodeMap)

	pathCodeDirMap = make(map[string]string, len(pathDirCodeMap))
	for dir, code := range pathDirCodeMap {
		pathCodeDirMap[code] = dir
	}
}

func getDesktopIdByFilePath(path string) string {
	var desktopId string
	for dir, _ := range pathDirCodeMap {
		if strings.HasPrefix(path, dir) {
			desktopId = path[len(dir):]
			desktopId = strings.Replace(desktopId, "/", "-", -1)
		}
	}
	return desktopId
}

func addDirTrailingSlash(dir string) string {
	if len(dir) == 0 {
		panic("length of dir is 0")
	}
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}
	return dir
}

func addDesktopExt(str string) string {
	if strings.HasSuffix(str, desktopExt) {
		return str
	}
	return str + desktopExt
}

func trimDesktopExt(str string) string {
	if strings.HasSuffix(str, desktopExt) {
		return str[:len(str)-len(desktopExt)]
	}
	return str
}

func zipDesktopPath(path string) string {
	for dir, code := range pathDirCodeMap {
		if strings.HasPrefix(path, dir) {
			path = code + path[len(dir):]
		}
	}
	return trimDesktopExt(path)
}

func unzipDesktopPath(path string) string {
	head := path[:3]
	for code, dir := range pathCodeDirMap {
		if code == head {
			path = dir + path[3:]
			break
		}
	}
	return addDesktopExt(path)
}
