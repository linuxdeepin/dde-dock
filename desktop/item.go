package desktop

import (
	"fmt"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/operations"
)

func getDefaultOpenApp(uri string) (*gio.AppInfo, error) {
	job := operations.NewGetDefaultLaunchAppJob(uri, false)
	job.Execute()
	if job.HasError() {
		return nil, job.GetError()
	}

	return job.Result().(*gio.AppInfo), nil
}

// Item is Normal Item, like file/directory/link.
type Item struct {
	uri      string
	uris     []string
	files    []*gio.File
	multiple bool
	app      *Application
	menu     *Menu
}

// NewItem creates new item.
func NewItem(app *Application, uris []string) *Item {
	return &Item{
		app:      app,
		uri:      uris[0],
		uris:     uris,
		multiple: len(uris) > 1,
	}
}

func (item *Item) settings() *Settings {
	return item.app.settings
}

// func (item *Item) App() *Application {
// 	return item.app
// }

func (item *Item) emitRequestDelete() {
	item.app.emitRequestDelete(item.uris)
}

func (item *Item) emitRequestRename() {
	item.app.emitRequestRename(item.uri)
}

func (item *Item) emitRequestEmptyTrash() {
	item.app.emitRequestEmptyTrash()
}

func (item *Item) emitRequestCreateFile() {
	item.app.emitRequestCreateFile()
}

func (item *Item) emitRequestCreateFileFromTemplate(template string) {
	item.app.emitRequestCreateFileFromTemplate(template)
}

func (item *Item) emitRequestCreateDirectory() {
	item.app.emitRequestCreateDirectory()
}

func (item *Item) showProperties() {
	item.app.showProperties(item.uris)
}

// ArchiveMimeTypes is a list of MIMEType for archive files.
var ArchiveMimeTypes = []string{
	"application/x-gtar",
	"application/x-zip",
	"application/x-zip-compressed",
	"application/zip",
	"application/x-zip",
	"application/x-tar",
	"application/x-7z-compressed",
	"application/x-rar",
	"application/x-rar-compressed",
	"application/x-jar",
	"application/x-java-archive",
	"application/x-war",
	"application/x-ear",
	"application/x-arj",
	"application/x-gzip",
	"application/gzip",
	"application/x-bzip-compressed-tar",
	"application/x-compressed-tar",
	"application/x-archive",
	"application/x-xz-compressed-tar",
	"application/x-bzip",
	"application/x-cbz",
	"application/x-xz",
	"application/x-lzma-compressed-tar",
	"application/x-ms-dos-executable",
	"application/x-lzma",
	"application/x-cd-image",
	"application/x-deb",
	"application/x-rpm",
	"application/x-stuffit",
	"application/x-tzo",
	"application/x-tarz",
	"application/x-tzo",
	"application/x-msdownload",
	"application/x-lha",
	"application/x-zoo",
}

func isArchived(f *gio.File) bool {
	info, err := f.QueryInfo(gio.FileAttributeStandardContentType, gio.FileQueryInfoFlagsNone, nil)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer info.Unref()
	contentType := info.GetContentType()
	for _, MIMEType := range ArchiveMimeTypes {
		if contentType == MIMEType {
			return true
		}
	}

	return false
}

func (item *Item) destroy() {
	for _, file := range item.files {
		file.Unref()
	}
}

