package softwarecenter

import (
	"dbus/com/linuxdeepin/softwarecenter"
	"fmt"
	"pkg.linuxdeepin.com/lib/dbus"
)

const (
	ActionStart  string = "action-start"
	ActionUpdate string = "action-update"
	ActionFinish string = "action-finish"
	ActionFailed string = "action-failed"

	ActionOperationInstall = iota + 1
	ActionOperationDelete
)

type Action struct {
	Name   string
	Detail dbus.Variant
}

type ActionStartDetail struct {
	PkgName   string
	Operation int32
}

type ActionUpdateDetail struct {
	PkgName     string
	Operation   int32
	Process     int32
	Description string
}

type PkgInfo struct {
	PkgName   string
	Deleted   bool
	Installed bool
	Upgraded  bool
}

type ActionFinishDetail struct {
	PkgName   string
	Operation int32
	Pkgs      []PkgInfo
}

type ActionFailedDetail struct {
	PkgName     string
	Operation   int32
	Pkgs        []PkgInfo
	Description string
}

func NewSoftwareCenter() (*softwarecenter.SoftwareCenter, error) {
	return softwarecenter.NewSoftwareCenter(
		"com.linuxdeepin.softwarecenter",
		"/com/linuxdeepin/softwarecenter",
	)
}

func makeActionStartDetail(detail []interface{}) dbus.Variant {
	pkgName := detail[0].(string)
	operation := detail[1].(int32)
	return dbus.MakeVariant(ActionStartDetail{
		PkgName:   pkgName,
		Operation: operation,
	})
}

func makeActionUpdateDetail(detail []interface{}) dbus.Variant {
	pkgName := detail[0].(string)
	operation := detail[1].(int32)
	process := detail[2].(int32)
	description := detail[3].(string)
	return dbus.MakeVariant(ActionUpdateDetail{
		PkgName:     pkgName,
		Operation:   operation,
		Process:     process,
		Description: description,
	})
}

func makePkgInfoList(infos interface{}) []PkgInfo {
	pkgInfo := make([]PkgInfo, 0)
	for _, v := range infos.([][]interface{}) {
		pkgName := v[0].(string)
		deleted := v[1].(bool)
		installed := v[2].(bool)
		upgraded := v[3].(bool)

		pkgInfo = append(pkgInfo, PkgInfo{
			pkgName,
			deleted,
			installed,
			upgraded,
		})
	}

	return pkgInfo
}

func makeActionFinishDetail(detail []interface{}) dbus.Variant {
	pkgName := detail[0].(string)
	operation := detail[1].(int32)
	return dbus.MakeVariant(ActionFinishDetail{
		PkgName:   pkgName,
		Operation: operation,
		Pkgs:      makePkgInfoList(detail[2]),
	})
}

func makeActionFailedDetail(detail []interface{}) dbus.Variant {
	pkgName := detail[0].(string)
	operation := detail[1].(int32)
	description := detail[3].(string)
	return dbus.MakeVariant(ActionFailedDetail{
		PkgName:     pkgName,
		Operation:   operation,
		Pkgs:        makePkgInfoList(detail[2]),
		Description: description,
	})
}

func UpdateSignalTranslator(message [][]interface{}) []Action {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	info := make([]Action, 0)
	for _, v := range message {
		actionName := v[0].(string)
		action := Action{}
		action.Name = actionName

		switch actionName {
		case ActionStart:
			detail := v[1].(dbus.Variant).Value().([]interface{})
			action.Detail = makeActionStartDetail(detail)
		case ActionUpdate:
			detail := v[1].(dbus.Variant).Value().([]interface{})
			action.Detail = makeActionUpdateDetail(detail)
		case ActionFinish:
			detail := v[1].(dbus.Variant).Value().([]interface{})
			action.Detail = makeActionFinishDetail(detail)
		case ActionFailed:
			detail := v[1].(dbus.Variant).Value().([]interface{})
			action.Detail = makeActionFailedDetail(detail)
		default:
			// logger.Warningf("\"%s\" is not handled", actionName)
		}

		info = append(info, action)
	}

	return info
}
