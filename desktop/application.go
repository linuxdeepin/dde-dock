package desktop

// #cgo pkg-config: glib-2.0
// #include <glib.h>
// #include <stdlib.h>
// int content_type_can_be_executable(char* type);
// int content_type_is(char* type, char* expected_type);
import "C"
import "unsafe"
import (
	"fmt"
	"os/exec"
	"path/filepath"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/glib-2.0"
	"pkg.deepin.io/lib/operations"
	"strings"
)

func GetUserSpecialDir(dir glib.UserDirectory) string {
	return glib.GetUserSpecialDir(dir)
}

func GetDesktopDir() string {
	return GetUserSpecialDir(glib.UserDirectoryDirectoryDesktop)
}

// IMenuable is the interface for something can generate menu.
type IMenuable interface {
	GenMenuContent() (*Menu, error)
	destroy()
}

// Application for desktop daemon.
type Application struct {
	desktop  *Desktop
	settings *Settings
	menuable IMenuable
	menu     *Menu

	ActivateFlagDisplay       int32
	ActivateFlagRunInTerminal int32
	ActivateFlagRun           int32

	RequestRename                 func(string)
	RequestDelete                 func([]string)
	RequestEmptyTrash             func()
	RequestSort                   func(string)
	RequestCleanup                func()
	ReqeustAutoArrange            func()
	RequestCreateFile             func()
	RequestCreateFileFromTemplate func(string)
	RequestCreateDirectory        func()
	ItemCut                       func([]string)
	RequestOpen                   func([]string)
	RequestDismissAppGroup        func(string)

	ActivateFileFailed      func(string)
	CreateAppGroupFailed    func(string)
	MergeIntoAppGroupFailed func(string)

	AppGroupCreated func(string, []string)
	AppGroupDeleted func(string)
	AppGroupMerged  func(string, []string)

	ItemDeleted  func(string)
	ItemCreated  func(string)
	ItemModified func(string)
}

// GetDBusInfo returns dbus info of Application.
func (app *Application) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.dde.daemon.Desktop",
		ObjectPath: "/com/deepin/dde/daemon/Desktop",
		Interface:  "com.deepin.dde.daemon.Desktop",
	}
}

// NewApplication creates a application, settings must not be nil.
func NewApplication(s *Settings) *Application {
	app := &Application{
		ActivateFlagRun:           0,
		ActivateFlagRunInTerminal: 1,
		ActivateFlagDisplay:       2,
	}
	app.desktop = NewDesktop(app)
	app.settings = s
	return app
}

func (app *Application) setSettings(s *Settings) {
	app.settings = s
}

func (app *Application) emitRequestRename(uri string) {
	dbus.Emit(app, "RequestRename", uri)
}

func (app *Application) emitRequestDelete(uris []string) {
	dbus.Emit(app, "RequestDelete", uris)
}

func (app *Application) emitRequestEmptyTrash() {
	dbus.Emit(app, "RequestEmptyTrash")
}

func (app *Application) emitRequestSort(sortPolicy string) {
	dbus.Emit(app, "RequestSort", sortPolicy)
}

func (app *Application) emitRequestCleanup() {
	dbus.Emit(app, "RequestCleanup")
}

func (app *Application) emitRequestAutoArrange() {
	dbus.Emit(app, "RequestAutoArrange")
}

func (app *Application) emitRequestCreateFile() {
	dbus.Emit(app, "RequestCreateFile")
}

func (app *Application) emitRequestCreateFileFromTemplate(template string) {
	dbus.Emit(app, "RequestCreateFileFromTemplate", template)
}

func (app *Application) emitRequestCreateDirectory() {
	dbus.Emit(app, "RequestCreateDirectory")
}

func (app *Application) emitItemCut(uris []string) {
	dbus.Emit(app, "ItemCut", uris)
}

func (app *Application) emitRequestOpen(uris []string) {
	dbus.Emit(app, "RequestOpen", uris)
}

func (app *Application) emitActivateFileFailed(reason string) {
	dbus.Emit(app, "ActivateFileFailed", reason)
}

func (app *Application) emitCreateAppGroupFailed(reason string) {
	dbus.Emit(app, "CreateAppGroupFailed", reason)
}

func (app *Application) emitMergeIntoAppGroupFailed(reason string) {
	dbus.Emit(app, "MergeIntoAppGroupFailed", reason)
}

func (app *Application) emitAppGroupCreated(group string) {
	dbus.Emit(app, "AppGroupCreated", group)
}

func (app *Application) emitAppGroupMerged(group string, files []string) {
	dbus.Emit(app, "AppGroupMerged", group, files)
}

func (app *Application) emitRequestDismissAppGroup(group string) {
	dbus.Emit(app, "RequestDismissAppGroup", group)
}

