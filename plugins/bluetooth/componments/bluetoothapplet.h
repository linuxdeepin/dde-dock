// Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef BLUETOOTHAPPLET_H
#define BLUETOOTHAPPLET_H

#include <QWidget>
#include <QScrollArea>
#include <QStandardItemModel>

#include <dtkwidget_global.h>
#include <DGuiApplicationHelper>

#include "org_deepin_dde_airplanemode1.h"

class QVBoxLayout;
class QHBoxLayout;

class Device;
class Adapter;
class BluetoothAdapterItem;
class AdaptersManager;
class HorizontalSeperator;

DWIDGET_BEGIN_NAMESPACE
class DLabel;
class DSwitchButton;
class DListView;
DWIDGET_END_NAMESPACE

DWIDGET_USE_NAMESPACE
DGUI_USE_NAMESPACE

using DBusAirplaneMode = org::deepin::dde::AirplaneMode1;

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
    void changeEvent(QEvent *event) override;

    void updateEnabledStatus();

private Q_SLOTS:
    void onThemeTypeChanged(DGuiApplicationHelper::ColorType themeType);

private:
    DLabel *m_label;
    QHBoxLayout *m_layout;
};

class BluetoothApplet : public QWidget
{
    Q_OBJECT
public:
    explicit BluetoothApplet(AdaptersManager *adapterManager, QWidget *parent = nullptr);
    bool poweredInitState();
    // 当前是否有蓝牙适配器
    bool hasAadapter();
    // 刷新蓝牙适配器搜索到的设备列表
    void setAdapterRefresh();
    // 设置当前所有蓝牙适配器的电源状态
    void setAdapterPowered(bool state);
    // 已连接蓝牙设备名称列表，用于获取鼠标悬浮在蓝牙插件上时tips显示内容
    QStringList connectedDevicesName();

    inline bool airplaneModeEnable() const { return m_airplaneModeEnable;}
    // 返回蓝牙适配器
    AdaptersManager *adaptersManager();

signals:
    void noAdapter();
    void justHasAdapter();
    void powerChanged(bool state);
    void deviceStateChanged(const Device *device);
    void requestHide();

public slots:
    // 蓝牙适配器增加
    void onAdapterAdded(Adapter *adapter);
    // 蓝牙适配器移除
    void onAdapterRemoved(Adapter *adapter);
    // 设置蓝牙适配器电源是否开启
    void onSetAdapterPower(Adapter *adapter, bool state);
    // 更新蓝牙适配器电源状态，用于更新任务栏蓝牙插件图标的显示状态
    void updateBluetoothPowerState();
    void setAirplaneModeEnabled(bool enable);

private:
    void initUi();
    void initConnect();
    // 更新蓝牙插件主界面大小
    void updateSize();
    void updateIconTheme();
    void initAdapters();

private:
    QScrollArea *m_scroarea;
    QWidget *m_contentWidget;
    AdaptersManager *m_adaptersManager;
    SettingLabel *m_settingLabel;
    QVBoxLayout *m_mainLayout;
    QVBoxLayout *m_contentLayout;
    HorizontalSeperator *m_seperator;

    QStringList m_connectDeviceName;
    QMap<QString, BluetoothAdapterItem *> m_adapterItems;           // 所有蓝牙适配器
    DBusAirplaneMode *m_airPlaneModeInter;
    bool m_airplaneModeEnable;
};

#endif // BLUETOOTHAPPLET_H
