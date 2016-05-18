package dock

import (
	"gir/glib-2.0"
	"strings"
)

type AppIdFilterGroup struct {
	KeyFileWmClass    *glib.KeyFile
	KeyFileWmInstance *glib.KeyFile
	KeyFileExecName   *glib.KeyFile
	KeyFileArgs       *glib.KeyFile
	KeyFileIconName   *glib.KeyFile
	KeyFileWmName     *glib.KeyFile
}

const (
	ddeDataDir           = "/usr/share/dde/data"
	filterWmClassFile    = ddeDataDir + "/filter_wmclass.ini"
	filterWmInstanceFile = ddeDataDir + "/filter_wminstance.ini"
	filterExecNameFile   = ddeDataDir + "/filter_execname.ini"
	filterArgsFile       = ddeDataDir + "/filter_arg.ini"
	filterIconNameFile   = ddeDataDir + "/filter_icon_name.ini"
	filterWmNameFile     = ddeDataDir + "/filter_wmname.ini"
)

func loadFilterFile(file string) (*glib.KeyFile, error) {
	kf := glib.NewKeyFile()
	_, err := kf.LoadFromFile(file, glib.KeyFileFlagsNone)
	if err != nil {
		return nil, err
	}
	return kf, nil
}

func NewAppIdFilterGroup() *AppIdFilterGroup {
	f := &AppIdFilterGroup{}
	var err error
	f.KeyFileWmClass, err = loadFilterFile(filterWmClassFile)
	if err != nil {
		logger.Warning(err)
	}
	f.KeyFileWmInstance, err = loadFilterFile(filterWmInstanceFile)
	if err != nil {
		logger.Warning(err)
	}
	f.KeyFileExecName, err = loadFilterFile(filterExecNameFile)
	if err != nil {
		logger.Warning(err)
	}
	f.KeyFileArgs, err = loadFilterFile(filterArgsFile)
	if err != nil {
		logger.Warning(err)
	}
	f.KeyFileIconName, err = loadFilterFile(filterIconNameFile)
	if err != nil {
		logger.Warning(err)
	}
	f.KeyFileWmName, err = loadFilterFile(filterWmNameFile)
	if err != nil {
		logger.Warning(err)
	}
	return f
}

func findAppIdByFilter(group, arg string, filter *glib.KeyFile) string {
	if group == "" || arg == "" || filter == nil {
		return ""
	}
	if filter.HasGroup(group) {
		_, keys, _ := filter.GetKeys(group)
		for _, key := range keys {
			if strings.Contains(arg, key) {
				value, _ := filter.GetString(group, key)
				return value
			}
		}
	}
	return ""
}
