/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     zhaolong <zhaolong@uniontech.com>
 *
 * Maintainer: zhaolong <zhaolong@uniontech.com>
 *
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

    m_devices[id] = device;
    divideDevice(device);

    emit deviceAdded(device);
}

void Adapter::removeDevice(const QString &deviceId)
{
    const Device *constDevice = m_devices.value(deviceId);
    auto device = const_cast<Device *>(constDevice);
    if (device) {
        m_devices.remove(deviceId);
        m_paredDev.remove(deviceId);
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
    }
}

//void Adapter::removeAllDevices()
//{
//    QMapIterator<QString, const Device *> iterator(m_devices);
//    while (iterator.hasNext()) {
//        iterator.next();
//        auto device = const_cast<Device *>(iterator.value());
//        if (device) {
//            m_devices.remove(device->id());
//            m_paredDev.remove(device->id());
//            delete device;
//        }
//    }
//}

const QMap<QString, const Device *> &Adapter::paredDevices() const
{
    return  m_paredDev;
}

//int Adapter::paredDevicesCount() const
//{
//    return  m_paredDev.size();
//}

void Adapter::divideDevice(const Device *device)
{
    if (device->paired()) {
        m_paredDev[device->id()] = device;
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

        m_devices[id] = device;
        divideDevice(device);
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

