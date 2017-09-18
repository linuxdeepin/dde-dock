/*
 * Copyright (C) 2015 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package dstore

import (
	"bufio"
	"encoding/json"
	"os"
)

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

type CategoryJSONInfo struct {
	ID      string
	Locales map[string]map[string]string
	Name    string
}

func GetAllInfos(file string) []CategoryJSONInfo {
	fallbackCategories := []CategoryJSONInfo{
		CategoryJSONInfo{
			ID:   OthersID,
			Name: OthersName,
		},
		CategoryJSONInfo{
			ID:   InternetID,
			Name: InternetName,
		},
		CategoryJSONInfo{
			ID:   OfficeID,
			Name: OfficeName,
		},
		CategoryJSONInfo{
			ID:   DevelopmentID,
			Name: DevelopmentName,
		},
		CategoryJSONInfo{
			ID:   ReadingID,
			Name: ReadingName,
		},
		CategoryJSONInfo{
			ID:   GraphicsID,
			Name: GraphicsName,
		},
		CategoryJSONInfo{
			ID:   GameID,
			Name: GameName,
		},
		CategoryJSONInfo{
			ID:   MusicID,
			Name: MusicName,
		},
		CategoryJSONInfo{
			ID:   SystemID,
			Name: SystemName,
		},
		CategoryJSONInfo{
			ID:   VideoID,
			Name: VideoName,
		},
		CategoryJSONInfo{
			ID:   ChatID,
			Name: ChatName,
		},
	}
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

	categoryInfos := make([]CategoryJSONInfo, len(jsonInfo))
	for i, info := range jsonInfo {
		categoryInfos[i] = info
	}

	return categoryInfos
}
