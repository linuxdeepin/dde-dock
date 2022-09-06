// SPDX-FileCopyrightText: 2016 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef BLUETOOTHADAPTERITEM_H
#define BLUETOOTHADAPTERITEM_H

#include "device.h"
#include "bluetoothapplet.h"

#include <QWidget>

#include <DListView>
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
class HorizontalSeperator;
class StateButton;

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
    void disconnectDevice();

private:
    void initActionList();
    void initConnect();

    DStyleHelper m_style;
    QString m_deviceIcon;

    const Device *m_device;
    DStandardItem *m_standarditem;
    DViewItemAction *m_labelAction;
    DViewItemAction *m_stateAction;
    DViewItemAction *m_connAction;
    DSpinner *m_loading;

    QWidget *m_iconWidget;
    StateButton *m_connButton;
};

class BluetoothAdapterItem : public QWidget
{
    Q_OBJECT
public:
    explicit BluetoothAdapterItem(Adapter *adapter, QWidget *parent = nullptr);
    ~BluetoothAdapterItem();
    Adapter *adapter() { return m_adapter; }
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

    QSize sizeHint() const override;

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

    Adapter *m_adapter;
    SettingLabel *m_adapterLabel;
    DSwitchButton *m_adapterStateBtn;
    DListView *m_deviceListview;
    DStyledItemDelegate *m_itemDelegate;
    QStandardItemModel *m_deviceModel;
    RefreshButton *m_refreshBtn;
    DBusBluetooth *m_bluetoothInter;
    bool m_showUnnamedDevices;

    QMap<QString, BluetoothDeviceItem *> m_deviceItems;
    HorizontalSeperator *m_seperator;
};

#endif // BLUETOOTHADAPTERITEM_H
