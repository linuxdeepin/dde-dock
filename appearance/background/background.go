package background

import (
	"fmt"
	"os"
	"path"

	"pkg.deepin.io/dde/api/thumbnails/images"
	"pkg.deepin.io/lib/glib-2.0"
	"pkg.deepin.io/lib/graphic"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	thumbWidth  int = 128
	thumbHeight     = 72

	wrapBgSchema    = "com.deepin.wrap.gnome.desktop.background"
	gsKeyBackground = "picture-uri"
)

type Background struct {
	URI string

	Deletable bool
}
type Backgrounds []*Background

func ListBackground() Backgrounds {
	var infos Backgrounds
	for _, file := range getBgFiles() {
		infos = append(infos, &Background{
			URI:       dutils.EncodeURI(file, dutils.SCHEME_FILE),
			Deletable: isDeletable(file),
		})
	}
	return infos
}

func IsBackgroundFile(file string) bool {
	return graphic.IsSupportedImage(dutils.DecodeURI(file))
}

func (infos Backgrounds) Set(uri string) (string, error) {
	uri = dutils.EncodeURI(uri, dutils.SCHEME_FILE)
	info := infos.Get(uri)
	if info != nil {
		return uri, doSetByURI(uri)
	}

	file := dutils.DecodeURI(uri)
	dest, err := getBgDest(file)
	if err != nil {
		return "", err
	}

	if !dutils.IsFileExist(dest) {
		err = dutils.CopyFile(file, dest)
		if err != nil {
			return "", err
		}
	}
	uri = dutils.EncodeURI(dest, dutils.SCHEME_FILE)

	return uri, doSetByURI(dest)
}

func (infos Backgrounds) GetURIs() []string {
	var uris []string
	for _, info := range infos {
		uris = append(uris, info.URI)
	}
	return uris
}

func (infos Backgrounds) Get(uri string) *Background {
	uri = dutils.EncodeURI(uri, dutils.SCHEME_FILE)
	for _, info := range infos {
		if uri == info.URI {
			return info
		}
	}
	return nil
}

func (infos Backgrounds) Delete(uri string) error {
	info := infos.Get(uri)
	if info == nil {
		return fmt.Errorf("Not found '%s'", uri)
	}

	return info.Delete()
}

func (infos Backgrounds) Thumbnail(uri string) (string, error) {
	info := infos.Get(uri)
	if info == nil {
		return "", fmt.Errorf("Not found '%s'", uri)
	}

	return info.Thumbnail()
}

func (info *Background) Delete() error {
	if !info.Deletable {
		return fmt.Errorf("Permission Denied")
	}

	return os.Remove(dutils.DecodeURI(info.URI))
}

func (info *Background) Thumbnail() (string, error) {
	return images.ThumbnailForTheme(info.URI, thumbWidth, thumbHeight, false)
}

func doSetByURI(uri string) error {
	uri = dutils.EncodeURI(uri, dutils.SCHEME_FILE)
	setting, err := dutils.CheckAndNewGSettings(wrapBgSchema)
	if err != nil {
		return err
	}
	defer setting.Unref()

	old := setting.GetString(gsKeyBackground)
	if old == uri {
		return nil
	}

	setting.SetString(gsKeyBackground, uri)
	return nil
}

func getBgDest(file string) (string, error) {
	id, ok := dutils.SumFileMd5(file)
	if !ok {
		return "", fmt.Errorf("Not found '%s'", file)
	}
	return path.Join(
		glib.GetUserSpecialDir(glib.UserDirectoryDirectoryPictures),
		"Wallpapers", id+".png"), nil
}