func (app *Application) emitRequestPaste(dest string) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return
	}

	obj := conn.Object("com.deepin.filemanager.Backend.Clipboard", "/com/deepin/filemanager/Backend/Clipboard")
	if obj != nil {
		obj.Call("com.deepin.filemanager.Backend.Clipboard.EmitPaste", 0, dest).Store()
	}
}

func (app *Application) showProperties(uris []string) {
	exec.Command("deepin-nautilus-properties", uris...).Start()
}

func isAppGroup(uri string) bool {
	return filepath.HasPrefix(filepath.Base(uri), AppGroupPrefix)
}

func isAllAppGroup(uris []string) bool {
	for _, uri := range uris {
		if !isAppGroup(uri) {
			return false
		}
	}
	return true
}

// IsAppGroup returns whether uri is a AppGroup
func (app *Application) IsAppGroup(uri string) bool {
	return isAppGroup(uri)
}

func isAllSpecificItemsAux(a, b string) bool {
	return strings.HasPrefix(a, "trash:") && strings.HasPrefix(b, "computer:")
}

func isAllSpecificItems(uris []string) bool {
	return isAllSpecificItemsAux(uris[0], uris[1]) || isAllSpecificItemsAux(uris[1], uris[0])
}

func (app *Application) getMenuable(uris []string) IMenuable {
	if len(uris) == 0 {
		return app.desktop
	}

	if len(uris) == 1 {
		uri := uris[0]
		if strings.HasPrefix(uri, "trash:") {
			return NewTrashItem(app, uri)
		} else if strings.HasPrefix(uri, "computer:") {
			return NewComputerItem(app, uri)
		}
	}

	if len(uris) == 2 && isAllSpecificItems(uris) {
		// TODO
	}

	if isAllAppGroup(uris) {
		return NewAppGroup(app, uris)
	}

	return NewItem(app, uris)
}

// GenMenuContent returns the menu content in json format used in DeepinMenu.
func (app *Application) GenMenuContent(uris []string) string {
	app.menuable = app.getMenuable(uris)
	menu, err := app.menuable.GenMenuContent()
	if err != nil {
		return ""
	}

	app.menu = menu
	return menu.ToJSON()
}

// HandleSelectedMenuItem will handle selected menu item according to passed id.
func (app *Application) HandleSelectedMenuItem(id string) {
	if app.menu == nil {
		return
	}
	app.menu.HandleAction(id)
}

// DestroyMenu destroys the useless menu.
func (app *Application) DestroyMenu() {
	if app.menu == nil {
		return
	}
	app.menuable.destroy()
	app.menu = nil
}

func filterDesktop(files []string) []string {
	availableFiles := []string{}
	for _, file := range files {
		if strings.HasSuffix(file, ".desktop") {
			availableFiles = append(availableFiles, file)
		}
	}
	return availableFiles
}

// RequestCreatingAppGroup creates app group according to the files, and emits AppGroupCreated signal when it's done.
func (app *Application) RequestCreatingAppGroup(files []string) {
	C.g_reload_user_special_dirs_cache()
	desktopDir := GetDesktopDir()

	availableFiles := filterDesktop(files)

	// get group name
	groupName := getGroupName(availableFiles)

	dirName := AppGroupPrefix + groupName

	// create app group.
	createJob := operations.NewCreateDirectoryJob(desktopDir, dirName, nil)
	createJob.Execute()

	if err := createJob.GetError(); err != nil {
		app.emitCreateAppGroupFailed(err.Error())
		return
	}

	// move files into app group.
	moveJob := operations.NewMoveJob(availableFiles, desktopDir, "", 0, nil)
	moveJob.Execute()
	if err := moveJob.GetError(); err != nil {
		app.emitCreateAppGroupFailed(err.Error())
	}

	app.emitAppGroupCreated(dirName)
}

// RequestMergeIntoAppGroup will merge files into existed AppGroup, and emits AppGroupMerged signal when it's done.
func (app *Application) RequestMergeIntoAppGroup(files []string, appGroup string) {
	availableFiles := filterDesktop(files)

	moveJob := operations.NewMoveJob(availableFiles, GetDesktopDir(), "", 0, nil)
	moveJob.Execute()
	if err := moveJob.GetError(); err != nil {
		app.emitMergeIntoAppGroupFailed(err.Error())
		return
	}

	app.emitAppGroupMerged(appGroup, availableFiles)
}

func (app *Application) doDisplayFile(file *gio.File, contentType string) {
	defaultApp := gio.AppInfoGetDefaultForType(contentType, false)
	if defaultApp == nil {
		app.emitActivateFileFailed("unknown default application")
		return
	}
	defer defaultApp.Unref()

	if _, err := defaultApp.Launch([]*gio.File{file}, nil); err != nil {
		app.emitActivateFileFailed(err.Error())
	}
}

