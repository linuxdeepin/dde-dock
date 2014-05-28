package main

import (
	"bytes"
	"flag"
	"fmt"
	"path"
	"text/template"
)

var funcMap = template.FuncMap{
	"ToCaplitalize":               ToCaplitalize,
	"ToSectionFuncBaseName":       ToSectionFuncBaseName,
	"ToKeyFuncBaseName":           ToKeyFuncBaseName,
	"ToKeyTypeRealData":           ToKeyTypeRealData,
	"ToKeyTypeDefaultValue":       ToKeyTypeDefaultValue,
	"IfNeedCheckValueLength":      IfNeedCheckValueLength,
	"GetAllVkeysRelatedSections":  GetAllVkeysRelatedSections,
	"GetVkeysOfSection":           GetVkeysOfSection,
	"ToKeyTypeShortName":          ToKeyTypeShortName,
	"ToKeyDisplayName":            ToKeyDisplayName,
	"ToKeyValue":                  ToKeyValue,
	"IsKeyUsedByFrontEnd":         IsKeyUsedByFrontEnd,
	"ToFrontEndWidget":            ToFrontEndWidget,
	"ToClassName":                 ToClassName,
	"GetAllKeysInPage":            GetAllKeysInPage,
	"GetKeyWidgetProp":            GetKeyWidgetProp,
	"ToKeyTypeInterfaceConverter": ToKeyTypeInterfaceConverter,
}

const (
	backEndDir            = ".."
	frontEndDir           = "../../../dss/modules/network/edit_autogen/"
	nmSettingsJSONFile    = "./nm_settings.json"
	nmSettingVkeyJSONFile = "./nm_setting_vkey.json"
	nmSettingPageJSONFile = "./nm_setting_vsection.json"
)

var (
	argWriteOutput       bool
	argBackEnd           bool
	argFrontEnd          bool
	nmSettingUtilsFile   = path.Join(backEndDir, "nm_setting_general_autogen.go")
	nmSettingVkeyFile    = path.Join(backEndDir, "nm_setting_virtual_key_autogen.go")
	frontEndConnPropFile = path.Join(frontEndDir, "BaseConnectionEdit.qml")
	nmSettings           []NMSettingStruct
	nmSettingVkeys       []NMSettingVkeyStruct
	nmSettingPages       []NMSettingPageStruct
)

type NMSettingStruct struct {
	SectionName  string // such as "NM_SETTING_CONNECTION_SETTING_NAME"
	SectionValue string // such as "connection"
	Keys         []NMSettingKeyStruct
}

type NMSettingKeyStruct struct {
	Name           string            // such as "NM_SETTING_CONNECTION_ID"
	Value          string            // such as "id"
	Type           string            // such as "ktypeString"
	Default        string            // such as "<default>", "<null>" or "true"
	UsedByBackEnd  bool              // determine if this key will be used by back-end(golang code)
	UsedByFrontEnd bool              // determine if this key will be used by front-end(qml code)
	LogicSet       bool              // determine if this key should to generate a logic setter
	DisplayName    string            // such as "Connection name"
	FrontEndWidget string            // such as "EditLinePasswordInput"
	WidgetProp     map[string]string // properties for front end widget, such as "WidgetProp":{"alwaysUpdate":"true"}
}

type NMSettingVkeyStruct struct {
	Name           string // such as "NM_SETTING_VK_802_1X_EAP"
	Value          string // such as "vk-eap"
	Type           string // such as "ktypeString"
	RelatedSection string // such as "NM_SETTING_802_1X_SETTING_NAME"
	RelatedKey     string // such as "NM_SETTING_802_1X_EAP"
	EnableWrapper  bool   // check if the virtual key is a wrapper just to enable target key
	UsedByFrontEnd bool   // check if is used by front-end
	Optional       bool   // if key is optional, will ignore error for it
	DisplayName    string
	FrontEndWidget string            // such as "EditLinePasswordInput"
	WidgetProp     map[string]string // properties for front end widget, such as "WidgetProp":{"alwaysUpdate":"true"}
}

