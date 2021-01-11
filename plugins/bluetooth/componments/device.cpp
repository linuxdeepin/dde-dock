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

#include "device.h"

#include <QDateTime>

QMap<QString,QString> Device::deviceType2Icon = {
    {"unknow","other"},
    {"computer","pc"},
    {"phone","phone"},
    {"video-display","vidicon"},
    {"multimedia-player","tv"},
    {"scanner","scaner"},
    {"input-keyboard","keyboard"},
    {"input-mouse","mouse"},
    {"input-gaming","other"},
    {"input-tablet","touchpad"},
    {"audio-card","pheadset"},
    {"network-wireless","lan"},
    {"camera-video","vidicon"},
    {"printer","print"},
    {"camera-photo","camera"},
    {"modem","other"}
};

Device::Device(QObject *parent)
    : QObject(parent)
    , m_paired(false)
    , m_trusted(false)
    , m_connecting(false)
    , m_rssi(0)
    , m_state(StateUnavailable)
    , m_connectState(false)
{
    m_time = static_cast<int>(QDateTime::currentDateTime().toTime_t());
}

Device::~Device()
{
}

void Device::updateDeviceTime()
{
    m_time = static_cast<int>(QDateTime::currentDateTime().toTime_t());
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

void Device::setTrusted(bool trusted)
{
    if (trusted != m_trusted) {
        m_trusted = trusted;
        Q_EMIT trustedChanged(trusted);
    }
}

void Device::setConnecting(bool connecting)
{
    if (connecting != m_connecting) {
        m_connecting = connecting;
        Q_EMIT connectingChanged(connecting);
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
    m_deviceType = deviceType2Icon[deviceType];
}

QDebug &operator<<(QDebug &stream, const Device *device)
{
    stream << "Device name:" << device->name()
           << " paired:" << device->paired()
           << " state:" << device->state();

    return stream;
}