// displayFile will display file using default app.
func (app *Application) displayFile(file string) {
	f := gio.FileNewForCommandlineArg(file)
	if f == nil {
		return
	}
	defer f.Unref()

	info, err := f.QueryInfo(gio.FileAttributeStandardContentType, gio.FileQueryInfoFlagsNone, nil)
	if err != nil {
		app.emitActivateFileFailed(err.Error())
		return
	}
	defer info.Unref()

	contentType := info.GetContentType()
	app.doDisplayFile(f, contentType)
}

// ActivateFile will activate file.
func (app *Application) ActivateFile(file string, args []string, isExecutable bool, flag int32) {
	isDesktopFile := strings.HasSuffix(file, ".desktop")
	if isDesktopFile && isExecutable {
		app.activateDesktopFile(file, args)
	} else {
		app.activateFile(file, args, isExecutable, flag)
	}
}

func (app *Application) activateDesktopFile(file string, args []string) {
	a := gio.NewDesktopAppInfoFromFilename(file)
	if a == nil {
		fmt.Println("XXXXX")
		return
	}
	defer a.Unref()

	a.LaunchUris(args, nil)
}

func (app *Application) activateFile(file string, args []string, isExecutable bool, flag int32) {
	f := gio.FileNewForCommandlineArg(file)
	if f == nil {
		app.emitActivateFileFailed("xxx")
		return
	}
	defer f.Unref()

	info, err := f.QueryInfo(gio.FileAttributeStandardContentType, gio.FileQueryInfoFlagsNone, nil)
	if err != nil {
		app.emitActivateFileFailed(err.Error())
		return
	}
	defer info.Unref()

	contentType := info.GetContentType()
	cContentType := C.CString(contentType)
	defer C.free(unsafe.Pointer(cContentType))

	cPlainType := C.CString("text/plain")
	defer C.free(unsafe.Pointer(cPlainType))

	if isExecutable && (C.int(C.content_type_can_be_executable(cContentType)) == 1 || strings.HasSuffix(file, ".bin")) {
		if C.int(C.content_type_is(cContentType, cPlainType)) == 1 { // runable file
			switch flag {
			case app.ActivateFlagRunInTerminal:
				runInTerminal("", file)
			case app.ActivateFlagRun:
				exec.Command(file).Run()
			case app.ActivateFlagDisplay:
				app.doDisplayFile(f, contentType)
			}
		} else { // binary file
			// FIXME: strange logic from dde-workspace, why args is not used on the other places.
			exec.Command(file, args...).Run()
		}
	} else {
		app.doDisplayFile(f, contentType)
	}
}

// TODO: move to filemanager.
type ItemInfo struct {
	DisplayName string
	BaseName    string
	URI         string
	MIME        string
	Icon        string
	Size        int64
	FileType    uint16
	IsBackup    bool
	IsHidden    bool
	IsReadOnly  bool
	IsSymlink   bool
	CanDelete   bool
	CanExecute  bool
	CanRead     bool
	CanRename   bool
	CanTrash    bool
	CanWrite    bool
}

func toItemInfo(p operations.ListProperty) ItemInfo {
	return ItemInfo{
		DisplayName: p.DisplayName,
		BaseName:    p.BaseName,
		URI:         p.URI,
		MIME:        p.MIME,
		Size:        p.Size,
		FileType:    p.FileType,
		IsBackup:    p.IsBackup,
		IsHidden:    p.IsHidden,
		IsReadOnly:  p.IsReadOnly,
		IsSymlink:   p.IsSymlink,
		CanDelete:   p.CanDelete,
		CanExecute:  p.CanExecute,
		CanRead:     p.CanRead,
		CanRename:   p.CanRename,
		CanTrash:    p.CanTrash,
		CanWrite:    p.CanWrite,
	}
}

func (app *Application) getItemInfo(p operations.ListProperty) ItemInfo {
	info := toItemInfo(p)
	info.Icon = operations.GetThemeIcon(p.URI, 48) //app.settings.iconSize)
	return info
}

func (app *Application) GetDesktopItems() (map[string]ItemInfo, error) {
	infos := map[string]ItemInfo{}
	var err error

	path := GetDesktopDir()
	job := operations.NewListDirJob(path, operations.ListJobFlagIncludeHidden)

	job.ListenProperty(func(p operations.ListProperty) {
		infos[p.URI] = app.getItemInfo(p)
	})

	job.ListenDone(func(e error) {
		if e != nil {
			err = e
			return
		}
	})

	job.Execute()

	return infos, err
}

func (app *Application) GetItemInfo(file string) (ItemInfo, error) {
	info := ItemInfo{}
	f := gio.FileNewForCommandlineArg(file)
	if f == nil {
		return info, fmt.Errorf("Invalid file: %q", file)
	}
	defer f.Unref()

	listProperty, err := operations.GetListProperty(f, nil)
	if err != nil {
		return info, err
	}

	return app.getItemInfo(listProperty), nil
}
