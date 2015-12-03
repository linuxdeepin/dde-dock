/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
	"bytes"
	"flag"
	"path"
	"text/template"
)

var funcMap = template.FuncMap{
	"ToCaplitalize":               ToCaplitalize,
	"ToSectionFuncBaseName":       ToSectionFuncBaseName,
	"ToKeyFuncBaseName":           ToKeyFuncBaseName,
	"ToKeyTypeRealData":           ToKeyTypeRealData,
	"ToKeyDefaultValue":           ToKeyDefaultValue,
	"IfNeedCheckValueLength":      IfNeedCheckValueLength,
	"GetAllVkeysRelatedSections":  GetAllVkeysRelatedSections,
	"GetVkeysOfSection":           GetVkeysOfSection,
	"ToKeyTypeShortName":          ToKeyTypeShortName,
	"ToKeyDisplayName":            ToKeyDisplayName,
	"ToKeyValue":                  ToKeyValue,
	"ToKeyAlwaysUpdate":           ToKeyAlwaysUpdate,
	"ToKeyUseValueRange":          ToKeyUseValueRange,
	"ToKeyMinValue":               ToKeyMinValue,
	"ToKeyMaxValue":               ToKeyMaxValue,
	"ToKeyRelatedSectionValue":    ToKeyRelatedSectionValue,
	"IsKeyUsedByFrontEnd":         IsKeyUsedByFrontEnd,
	"ToFrontEndWidget":            ToFrontEndWidget,
	"ToClassName":                 ToClassName,
	"ToVsClassName":               ToVsClassName,
	"GetAllKeysInVsection":        GetAllKeysInVsection,
	"GetKeyWidgetProps":           GetKeyWidgetProps,
	"ToKeyTypeInterfaceConverter": ToKeyTypeInterfaceConverter,
	"IsEnableWrapperVkey":         IsEnableWrapperVkey,
}

const (
	nmSettingsJSONFile        = "./nm_settings.json"
	nmSettingVkeyJSONFile     = "./nm_setting_vkey.json"
	nmSettingVsectionJSONFile = "./nm_setting_vsection.json"
)

var (
	argWriteOutput        bool
	argGenBackEndGo       bool
	argGenFrontEndQt      bool
	argGenFrontEndQml     bool
	argBackEndGoDir       string
	argFrontEndQtDir      string
	argFrontEndQmlDir     string
	nmSettingUtilsFile    string
	nmSettingVkeyFile     string
	nmSettingVsectionFile string
	backEndFile           string
	frontEndConnPropFile  string
	nmSections            []NMSectionStruct
	nmVkeys               []NMVkeyStruct
	nmVsections           []NMVsectionStruct
)

type NMSectionStruct struct {
	Name  string // such as "NM_SETTING_CONNECTION_SETTING_NAME"
	Value string // such as "connection"
	Keys  []NMKeyStruct
}

type NMVsectionStruct struct {
	Ignore          bool
	Name            string // such as "NM_SETTING_VS_GENERAL"
	Value           string // such as "vs-general"
	DisplayName     string
	RelatedSections []string
}

type NMKeyStruct struct {
	Name           string            // such as "NM_SETTING_CONNECTION_ID"
	Value          string            // such as "id"
	Type           string            // such as "ktypeString"
	Default        string            // such as "<default>", "<null>" or "true"
	UsedByBackEnd  bool              // determine if this key will be used by back-end(golang code)
	UsedByFrontEnd bool              // determine if this key will be used by front-end(qml code)
	LogicSet       bool              // determine if this key should to generate a logic setter
	DisplayName    string            // such as "Connection name"
	FrontEndWidget string            // front-end widget name, such as "EditLinePasswordInput"
	AlwaysUpdate   bool              // mark if front-end widget should be re-get value when other keys changed
	UseValueRange  bool              // mark custom value range will be used for integer keys
	MinValue       int               // minimize value
	MaxValue       int               // maximum value
	WidgetProps    map[string]string // properties for front-end widget, such as "WidgetProps":{"alwaysUpdate":"true"}
}

type NMVkeyStruct struct {
	Name           string            // such as "NM_SETTING_VK_802_1X_EAP"
	Value          string            // such as "vk-eap"
	Type           string            // such as "ktypeString"
	VkType         string            // could be "vkTypeWrapper", "vkTypeEnableWrapper" or "vkTypeController"
	RelatedSection string            // such as "NM_SETTING_802_1X_SETTING_NAME"
	RelatedKeys    []string          // such as "NM_SETTING_802_1X_EAP"
	UsedByFrontEnd bool              // check if is used by front-end
	ChildKey       bool              // such as ip address, mask and gateway
	Optional       bool              // if key is optional, will ignore error for it
	DisplayName    string            // display name for font-end, need i18n
	FrontEndWidget string            // such as "EditLinePasswordInput"
	AlwaysUpdate   bool              // mark if front-end widget should be re-get value when other keys changed
	UseValueRange  bool              // mark custom value range will be used for integer keys
	MinValue       int               // minimize value
	MaxValue       int               // maximum value
	WidgetProps    map[string]string // properties for front end widget, such as "WidgetProps":{"alwaysUpdate":"true"}
}