type NMSettingPageStruct struct {
	Ignore          bool
	Name            string
	DisplayName     string
	RelatedSections []string
}

func genNMSettingCode(nmSetting NMSettingStruct) (content string) {
	content = fileHeader
	content += genTpl(nmSetting, tplGetKeyType)            // get key type
	content += genTpl(nmSetting, tplIsKeyInSettingSection) // check is key in current section
	content += genTpl(nmSetting, tplGetDefaultValue)       // get default value
	content += genTpl(nmSetting, tplGeneralGetterJSON)     // general json getter
	content += genTpl(nmSetting, tplGeneralSetterJSON)     // general json setter
	content += genTpl(nmSetting, tplCheckExists)           // check if key exists
	content += genTpl(nmSetting, tplEnsureNoEmpty)         // ensure section and key exists and not empty
	content += genTpl(nmSetting, tplGetter)                // getter
	content += genTpl(nmSetting, tplSetter)                // setter
	content += genTpl(nmSetting, tplJSONGetter)            // json getter
	content += genTpl(nmSetting, tplJSONSetter)            // json setter
	content += genTpl(nmSetting, tplLogicJSONSetter)       // logic json setter
	content += genTpl(nmSetting, tplRemover)               // remover
	return
}

func genNMSettingGeneralUtilsCode(nmSettings []NMSettingStruct) (content string) {
	content = genTpl(nmSettings, tplGeneralSettingUtils) // general setting utils
	return
}

func genNMSettingVkeyCode(nmSettings []NMSettingStruct, nmSettingVkeys []NMSettingVkeyStruct) (content string) {
	content = genTpl(nmSettingVkeys, tplVkey) // general setting utils
	return
}

func genFrontEndConnPropCode(nmPages []NMSettingPageStruct) (content string) {
	content = genTpl(nmPages, tplFrontEndConnProp)
	return
}

func genFrontEndSectionCode(nmPage NMSettingPageStruct) (content string) {
	content = genTpl(nmPage, tplFrontEndSection)
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

func genBackEndCode() {
	// back-end code, echo nm setting sections
	for _, nmSetting := range nmSettings {
		autogenContent := genNMSettingCode(nmSetting)
		backEndFile := getBackEndFilePath(nmSetting.SectionName)
		writeOrDisplayResultForBackEnd(backEndFile, autogenContent)
	}

	// back-end code, general setting utils
	autogenContent := genNMSettingGeneralUtilsCode(nmSettings)
	writeOrDisplayResultForBackEnd(nmSettingUtilsFile, autogenContent)

	// back-end code, virtual key
	autogenContent = genNMSettingVkeyCode(nmSettings, nmSettingVkeys)
	writeOrDisplayResultForBackEnd(nmSettingVkeyFile, autogenContent)
}

func genFrontEndCode() {
	// front-end code, BaseConnectionProperties.qml
	autogenContent := genFrontEndConnPropCode(nmSettingPages)
	writeOrDisplayResultForFrontEnd(frontEndConnPropFile, autogenContent)

	for _, nmPage := range nmSettingPages {
		if nmPage.Ignore {
			continue
		}
		autogenContent = genFrontEndSectionCode(nmPage)
		frontEndFile := getFrontEndFilePath(nmPage.Name)
		writeOrDisplayResultForFrontEnd(frontEndFile, autogenContent)
	}
}

func main() {
	flag.BoolVar(&argWriteOutput, "w", false, "write to file")
	flag.BoolVar(&argBackEnd, "b", false, "generate back-end code")
	flag.BoolVar(&argFrontEnd, "f", false, "generate front-end code")
	flag.Parse()

	unmarshalJSONFile(nmSettingsJSONFile, &nmSettings)
	unmarshalJSONFile(nmSettingVkeyJSONFile, &nmSettingVkeys)
	unmarshalJSONFile(nmSettingPageJSONFile, &nmSettingPages)
	if argBackEnd {
		genBackEndCode()
	}
	if argFrontEnd {
		genFrontEndCode()
	}
}
