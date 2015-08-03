package fonts

type Font struct {
	Id         string
	Name       string
	Family     string
	FamilyName string
	File       string

	Styles []string
	Lang   []string

	Monospace bool
	Deletable bool
}
type Fonts []*Font

func ListFont() Fonts {
	return fcInfosToFonts()
}

func (infos Fonts) ListStandard() Fonts {
	var ret Fonts
	for _, info := range infos {
		if !info.supportedCurLang() {
			continue
		}
		ret = append(ret, info)
	}
	return ret
}

func (infos Fonts) ListMonospace() Fonts {
	var ret Fonts
	for _, info := range infos {
		if !info.Monospace {
			continue
		}
		ret = append(ret, info)
	}
	return ret
}

func (infos Fonts) Get(id string) *Font {
	for _, info := range infos {
		if info.Id == id {
			return info
		}
	}
	return nil
}

func (infos Fonts) convertToFamilies() Families {
	var ret Families
	for _, info := range infos {
		ret = ret.add(&Family{
			Id:     info.Family,
			Name:   info.FamilyName,
			Styles: info.Styles,
			//Files:  []string{info.File},
		})
	}
	return ret
}

func (info *Font) supportedCurLang() bool {
	lang := getCurLang()
	for _, v := range info.Lang {
		if lang == v {
			return true
		}
	}
	return false
}
