/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#ifndef WIRELESSITEM_H
#define WIRELESSITEM_H

#include "constants.h"
#include "deviceitem.h"
#include "applet/wirelesslist.h"

#include <QHash>
#include <QLabel>

#include <WirelessDevice>
#include <NetworkModel>

class TipsWidget;
class WirelessItem : public DeviceItem
{
    Q_OBJECT

public:
    enum WirelessStatus {
        Unknown              = 0,
        Enabled              = 1,
        Disabled             = 2,
        Connected            = 4,
        Disconnected         = 8,
        Connecting           = 16,
        Authenticating       = 32,
        ObtainingIP          = 64,
        ObtainIpFailed       = 128,
        ConnectNoInternet    = 256,
        Failed               = 512};
    Q_ENUM(WirelessStatus)

public:
    explicit WirelessItem(dde::network::WirelessDevice *device, dde::network::NetworkModel *model);
    ~WirelessItem();

    QWidget *itemApplet();
    int APcount();
    bool deviceEanbled();
    void setDeviceEnabled(bool enable);
    WirelessStatus getDeviceState();
    QJsonObject &getConnectedApInfo();
    QJsonObject getActiveWirelessConnectionInfo();
    inline int deviceInfo() { return m_index; }
    void setControlPanelVisible(bool visible);

private:
    /**
     * @def initConnect
     * @brief 初始化信号槽
     */
    void initConnect();


public Q_SLOTS:
    // set the device name displayed
    // in the top-left corner of the applet
    void setDeviceInfo(const int index);

Q_SIGNALS:
    void requestActiveAP(const QString &devPath, const QString &apPath, const QString &uuid) const;
    void requestDeactiveAP(const QString &devPath) const;
    void requestWirelessScan();
    void deviceStateChanged();
    void activeApInfoChanged();

//protected:
//    bool eventFilter(QObject *o, QEvent *e);

private slots:
    void init();
    void adjustHeight(bool visibel);

private:
    int m_index;
    QWidget *m_wirelessApplet;

    WirelessList *m_APList;
    QJsonObject m_activeApInfo;
    dde::network::NetworkModel *m_model;
};

#endif // WIRELESSITEM_H
