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

#ifndef WIREDITEM_H
#define WIREDITEM_H

#include "deviceitem.h"

#include <WiredDevice>

#include <DGuiApplicationHelper>
#include <DSpinner>

using namespace dde::network;
DGUI_USE_NAMESPACE
DWIDGET_USE_NAMESPACE

class QLabel;
class TipsWidget;
class HorizontalSeperator;
class StateButton;
class WiredItem : public DeviceItem
{
    Q_OBJECT

public:
    enum WiredStatus {
        Unknow              = 0,
        Enabled             = 0x00000001,
        Disabled            = 0x00000002,
        Connected           = 0x00000004,
        Disconnected        = 0x00000008,
        Connecting          = 0x00000010,
        Authenticating      = 0x00000020,
        ObtainingIP         = 0x00000040,
        ObtainIpFailed      = 0x00000080,
        ConnectNoInternet   = 0x00000100,
        Nocable             = 0x00000200,
        Failed              = 0x00000400,
    };
    Q_ENUM(WiredStatus)

public:
    explicit WiredItem(dde::network::WiredDevice *device, const QString &deviceName, QWidget *parent = nullptr);
    void setTitle(const QString &name);
    bool deviceEabled();
    void setDeviceEnabled(bool enabled);
    WiredStatus getDeviceState();
    QJsonObject getActiveWiredConnectionInfo();
    inline QString &deviceName() { return m_deviceName; }
    void setThemeType(DGuiApplicationHelper::ColorType themeType);

signals:
    void requestActiveConnection(const QString &devPath, const QString &uuid);
    void wiredStateChanged();
    void enableChanged();
    void activeConnectionChanged();

private slots:
    void deviceStateChanged(NetworkDevice::DeviceStatus state);
    void changedActiveWiredConnectionInfo(const QJsonObject &connInfo);

private:
    QString m_deviceName;
    QLabel *m_connectedName;
    QLabel *m_wiredIcon;
    StateButton *m_stateButton;
    DSpinner *m_loadingStat;

    HorizontalSeperator *m_line;
};

#endif // WIREDITEM_H
