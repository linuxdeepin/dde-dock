/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

package main

import (
	"bytes"
	"flag"
	"fmt"
	"text/template"
)

const (
	nmConstsYamlFile             = "./nm_consts_gen.yml"
	nmConstsKeysOverrideYamlFile = "./nm_consts_keys_override.yml"
	nmVpnAliasSettingsYamlFile   = "./nm_vpn_alias_settings.yml"
	nmVirtualSettingYamlFile     = "./nm_virtual_sections.yml"
	nmLogicSetKeysYamlFile       = "./nm_logicset_keys.yml"
)

var argOutputConstsFile string
var argOutputBeansFile string
var argTestMode bool

var funcMap = template.FuncMap{
	"UnwrapInterface":              UnwrapInterface,
	"GetKeyTypeGoSyntax":           GetKeyTypeGoSyntax,
	"GetKeyTypeGoIfcConverterFunc": GetKeyTypeGoIfcConverterFunc,
	"GetKeyDefaultValue":           GetKeyDefaultValue,
	"IfNeedCheckValueLength":       IfNeedCheckValueLength,
	"IsLogicSetKey":                IsLogicSetKey,
	"GetKeyFuncBaseName":           GetKeyFuncBaseName,
	"GetKeyTypeShortName":          GetKeyTypeShortName,
	"GetVsRelatedSettings":         GetVsRelatedSettings,
}

var nmConsts nmConstsStruct
var nmOverrideKeys []nmSettingKey

type nmConstsStruct struct {
	NMEnums []struct {
		EnumClass string `yaml:"EnumClass"`
		Members   []struct {
			Name  string      `yaml:"Name"`
			Value interface{} `yaml:"Value"`
		} `yaml:"Members"`
	} `yaml:"NMEnums"`
	NMSettings []nmSetting `yaml:"NMSettings"`
}

type nmSetting struct {
	SettingClass    string          `yaml:"SettingClass"`
	Name            string          `yaml:"Name"`
	RealSettingName string          `yaml:"RealSettingName,omitempty"` // only used for alias settings
	Value           string          `yaml:"Value"`
	Keys            []*nmSettingKey `yaml:"Keys"`
}
type nmSettingKey struct {
	KeyName      string `yaml:"KeyName"`
	Value        string `yaml:"Value"`
	CapcaseName  string `yaml:"CapcaseName"`
	Type         string `yaml:"Type"`
	DefaultValue string `yaml:"DefaultValue,omitempty"`
}

var nmVpnAliasSettings []nmSetting

var nmVirtualSections []nmVirtualSection

type nmVirtualSection struct {
	VirtaulSectionName string `yaml:"VirtaulSectionName"`
	Value              string `yaml:"Value"`
	DisplayName        string `yaml:"DisplayName"`
	Expanded           bool   `yaml:"Expanded"`
	Keys               []struct {
		KeyValue      string `yaml:"KeyValue"`
		Section       string `yaml:"Section"`
		DisplayName   string `yaml:"DisplayName"`
		WidgetType    string `yaml:"WidgetType"`
		AlwaysUpdate  bool   `yaml:"AlwaysUpdate,omitempty"`
		UseValueRange bool   `yaml:"UseValueRange,omitempty"`
		MinValue      int    `yaml:"MinValue,omitempty"`
		MaxValue      int    `yaml:"MaxValue,omitempty"`
		VKeyInfo      struct {
			VirtualKeyName string   `yaml:"VirtualKeyName"`
			Type           string   `yaml:"Type"`
			VkType         string   `yaml:"VkType"`
			RelatedKeys    []string `yaml:"RelatedKeys"`
			ChildKey       bool     `yaml:"ChildKey"`
			Optional       bool     `yaml:"Optional"`
		} `yaml:"VKeyInfo,omitempty"`
	} `yaml:"Keys"`
}

var nmLogicSetKeys []string

func genNMConstsCode() (content string) {
	content = nmConstsHeader
	content += genTpl(nmVirtualSections, tplNMVirtualConsts)
	content += genTpl(nmConsts, tplNMConsts)
	return
}

func genNMBeansCode() (content string) {
	content = nmSettingBeansHeader
	content += genTpl(nmConsts, tplNMBeans)
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

func main() {
	flag.StringVar(&argOutputConstsFile, "output-consts", "../nm/nm_consts_gen.go", "generate to networkmanager constss .go file")
	flag.StringVar(&argOutputBeansFile, "output-beans", "../nm_setting_beans_gen.go", "generate to networkmanager setting getter and setter beans .go file")
	flag.BoolVar(&argTestMode, "test", false, "test mode, output console instead writing to file")
	flag.Parse()

	yamlUnmarshalFile(nmConstsYamlFile, &nmConsts)
	yamlUnmarshalFile(nmConstsKeysOverrideYamlFile, &nmOverrideKeys)
	mergeOverrideKeys()

	yamlUnmarshalFile(nmVpnAliasSettingsYamlFile, &nmVpnAliasSettings)
	for _, setting := range nmVpnAliasSettings {
		nmConsts.NMSettings = append(nmConsts.NMSettings, setting)
	}

	yamlUnmarshalFile(nmVirtualSettingYamlFile, &nmVirtualSections)
	yamlUnmarshalFile(nmLogicSetKeysYamlFile, &nmLogicSetKeys)
	yamlUnmarshalFile(nmLogicSetKeysYamlFile, &nmLogicSetKeys)

	content := genNMConstsCode()
	if argTestMode {
		fmt.Println(content)
	} else {
		writeOutputFile(argOutputConstsFile, content)
	}

	content = genNMBeansCode()
	if argTestMode {
		fmt.Println(content)
	} else {
		writeOutputFile(argOutputBeansFile, content)
	}
}
