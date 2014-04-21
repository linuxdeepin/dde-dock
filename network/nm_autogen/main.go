package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"text/template"
)

var funcMap = template.FuncMap{
	"ToFieldFuncBaseName":       ToFieldFuncBaseName,
	"ToKeyFuncBaseName":         ToKeyFuncBaseName,
	"ToKeyTypeRealData":         ToKeyTypeRealData,
	"ToKeyTypeDefaultValueJSON": ToKeyTypeDefaultValueJSON,
	"IfNeedCheckValueLength":    IfNeedCheckValueLength,
}

const nmSettingsJSONFile = "./nm_settings.json"

type NMSettingStruct struct {
	FieldName   string // such as "NM_SETTING_CONNECTION_SETTING_NAME"
	BackEndFile string
	Keys        []struct {
		Name         string // such as "NM_SETTING_CONNECTION_ID"
		Value        string // such as "id"
		Type         string // such as "ktypeString"
		Default      string // such as "<default>", "<null>" or "true"
		BackEndUsed  bool   // determine if this key will be used by back-end(golang code)
		FrontEndUsed bool   // determine if this key will be used by front-end(qml code)
		LogicSet     bool   // determine if this key should to generate a logic setter
	}
}

func generateNMSettingCode(nmSetting NMSettingStruct) (content string) {
	content = fileHeader
	content += generateTemplate(nmSetting, tplGetKeyType)          // get key type
	content += generateTemplate(nmSetting, tplIsKeyInSettingField) // check is key in current field
	content += generateTemplate(nmSetting, tplGetDefaultValueJSON) // get default json value
	content += generateTemplate(nmSetting, tplGeneralGetterJSON)   // general json getter
	content += generateTemplate(nmSetting, tplGeneralSetterJSON)   // general json setter
	content += generateTemplate(nmSetting, tplCheckExists)         // check if key exists
	content += generateTemplate(nmSetting, tplEnsureNoEmpty)       // ensure field and key exists and not empty
	content += generateTemplate(nmSetting, tplGetter)              // getter
	content += generateTemplate(nmSetting, tplSetter)              // setter
	content += generateTemplate(nmSetting, tplJSONGetter)          // json getter
	content += generateTemplate(nmSetting, tplJSONSetter)          // json setter
	content += generateTemplate(nmSetting, tplRemover)             // remover

	// TODO logic setter
	// TODO logic json setter
	// TODO get avaiable values

	return
}

func generateTemplate(nmSetting NMSettingStruct, tplstr string) (content string) {
	templator := template.New("nm autogen").Funcs(funcMap)
	tpl, err := templator.Parse(tplstr)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := bytes.NewBufferString("")
	err = tpl.Execute(buf, nmSetting)
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

	for _, nmSetting := range nmSettings {
		autogenContent := generateNMSettingCode(nmSetting)
		if writeOutput {
			// write to file and execute gofmt
			err = ioutil.WriteFile(nmSetting.BackEndFile, []byte(autogenContent), 0644)
			if err != nil {
				fmt.Println("error, write file failed:", err)
				continue
			}
			execAndWait(10, "gofmt", "-w", nmSetting.BackEndFile)
		} else {
			fmt.Println(autogenContent)
			fmt.Println()
		}
	}
}
