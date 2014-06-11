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
	"ToKeyRelatedSectionValue":    ToKeyRelatedSectionValue,
	"IsKeyUsedByFrontEnd":         IsKeyUsedByFrontEnd,
	"ToFrontEndWidget":            ToFrontEndWidget,
	"ToClassName":                 ToClassName,
	"ToVsClassName":               ToVsClassName,
	"GetAllKeysInVsection":        GetAllKeysInVsection,
	"GetKeyWidgetProp":            GetKeyWidgetProp,
	"ToKeyTypeInterfaceConverter": ToKeyTypeInterfaceConverter,
	"IsEnableWrapperVkey":         IsEnableWrapperVkey,
}

const (
	backEndDir                = ".."
	frontEndDir               = "../../../dss/modules/network/edit_autogen/"
	nmSettingsJSONFile        = "./nm_settings.json"
	nmSettingVkeyJSONFile     = "./nm_setting_vkey.json"
	nmSettingVsectionJSONFile = "./nm_setting_vsection.json"
)

var (
	argWriteOutput        bool
	argBackEnd            bool
	argFrontEnd           bool
	nmSettingUtilsFile    = path.Join(backEndDir, "nm_setting_general_autogen.go")
	nmSettingVkeyFile     = path.Join(backEndDir, "nm_setting_virtual_key_autogen.go")
	nmSettingVsectionFile = path.Join(backEndDir, "nm_setting_virtual_section_autogen.go")
	frontEndConnPropFile  = path.Join(frontEndDir, "BaseConnectionEdit.qml")
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
	FrontEndWidget string            // such as "EditLinePasswordInput"
	WidgetProp     map[string]string // properties for front end widget, such as "WidgetProp":{"alwaysUpdate":"true"}
}

type NMVkeyStruct struct {
	Name           string   // such as "NM_SETTING_VK_802_1X_EAP"
	Value          string   // such as "vk-eap"
	Type           string   // such as "ktypeString"
	VkType         string   // could be "vkTypeWrapper", "vkTypeEnableWrapper" or "vkTypeController"
	RelatedSection string   // such as "NM_SETTING_802_1X_SETTING_NAME"
	RelatedKeys    []string // such as "NM_SETTING_802_1X_EAP"
	UsedByFrontEnd bool     // check if is used by front-end
	Optional       bool     // if key is optional, will ignore error for it
	DisplayName    string
	FrontEndWidget string            // such as "EditLinePasswordInput"
	WidgetProp     map[string]string // properties for front end widget, such as "WidgetProp":{"alwaysUpdate":"true"}
}

func genNMSettingCode(nmSection NMSectionStruct) (content string) {
	content = fileHeader
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
	content = genTpl(nmSections, tplGeneralSettingUtils) // general setting utils
	return
}

func genNMVkeyCode(nmSections []NMSectionStruct, nmVkeys []NMVkeyStruct) (content string) {
	content = genTpl(nmVkeys, tplVkey)
	return
}

func genNMVsectionCode(nmSections []NMSectionStruct, nmVkeys []NMVsectionStruct) (content string) {
	content = genTpl(nmVkeys, tplVsection)
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
	// back-end code, echo nm setting sections
	for _, nmSection := range nmSections {
		autogenContent := genNMSettingCode(nmSection)
		backEndFile := getBackEndFilePath(nmSection.Name)
		writeOrDisplayResultForBackEnd(backEndFile, autogenContent)
	}

	// back-end code, general setting utils
	autogenContent := genNMGeneralUtilsCode(nmSections)
	writeOrDisplayResultForBackEnd(nmSettingUtilsFile, autogenContent)

	// back-end code, virtual key
	autogenContent = genNMVkeyCode(nmSections, nmVkeys)
	writeOrDisplayResultForBackEnd(nmSettingVkeyFile, autogenContent)

	// back-end code, virtual sections
	autogenContent = genNMVsectionCode(nmSections, nmVsections)
	writeOrDisplayResultForBackEnd(nmSettingVsectionFile, autogenContent)
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
	flag.BoolVar(&argBackEnd, "b", false, "generate back-end code")
	flag.BoolVar(&argFrontEnd, "f", false, "generate front-end code")
	flag.Parse()

	unmarshalJSONFile(nmSettingsJSONFile, &nmSections)
	unmarshalJSONFile(nmSettingVkeyJSONFile, &nmVkeys)
	unmarshalJSONFile(nmSettingVsectionJSONFile, &nmVsections)
	if argBackEnd {
		genBackEndCode()
	}
	if argFrontEnd {
		genFrontEndCode()
	}
}
