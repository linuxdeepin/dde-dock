/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     chenwei <chenwei@uniontech.com>
 *
 * Maintainer: chenwei <chenwei@uniontech.com>
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

#ifndef BLUETOOTHADAPTERITEM_H
#define BLUETOOTHADAPTERITEM_H

#include "componments/device.h"

#include <QWidget>

#include <DStyleHelper>
#include <DApplicationHelper>

#include <com_deepin_daemon_bluetooth.h>

using  DBusBluetooth = com::deepin::daemon::Bluetooth;

DWIDGET_USE_NAMESPACE

DWIDGET_BEGIN_NAMESPACE
class DSwitchButton;
class DStandardItem;
class DListView;
class DSpinner;
DWIDGET_END_NAMESPACE

class Adapter;
class SettingLabel;
class QStandardItemModel;
class RefreshButton;

const QString LightString = QString(":/light/buletooth_%1_light.svg");
const QString DarkString = QString(":/dark/buletooth_%1_dark.svg");

class BluetoothDeviceItem : public QObject
{
    Q_OBJECT
public:
    explicit BluetoothDeviceItem(QStyle *style = nullptr, const Device *device = nullptr, DListView *parent = nullptr);
    virtual ~BluetoothDeviceItem();

    DStandardItem *standardItem() { return m_standarditem; }
    const Device *device() { return m_device; }

public slots:
    // 系统主题发生改变时更新蓝牙图标
    void updateIconTheme(DGuiApplicationHelper::ColorType type);
    // 更新蓝牙设备的连接状态
    void updateDeviceState(Device::State state);

signals:
    void requestTopDeviceItem(DStandardItem *item);
    void deviceStateChanged(const Device *device);

private:
    void initActionList();
    void initConnect();

    DStyleHelper m_style;
    QString m_deviceIcon;

    const Device *m_device = nullptr;
    DStandardItem *m_standarditem = nullptr;
    DViewItemAction *m_labelAction = nullptr;
    DViewItemAction *m_stateAction = nullptr;
    DSpinner *m_loading = nullptr;
};

class BluetoothAdapterItem : public QWidget
{
    Q_OBJECT
public:
    explicit BluetoothAdapterItem(Adapter *adapter, QWidget *parent = nullptr);
    ~BluetoothAdapterItem();
    Adapter *adapter() { return m_adapter; }
    int currentDeviceCount();
    QStringList connectedDevicesName();

public slots:
    // 添加蓝牙设备
    void onDeviceAdded(const Device *device);
    // 移除蓝牙设备
    void onDeviceRemoved(const Device *device);
    // 蓝牙设备名称更新
    void onDeviceNameUpdated(const Device *device);
    // 连接蓝牙设备
    void onConnectDevice(const QModelIndex &index);
    // 将已连接的蓝牙设备放到列表第一个
    void onTopDeviceItem(DStandardItem *item);
    // 设置蓝牙适配器名称
    void onAdapterNameChanged(const QString name);
    void updateIconTheme(DGuiApplicationHelper::ColorType type);

signals:
    void adapterPowerChanged();
    void requestSetAdapterPower(Adapter *adapter, bool state);
    void requestRefreshAdapter(Adapter *adapter);
    void connectDevice(const Device *device, Adapter *adapter);
    void deviceCountChanged();
    void deviceStateChanged(const Device *device);

private:
    void initData();
    void initUi();
    void initConnect();
    void setUnnamedDevicesVisible(bool isShow);

    Adapter *m_adapter = nullptr;
    SettingLabel *m_adapterLabel = nullptr;
    DSwitchButton *m_adapterStateBtn = nullptr;
    DListView *m_deviceListview = nullptr;
    QStandardItemModel *m_deviceModel = nullptr;
    RefreshButton *m_refreshBtn = nullptr;
    DBusBluetooth *m_bluetoothInter;
    bool m_showUnnamedDevices;

    QMap<QString, BluetoothDeviceItem *> m_deviceItems;
};

#endif // BLUETOOTHADAPTERITEM_H
