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
#include <QList>
#include <QPointer>

#include <WirelessDevice>
#include <NetworkModel>

#include <DSpinner>

#include <com_deepin_daemon_airplanemode.h>

using AirplanInter = com::deepin::daemon::AirplaneMode;


DWIDGET_USE_NAMESPACE

class AccessPointWidget;
class QVBoxLayout;
class QTimer;
class WirelessList : public QScrollArea
{
    Q_OBJECT

public:
    explicit WirelessList(dde::network::WirelessDevice *deviceIter, dde::network::NetworkModel *model, QWidget *parent = 0);
    ~WirelessList();

    QWidget *controlPanel();
    int APcount();
    /**
     * @def setModel
     * @brief 设置model
     * @param model
     */
    void setModel(dde::network::NetworkModel *model);

public Q_SLOTS:
    void setDeviceInfo(const int index);
    void onEnableButtonToggle(const bool enable);

signals:
    void requestSetDeviceEnable(const QString &path, const bool enable) const;
    void requestActiveAP(const QString &devPath, const QString &apPath, const QString &uuid) const;
    void requestDeactiveAP(const QString &devPath) const;
    void requestUpdatePopup();

private:
    /**
     * @def initUI
     * @brief 初始化ui界面
     */
    void initUI();
    /**
     * @def initConnect
     * @brief 初始化信号和槽
     */
    void initConnect();

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
    dde::network::NetworkModel *m_model;

    AccessPoint m_activeAP;
    AccessPoint m_activatingAP;
    AccessPoint m_activeHotspotAP;
    QList<AccessPoint> m_apList;
    QList<AccessPointWidget *> m_apwList;

    QTimer *m_updateAPTimer;
    DSpinner *m_loadingStat;

    QVBoxLayout *m_centralLayout;
    QWidget *m_centralWidget;
    //无线网顶部 --（名称  刷新  开关）在这个里面
    DeviceControlWidget *m_controlPanel;

    AccessPointWidget *m_clickedAPW;

    AirplanInter *m_airplaninter;

public:
    bool isHotposActive;
};

#endif // WIRELESSAPPLET_H
