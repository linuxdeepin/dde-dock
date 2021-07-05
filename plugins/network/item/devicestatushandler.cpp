/*
 * Copyright (C) 2011 ~ 2021 Deepin Technology Co., Ltd.
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

#include "devicestatushandler.h"

#include <wireddevice.h>
#include <wirelessdevice.h>
#include <networkcontroller.h>

DeviceStatusHandler::DeviceStatusHandler(QObject *parent)
    : QObject(parent)
{
}

DeviceStatusHandler::~DeviceStatusHandler()
{
}

PluginState DeviceStatusHandler::pluginState()
{
    QList<NetworkDeviceBase *> devices = NetworkController::instance()->devices();
    // 筛选出所有的有线和无线的状态
    QList<WiredDevice *> wiredDevices;
    QList<WirelessDevice *> wirelessDevice;
    for (NetworkDeviceBase *deviceBase : devices) {
        if (deviceBase->deviceType() == DeviceType::Wired) {
            WiredDevice *device = static_cast<WiredDevice *>(deviceBase);
            wiredDevices << device;
        } else if (deviceBase->deviceType() == DeviceType::Wireless) {
            WirelessDevice *device = static_cast<WirelessDevice *>(deviceBase);
            wirelessDevice << device;
        }
    }

    // 计算有线网络和无线网络的合并状态
    NetDeviceStatus wiredStat = wiredStatus(wiredDevices);
    NetDeviceStatus wirelessStat = wirelessStatus(wirelessDevice);
    return plugState(wiredStat, wirelessStat);
}

NetDeviceStatus DeviceStatusHandler::wiredStatus(WiredDevice *device)
{
    // 如果当前网卡是禁用，直接返回禁用
    if (!device->isEnabled())
        return NetDeviceStatus::Disabled;

    // 网络是已连接，但是当前的连接状态不是Full，则认为网络连接成功，但是无法上网
    if (device->deviceStatus() == DeviceStatus::Activated
            && NetworkController::instance()->connectivity() != Connectivity::Full) {
        return NetDeviceStatus::ConnectNoInternet;
    }

    // 获取IP地址失败
    if (!device->IPValid())
        return NetDeviceStatus::ObtainIpFailed;

    // 根据设备状态来直接获取返回值
    switch (device->deviceStatus()) {
    case DeviceStatus::Unmanaged:
    case DeviceStatus::Unavailable:    return NetDeviceStatus::Nocable;
    case DeviceStatus::Disconnected:   return NetDeviceStatus::Disconnected;
    case DeviceStatus::Prepare:
    case DeviceStatus::Config:         return NetDeviceStatus::Connecting;
    case DeviceStatus::Needauth:       return NetDeviceStatus::Authenticating;
    case DeviceStatus::IpConfig:
    case DeviceStatus::IpCheck:
    case DeviceStatus::Secondaries:    return NetDeviceStatus::ObtainingIP;
    case DeviceStatus::Activated:      return NetDeviceStatus::Connected;
    case DeviceStatus::Deactivation:
    case DeviceStatus::Failed:         return NetDeviceStatus::ConnectFailed;
    default:                           return NetDeviceStatus::Unknown;
    }

    Q_UNREACHABLE();
    return NetDeviceStatus::Unknown;
}

NetDeviceStatus DeviceStatusHandler::wiredStatus(const QList<WiredDevice *> &devices)
{
    QList<NetDeviceStatus> deviceStatus;
    for (WiredDevice *device : devices)
        deviceStatus << wiredStatus(device);

    // 显示的规则:从allDeviceStatus列表中按照顺序遍历所有的状态，
    // 再遍历所有的设备的状态，只要其中一个设备的状态满足当前的状态，就返回当前状态
    static QList<NetDeviceStatus> allDeviceStatus =
        { NetDeviceStatus::Authenticating, NetDeviceStatus::ObtainingIP, NetDeviceStatus::Connected,
        NetDeviceStatus::ConnectNoInternet, NetDeviceStatus::Connecting, NetDeviceStatus::Disconnected,
        NetDeviceStatus::Disabled, NetDeviceStatus::Nocable, NetDeviceStatus::Unknown };
    for (int i = 0; i < allDeviceStatus.size(); i++) {
        NetDeviceStatus status = allDeviceStatus[i];
        if (deviceStatus.contains(status))
            return status;
    }

    return NetDeviceStatus::Unknown;
}

NetDeviceStatus DeviceStatusHandler::wirelessStatus(WirelessDevice *device)
{
    if (!device->isEnabled())
        return NetDeviceStatus::Disabled;

    if (device->deviceStatus() == DeviceStatus::Activated
            && device->connectivity() != Connectivity::Full) {
        return NetDeviceStatus::ConnectNoInternet;
    }

    if (!device->IPValid())
        return NetDeviceStatus::ObtainIpFailed;

    DeviceStatus status = device->deviceStatus();
    switch (status) {
    case DeviceStatus::Unmanaged:
    case DeviceStatus::Unavailable:
    case DeviceStatus::Disconnected:  return NetDeviceStatus::Disconnected;
    case DeviceStatus::Prepare:
    case DeviceStatus::Config:        return NetDeviceStatus::Connecting;
    case DeviceStatus::Needauth:      return NetDeviceStatus::Authenticating;
    case DeviceStatus::IpConfig:
    case DeviceStatus::IpCheck:
    case DeviceStatus::Secondaries:   return NetDeviceStatus::ObtainingIP;
    case DeviceStatus::Activated:     return NetDeviceStatus::Connected;
    case DeviceStatus::Deactivation:
    case DeviceStatus::Failed:        return NetDeviceStatus::ConnectFailed;
    default:                          return NetDeviceStatus::Unknown;
    }

    Q_UNREACHABLE();
    return NetDeviceStatus::Unknown;
}

NetDeviceStatus DeviceStatusHandler::wirelessStatus(const QList<WirelessDevice *> &devices)
{
    // 所有设备状态叠加
    QList<NetDeviceStatus> devStatus;
    for (WirelessDevice *device : devices)
        devStatus << wirelessStatus(device);

    static QList<NetDeviceStatus> allDeviceStatus =
        { NetDeviceStatus::Authenticating, NetDeviceStatus::ObtainingIP, NetDeviceStatus::Connected,
        NetDeviceStatus::ConnectNoInternet, NetDeviceStatus::Connecting, NetDeviceStatus::Disconnected,
        NetDeviceStatus::Disabled, NetDeviceStatus::Unknown};

    for (int i = 0; i < allDeviceStatus.size(); i++) {
        NetDeviceStatus status = allDeviceStatus[i];
        if (devStatus.contains(status))
            return status;
    }

    return NetDeviceStatus::Unknown;
}

PluginState DeviceStatusHandler::plugState(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    if (isUnknow(wiredStatus, wirelessStatus))
        return PluginState::Unknow;

    if (isDisabled(wiredStatus, wirelessStatus))
        return PluginState::Disabled;

    if (isWiredDisconnected(wiredStatus, wirelessStatus))
        return PluginState::WiredDisconnected;

    if (isWiredDisabled(wiredStatus, wirelessStatus))
        return PluginState::WiredDisabled;

    if (isWiredConnected(wiredStatus, wirelessStatus))
        return PluginState::WiredConnected;

    if (isWiredConnecting(wiredStatus, wirelessStatus))
        return PluginState::WiredConnecting;

    if (isWiredConnectNoInternet(wiredStatus, wirelessStatus))
        return PluginState::WiredConnectNoInternet;

    if (isNocable(wiredStatus, wirelessStatus))
        return PluginState::Nocable;

    if (isWiredFailed(wiredStatus, wirelessStatus))
        return PluginState::WiredFailed;

    if (isWirelessDisconnected(wiredStatus, wirelessStatus))
        return PluginState::WirelessDisconnected;

    if (isWirelessDisabled(wiredStatus, wirelessStatus))
        return PluginState::WirelessDisabled;

    if (isWirelessConnected(wiredStatus, wirelessStatus))
        return PluginState::WirelessConnected;

    if (isWirelessConnecting(wiredStatus, wirelessStatus))
        return PluginState::WirelessConnecting;

    if (isWirelessConnectNoInternet(wiredStatus, wirelessStatus))
        return PluginState::WirelessConnectNoInternet;

    if (isWirelessFailed(wiredStatus, wirelessStatus))
        return PluginState::WirelessFailed;

    if (isDisconnected(wiredStatus, wirelessStatus))
        return PluginState::Disconnected;

    if (isConnected(wiredStatus, wirelessStatus))
        return PluginState::Connected;

    if (isConnecting(wiredStatus, wirelessStatus))
        return PluginState::Connecting;

    if (isConnectNoInternet(wiredStatus, wirelessStatus))
        return PluginState::ConnectNoInternet;

    return PluginState::Failed;
}

bool DeviceStatusHandler::isUnknow(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 无线和有线都是未知状态，则认为是未知状态(都没有网卡)
    return (wiredStatus == NetDeviceStatus::Unknown
            && wirelessStatus == NetDeviceStatus::Unknown);
}

bool DeviceStatusHandler::isDisabled(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 有线无线都禁用，则认为是禁用状态
    return (wiredStatus == NetDeviceStatus::Disabled
            && wirelessStatus == NetDeviceStatus::Disabled);
}

bool DeviceStatusHandler::isWiredDisconnected(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 没有无线或者无线禁用的情况下, 有线设备开启、有线设备断开连接，有线设备获取IP失败，认为有线连接失败
    return ((wirelessStatus == NetDeviceStatus::Unknown && wiredStatus == NetDeviceStatus::Enabled)
            || (wirelessStatus == NetDeviceStatus::Unknown && wiredStatus == NetDeviceStatus::Disconnected)
            || (wirelessStatus == NetDeviceStatus::Unknown && wiredStatus == NetDeviceStatus::ObtainIpFailed)
            || (wirelessStatus == NetDeviceStatus::Disabled && wiredStatus == NetDeviceStatus::Enabled)
            || (wirelessStatus == NetDeviceStatus::Disabled && wiredStatus == NetDeviceStatus::Disconnected)
            || (wirelessStatus == NetDeviceStatus::Disabled && wiredStatus == NetDeviceStatus::ObtainIpFailed));
}

bool DeviceStatusHandler::isWiredDisabled(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 有线禁用了，没有无线网卡
    return (wiredStatus == NetDeviceStatus::Disabled
            && wirelessStatus == NetDeviceStatus::Unknown);
}

bool DeviceStatusHandler::isWiredConnected(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 如果有线是连接状态，没有无线，无线启用，无线禁用，无线断开连接，无线获取IP失败，无线连接成功但是无网络
    // 无线连接失败，则认为是有线连接成功
    static QList<NetDeviceStatus> wirelessFailuredStatus =
        { NetDeviceStatus::Unknown, NetDeviceStatus::Enabled
        , NetDeviceStatus::Disabled, NetDeviceStatus::Disconnected
        , NetDeviceStatus::ObtainIpFailed, NetDeviceStatus::ConnectNoInternet
        , NetDeviceStatus::ConnectFailed };

    return ((wiredStatus == NetDeviceStatus::Connected)
            && (wirelessFailuredStatus.contains(wirelessStatus)));
}

bool DeviceStatusHandler::isWiredConnecting(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 有线正在连接，正在认证，正在获取IP
    // 没有无线，无线开启和禁用、连接成功，断开连接，获取IP失败，连接但是没有网络，连接失败
    // 这种情况认为是有线正在连接
    static QList<NetDeviceStatus> wiredConnectingStatus =
        { NetDeviceStatus::Connecting, NetDeviceStatus::Authenticating, NetDeviceStatus::ObtainingIP };
    static QList<NetDeviceStatus> wirelessConnecting =
        { NetDeviceStatus::Unknown, NetDeviceStatus::Enabled
        , NetDeviceStatus::Disabled, NetDeviceStatus::Connected
        , NetDeviceStatus::Disconnected, NetDeviceStatus::ObtainIpFailed
        , NetDeviceStatus::ConnectNoInternet, NetDeviceStatus::ConnectFailed };

    return (wiredConnectingStatus.contains(wiredStatus)
            && wirelessConnecting.contains(wirelessStatus));
}

bool DeviceStatusHandler::isWiredConnectNoInternet(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 没有无线，无线开启或禁用，无线断开连接，无线获取IP失败，无线连接失败
    // 有线连接但是没有网络
    // 这种情况认为是有线连接无网络
    static QList<NetDeviceStatus> wirelessNoConnectStatus =
        { NetDeviceStatus::Unknown, NetDeviceStatus::Enabled
        , NetDeviceStatus::Disabled, NetDeviceStatus::Disconnected
        , NetDeviceStatus::ObtainIpFailed, NetDeviceStatus::ConnectFailed };

    return (wiredStatus == NetDeviceStatus::ConnectNoInternet
            && wirelessNoConnectStatus.contains(wirelessStatus));
}

bool DeviceStatusHandler::isNocable(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 有线不可用，没有无线或无线禁用，这种情况认为是网络不可用
    return (wiredStatus == NetDeviceStatus::Nocable
            && (wirelessStatus == NetDeviceStatus::Unknown
                || wirelessStatus == NetDeviceStatus::Disabled ));
}

bool DeviceStatusHandler::isWiredFailed(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 有线连接失败，没有无线或无线禁用，这种情况认为是有线不可用
    return (wiredStatus == NetDeviceStatus::ConnectFailed
         && (wirelessStatus == NetDeviceStatus::Unknown
             || wirelessStatus == NetDeviceStatus::Disabled ));
}

bool DeviceStatusHandler::isWirelessDisconnected(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 无线启用，断开连接，获取IP失败
    // 没有有线，有线禁用，有线无效，有线失败
    // 这种情况认为是无线断开连接
    static QList<NetDeviceStatus> wirelessStatusOfDis =
        { NetDeviceStatus::Enabled, NetDeviceStatus::Disconnected
        , NetDeviceStatus::ObtainIpFailed };
    static QList<NetDeviceStatus> wiredStatusOfDis =
        { NetDeviceStatus::Unknown, NetDeviceStatus::Disabled
        , NetDeviceStatus::Nocable, NetDeviceStatus::ConnectFailed };

    return ((wirelessStatusOfDis.contains(wirelessStatus)
            && wiredStatusOfDis.contains(wiredStatus))
            || (wirelessStatus == NetDeviceStatus::ConnectFailed
                && wiredStatus == NetDeviceStatus::Nocable));
}

bool DeviceStatusHandler::isWirelessDisabled(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 无线禁用，没有有线，这种情况认为是无线禁用
    return (wirelessStatus == NetDeviceStatus::Disabled
            && wiredStatus == NetDeviceStatus::Unknown);
}

bool DeviceStatusHandler::isWirelessConnected(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 没有有线网卡，有线启用禁用，断开连接，获取IP失败，已经连接网络但是无法上网，无效和失败
    // 无线网络已连接，这种情况认为是无线网络已连接
    static QList<NetDeviceStatus> wiredFailusStatus =
        { NetDeviceStatus::Unknown, NetDeviceStatus::Enabled
        , NetDeviceStatus::Disabled, NetDeviceStatus::Disconnected
        , NetDeviceStatus::ObtainIpFailed, NetDeviceStatus::ConnectNoInternet
        , NetDeviceStatus::Nocable, NetDeviceStatus::ConnectFailed };

    return (wiredFailusStatus.contains(wiredStatus)
            && wirelessStatus == NetDeviceStatus::Connected);
}

bool DeviceStatusHandler::isWirelessConnecting(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 无线状态：正在连接，正在认证, 正在获取IP
    // 有线状态：没有有线，禁用，启用，连接成功，断开连接，获取IP失败，连接成功但是无法上网，无效，连接失败
    // 这种情况认为是无线正在连接
    static QList<NetDeviceStatus> wirelessConnecting =
        { NetDeviceStatus::Connecting, NetDeviceStatus::Authenticating
        , NetDeviceStatus::ObtainingIP };

    static QList<NetDeviceStatus> wiredOfConnecting =
        { NetDeviceStatus::Unknown, NetDeviceStatus::Enabled
        , NetDeviceStatus::Disabled, NetDeviceStatus::Connected
        , NetDeviceStatus::Disconnected, NetDeviceStatus::ObtainIpFailed
        , NetDeviceStatus::ConnectNoInternet, NetDeviceStatus::Nocable
        , NetDeviceStatus::ConnectFailed };

    return (wirelessConnecting.contains(wirelessStatus)
            && wiredOfConnecting.contains(wiredStatus));
}

bool DeviceStatusHandler::isWirelessConnectNoInternet(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 无线状态：没有无线，启用，禁用，断开连接，获取IP失败，无效，失败
    // 无线是已连接但是无法上网，这种情况认为是无线已连接但是无法上网
    static QList<NetDeviceStatus> allWiredStatus =
        { NetDeviceStatus::Unknown, NetDeviceStatus::Enabled
        , NetDeviceStatus::Disabled, NetDeviceStatus::Disconnected
        , NetDeviceStatus::ObtainIpFailed, NetDeviceStatus::Nocable
        , NetDeviceStatus::ConnectFailed };

    return (allWiredStatus.contains(wiredStatus)
            && wirelessStatus == NetDeviceStatus::ConnectNoInternet);
}

bool DeviceStatusHandler::isWirelessFailed(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 无线连接失败，在没有有线和有线禁用的情况下，认为无线连接失败
    return (wirelessStatus == NetDeviceStatus::ConnectFailed
            && (wiredStatus == NetDeviceStatus::Unknown
                || wiredStatus == NetDeviceStatus::Disabled));
}

bool DeviceStatusHandler::isDisconnected(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 无线启用，断开连接，获取IP失败，连接失败
    // 有线启用，断开连接，获取IP失败
    // 这种情况认为网络连接断开
    static QList<NetDeviceStatus> disconectStatusWireless =
        { NetDeviceStatus::Enabled, NetDeviceStatus::Disconnected
        , NetDeviceStatus::ObtainIpFailed, NetDeviceStatus::ConnectFailed };

    static QList<NetDeviceStatus> disconnectStatusWired =
        { NetDeviceStatus::Enabled, NetDeviceStatus::Disconnected, NetDeviceStatus::ObtainIpFailed };

    return (disconectStatusWireless.contains(wirelessStatus)
            && disconnectStatusWired.contains(wiredStatus));
}

bool DeviceStatusHandler::isConnected(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 有线和无线都连接成功，这种情况认为连接成功
    return (wirelessStatus == NetDeviceStatus::Connected
            && wiredStatus == NetDeviceStatus::Connected);
}

bool DeviceStatusHandler::isConnecting(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 有线正在连接，正在认证，正在获取IP
    // 无线正在连接，正在认证，正在获取IP
    // 这种情况认为是正在连接
    static QList<NetDeviceStatus> connectingWired =
        { NetDeviceStatus::Connecting, NetDeviceStatus::Authenticating
        , NetDeviceStatus::ObtainingIP };

    static QList<NetDeviceStatus> connectingWireless =
        { NetDeviceStatus::Connecting, NetDeviceStatus::Authenticating
          , NetDeviceStatus::ObtainingIP };

    return (connectingWired.contains(wiredStatus)
            && connectingWireless.contains(wirelessStatus));
}

bool DeviceStatusHandler::isConnectNoInternet(const NetDeviceStatus &wiredStatus, const NetDeviceStatus &wirelessStatus)
{
    // 有线和无线都已经连接但是无法上网，这种情况认为是已连接但是无法上网
    return (wirelessStatus == NetDeviceStatus::ConnectNoInternet
            && wiredStatus == NetDeviceStatus::ConnectNoInternet);
}
