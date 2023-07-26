// Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "device.h"

#include <DStyleHelper>
#include <QDateTime>

Device::Device(QObject *parent)
    : QObject(parent)
    , m_paired(false)
    , m_trusted(false)
    , m_connecting(false)
    , m_rssi(0)
    , m_state(StateUnavailable)
    , m_connectState(false)
    , m_battery(0)
{
}

Device::~Device()
{
}

void Device::setId(const QString &id)
{
    m_id = id;
}

void Device::setName(const QString &name)
{
    if (name != m_name) {
        m_name = name;
        Q_EMIT nameChanged(name);
    }
}

void Device::setAlias(const QString &alias)
{
    if (alias != m_alias) {
        m_alias = alias;
        Q_EMIT aliasChanged(alias);
    }
}

void Device::setPaired(bool paired)
{
    if (paired != m_paired) {
        m_paired = paired;
        Q_EMIT pairedChanged(paired);
    }
}

void Device::setState(const State &state)
{
    if (state != m_state) {
        m_state = state;
        Q_EMIT stateChanged(state);
    }
}

void Device::setConnectState(const bool connectState)
{
    if (connectState != m_connectState) {
        m_connectState = connectState;
        Q_EMIT connectStateChanged(connectState);
    }
}

void Device::setRssi(int rssi)
{
    if (m_rssi != rssi) {
        m_rssi = rssi;
        Q_EMIT rssiChanged(rssi);
    }
}

void Device::setDeviceType(const QString &deviceType)
{
    m_deviceType = deviceType;
}

void Device::setBattery(int battery)
{
    if (m_battery != battery) {
        m_battery = battery;
        Q_EMIT batteryChanged(battery);
    }
}

QDebug &operator<<(QDebug &stream, const Device *device)
{
    stream << "Device name:" << device->name()
           << " paired:" << device->paired()
           << " state:" << device->state();

    return stream;
}
