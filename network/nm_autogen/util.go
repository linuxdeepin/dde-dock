package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

func writeOrDisplayResultForBackEnd(file, content string) {
	if argWriteOutput {
		writeBackEndFile(file, content)
	} else {
		fmt.Println(content)
		fmt.Println()
	}
}

func writeOrDisplayResultForFrontEnd(file, content string) {
	if argWriteOutput {
		writeFrontEndFile(file, content)
	} else {
		fmt.Println(content)
		fmt.Println()
	}
}

func writeBackEndFile(file, content string) {
	// write to .go file and execute gofmt
	err := ioutil.WriteFile(file, []byte(content), 0644)
	if err != nil {
		fmt.Println("error, write file failed:", err)
		return
	}
	execAndWait(10, "gofmt", "-w", file)
	fmt.Println(file)
}

func writeFrontEndFile(file, content string) {
	// write to .go file and execute gofmt
	err := ioutil.WriteFile(file, []byte(content), 0644)
	if err != nil {
		fmt.Println("error, write file failed:", err)
		return
	}
	fmt.Println(file)
}

func GetAllKeysInVsection(vsectionName string) (keys []string) {
	vsectionInfo := getVsectionInfo(vsectionName)
	for _, section := range vsectionInfo.RelatedSections {
		keys = appendStrArrayUnique(keys, GetAllKeysInSection(section)...)
	}
	return
}

// GetAllKeysInSection return all keys that will be used by front-end.
func GetAllKeysInSection(sectionName string) (keys []string) {
	// virtual keys in section that with none related key
	for _, vk := range nmVkeys {
		if vk.RelatedSection == sectionName && isStringInArray("NM_SETTING_VK_NONE_RELATED_KEY", vk.RelatedKeys) {
			keys = append(keys, vk.Name)
		}
	}
	for _, nmSection := range nmSections {
		if nmSection.Name == sectionName {
			for _, k := range nmSection.Keys {
				vksNames := getRelatedVks(k.Name)
				if len(vksNames) > 0 {
					// if virtual key is a enable wrapper, both
					// virtual key and real key will be appended.
					for _, vkName := range vksNames {
						keys = appendStrArrayUnique(keys, vkName)
						if getVkInfo(vkName).EnableWrapper {
							keys = appendStrArrayUnique(keys, k.Name)
						}
					}
				} else {
					keys = append(keys, k.Name)
				}
			}
			break
		}
	}
	return
}

// get all related virtual keys of real key
func getRelatedVks(keyName string) (vks []string) {
	for _, vk := range nmVkeys {
		if isStringInArray(keyName, vk.RelatedKeys) {
			vks = append(vks, vk.Name)
		}
	}
	return
}

func ToKeyDisplayName(keyName string) (displayName string) {
	var keyValue string
	if isVk(keyName) {
		vkInfo := getVkInfo(keyName)
		displayName = vkInfo.DisplayName
		keyValue = vkInfo.Value
		if displayName == "<default>" {
			keyInfo := getKeyInfo(vkInfo.RelatedKeys[0])
			displayName = keyInfo.DisplayName
		}
	} else {
		keyInfo := getKeyInfo(keyName)
		displayName = keyInfo.DisplayName
		keyValue = keyInfo.Value
	}
	if displayName != "<default>" {
		return
	}
	return "!!" + keyValue
}

// "NM_SETTING_802_1X_EAP" -> "eap"
func ToKeyValue(keyName string) (keyValue string) {
	if isVk(keyName) {
		keyValue = getVkInfo(keyName).Value
	} else {
		keyValue = getKeyInfo(keyName).Value
	}
	return
}

// "NM_SETTING_802_1X_EAP" -> "802-1x"
func ToKeyRelatedSectionValue(keyName string) (sectionValue string) {
	if isVk(keyName) {
		sectionName := getVkInfo(keyName).RelatedSection
		sectionValue = ToSectionValue(sectionName)
	} else {
		sectionValue = getKeyRelatedSectionInfo(keyName).Value
	}
	return
}
func getKeyRelatedSectionInfo(keyName string) (sectionInfo NMSectionStruct) {
	for _, section := range nmSections {
		for _, key := range section.Keys {
			if key.Name == keyName {
				sectionInfo = section
				return
			}
		}
	}
	fmt.Println("invalid key name", keyName)
	os.Exit(1)
	return
}

// "NM_SETTING_802_1X_SETTING_NAME" -> "802-1x"
func ToSectionValue(sectionName string) (sectionValue string) {
	// TODO NM_SETTING_VK_NONE_RELATED_FIELD
	sectionInfo := getSectionInfo(sectionName)
	return sectionInfo.Value
}

func GetKeyWidgetProp(keyName string) (prop map[string]string) {
	if isVk(keyName) {
		prop = getVkInfo(keyName).WidgetProp
	} else {
		prop = getKeyInfo(keyName).WidgetProp
	}
	return
}

func IsKeyUsedByFrontEnd(keyName string) (used bool) {
	if isVk(keyName) {
		used = getVkInfo(keyName).UsedByFrontEnd
	} else {
		used = getKeyInfo(keyName).UsedByFrontEnd
	}
	return
}

