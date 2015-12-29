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
