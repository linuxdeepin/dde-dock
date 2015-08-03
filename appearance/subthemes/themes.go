package subthemes

import (
	"fmt"
	"os"
	"path"
	"pkg.deepin.io/dde/api/themes"
	"pkg.deepin.io/dde/api/thumbnails/cursor"
	"pkg.deepin.io/dde/api/thumbnails/gtk"
	"pkg.deepin.io/dde/api/thumbnails/icon"
	dutils "pkg.deepin.io/lib/utils"
	"strings"
)

const (
	thumbWidth  int = 128
	thumbHeight     = 72

	thumbDir = "/usr/share/personalization/thumbnail"
)

type Theme struct {
	Id   string
	Path string

	Deletable bool
}
type Themes []*Theme

func ListGtkTheme() Themes {
	return getThemes(themes.ListGtkTheme())
}

func ListIconTheme() Themes {
	return getThemes(themes.ListIconTheme())
}

func ListCursorTheme() Themes {
	return getThemes(themes.ListCursorTheme())
}

func IsGtkTheme(id string) bool {
	return themes.IsThemeInList(id, themes.ListGtkTheme())
}

func IsIconTheme(id string) bool {
	return themes.IsThemeInList(id, themes.ListIconTheme())
}

func IsCursorTheme(id string) bool {
	return themes.IsThemeInList(id, themes.ListCursorTheme())
}

func SetGtkTheme(id string) error {
	return themes.SetGtkTheme(id)
}

func SetIconTheme(id string) error {
	return themes.SetIconTheme(id)
}

func SetCursorTheme(id string) error {
	return themes.SetCursorTheme(id)
}

func GetGtkThumbnail(id string) (string, error) {
	info := ListGtkTheme().Get(id)
	if info == nil {
		return "", fmt.Errorf("Not found '%s'", id)
	}

	var thumb = path.Join(thumbDir, "WindowThemes", id+"-thumbnail.png")
	if dutils.IsFileExist(thumb) {
		return thumb, nil
	}

	return gtk.GenThumbnail(path.Join(info.Path, "index.theme"), "",
		thumbWidth, thumbHeight)
}

func GetIconThumbnail(id string) (string, error) {
	info := ListIconTheme().Get(id)
	if info == nil {
		return "", fmt.Errorf("Not found '%s'", id)
	}

	var thumb = path.Join(thumbDir, "IconThemes", id+"-thumbnail.png")
	if dutils.IsFileExist(thumb) {
		return thumb, nil
	}
	return icon.GenThumbnail(path.Join(info.Path, "index.theme"), "",
		thumbWidth, thumbHeight)
}

func GetCursorThumbnail(id string) (string, error) {
	info := ListCursorTheme().Get(id)
	if info == nil {
		return "", fmt.Errorf("Not found '%s'", id)
	}

	var thumb = path.Join(thumbDir, "CursorThemes", id+"-thumbnail.png")
	if dutils.IsFileExist(thumb) {
		return thumb, nil
	}
	return cursor.GenThumbnail(path.Join(info.Path, "cursor.theme"), "",
		thumbWidth, thumbHeight)
}

func (infos Themes) GetIds() []string {
	var ids []string
	for _, info := range infos {
		ids = append(ids, info.Id)
	}
	return ids
}

func (infos Themes) Get(id string) *Theme {
	for _, info := range infos {
		if id == info.Id {
			return info
		}
	}
	return nil
}

func (infos Themes) Delete(id string) error {
	info := infos.Get(id)
	if info == nil {
		return fmt.Errorf("Not found '%s'", id)
	}
	return info.Delete()
}

func (info *Theme) Delete() error {
	if !info.Deletable {
		return fmt.Errorf("Permission Denied")
	}
	return os.RemoveAll(info.Path)
}

func getThemes(files []string) Themes {
	var infos Themes
	for _, v := range files {
		infos = append(infos, &Theme{
			Id:        path.Base(v),
			Path:      v,
			Deletable: isDeletable(v),
		})
	}
	return infos
}

func isDeletable(file string) bool {
	if strings.Contains(file, os.Getenv("HOME")) {
		return true
	}
	return false
}