func setupDirs() {
	nmSettingUtilsFile = path.Join(argBackEndGoDir, "nm_setting_general_gen.go")
	nmSettingVkeyFile = path.Join(argBackEndGoDir, "nm_setting_virtual_key_gen.go")
	nmSettingVsectionFile = path.Join(argBackEndGoDir, "nm_setting_virtual_section_gen.go")
	backEndFile = path.Join(argBackEndGoDir, "nm_setting_bean_gen.go")
	frontEndConnPropFile = path.Join(argFrontEndQmlDir, "BaseConnectionEdit.qml")
}

func genNMSettingCode(nmSection NMSectionStruct) (content string) {
	content += "\n// Origin file name " + getBackEndFilePath(nmSection.Name)
	content += genTpl(nmSection, tplGetKeyType)            // get key type
	content += genTpl(nmSection, tplIsKeyInSettingSection) // check is key in current section
	content += genTpl(nmSection, tplGetDefaultValue)       // get default value
	content += genTpl(nmSection, tplGeneralGetterJSON)     // general json getter
	content += genTpl(nmSection, tplGeneralSetterJSON)     // general json setter
	content += genTpl(nmSection, tplCheckExists)           // check if key exists
	content += genTpl(nmSection, tplEnsureNoEmpty)         // ensure section and key exists and not empty
	content += genTpl(nmSection, tplGetter)                // getter
	content += genTpl(nmSection, tplSetter)                // setter
	content += genTpl(nmSection, tplJSONGetter)            // json getter
	content += genTpl(nmSection, tplJSONSetter)            // json setter
	content += genTpl(nmSection, tplLogicJSONSetter)       // logic json setter
	content += genTpl(nmSection, tplRemover)               // remover
	return
}

func genNMGeneralUtilsCode(nmSections []NMSectionStruct) (content string) {
	content += "\n// Origin file name " + nmSettingUtilsFile
	content += genTpl(nmSections, tplGeneralSettingUtils) // general setting utils
	return
}

func genNMVkeyCode(nmSections []NMSectionStruct, nmVkeys []NMVkeyStruct) (content string) {
	content += "\n// Origin file name " + nmSettingVkeyFile
	content += genTpl(nmVkeys, tplVkey)
	return
}

func genNMVsectionCode(nmSections []NMSectionStruct, nmVkeys []NMVsectionStruct) (content string) {
	content += "\n// Origin file name " + nmSettingVsectionFile
	content += genTpl(nmVkeys, tplVsection)
	return
}

func genFrontEndConnPropCode(nmVsections []NMVsectionStruct) (content string) {
	content = genTpl(nmVsections, tplFrontEndConnProp)
	return
}

func genFrontEndSectionCode(nmVsection NMVsectionStruct) (content string) {
	content = genTpl(nmVsection, tplFrontEndSection)
	return
}

func genTpl(data interface{}, tplstr string) (content string) {
	templator := template.New("nm autogen").Funcs(funcMap)
	tpl, err := templator.Parse(tplstr)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBufferString("")
	err = tpl.Execute(buf, data)
	if err != nil {
		panic(err)
	}
	content = string(buf.Bytes())
	return
}

func genBackEndCode() {
	autogenContent := fileHeader

	// back-end code, virtual sections
	autogenContent += genNMVsectionCode(nmSections, nmVsections)

	// back-end code, virtual key
	autogenContent += genNMVkeyCode(nmSections, nmVkeys)

	// back-end code, general setting utils
	autogenContent += genNMGeneralUtilsCode(nmSections)

	// back-end code, echo nm setting sections
	for _, nmSection := range nmSections {
		autogenContent += genNMSettingCode(nmSection)
	}

	writeOrDisplayResultForBackEnd(backEndFile, autogenContent)
}

func genFrontEndCode() {
	// front-end code, BaseConnectionProperties.qml
	autogenContent := genFrontEndConnPropCode(nmVsections)
	writeOrDisplayResultForFrontEnd(frontEndConnPropFile, autogenContent)

	for _, nmVsection := range nmVsections {
		if nmVsection.Ignore {
			continue
		}
		autogenContent = genFrontEndSectionCode(nmVsection)
		frontEndFile := getFrontEndFilePath(nmVsection.Name)
		writeOrDisplayResultForFrontEnd(frontEndFile, autogenContent)
	}
}

func main() {
	flag.BoolVar(&argWriteOutput, "w", false, "write to file")
	flag.BoolVar(&argWriteOutput, "write", false, "write to file")
	flag.BoolVar(&argGenBackEndGo, "gen-back-end-go", false, "generate back-end go code")
	flag.BoolVar(&argGenFrontEndQt, "gen-front-end-qt", false, "generate front-end qt code")
	flag.BoolVar(&argGenFrontEndQml, "gen-front-end-qml", false, "generate front-end qml code")
	flag.StringVar(&argBackEndGoDir, "back-end-go-dir", "..", "go back-end directory")
	flag.StringVar(&argFrontEndQtDir, "front-end-qt-dir", "./front_end_qt_gen", "qt front-end directory")
	flag.StringVar(&argFrontEndQmlDir, "front-end-qml-dir", "./front_end_qml_gen", "qml front-end directory")
	flag.Parse()

	setupDirs()

	unmarshalJSONFile(nmSettingsJSONFile, &nmSections)
	unmarshalJSONFile(nmSettingVkeyJSONFile, &nmVkeys)
	unmarshalJSONFile(nmSettingVsectionJSONFile, &nmVsections)
	if argGenBackEndGo {
		genBackEndCode()
	}
	if argGenFrontEndQml {
		genFrontEndCode()
	}
}
