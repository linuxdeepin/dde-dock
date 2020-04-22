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

namespace Dock {
class TipsWidget;
}

class WirelessItem : public DeviceItem
{
    Q_OBJECT

public:
    enum WirelessStatus
    {
        Unknow              = 0,
        Enabled             = 0x00010000,
        Disabled            = 0x00020000,
        Connected           = 0x00040000,
        Disconnected        = 0x00080000,
        Connecting          = 0x00100000,
        Authenticating      = 0x00200000,
        ObtainingIP         = 0x00400000,
        ObtainIpFailed      = 0x00800000,
        ConnectNoInternet   = 0x01000000,
        Failed              = 0x02000000,
    };
    Q_ENUM(WirelessStatus)

public:
    explicit WirelessItem(dde::network::WirelessDevice *device);
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

public Q_SLOTS:
    // set the device name displayed
    // in the top-left corner of the applet
    void setDeviceInfo(const int index);

Q_SIGNALS:
    void requestActiveAP(const QString &devPath, const QString &apPath, const QString &uuid) const;
    void requestDeactiveAP(const QString &devPath) const;
    void feedSecret(const QString &connectionPath, const QString &settingName, const QString &password, const bool autoConnect);
    void cancelSecret(const QString &connectionPath, const QString &settingName);
    void queryActiveConnInfo();
    void requestWirelessScan();
    void createApConfig(const QString &devPath, const QString &apPath);
    void queryConnectionSession( const QString &devPath, const QString &uuid );
    void deviceStateChanged();

protected:
    bool eventFilter(QObject *o, QEvent *e);

private slots:
    void init();
    void adjustHeight(bool visibel);

private:
    int m_index;
    QTimer *m_refreshTimer;
    QWidget *m_wirelessApplet;

    WirelessList *m_APList;
    QJsonObject m_activeApInfo;
};

#endif // WIRELESSITEM_H