// GenMenuContent generates json format menu content used in DeepinMenu for normal itself.
func (item *Item) GenMenuContent() (*Menu, error) {
	item.menu = NewMenu()
	item.files = make([]*gio.File, len(item.uris))
	for i, uri := range item.uris {
		item.files[i] = gio.FileNewForCommandlineArg(uri)
		if item.files[i] == nil {
			return nil, fmt.Errorf("No such a file or directory: %s", item.uri)
		}
	}

	menu := item.menu.AppendItem(NewMenuItem(Tr("_Open"), func() {
		// be care of different open app.
		item.app.emitRequestOpen(item.uris)
	}, true))

	// same type?
	if true {
		openWithMenuItem := NewMenuItem(Tr("Open _with"), func() {}, true)
		menu.AppendItem(openWithMenuItem)

		openWithSubMenu := NewMenu()
		openWithSubMenu.SetIDGenerator(menu.genID)
		openWithMenuItem.subMenu = openWithSubMenu

		job := operations.NewGetRecommendedLaunchAppsJob(item.uri)
		job.Execute()
		if !job.HasError() {
			recommendedApps := job.Result().([]*gio.AppInfo)
			if len(recommendedApps) > 0 {
				for _, app := range recommendedApps {
					openWithSubMenu.AppendItem(NewMenuItem(app.GetName(), func(id string) func() {
						return func() {
							fmt.Println("open with", id)
							app := gio.NewDesktopAppInfo(id)
							if app == nil {
								fmt.Println("get app failed:", id)
								return
							}
							defer app.Unref()
							defer func() {
								for _, file := range item.files {
									file.Unref()
								}
							}()
							// be care of different open app.
							app.Launch(item.files, nil)
						}
					}(app.GetId()), true))
					app.Unref()
				}

				openWithSubMenu.AddSeparator()
			}
		}

		openWithSubMenu.AppendItem(NewMenuItem(Tr("_Chose"), func() {
			// TODO:
			fmt.Println("chose open with")
		}, true))

	}

	menu.AddSeparator()

	// TODO: use plugin, remove useless function.
	if true {
		runFileRoller := func(cmd string, files []*gio.File) error {
			app, err := gio.AppInfoCreateFromCommandline(cmd, "", gio.AppInfoCreateFlagsSupportsStartupNotification)
			if err != nil {
				return err
			}
			defer app.Unref()
			_, err = app.Launch(files, nil)
			return err
		}

		menu.AppendItem(NewMenuItem(Tr("Co_mpress"), func() {
			err := runFileRoller("file-roller -d %U", item.files)
			if err != nil {
				fmt.Println(err)
			}
		}, true))

		allIsArchived := true
		for _, file := range item.files {
			if !isArchived(file) {
				allIsArchived = false
				break
			}
		}

		if allIsArchived {
			menu.AppendItem(NewMenuItem(Tr("_Extract"), func() {
				err := runFileRoller("file-roller -f", item.files)
				if err != nil {
					fmt.Println(err)
				}
			}, true)).AppendItem(NewMenuItem(Tr("Extract _Here"), func() {
				err := runFileRoller("file-roller -h", item.files)
				if err != nil {
					fmt.Println(err)
				}
			}, true)).AddSeparator()
		}
	}

	menu.AppendItem(NewMenuItem(Tr("Cu_t"), func() {
		operations.CutToClipboard(item.uris)
		item.app.emitItemCut(item.uris)
	}, true)).AppendItem(NewMenuItem(Tr("_Copy"), func() {
		operations.CopyToClipboard(item.uris)
	}, true))

	if !item.multiple {
		fileType := item.files[0].QueryFileType(gio.FileQueryInfoFlagsNone, nil)
		if fileType == gio.FileTypeDirectory {
			menu.AppendItem(NewMenuItem(Tr("Paste _Into"), func() {
				item.app.emitRequestPaste(item.uri)
			}, operations.CanPaste(item.uri))).AddSeparator().AppendItem(NewMenuItem(Tr("Open in _terminal"), func() {
				runInTerminal(item.uri, "")
			}, !item.multiple))
		}
	}

	menu.AddSeparator()

	menu.AppendItem(NewMenuItem(Tr("_Rename"), func() {
		item.emitRequestRename()
	}, !item.multiple)).AppendItem(NewMenuItem(Tr("_Delete"), func() {
		item.emitRequestDelete()
	}, true))

	menu.AddSeparator()

	return menu.AppendItem(NewMenuItem(Tr("_Properties"), func() {
		item.showProperties()
	}, true)), nil
}
