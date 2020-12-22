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

#include "wirelesslist.h"
#include "accesspointwidget.h"
#include "constants.h"

#include <QJsonDocument>
#include <QScreen>
#include <QDebug>
#include <QGuiApplication>
#include <QVBoxLayout>
#include <QScrollBar>
#include <QTimer>

#include <dinputdialog.h>
#include <DDBusSender>

DWIDGET_USE_NAMESPACE

using namespace dde::network;

extern const int ItemWidth = 250;
extern const int ItemMargin = 10;
extern const int ItemHeight;

WirelessList::WirelessList(WirelessDevice *deviceIter, NetworkModel *model, QWidget *parent)
    : QScrollArea(parent)
    , m_device(deviceIter)
    , m_model(model)
    , m_activeAp(nullptr)
    , m_centralLayout(new QVBoxLayout)
    , m_activeHotspotAP(AccessPoint())
    , m_centralWidget(new QWidget(this))
    , m_controlPanel(new DeviceControlWidget(this))
    , m_updateTimer(new QTimer(this))
    , m_clickIntervalTimer(new QTimer(this))
{
    initUI();
    initConnect();
    //由于信号和槽函数会在连接之前将数据发过来，则刚启动的时候可能已经更新完成数据了，所以导致页面上没有数据
    m_device->initWirelessData();
    //更新开关状态
    Q_EMIT m_model->initDeviceEnable(m_device->path());
    m_updateTimer->setInterval(0);
    m_updateTimer->setSingleShot(true);
    m_clickIntervalTimer->setInterval(500);
    m_clickIntervalTimer->setSingleShot(true);
}

WirelessList::~WirelessList()
{
}

void WirelessList::initUI()
{
    //设置固定高度
    setFixedHeight(ItemHeight);

    m_centralWidget->setFixedWidth(ItemWidth - 2 * ItemMargin);
    m_centralWidget->setLayout(m_centralLayout);
    m_centralLayout->addWidget(m_controlPanel);
    m_centralLayout->setSpacing(0);
    m_centralLayout->setMargin(0);

    setWidget(m_centralWidget);
    setFrameShape(QFrame::NoFrame);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_centralWidget->setAutoFillBackground(false);
    viewport()->setAutoFillBackground(false);

    isHotposActive = false;
}

void WirelessList::initConnect()
{
    //后端开关wifi
    connect(m_device, &WirelessDevice::enableChanged, this,
            [ = ](){    qDebug() << "Device change signal from daemon, state:" << m_device->enabled();
                        onDeviceEnableChanged(m_device->enabled());

                        if (m_device->enabled())
                            //这里会在页面创建的时候去初始化一次，所以无需在构造函数中再调用
                            m_device->initWirelessData();});
    //开关wifi
    connect(m_controlPanel, &DeviceControlWidget::enableButtonToggled, this, &WirelessList::onEnableButtonToggle);
    //刷新wifi按钮
    connect(m_controlPanel, &DeviceControlWidget::requestRefresh, m_model, &NetworkModel::updateApList);
    connect(m_updateTimer, &QTimer::timeout, this, &WirelessList::updateView);
    //修改wifi数据的信号
    connect(m_device, &WirelessDevice::apAdded, this, &WirelessList::APAdded);
    connect(m_device, &WirelessDevice::apRemoved, this, &WirelessList::ApRemoved);
    connect(m_device, &WirelessDevice::apInfoChanged, this, &WirelessList::ApInfoChange);
    //当wifi连接时，需要先去后端验证是否是企业wifi再做处理，所以要关联该信号
    connect(m_device, &WirelessDevice::activateAccessPointFailed, this, &WirelessList::onActivateApFailed);
    //活动的wifi进行处理
    connect(m_device, &WirelessDevice::activeConnectionsChanged, this, &WirelessList::ActiveConnectChange);
    //热点
    connect(m_device, &WirelessDevice::hotspotEnabledChanged, this, &WirelessList::onHotspotEnabledChanged);
}

