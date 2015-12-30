package category

import (
	"bufio"
	"encoding/json"
	"os"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
)

type CategoryJSONInfo struct {
	ID      string
	Locales map[string]map[string]string
	Name    string
}

func GetAllInfos(file string) []CategoryInfo {
	fallbackCategories := []CategoryInfo{
		NewInfo(OthersID, OthersName),
		NewInfo(InternetID, InternetName),
		NewInfo(OfficeID, OfficeName),
		NewInfo(DevelopmentID, DevelopmentName),
		NewInfo(ReadingID, ReadingName),
		NewInfo(GraphicsID, GraphicsName),
		NewInfo(GameID, GameName),
		NewInfo(MusicID, MusicName),
		NewInfo(SystemID, SystemName),
		NewInfo(VideoID, VideoName),
		NewInfo(ChatID, ChatName),
	}

	var categoryInfos []CategoryInfo
	f, err := os.Open(file)
	if err != nil {
		return fallbackCategories
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	var jsonInfo []CategoryJSONInfo
	if err := decoder.Decode(&jsonInfo); err != nil {
		return fallbackCategories
	}

	categoryInfos = make([]CategoryInfo, len(jsonInfo))
	for i, info := range jsonInfo {
		cid, _ := getCategoryID(info.ID)
		categoryInfos[i] = NewInfo(cid, info.Name)
	}

	return categoryInfos
}

type AppJSONInfo struct {
	Id         string
	Category   string
	Name       string
	LocaleName map[string]string
}

func getCategoryInfo(file string) (map[string]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := json.NewDecoder(bufio.NewReader(f))
	jsonData := map[string]AppJSONInfo{}
	if err := decoder.Decode(&jsonData); err != nil {
		return nil, err
	}

	infos := map[string]string{}
	for k, v := range jsonData {
		infos[k] = v.Category
	}
	return infos, nil
}

func getXCategoryInfo(file string) (map[string]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	infos := map[string]string{}
	decoder := json.NewDecoder(bufio.NewReader(f))
	if err := decoder.Decode(&infos); err != nil {
		return nil, err
	}

	return infos, nil
}