func getSectionInfo(sectionName string) (sectionInfo NMSectionStruct) {
	for _, section := range nmSections {
		if section.Name == sectionName {
			sectionInfo = section
			return
		}
	}
	fmt.Println("invalid section name", sectionName)
	os.Exit(1)
	return
}

func getVsectionInfo(vsectionName string) (vsectionInfo NMVsectionStruct) {
	for _, vsection := range nmVsections {
		if vsection.Name == vsectionName {
			vsectionInfo = vsection
			return
		}
	}
	fmt.Println("invalid vsection name", vsectionName)
	os.Exit(1)
	return
}

func getKeyInfo(keyName string) (keyInfo NMKeyStruct) {
	for _, section := range nmSections {
		for _, key := range section.Keys {
			if key.Name == keyName {
				keyInfo = key
				return
			}
		}
	}
	fmt.Println("invalid key name", keyName)
	os.Exit(1)
	return
}

func getVkInfo(vkName string) (vkInfo NMVkeyStruct) {
	for _, vk := range nmVkeys {
		if vk.Name == vkName {
			vkInfo = vk
			return
		}
	}
	fmt.Println("invalid key name", vkName)
	os.Exit(1)
	return
}

// check if target key is a virtual key
func isVk(keyName string) (ok bool) {
	for _, vk := range nmVkeys {
		if vk.Name == keyName {
			return true
		}
	}
	return false
}

// NM_SETTING_CONNECTION_SETTING_NAME -> ConnectionSetting, NM_SETTING_VK_VPN_L2TP_SETTING_NAME -> VpnL2tp
func ToSectionFuncBaseName(name string) (funcName string) {
	name = strings.Replace(name, "NM_SETTING_ALIAS_", "NM_SETTING_", -1) // remove virtual section tag
	funcName = strings.TrimPrefix(name, "NM_")
	funcName = strings.TrimSuffix(funcName, "_SETTING_NAME")
	funcName = strings.Replace(funcName, "_", " ", -1)
	funcName = ToCaplitalize(funcName)
	funcName = strings.Replace(funcName, " ", "", -1)
	return
}

// NM_SETTING_CONNECTION_ID -> SettingConnectionId
func ToKeyFuncBaseName(name string) (funcName string) {
	funcName = strings.TrimPrefix(name, "NM_")
	funcName = strings.Replace(funcName, "_", " ", -1)
	funcName = ToCaplitalize(funcName)
	funcName = strings.Replace(funcName, " ", "", -1)
	return
}

// "hello world" -> "Hello World", "HELLO WORLD" -> "Hello World"
func ToCaplitalize(str string) (capstr string) {
	scaner := bufio.NewScanner(strings.NewReader(str))
	scaner.Split(bufio.ScanWords)
	for scaner.Scan() {
		word := scaner.Text()
		if len(word) > 1 {
			capstr = capstr + " " + strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		} else if len(word) == 1 {
			capstr = capstr + " " + strings.ToUpper(word)
		}
	}
	capstr = strings.TrimSpace(capstr)
	return
}

// "vk-key-mgmt" -> "VkKeyMgmt", "vk key mgmt" -> "VkKeyMgmt", "vk_key_mgmt" -> "VkKeyMgmt"
func ToClassName(str string) (id string) {
	id = strings.Replace(str, "_", " ", -1)
	id = strings.Replace(id, "-", " ", -1)
	id = ToCaplitalize(id)
	id = strings.Replace(id, " ", "", -1)
	return
}

// "vs-section" -> "Section"
func ToVsClassName(str string) (id string) {
	str = strings.TrimPrefix(str, "vs-")
	return ToClassName(str)
}

func unmarshalJSONFile(jsonFile string, data interface{}) {
	fileContent, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		fmt.Println("error, open file failed:", err)
		os.Exit(1)
	}

	err = json.Unmarshal(fileContent, data)
	if err != nil {
		fmt.Printf("error, unmarshal json %s failed: %v\n", jsonFile, err)
		os.Exit(1)
	}
}

func execAndWait(timeout int, name string, arg ...string) (stdout, stderr string, err error) {
	cmd := exec.Command(name, arg...)
	var bufStdout, bufStderr bytes.Buffer
	cmd.Stdout = &bufStdout
	cmd.Stderr = &bufStderr
	err = cmd.Start()
	if err != nil {
		return
	}

	// wait for process finished
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		if err = cmd.Process.Kill(); err != nil {
			return
		}
		<-done
		err = fmt.Errorf("time out and process was killed")
	case err = <-done:
		stdout = bufStdout.String()
		stderr = bufStderr.String()
		if err != nil {
			return
		}
	}
	return
}

func isStringInArray(s string, list []string) bool {
	for _, i := range list {
		if i == s {
			return true
		}
	}
	return false
}

func appendStrArrayUnique(a1 []string, a2 ...string) (a []string) {
	a = a1
	for _, s := range a2 {
		if !isStringInArray(s, a) {
			a = append(a, s)
		}
	}
	return
}
