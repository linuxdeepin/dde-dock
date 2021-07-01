/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
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

#ifndef DEVICESTATUSHANDLER_H
#define DEVICESTATUSHANDLER_H

#include <QObject>

namespace dde {
  namespace network {
    class UWiredDevice;
    class UWirelessDevice;
  }
}

using namespace dde::network;

enum class NetDeviceStatus {
    Unknown = 0,
    Enabled,
    Disabled,
    Connected,
    Disconnected,
    Connecting,
    Authenticating,
    ObtainingIP,
    ObtainIpFailed,
    ConnectNoInternet,
    Nocable,
    ConnectFailed
};

enum class PluginState
{
    Unknow = 0,
    Disabled,
    Connected,
    Disconnected,
    Connecting,
    Failed,
    ConnectNoInternet,
    WirelessDisabled,
    WiredDisabled,
    WirelessConnected,
    WiredConnected,
    WirelessDisconnected,
    WiredDisconnected,
    WirelessConnecting,
    WiredConnecting,
    WirelessConnectNoInternet,
    WiredConnectNoInternet,
    WirelessFailed,
    WiredFailed,
    Nocable
};

#define DECLARE_STATIC_CHECKSTATUS(method) static bool method(const NetDeviceStatus &, const NetDeviceStatus &);

class DeviceStatusHandler : public QObject
{
    Q_OBJECT

public:
    // 获取当前所有的设备列表的状态
    static PluginState pluginState();

private:
    explicit DeviceStatusHandler(QObject *parent = Q_NULLPTR);
    ~DeviceStatusHandler();

    static NetDeviceStatus wiredStatus(UWiredDevice * device);
    static NetDeviceStatus wiredStatus(QList<UWiredDevice *> devices);
    static NetDeviceStatus wirelessStatus(UWirelessDevice *device);
    static NetDeviceStatus wirelessStatus(QList<UWirelessDevice *>devices);
    static PluginState plugState(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus);

private:
    DECLARE_STATIC_CHECKSTATUS(isUnknow)
    DECLARE_STATIC_CHECKSTATUS(isDisabled)
    DECLARE_STATIC_CHECKSTATUS(isWiredDisconnected)
    DECLARE_STATIC_CHECKSTATUS(isWiredDisabled)
    DECLARE_STATIC_CHECKSTATUS(isWiredConnected)
    DECLARE_STATIC_CHECKSTATUS(isWiredConnecting)
    DECLARE_STATIC_CHECKSTATUS(isWiredConnectNoInternet)
    DECLARE_STATIC_CHECKSTATUS(isNocable)
    DECLARE_STATIC_CHECKSTATUS(isWiredFailed)
    DECLARE_STATIC_CHECKSTATUS(isWirelessDisconnected)
    DECLARE_STATIC_CHECKSTATUS(isWirelessDisabled)
    DECLARE_STATIC_CHECKSTATUS(isWirelessConnected)
    DECLARE_STATIC_CHECKSTATUS(isWirelessConnecting)
    DECLARE_STATIC_CHECKSTATUS(isWirelessConnectNoInternet)
    DECLARE_STATIC_CHECKSTATUS(isWirelessFailed)
    DECLARE_STATIC_CHECKSTATUS(isDisconnected)
    DECLARE_STATIC_CHECKSTATUS(isConnected)
    DECLARE_STATIC_CHECKSTATUS(isConnecting)
    DECLARE_STATIC_CHECKSTATUS(isConnectNoInternet)
};

#endif // DEVICESTATUSHANDLER_H
