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

#ifndef BLUETOOTHAPPLET_H
#define BLUETOOTHAPPLET_H

#include <QWidget>
#include <QStandardItemModel>

#include <dtkwidget_global.h>

class QVBoxLayout;
class QHBoxLayout;

class Device;
class Adapter;
class BluetoothAdapterItem;
class AdaptersManager;

DWIDGET_BEGIN_NAMESPACE
class DLabel;
class DSwitchButton;
class DListView;
DWIDGET_END_NAMESPACE

DWIDGET_USE_NAMESPACE

class SettingLabel : public QWidget
{
    Q_OBJECT
public:
    explicit SettingLabel(QString text, QWidget *parent = nullptr);
    void addButton(QWidget *button, int space);
    DLabel *label() { return m_label; }
signals:
    void clicked();

protected:
    void mousePressEvent(QMouseEvent *ev) override;
    void paintEvent(QPaintEvent *event) override;

private:
    DLabel *m_label = nullptr;
    QHBoxLayout *m_layout = nullptr;
};

class BluetoothApplet : public QWidget
{
    Q_OBJECT
public:
    explicit BluetoothApplet(QWidget *parent = nullptr);
    bool poweredInitState();
    // 当前是否有蓝牙适配器
    bool hasAadapter();
    // 刷新蓝牙适配器搜索到的设备列表
    void setAdapterRefresh();
    // 设置当前所有蓝牙适配器的电源状态
    void setAdapterPowered(bool state);
    // 已连接蓝牙设备名称列表，用于获取鼠标悬浮在蓝牙插件上时tips显示内容
    QStringList connectedDevicesName();

signals:
    void noAdapter();
    void justHasAdapter();
    void powerChanged(bool state);
    void deviceStateChanged(const Device* device);

public slots:
    // 蓝牙适配器增加
    void onAdapterAdded(Adapter *adapter);
    // 蓝牙适配器移除
    void onAdapterRemoved(Adapter *adapter);
    // 设置蓝牙适配器电源是否开启
    void onSetAdapterPower(Adapter *adapter, bool state);
    // 更新蓝牙适配器电源状态，用于更新任务栏蓝牙插件图标的显示状态
    void updateBluetoothPowerState();

private:
    void initUi();
    void initConnect();
    // 更新蓝牙插件主界面大小
    void updateSize();
    QWidget *m_contentWidget = nullptr;
    AdaptersManager *m_adaptersManager = nullptr;
    SettingLabel *m_settingLabel = nullptr;
    QVBoxLayout *m_mainLayout = nullptr;
    QVBoxLayout *m_contentLayout = nullptr;

    QStringList m_connectDeviceName;
    QMap<QString, BluetoothAdapterItem *> m_adapterItems;
};

#endif // BLUETOOTHAPPLET_H
