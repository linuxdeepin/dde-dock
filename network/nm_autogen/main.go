package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

var funcMap = template.FuncMap{
	"ToFieldFuncBaseName":       ToFieldFuncBaseName,
	"ToKeyFuncBaseName":         ToKeyFuncBaseName,
	"ToKeyTypeRealData":         ToKeyTypeRealData,
	"ToKeyTypeDefaultValueJSON": ToKeyTypeDefaultValueJSON,
	"IfNeedCheckValueLength":    IfNeedCheckValueLength,
}

var jsonFiles = []string{
	"./nm_setting_802_1x.json",
	"./nm_setting_connection.json",
	"./nm_setting_ipv4.json",
	"./nm_setting_ipv6.json",
	"./nm_setting_wired.json",
	"./nm_setting_wireless.json",
	"./nm_setting_wireless_security.json",
	"./nm_setting_pppoe.json",
	"./nm_setting_ppp.json",
}

type NMSettingStruct struct {
	FieldName  string
	OutputFile string
	Keys       []struct {
		Name    string
		Type    string
		Default string
	}
}

func generateNMSettingCode(nmSetting NMSettingStruct) (content string) {
	content = fileHeader
	content += generateTemplate(nmSetting, tplGetKeyType)          // get key type
	content += generateTemplate(nmSetting, tplIsKeyInSettingField) // check is key in current field
	content += generateTemplate(nmSetting, tplGetDefaultValueJSON) // get default json value
	content += generateTemplate(nmSetting, tplGeneralGetterJSON)   // general json getter
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

// "ktypeString" -> "string", "ktypeArrayByte" -> "[]byte"
func ToKeyTypeRealData(ktype string) (realData string) {
	switch ktype {
	default:
		fmt.Println("invalid ktype:", ktype)
		os.Exit(1)
	case "ktypeString":
		realData = "string"
	case "ktypeByte":
		realData = "byte"
	case "ktypeInt32":
		realData = "int32"
	case "ktypeUint32":
		realData = "uint32"
	case "ktypeUint64":
		realData = "uint64"
	case "ktypeBoolean":
		realData = "bool"
	case "ktypeArrayByte":
		realData = "[]byte"
	case "ktypeArrayString":
		realData = "[]string"
	case "ktypeArrayUint32":
		realData = "[]uint32"
	case "ktypeArrayArrayByte":
		realData = "[][]byte"
	case "ktypeArrayArrayUint32":
		realData = "[][]uint32"
	case "ktypeDictStringString":
		realData = "map[string]string"
	case "ktypeIpv6Addresses":
		realData = "Ipv6Addresses"
	case "ktypeIpv6Routes":
		realData = "Ipv6Routes"
	case "ktypeWrapperString":
		realData = "[]byte"
	case "ktypeWrapperMacAddress":
		realData = "[]byte"
	case "ktypeWrapperIpv4Dns":
		realData = "[]uint32"
	case "ktypeWrapperIpv4Addresses":
		realData = "[][]uint32"
	case "ktypeWrapperIpv4Routes":
		realData = "[][]uint32"
	case "ktypeWrapperIpv6Dns":
		realData = "[][]byte"
	case "ktypeWrapperIpv6Addresses":
		realData = "Ipv6Addresses"
	case "ktypeWrapperIpv6Routes":
		realData = "Ipv6Routes"
	}
	return
}

// "ktypeString" -> `""`, "ktypeBool" -> `false`
func ToKeyTypeDefaultValueJSON(ktype, customValue string) (valueJSON string) {
	if customValue == "<null>" {
		if ktype == "ktypeString" {
			return `""`
		} else {
			return "null"
		}
	} else if customValue != "<default>" {
		return customValue
	}
	switch ktype {
	default:
		fmt.Println("invalid ktype:", ktype)
		os.Exit(1)
	case "ktypeString":
		valueJSON = `""`
	case "ktypeByte":
		valueJSON = `0`
	case "ktypeInt32":
		valueJSON = `0`
	case "ktypeUint32":
		valueJSON = `0`
	case "ktypeUint64":
		valueJSON = `0`
	case "ktypeBoolean":
		valueJSON = `false`
	case "ktypeArrayByte":
		valueJSON = `""`
	case "ktypeArrayString":
		valueJSON = `null`
	case "ktypeArrayUint32":
		valueJSON = `null`
	case "ktypeArrayArrayByte":
		valueJSON = `null`
	case "ktypeArrayArrayUint32":
		valueJSON = `null`
	case "ktypeDictStringString":
		valueJSON = `null`
	case "ktypeIpv6Addresses":
		valueJSON = `null`
	case "ktypeIpv6Routes":
		valueJSON = `null`
	case "ktypeWrapperString":
		valueJSON = `""`
	case "ktypeWrapperMacAddress":
		valueJSON = `""`
	case "ktypeWrapperIpv4Dns":
		valueJSON = `null`
	case "ktypeWrapperIpv4Addresses":
		valueJSON = `null`
	case "ktypeWrapperIpv4Routes":
		valueJSON = `null`
	case "ktypeWrapperIpv6Dns":
		valueJSON = `null`
	case "ktypeWrapperIpv6Addresses":
		valueJSON = `null`
	case "ktypeWrapperIpv6Routes":
		valueJSON = `null`
	}
	return
}

// test if need check value length to ensure value not empty
func IfNeedCheckValueLength(ktype string) (need string) {
	switch ktype {
	default:
		fmt.Println("invalid ktype:", ktype)
		os.Exit(1)
	case "ktypeString":
		need = "t"
	case "ktypeByte":
		need = ""
	case "ktypeInt32":
		need = ""
	case "ktypeUint32":
		need = ""
	case "ktypeUint64":
		need = ""
	case "ktypeBoolean":
		need = ""
	case "ktypeArrayByte":
		need = "t"
	case "ktypeArrayString":
		need = "t"
	case "ktypeArrayUint32":
		need = "t"
	case "ktypeArrayArrayByte":
		need = "t"
	case "ktypeArrayArrayUint32":
		need = "t"
	case "ktypeDictStringString":
		need = "t"
	case "ktypeIpv6Addresses":
		need = "t"
	case "ktypeIpv6Routes":
		need = "t"
	case "ktypeWrapperString":
		need = "t"
	case "ktypeWrapperMacAddress":
		need = "t"
	case "ktypeWrapperIpv4Dns":
		need = "t"
	case "ktypeWrapperIpv4Addresses":
		need = "t"
	case "ktypeWrapperIpv4Routes":
		need = "t"
	case "ktypeWrapperIpv6Dns":
		need = "t"
	case "ktypeWrapperIpv6Addresses":
		need = "t"
	case "ktypeWrapperIpv6Routes":
		need = "t"
	}
	return
}

// NM_SETTING_CONNECTION_SETTING_NAME -> ConnectionSetting
func ToFieldFuncBaseName(name string) (funcName string) {
	funcName = strings.TrimPrefix(name, "NM_")
	funcName = strings.TrimSuffix(funcName, "_SETTING_NAME")
	funcName = strings.Replace(funcName, "_", " ", -1)
	funcName = caplitalizeString(funcName)
	funcName = strings.Replace(funcName, " ", "", -1)
	return
}

// NM_SETTING_CONNECTION_ID -> SettingConnectionId
func ToKeyFuncBaseName(name string) (funcName string) {
	funcName = strings.TrimPrefix(name, "NM_")
	funcName = strings.Replace(funcName, "_", " ", -1)
	funcName = caplitalizeString(funcName)
	funcName = strings.Replace(funcName, " ", "", -1)
	return
}

// "hello world" -> "Hello World", "HELLO WORLD" -> "Hello World"
func caplitalizeString(str string) (capstr string) {
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

func main() {
	var writeOutput bool

	flag.BoolVar(&writeOutput, "w", false, "write to output file")
	flag.Parse()

	for _, f := range jsonFiles {
		fileContent, err := ioutil.ReadFile(f)
		if err != nil {
			fmt.Println("error, open file failed:", err)
			continue
		}

		var nmSetting NMSettingStruct
		err = json.Unmarshal(fileContent, &nmSetting)
		if err != nil {
			fmt.Printf("error, unmarshal json %s failed: %v\n", f, err)
			continue
		}

		autogenContent := generateNMSettingCode(nmSetting)
		if writeOutput {
			// write to file and execute gofmt
			err = ioutil.WriteFile(nmSetting.OutputFile, []byte(autogenContent), 0644)
			if err != nil {
				fmt.Println("error, write file failed:", err)
				continue
			}
			execAndWait(10, "gofmt", "-w", nmSetting.OutputFile)
		} else {
			fmt.Println(autogenContent)
			fmt.Println()
		}
	}
}
