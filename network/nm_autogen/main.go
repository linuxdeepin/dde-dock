package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"path"
	"text/template"
)

var funcMap = template.FuncMap{
	"ToFieldFuncBaseName":       ToFieldFuncBaseName,
	"ToKeyFuncBaseName":         ToKeyFuncBaseName,
	"ToKeyTypeRealData":         ToKeyTypeRealData,
	"ToKeyTypeDefaultValueJSON": ToKeyTypeDefaultValueJSON,
	"IfNeedCheckValueLength":    IfNeedCheckValueLength,
}

const (
	backEndDir         = ".."
	frontEndDir        = "../../../dss/modules/network/_components_autogen/"
	nmSettingsJSONFile = "./nm_settings.json"
)

var nmSettingUtilsFile = path.Join(backEndDir, "nm_setting_utils_autogen.go")

type NMSettingStruct struct {
	FieldName string // such as "NM_SETTING_CONNECTION_SETTING_NAME"
	Keys      []struct {
		Name         string // such as "NM_SETTING_CONNECTION_ID"
		Value        string // such as "id"
		Type         string // such as "ktypeString"
		Default      string // such as "<default>", "<null>" or "true"
		BackEndUsed  bool   // determine if this key will be used by back-end(golang code)
		FrontEndUsed bool   // determine if this key will be used by front-end(qml code)
		LogicSet     bool   // determine if this key should to generate a logic setter
		// DisplayName   string	// TODO
	}
}

type NMPageStruct struct {
	Name          string
	DisplayName   string
	RelatedFields []string
	// FrontEndFile  string // TODO
}

func genNMSettingCode(nmSetting NMSettingStruct) (content string) {
	content = fileHeader
	content += genTpl(nmSetting, tplGetKeyType)          // get key type
	content += genTpl(nmSetting, tplIsKeyInSettingField) // check is key in current field
	content += genTpl(nmSetting, tplGetDefaultValueJSON) // get default json value
	content += genTpl(nmSetting, tplGeneralGetterJSON)   // general json getter
	content += genTpl(nmSetting, tplGeneralSetterJSON)   // general json setter
	content += genTpl(nmSetting, tplCheckExists)         // check if key exists
	content += genTpl(nmSetting, tplEnsureNoEmpty)       // ensure field and key exists and not empty
	content += genTpl(nmSetting, tplGetter)              // getter
	content += genTpl(nmSetting, tplSetter)              // setter
	content += genTpl(nmSetting, tplJSONGetter)          // json getter
	content += genTpl(nmSetting, tplJSONSetter)          // json setter
	content += genTpl(nmSetting, tplRemover)             // remover

	// TODO logic setter
	// TODO logic json setter
	// TODO get avaiable values

	return
}

func genNMSettingGeneralUtilsCode(nmSettings []NMSettingStruct) (content string) {
	content = genTpl(nmSettings, tplGeneralSettingUtils) // general setting utils
	return
}

func genTpl(data interface{}, tplstr string) (content string) {
	templator := template.New("nm autogen").Funcs(funcMap)
	tpl, err := templator.Parse(tplstr)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := bytes.NewBufferString("")
	err = tpl.Execute(buf, data)
	if err != nil {
		fmt.Println(err)
		return
	}
	content = string(buf.Bytes())
	return
}

func main() {
	var writeOutput bool

	flag.BoolVar(&writeOutput, "w", false, "write to output file")
	flag.Parse()

	fileContent, err := ioutil.ReadFile(nmSettingsJSONFile)
	if err != nil {
		fmt.Println("error, open file failed:", err)
		return
	}

	var nmSettings []NMSettingStruct
	err = json.Unmarshal(fileContent, &nmSettings)
	if err != nil {
		fmt.Printf("error, unmarshal json %s failed: %v\n", nmSettingsJSONFile, err)
		return
	}

	// back-end code
	for _, nmSetting := range nmSettings {
		autogenContent := genNMSettingCode(nmSetting)
		if writeOutput {
			// write to file and execute gofmt
			backEndFile := getBackEndFilePath(nmSetting.FieldName)
			writeBackendFile(backEndFile, autogenContent)
		} else {
			fmt.Println(autogenContent)
			fmt.Println()
		}
	}

	// back-end code, general setting utils
	autogenContent := genNMSettingGeneralUtilsCode(nmSettings)
	if writeOutput {
		writeBackendFile(nmSettingUtilsFile, autogenContent)
	} else {
		fmt.Println(autogenContent)
		fmt.Println()
	}
}