QWidget *WirelessList::controlPanel()
{
    return m_controlPanel;
}

int WirelessList::APcount()
{
    return m_apwList.size();
}

void WirelessList::APAdded(const QJsonObject &apInfo)
{
    //wifi信号强度
    int strength = apInfo["Strength"].toInt();
    QString Ssid = apInfo["Ssid"].toString();
    //当信号小于等于5则不加入列表中
    if (strength < 5) {
        qDebug() << "append AP fail, " << "Ssid = " << Ssid  << ", Strength =" << strength;
        return;
    }
    AccessPointWidget *apw = accessPointWidgetBySsid(Ssid);
    if (apw) {
        ApInfoChange(apInfo);
    } else {
        apw = new AccessPointWidget(apInfo);
        m_apwList.append(apw);
        //加入到layout里面
        m_centralLayout->addWidget(apw);
    }

    //刚添加的时候给一个状态,防止重启刚打开的时候状态错误
    if (m_device->activeApSsid() == Ssid) {
        ActiveConnectChange(m_device->activeApInfo());
    }

    //这里做一个刷新，防止wifi数据在添加后没有刷新,信号槽后面绑定，防止出现多次刷新的情况
    m_updateTimer->start();

    //关联wifi刷新信号
    connect(apw, &AccessPointWidget::apChange, this, [ = ](){
        if (!m_updateTimer->isActive())
            m_updateTimer->start();});
    //关联点击连接信号
    connect(apw, &AccessPointWidget::requestConnectAP, m_model,
            [ = ](const QString &apPath, const QString &uuid) {
                if (m_clickIntervalTimer->isActive()) return;
                m_clickIntervalTimer->start();
                m_clickAp = apw;
                Q_EMIT m_model->requestConnectAp(m_device->path(), apPath, uuid);});
    //关联断开连接信号
    connect(apw, &AccessPointWidget::requestDisconnectAP, m_model, &NetworkModel::requestDisconnctAP);
}

void WirelessList::ApInfoChange(const QJsonObject &apInfo)
{
    int strength = apInfo["Strength"].toInt();
    QString ssid = apInfo["Ssid"].toString();
    AccessPointWidget *apw = accessPointWidgetBySsid(ssid);
    //当apw等于空，可能是之前在添加的时候就wifi强度就小于5,然后wifi强度变成大于5了
    if (!apw) {
        if (strength >= 5) {
            qDebug() << "Network strength increased to more than 5 " << "ssid = " << ssid << ", Current intensity = " << strength;         ;
            APAdded(apInfo);
        }
        return;
    }
    if (strength < 5) {
       ApRemoved(apInfo);
    } else {
        apw->updateApInfo(apInfo);
    }
}

void WirelessList::ApRemoved(const QJsonObject &apInfo)
{
    QString ssid = apInfo["Ssid"].toString();
    AccessPointWidget *apw = accessPointWidgetBySsid(ssid);
    if (!apw) return;
    if (m_activeAp == apw)
        m_activeAp = nullptr;
    const int mIndex = m_apwList.indexOf(apw);
    m_centralLayout->removeWidget(apw);
    m_apwList.removeAt(mIndex);
    delete apw;
    apw = nullptr;
    m_updateTimer->start();
}

void WirelessList::setDeviceInfo(const int index)
{
    if (m_device.isNull()) {
        return;
    }

    // set device enable state
    m_controlPanel->setDeviceEnabled(m_device->enabled());
    // set device name
    if (index == -1)
        m_controlPanel->setDeviceName(tr("Wireless Network"));
    else
        m_controlPanel->setDeviceName(tr("Wireless Network %1").arg(index));
}

