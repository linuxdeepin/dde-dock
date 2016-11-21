#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import json
import time

from dbus_gen.com_deepin_daemon_Network import Network
from dbus_gen.com_deepin_daemon_Network_ConnectionSession import ConnectionSession

dbus_network = Network('com.deepin.daemon.Network', '/com/deepin/daemon/Network')

def get_network_devices():
    return json.loads(dbus_network.Devices)

def get_active_connections():
    return json.loads(dbus_network.ActiveConnections)

def get_default_wired_device():
    devices = get_network_devices()
    wired_devices = devices.get('wired')
    if wired_devices and len(wired_devices) > 0:
        return wired_devices[0]['Path']
    else:
        return None

def get_default_wireless_device():
    devices = get_network_devices()
    wireless_devices = devices.get('wireless')
    if wireless_devices and len(wireless_devices) > 0:
        return wireless_devices[0]['Path']
    else:
        return None

def create_connection(conn_type, device_path):
    return dbus_network.CreateConnection(conn_type, device_path)

def is_connection_connected(uuid):
    active_connections = get_active_connections()
    for active_path, active_value in active_connections.items():
        if active_value.get('Uuid') == uuid and active_value.get('State') == 2:
            return True
    return False

def disconnect_default_wired_device():
    dbus_network.DisconnectDevice(get_default_wired_device())
    time.sleep(2)

def connect_default_wired_device():
    dbus_network.ActivateConnection(dbus_network.GetWiredConnectionUuid(get_default_wired_device()),
                                    get_default_wired_device())
    time.sleep(2)

def disconnect_default_wireless_device():
    dbus_network.DisconnectDevice(get_default_wireless_device())
    time.sleep(2)

def test_active_connection(testcase, uuid, device_path, delete_conn = True):
    path = dbus_network.ActivateConnection(uuid, "/")
    testcase.assertIsNotNone(path)
    time.sleep(10) # wait for connection connected
    testcase.assertTrue(is_connection_connected(uuid))
    dbus_network.DeactivateConnection(uuid)
    if delete_conn:
        dbus_network.DeleteConnection(uuid)
    time.sleep(5) # wait for connection deleted
