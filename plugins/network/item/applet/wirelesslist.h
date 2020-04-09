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

#ifndef WIRELESSAPPLET_H
#define WIRELESSAPPLET_H

#include "devicecontrolwidget.h"
#include "accesspoint.h"

#include <QScrollArea>
#include <QVBoxLayout>
#include <QList>
#include <QTimer>
#include <QPointer>

//#include <dpicturesequenceview.h>
#include <WirelessDevice>
#include <DSpinner>

#include <com_deepin_daemon_airplanemode.h>

using AirplanInter = com::deepin::daemon::AirplaneMode;


DWIDGET_USE_NAMESPACE

class AccessPointWidget;
class WirelessList : public QScrollArea
{
    Q_OBJECT

public:
    explicit WirelessList(dde::network::WirelessDevice *deviceIter, QWidget *parent = 0);
    ~WirelessList();

    QWidget *controlPanel();
    int APcount();

public Q_SLOTS:
    void setDeviceInfo(const int index);
    void onEnableButtonToggle(const bool enable);

signals:
    void requestSetDeviceEnable(const QString &path, const bool enable) const;
    void requestActiveAP(const QString &devPath, const QString &apPath, const QString &uuid) const;
    void requestDeactiveAP(const QString &devPath) const;
    void requestWirelessScan();
    void requestUpdatePopup();

private slots:
    void loadAPList();
    void APAdded(const QJsonObject &apInfo);
    void APRemoved(const QJsonObject &apInfo);
    void APPropertiesChanged(const QJsonObject &apInfo);
    void updateAPList();
    void onDeviceEnableChanged(const bool enable);
    void activateAP(const QString &apPath, const QString &ssid);
    void deactiveAP();
    void updateIndicatorPos();
    void onActiveConnectionInfoChanged();
    void onActivateApFailed(const QString &apPath, const QString &uuid);
    void onHotspotEnabledChanged(const bool enabled);

private:
    AccessPoint accessPointBySsid(const QString &ssid);
    AccessPointWidget *accessPointWidgetByAp(const AccessPoint ap);

private:
    QPointer<dde::network::WirelessDevice> m_device;

    AccessPoint m_activeAP;
    AccessPoint m_activatingAP;
    AccessPoint m_activeHotspotAP;
    QList<AccessPoint> m_apList;
    QList<AccessPointWidget*> m_apwList;

    QTimer *m_updateAPTimer;
//    Dtk::Widget::DPictureSequenceView *m_indicator;
    DSpinner *m_loadingStat;

    QVBoxLayout *m_centralLayout;
    QWidget *m_centralWidget;
    DeviceControlWidget *m_controlPanel;

    AccessPointWidget *m_clickedAPW;

    AirplanInter *m_airplaninter;

public:
    bool isHotposActive;

};

#endif // WIRELESSAPPLET_H
