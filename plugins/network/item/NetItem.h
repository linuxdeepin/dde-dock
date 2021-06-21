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

#ifndef NETITEM_H
#define NETITEM_H

#include <DStandardItem>

#include <QAbstractListModel>
#include <QModelIndex>
#include <QJsonObject>
#include <QStyledItemDelegate>

#include <unetworkconst.h>

DWIDGET_USE_NAMESPACE
DGUI_USE_NAMESPACE

class DeviceController;
class NetworkDevice;
class NetItem;
class QLabel;
class QPushButton;
class UNetworkDeviceBase;
class UWiredDevice;
class UWirelessDevice;
class UAccessPoints;
class UWiredConnection;

namespace Dtk {
  namespace Widget {
    class DListView;
    class DSwitchButton;
    class DViewItemAction;
    class DLoadingIndicator;
    class DSpinner;
  }
}

enum NetItemRole
{
    TypeRole = Qt::UserRole + 100,
    DeviceDataRole,
    DataRole
};

#define PANELWIDTH 360

enum NetItemType
{
    DeviceControllViewItem = 0,     // 总控开关
    WirelessControllViewItem,       // 无线网卡开关
    WirelessViewItem,               // 无线列表
    WiredControllViewItem,          // 有线网卡开关
    WiredViewItem                   // 有线列表
};

class NetItem : public QObject
{
    Q_OBJECT

public:
    NetItem(QWidget *parent);
    virtual ~NetItem();

    virtual DStandardItem *standardItem();
    virtual void updateView() {}
    virtual NetItemType itemType() = 0;

protected:
    QWidget *parentWidget();

private:
    DStandardItem *m_standardItem;
    QWidget *m_parentWidget;
};

class DeviceControllItem : public NetItem
{
    Q_OBJECT

public:
    DeviceControllItem(QWidget *parent, UDeviceType deviceType);
    ~DeviceControllItem() Q_DECL_OVERRIDE;

    void setDevices(const QList<UNetworkDeviceBase *> &devices);
    UDeviceType deviceType();
    void updateView() Q_DECL_OVERRIDE;
    NetItemType itemType() Q_DECL_OVERRIDE;

private:
    void initItemText();
    void initSwitcher();
    void initConnection();

private Q_SLOTS:
    void onSwitchDevices(bool on);

private:
    QList<UNetworkDeviceBase *> m_devices;
    UDeviceType m_deviceType;
    DSwitchButton *m_switcher;
};

class WiredControllItem : public NetItem
{
    Q_OBJECT

public:
    WiredControllItem(QWidget *parent, UWiredDevice *device);
    ~WiredControllItem() Q_DECL_OVERRIDE;

    UWiredDevice *device();
    void updateView() Q_DECL_OVERRIDE;
    NetItemType itemType() Q_DECL_OVERRIDE;

protected Q_SLOTS:
    void onSwitchDevices(bool on);

private:
    UWiredDevice *m_device;
    DSwitchButton *m_switcher;
};

class WirelessControllItem : public NetItem
{
    Q_OBJECT

public:
    WirelessControllItem(QWidget *parent, UWirelessDevice *device);
    ~WirelessControllItem() Q_DECL_OVERRIDE;

    UWirelessDevice *device();
    void updateView() Q_DECL_OVERRIDE;
    NetItemType itemType() Q_DECL_OVERRIDE;

protected:
    bool eventFilter(QObject *object, QEvent *event) Q_DECL_OVERRIDE;

    QString iconFile();

protected Q_SLOTS:
    void onSwitchDevices(bool on);

private:
    UWirelessDevice *m_device;
    QWidget *m_widget;
    DSwitchButton *m_switcher;
    DLoadingIndicator *m_loadingIndicator;
};

class WiredItem : public NetItem
{
    Q_OBJECT

public:
    WiredItem(QWidget *parent, UWiredDevice *device, UWiredConnection *connection);
    ~WiredItem() Q_DECL_OVERRIDE;

    UWiredConnection *connection();
    void updateView() Q_DECL_OVERRIDE;
    NetItemType itemType() Q_DECL_OVERRIDE;

protected:
    void initUi();
    void initConnection();

    bool eventFilter(QObject *object, QEvent *event) Q_DECL_OVERRIDE;

protected Q_SLOTS:
    void onConnectionClicked();

private:
    UWiredConnection *m_connection;
    DViewItemAction *m_connectionItem;
    UWiredDevice *m_device;
    QPushButton *m_button;

    DViewItemAction *m_connIconAction;
};

class WirelessItem : public NetItem
{
    Q_OBJECT

public:
    WirelessItem(QWidget *parent, UWirelessDevice *device, UAccessPoints *ap);
    ~WirelessItem() Q_DECL_OVERRIDE;

    const UAccessPoints *accessPoint();
    void updateView() Q_DECL_OVERRIDE;
    NetItemType itemType() Q_DECL_OVERRIDE;
    static QString getStrengthStateString(int strength);

protected:
    bool eventFilter(QObject *object, QEvent *event) Q_DECL_OVERRIDE;

private:
    void initUi();
    void initConnection();
    void updateSrcirityIcon();
    void updateWifiIcon();
    void updateConnectionStatus();

private Q_SLOTS:
    void onConnection();

private:
    UAccessPoints *m_accessPoint;
    DViewItemAction *m_connLabel;
    QPushButton *m_button;
    UWirelessDevice *m_device;
    DViewItemAction *m_securityAction;
    DViewItemAction *m_wifiLabel;
    DSpinner *m_loadingStat;
};

#endif //  NETWORKAPPLETMODEL_H