void WirelessList::updateView()
{
    //代码调试的时候防止直接调用刷新
    Q_ASSERT(sender() == m_updateTimer);
    //当适配器消失，则直接退出
    if (m_device.isNull()) {
        return;
    }
    qDebug() << "m_device->enabled()" << m_device->enabled();
    if (m_device->enabled()) {
        std::sort(m_apwList.begin(), m_apwList.end(), [&](const AccessPointWidget *apw1, const AccessPointWidget *apw2) {
            if (apw1 == m_activeAp)
                return true;

            if (apw2 == m_activeAp)
                return false;

            if (apw1->strength() != apw2->strength())
                return apw1->strength() > apw2->strength();

            return  apw1->ssid() > apw2->ssid();
        });
        for (AccessPointWidget *apw: m_apwList) {
            m_centralLayout->removeWidget(apw);
            m_centralLayout->addWidget(apw);    
        }

    }

    const int contentHeight = APcount() * ItemHeight;
    m_centralWidget->setFixedHeight(contentHeight);
    setFixedHeight(contentHeight);
    emit requestUpdatePopup();

    return;
}

AccessPointWidget *WirelessList::accessPointWidgetBySsid(const QString &ssid)
{
    for (auto apw: m_apwList) {
        if (ssid == apw->ssid())
            return apw;
    }
    return nullptr;
}

void WirelessList::onEnableButtonToggle(const bool enable)
{
    if (m_device.isNull()) {
        return;
    }
    qDebug() <<"click enable , set Enable =" << enable;
    //直接调用dde::network::networkModel中的接口，防止数据出现延迟之类的问题
    m_model->onDeviceEnable(m_device->path(), enable);
    onDeviceEnableChanged(enable);
}

void WirelessList::onDeviceEnableChanged(const bool enable)
{
    m_controlPanel->setDeviceEnabled(enable);
    m_centralLayout->setEnabled(enable);
}

void WirelessList::ActiveConnectChange(const QJsonObject &activeAp)
{
    if (activeAp.isEmpty()) return;
    //ps: 这里的activeAp可能是空的，所以不要在判断m_activeAp为空之前使用
    // 0:Unknow, 1:Activating, 2:Activated, 3:Deactivating, 4:Deactivated
    QString activeSsid = activeAp["Id"].toString();
    const int state = activeAp["State"].toInt(0);
    if (state == AccessPointWidget::ApState::Deactivated) {
        //这里防止出现网络热点消失，然而dock上的网络图标并没有重置的问题
        if (m_activeAp)
            m_activeAp->setActiveState(AccessPointWidget::ApState(state));

        m_activeAp = nullptr;
    } else {
        m_activeAp = accessPointWidgetBySsid(activeSsid);
    }

    if (!m_activeAp) {
        m_updateTimer->start();
        qDebug() << "The network in the current connection cannot be found in the list.";
        qDebug() << "Ssid = " << activeSsid;
        return;
    }

    m_activeAp->setActiveState(AccessPointWidget::ApState(state));
}

void WirelessList::onActivateApFailed(const QString &apPath, const QString &uuid)
{
    if (m_device.isNull() && !m_clickAp) {
        return;
    }
    if (m_clickAp->path() == apPath) {
        qDebug() << "wireless connect failed and may require more configuration,"
                 << "path:" << m_clickAp->path() << "ssid" << m_clickAp->ssid() << "devicePath:" << m_device->path()
                 << "secret:" << m_clickAp->secured() << "strength:" << m_clickAp->strength() << "uuid:" << uuid;
        //打开网络相关的无线网页面
        DDBusSender()
        .service("com.deepin.dde.ControlCenter")
        .interface("com.deepin.dde.ControlCenter")
        .path("/com/deepin/dde/ControlCenter")
        .method("ShowPage")
        .arg(QString("network"))
        .arg(m_device->path())
        .call();
    }
}

void WirelessList::onHotspotEnabledChanged(const bool enabled)
{
    // Note: the obtained hotspot info is not complete
    m_activeHotspotAP = enabled ? AccessPoint(m_device->activeHotspotInfo().value("Hotspot").toObject())
                        : AccessPoint();
    isHotposActive = enabled;
}
