package main

import (
	"bytes"
	"flag"
	"fmt"
	"path"
	"text/template"
)

var funcMap = template.FuncMap{
	"ToFieldFuncBaseName":       ToFieldFuncBaseName,
	"ToKeyFuncBaseName":         ToKeyFuncBaseName,
	"ToKeyTypeRealData":         ToKeyTypeRealData,
	"ToKeyTypeDefaultValueJSON": ToKeyTypeDefaultValueJSON,
	"IfNeedCheckValueLength":    IfNeedCheckValueLength,
	"GetAllVkFields":            GetAllVkFields,
	"GetAllVkFieldKeys":         GetAllVkFieldKeys,
	"IsVkNeedLogicSetter":       IsVkNeedLogicSetter,
	"ToKeyTypeShortName":        ToKeyTypeShortName,
}

const (
	backEndDir          = ".."
	frontEndDir         = "../../../dss/modules/network/_components_autogen/"
	nmSettingsJSONFile  = "./nm_settings.json"
	nmSettingVkJSONFile = "./nm_setting_vk.json"
)

var (
	nmSettingUtilsFile = path.Join(backEndDir, "nm_setting_utils_autogen.go")
	nmSettingVkFile    = path.Join(backEndDir, "nm_setting_virtual_key_autogen.go")
)

type NMSettingStruct struct {
	FieldName  string // such as "NM_SETTING_CONNECTION_SETTING_NAME"
	FieldValue string // such as "connection"
	Keys       []struct {
		Name           string // such as "NM_SETTING_CONNECTION_ID"
		Value          string // such as "id"
		Type           string // such as "ktypeString"
		Default        string // such as "<default>", "<null>" or "true"
		UsedByBackEnd  bool   // determine if this key will be used by back-end(golang code)
		UsedByFrontEnd bool   // determine if this key will be used by front-end(qml code)
		LogicSet       bool   // determine if this key should to generate a logic setter
		// DisplayName   string	// TODO
	}
}

type NMSettingVkStruct struct {
	Name           string // such as "NM_SETTING_VK_802_1X_EAP"
	Value          string // such as "vk-eap"
	Type           string // such as "ktypeString"
	RelatedField   string // such as "NM_SETTING_802_1X_SETTING_NAME"
	RelatedKey     string // such as "NM_SETTING_802_1X_EAP"
	UsedByFrontEnd bool   // check if is used by front-end
	Optional       bool   // if key is optional, will ignore error for it
	// DisplayName   string	// TODO
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

func genNMSettingVirtualKeyCode(nmSettings []NMSettingStruct, nmSettingVks []NMSettingVkStruct) (content string) {
	content = genTpl(nmSettingVks, tplVirtualKey) // general setting utils
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

	var nmSettings []NMSettingStruct
	unmarshalJSONFile(nmSettingsJSONFile, &nmSettings)

	// back-end code, echo nm setting fields
	for _, nmSetting := range nmSettings {
		autogenContent := genNMSettingCode(nmSetting)
		if writeOutput {
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

	// back-end code, virtual key
	var nmSettingVks []NMSettingVkStruct
	unmarshalJSONFile(nmSettingVkJSONFile, &nmSettingVks)
	autogenContent = genNMSettingVirtualKeyCode(nmSettings, nmSettingVks)
	if writeOutput {
		writeBackendFile(nmSettingVkFile, autogenContent)
	} else {
		fmt.Println(autogenContent)
		fmt.Println()
	}
}
