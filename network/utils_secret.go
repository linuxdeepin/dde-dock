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

package network

import (
	"dbus/org/freedesktop/secret"
	"pkg.deepin.io/lib/dbus"
)

// keep keyring tags same with nm-applet
const (
	keyringTagUUID = "connection-uuid"
	keyringTagSN   = "setting-name"
	keyringTagSK   = "setting-key"
)

type itemSecretStruct struct {
	Session     dbus.ObjectPath
	Parameters  []byte
	Value       []byte
	ContentType string `dbus:"content_type"`
}

// convert []interface{} from DBus to the really item secret struct.
// A []interface{} example is: []interface{}{
// "/org/freedesktop/secrets/session/s16", []byte{}, []byte{0x70,
// 0x61}, "text/plain" }⏎
func convertToItemSecret(ifcArray []interface{}) (secret itemSecretStruct, ok bool) {
	if len(ifcArray) < 4 {
		ok = false
		return
	}
	secret = itemSecretStruct{}
	if secret.Session, ok = ifcArray[0].(dbus.ObjectPath); !ok {
		return
	}
	if secret.Parameters, ok = ifcArray[1].([]byte); !ok {
		return
	}
	if secret.Value, ok = ifcArray[2].([]byte); !ok {
		return
	}
	if secret.ContentType, ok = ifcArray[3].(string); !ok {
		return
	}
	return
}

func secretNewService() (service *secret.Service, err error) {
	service, err = secret.NewService("org.freedesktop.secrets", "/org/freedesktop/secrets")
	if err != nil {
		logger.Error(err)
	}
	return
}

func secretNewCollection(path dbus.ObjectPath) (collection *secret.Collection, err error) {
	collection, err = secret.NewCollection("org.freedesktop.secrets", path)
	if err != nil {
		logger.Error(err)
	}
	return
}

func secretNewSession(path dbus.ObjectPath) (session *secret.Session, err error) {
	session, err = secret.NewSession("org.freedesktop.secrets", path)
	if err != nil {
		logger.Error(err)
	}
	return
}

func secretNewItem(path dbus.ObjectPath) (item *secret.Item, err error) {
	item, err = secret.NewItem("org.freedesktop.secrets", path)
	if err != nil {
		logger.Error(err)
	}
	return
}

func secretServiceOpenSession() (sessionPath dbus.ObjectPath, err error) {
	service, err := secretNewService()
	if err != nil {
		return
	}
	defer secret.DestroyService(service)

	input := dbus.MakeVariant("")
	_, sessionPath, err = service.OpenSession("plain", input)
	if err != nil {
		logger.Error(err)
	}
	return
}

var secretSessionInstance dbus.ObjectPath

func secretGetSessionInstance() dbus.ObjectPath {
	if len(secretSessionInstance) != 0 {
		return secretSessionInstance
	}
	secretSessionInstance, err := secretServiceOpenSession()
	if err != nil {
		return ""
	}
	return secretSessionInstance
}

func secretCloseSession() {
	if len(secretSessionInstance) == 0 {
		return
	}
	session, err := secretNewSession(secretSessionInstance)
	if err != nil {
		return
	}

	err = session.Close()
	if err != nil {
		logger.Error(err)
	}
}

// secretGet get secret value from keyring for target key in network-manager configuration
func secretGet(uuid, settingName, settingKey string) (value string, ok bool) {
	if values, okNest := secretGetAll(uuid, settingName); okNest {
		value, ok = values[settingKey]
	}
	return
}

// secretGetAll get all secret values from keyring for network-manager configuration
func secretGetAll(uuid, settingName string) (keyValues map[string]string, ok bool) {
	sessionPath := secretGetSessionInstance()

	service, err := secretNewService()
	if err != nil {
		ok = false
		return
	}
	defer secret.DestroyService(service)

	keyValues = make(map[string]string)

	attributes := map[string]string{
		keyringTagUUID: uuid,
		keyringTagSN:   settingName,
	}

	unlockedPaths, lockedPaths, err := service.SearchItems(attributes)
	if err != nil {
		logger.Error(err)
	}
	reUnlockedPaths, _, _ := service.Unlock(lockedPaths)

	doSecretGetAll(service, sessionPath, unlockedPaths, keyValues)
	doSecretGetAll(service, sessionPath, reUnlockedPaths, keyValues)
	if len(keyValues) > 0 {
		ok = true
	}
	return
}
func doSecretGetAll(service *secret.Service, sessionPath dbus.ObjectPath, itemPaths []dbus.ObjectPath, values map[string]string) {
	secretsMap, err := service.GetSecrets(itemPaths, sessionPath)
	if err != nil {
		logger.Error(err)
	}
	for itemPath, itemSecretIfcArray := range secretsMap {
		if itemSecret, ok := convertToItemSecret(itemSecretIfcArray); ok {
			item, err := secretNewItem(itemPath)
			if err != nil {
				continue
			}
			defer secret.DestroyItem(item)

			attributes := item.Attributes.Get()
			if keyName, ok := attributes[keyringTagSK]; ok {
				values[keyName] = string(itemSecret.Value)
			}
		}
	}
	return
}

// secretSet update or new secret item to keyring for network-manager configuration
func secretSet(uuid, settingName, settingKey string, value string) {
	sessionPath := secretGetSessionInstance()

	service, err := secretNewService()
	if err != nil {
		return
	}
	defer secret.DestroyService(service)

	collectionPath, err := service.ReadAlias("default")
	if err != nil {
		return
	}

	unlockedPaths, _, err := service.Unlock([]dbus.ObjectPath{collectionPath})
	if err != nil {
		return
	}
	collectionPath = unlockedPaths[0]

	collection, err := secretNewCollection(collectionPath)
	if err != nil {
		return
	}
	defer secret.DestroyCollection(collection)

	properties := map[string]dbus.Variant{
		"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("NetworkManager password secret"),
		"org.freedesktop.Secret.Item.Type":  dbus.MakeVariant("org.freedesktop.Secret.Generic"),
		"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
			keyringTagUUID: uuid,
			keyringTagSN:   settingName,
			keyringTagSK:   settingKey,
		}),
	}

	itemSecret := itemSecretStruct{
		Session:     sessionPath,
		Parameters:  []byte{},
		Value:       []byte(value),
		ContentType: "text/plain",
	}

	_, _, err = collection.CreateItem(properties, itemSecret, true)
	if err != nil {
		logger.Error(err)
	}

	return
}

// secretDeleteAll delete all secret keyring items for network-manager configuration
func secretDeleteAll(uuid string) {
	sessionPath := secretGetSessionInstance()

	service, err := secretNewService()
	if err != nil {
		return
	}
	defer secret.DestroyService(service)

	attributes := map[string]string{
		keyringTagUUID: uuid,
	}

	unlockedPaths, lockedPaths, err := service.SearchItems(attributes)
	if err != nil {
		logger.Error(err)
	}
	reUnlockedPaths, _, _ := service.Unlock(lockedPaths)

	doSecretDeleteAll(service, sessionPath, unlockedPaths)
	doSecretDeleteAll(service, sessionPath, reUnlockedPaths)
}
func doSecretDeleteAll(service *secret.Service, sessionPath dbus.ObjectPath, itemPaths []dbus.ObjectPath) {
	for _, itemPath := range itemPaths {
		item, err := secretNewItem(itemPath)
		if err != nil {
			continue
		}
		defer secret.DestroyItem(item)
		item.Delete()
	}
}
