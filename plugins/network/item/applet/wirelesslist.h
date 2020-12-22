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
//    void requestActiveAP(const QString &devPath, const QString &apPath, const QString &uuid) const;
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
//    void loadAPList();
    void APAdded(const QJsonObject &apInfo);
    void ApRemoved(const QJsonObject &apInfo);
    void ApInfoChange(const QJsonObject &apInfo);
    //刷新前端页面数据
    void updateView();
    void onDeviceEnableChanged(const bool enable);
    void ActiveConnectChange(const QJsonObject &activeAp);
    void onActivateApFailed(const QString &apPath, const QString &uuid);
    //热点功能目前展示出来
    void onHotspotEnabledChanged(const bool enabled);

private:
    AccessPointWidget *accessPointWidgetBySsid(const QString &ssid);

private:
    QPointer<dde::network::WirelessDevice> m_device;
    dde::network::NetworkModel *m_model;

    //表示连接中和连接成功的wifi
    AccessPointWidget *m_activeAp;
    //这里使用一个click值保存，由于企业wifi在被连接时，
    //并不会立即将其变成活动中的wifi，而会ActivateAccessPoint在这个接口中表现出来
    //所以不能直接用m_activeAp保存状态，这样是不对的4
    AccessPointWidget *m_clickAp;
    AccessPoint m_activeHotspotAP;
    QList<AccessPointWidget *> m_apwList;
    QTimer *m_clickIntervalTimer;


    //刷新操作的函数
    QTimer *m_updateTimer;
    QVBoxLayout *m_centralLayout;
    QWidget *m_centralWidget;
    //无线网顶部 --（名称  刷新  开关）在这个里面
    DeviceControlWidget *m_controlPanel;

public:
    bool isHotposActive;
};

#endif // WIRELESSAPPLET_H
