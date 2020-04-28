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

#ifndef ADAPTERITEM_H
#define ADAPTERITEM_H

#include "device.h"

#include <QScrollArea>
#include <QMap>
#include <QVBoxLayout>
#include <QLabel>

class HorizontalSeparator;
class Adapter;
class SwitchItem;
class DeviceItem;
class AdaptersManager;
class MenueItem;
class AdapterItem : public QScrollArea
{
    Q_OBJECT
public:
    explicit AdapterItem(AdaptersManager *a, Adapter *adapter, QWidget *parent = nullptr);
    int deviceCount();
    void setPowered(bool powered);
    bool isPowered();
    int viewHeight();
    inline Device::State initDeviceState() { return  m_initDeviceState; }
    inline Device::State currentDeviceState() { return m_currentDeviceState; }

signals:
    void deviceStateChanged(const Device::State state);
    void powerChanged(bool powered);
    void sizeChange();

private slots:
    void deviceItemPaired(const bool paired);
    void deviceRssiChanged();
    void removeDeviceItem(const Device *device);
    void showAndConnect(bool change);
    void addDeviceItem(const Device *constDevice);
    void deviceChangeState(const Device::State state);
    void moveDeviceItem(Device::State state, DeviceItem *item);

private:
    void createDeviceItem(Device *device);
    void updateView();
    void showDevices(bool change);

private:
    QWidget *m_centralWidget;
    HorizontalSeparator *m_line;
    QVBoxLayout *m_deviceLayout;
    MenueItem *m_openControlCenter;

    AdaptersManager *m_adaptersManager;

    Adapter *m_adapter;
    SwitchItem *m_switchItem;
    QMap<QString, DeviceItem*> m_deviceItems;
    Device::State m_initDeviceState;
    Device::State m_currentDeviceState;
    QList<DeviceItem *> m_sortConnected;
    QList<DeviceItem *> m_sortUnConnect;
};

#endif // ADAPTERITEM_H
