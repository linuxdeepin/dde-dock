// Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "adapter.h"
#include "device.h"

#include <QJsonObject>
#include <QJsonDocument>
#include <QJsonArray>
#include <QJsonValue>

Adapter::Adapter(QObject *parent)
    : QObject(parent)
    , m_id("")
    , m_name("")
    , m_powered(false)
    , m_current(false)
    , m_discover(false)
{
}

void Adapter::setName(const QString &name)
{
    if (name != m_name) {
        m_name = name;
        Q_EMIT nameChanged(name);
    }
}

void Adapter::addDevice(const QJsonObject &deviceObj)
{
    const QString id = deviceObj["Path"].toString();
    const QString name = deviceObj["Name"].toString();
    const QString alias = deviceObj["Alias"].toString();
    const bool paired = deviceObj["Paired"].toBool();
    const int rssi = deviceObj["RSSI"].toInt();
    const Device::State state = Device::State(deviceObj["State"].toInt());
    const bool connectState = deviceObj["ConnectState"].toBool();
    const QString bluetoothDeviceType = deviceObj["Icon"].toString();
    const int battery = deviceObj["Battery"].toInt();

    removeDevice(id);

    auto device = new Device(this);

    device->setId(id);
    device->setName(name);
    device->setAlias(alias);
    device->setPaired(paired);
    device->setState(state);
    device->setConnectState(connectState);
    device->setRssi(rssi);
    device->setAdapterId(m_id);
    device->setDeviceType(bluetoothDeviceType);
    device->setBattery(battery);

    m_devices[id] = device;

    emit deviceAdded(device);
}

void Adapter::removeDevice(const QString &deviceId)
{
    const Device *constDevice = m_devices.value(deviceId);
    auto device = const_cast<Device *>(constDevice);
    if (device) {
        m_devices.remove(deviceId);
        emit deviceRemoved(device);
        delete device;
    }
}

void Adapter::updateDevice(const QJsonObject &dviceJson)
{
    const QString id = dviceJson["Path"].toString();
    const QString name = dviceJson["Name"].toString();
    const QString alias = dviceJson["Alias"].toString();
    const bool paired = dviceJson["Paired"].toBool();
    const int rssi = dviceJson["RSSI"].toInt();
    const Device::State state = Device::State(dviceJson["State"].toInt());
    const bool connectState = dviceJson["ConnectState"].toBool();
    const QString bluetoothDeviceType = dviceJson["Icon"].toString();
    const int battery = dviceJson["Battery"].toInt();

    // FIXME: Solve the problem that the device name in the Bluetooth list is blank
    if (name.isEmpty() && alias.isEmpty())
        return ;

    const Device *constdevice = m_devices.value(id);
    auto device = const_cast<Device *>(constdevice);
    if (device) {
        device->setId(id);
        device->setName(name);
        device->setAlias(alias);
        device->setPaired(paired);
        device->setRssi(rssi);
        //setState放后面，是因为用到了connectState,fix bug 55245
        device->setConnectState(connectState);
        device->setState(state);
        device->setDeviceType(bluetoothDeviceType);
        device->setBattery(battery);
        emit deviceNameUpdated(device);
    }
}

void Adapter::setPowered(bool powered)
{
    if (powered != m_powered) {
        m_powered = powered;
        Q_EMIT poweredChanged(powered);
    }
}

void Adapter::initDevicesList(const QJsonDocument &doc)
{
    QJsonArray arr = doc.array();
    for (QJsonValue val : arr) {
        QJsonObject deviceObj = val.toObject();
        const QString adapterId = deviceObj["AdapterPath"].toString();
        const QString id = deviceObj["Path"].toString();
        const QString name = deviceObj["Name"].toString();
        const QString alias = deviceObj["Alias"].toString();
        const bool paired = deviceObj["Paired"].toBool();
        const int rssi = deviceObj["RSSI"].toInt();
        const Device::State state = Device::State(deviceObj["State"].toInt());
        const bool connectState = deviceObj["ConnectState"].toBool();
        const QString bluetoothDeviceType = deviceObj["Icon"].toString();
        const int battery = deviceObj["Battery"].toInt();

        auto device = new Device(this);
        device->setId(id);
        device->setName(name);
        device->setAlias(alias);
        device->setPaired(paired);
        device->setState(state);
        device->setConnectState(connectState);
        device->setRssi(rssi);
        device->setAdapterId(adapterId);
        device->setDeviceType(bluetoothDeviceType);
        device->setBattery(battery);

        m_devices[id] = device;
    }
}

QMap<QString, const Device *> Adapter::devices() const
{
    return m_devices;
}

const Device *Adapter::deviceById(const QString &id) const
{
    return m_devices.keys().contains(id) ? m_devices[id] : nullptr;
}

void Adapter::setId(const QString &id)
{
    m_id = id;
}

void Adapter::setDiscover(bool discover)
{
    if (discover != m_discover) {
        m_discover = discover;
        Q_EMIT discoveringChanged(discover);
    }
}

