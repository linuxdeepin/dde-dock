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
	secret "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.secrets"

	"pkg.deepin.io/lib/dbus1"
)

// keep keyring tags same with nm-applet
const (
	keyringTagUUID = "connection-uuid"
	keyringTagSN   = "setting-name"
	keyringTagSK   = "setting-key"
)

func secretNewService() (service *secret.Service, err error) {
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	service = secret.NewService(sessionBus)
	return
}

func secretNewCollection(path dbus.ObjectPath) (collection *secret.Collection, err error) {
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	collection, err = secret.NewCollection(sessionBus, path)
	if err != nil {
		logger.Error(err)
	}
	return
}

func secretNewSession(path dbus.ObjectPath) (session *secret.Session, err error) {
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	session, err = secret.NewSession(sessionBus, path)
	if err != nil {
		logger.Error(err)
	}
	return
}

func secretNewItem(path dbus.ObjectPath) (item *secret.Item, err error) {
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}
	item, err = secret.NewItem(sessionBus, path)
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

	input := dbus.MakeVariant("")
	_, sessionPath, err = service.OpenSession(0, "plain", input)
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

	err = session.Close(0)
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

	keyValues = make(map[string]string)

	attributes := map[string]string{
		keyringTagUUID: uuid,
		keyringTagSN:   settingName,
	}

	unlockedPaths, lockedPaths, err := service.SearchItems(0, attributes)
	if err != nil {
		logger.Error(err)
	}
	reUnlockedPaths, _, _ := service.Unlock(0, lockedPaths)

	doSecretGetAll(service, sessionPath, unlockedPaths, keyValues)
	doSecretGetAll(service, sessionPath, reUnlockedPaths, keyValues)
	if len(keyValues) > 0 {
		ok = true
	}
	return
}
func doSecretGetAll(service *secret.Service, sessionPath dbus.ObjectPath, itemPaths []dbus.ObjectPath, values map[string]string) {
	secretsMap, err := service.GetSecrets(0, itemPaths, sessionPath)
	if err != nil {
		logger.Error(err)
	}
	for itemPath, itemSecret := range secretsMap {
		item, err := secretNewItem(itemPath)
		if err != nil {
			continue
		}

		attributes, _ := item.Attributes().Get(0)
		if keyName, ok := attributes[keyringTagSK]; ok {
			values[keyName] = string(itemSecret.Value)
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

	collectionPath, err := service.ReadAlias(0, "default")
	if err != nil {
		return
	}

	unlockedPaths, _, err := service.Unlock(0, []dbus.ObjectPath{collectionPath})
	if err != nil {
		return
	}
	collectionPath = unlockedPaths[0]

	collection, err := secretNewCollection(collectionPath)
	if err != nil {
		return
	}

	properties := map[string]dbus.Variant{
		"org.freedesktop.Secret.Item.Label": dbus.MakeVariant("NetworkManager password secret"),
		"org.freedesktop.Secret.Item.Type":  dbus.MakeVariant("org.freedesktop.Secret.Generic"),
		"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
			keyringTagUUID: uuid,
			keyringTagSN:   settingName,
			keyringTagSK:   settingKey,
		}),
	}

	itemSecret := secret.Secret{
		Session:     sessionPath,
		Parameters:  []byte{},
		Value:       []byte(value),
		ContentType: "text/plain",
	}

	_, _, err = collection.CreateItem(0, properties, itemSecret, true)
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

	attributes := map[string]string{
		keyringTagUUID: uuid,
	}

	unlockedPaths, lockedPaths, err := service.SearchItems(0, attributes)
	if err != nil {
		logger.Error(err)
	}
	reUnlockedPaths, _, _ := service.Unlock(0, lockedPaths)

	doSecretDeleteAll(service, sessionPath, unlockedPaths)
	doSecretDeleteAll(service, sessionPath, reUnlockedPaths)
}
func doSecretDeleteAll(service *secret.Service, sessionPath dbus.ObjectPath, itemPaths []dbus.ObjectPath) {
	for _, itemPath := range itemPaths {
		item, err := secretNewItem(itemPath)
		if err != nil {
			continue
		}
		item.Delete(0)
	}
}
