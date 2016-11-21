#!/usr/bin/env python3
# -*- coding: utf-8 -*-

# This example will create connection for target wifi and try to connect to it.

import json

import utils_dbus

from dbus_gen.com_deepin_daemon_Network import Network
from dbus_gen.com_deepin_daemon_Network_ConnectionSession import ConnectionSession

dbus_network = Network('com.deepin.daemon.Network', '/com/deepin/daemon/Network')

wifi_ssid = "test"
wifi_psk = "12345678"

# ensure target network device enabled
dbus_network.EnableDevice(utils_dbus.get_default_wireless_device(), True)

session_path = utils_dbus.create_connection('wireless', utils_dbus.get_default_wireless_device())
dbus_session = ConnectionSession('com.deepin.daemon.Network', session_path)
uuid = dbus_session.Uuid
dbus_session.SetKey('802-11-wireless', 'ssid', json.dump(wifi_ssid))
dbus_session.SetKey('802-11-wireless-security', 'vk-key-mgmt', json.dumps("wpa-psk"))
dbus_session.SetKey('802-11-wireless-security', 'psk', json.dump(wifi_psk))
dbus_session.Save()
path = dbus_network.ActivateConnection(uuid, utils_dbus.get_default_wireless_device())
print(path)
