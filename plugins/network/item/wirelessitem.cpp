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

#include "wirelessitem.h"
#include "networkplugin.h"
#include "../frame/util/imageutil.h"
#include "../widgets/tipswidget.h"
#include <DGuiApplicationHelper>

#include <QPainter>
#include <QMouseEvent>
#include <QApplication>
#include <QIcon>

using namespace dde::network;
DGUI_USE_NAMESPACE

WirelessItem::WirelessItem(WirelessDevice *device)
    : DeviceItem(device),
      m_refreshTimer(new QTimer(this)),
      m_wirelessApplet(new QWidget),
      m_APList(nullptr)
{
    m_refreshTimer->setSingleShot(true);
    m_refreshTimer->setInterval(100);

    m_wirelessApplet->setVisible(false);

    connect(m_refreshTimer, &QTimer::timeout, [&]{
        if (m_device.isNull()) {
            return;
        }
        WirelessDevice *dev = static_cast<WirelessDevice *>(m_device.data());
        // the status is Activated and activeApInfo is empty if the hotspot status of this wireless device is enabled
        if (m_device->status() == NetworkDevice::Activated && dev->activeApInfo().isEmpty() && !dev->hotspotEnabled()) {
            Q_EMIT queryActiveConnInfo();
            return;
        }
    });
    connect(m_device, static_cast<void (NetworkDevice::*)(const QString &statStr) const>(&NetworkDevice::statusChanged), this, &WirelessItem::deviceStateChanged);
    connect(static_cast<WirelessDevice *>(m_device.data()), &WirelessDevice::activeApInfoChanged, m_refreshTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(static_cast<WirelessDevice *>(m_device.data()), &WirelessDevice::activeWirelessConnectionInfoChanged, m_refreshTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, [ = ] {
        update();
    });

    connect(static_cast<WirelessDevice *>(m_device.data()), &WirelessDevice::apInfoChanged, this, [ = ](QJsonObject info) {
        const auto &activeApInfo = static_cast<WirelessDevice *>(m_device.data())->activeApInfo();
        if (activeApInfo.value("Ssid").toString() == info.value("Ssid").toString()) {
            m_activeApInfo = info;
        }
        update();
    });

    init();
}

WirelessItem::~WirelessItem()
{
    m_APList->deleteLater();
    m_APList->controlPanel()->deleteLater();
}

QWidget *WirelessItem::itemApplet()
{
    return m_wirelessApplet;
}

int WirelessItem::APcount()
{
    return m_APList->APcount();
}

bool WirelessItem::deviceActivated()
{
    switch (m_device->status()) {
    case NetworkDevice::Unknow:
    case NetworkDevice::Unmanaged:
    case NetworkDevice::Unavailable:
    case NetworkDevice::Disconnected:
    case NetworkDevice::Deactivation:
    case NetworkDevice::Failed: {
        return false;
    }
    case NetworkDevice::Prepare:
    case NetworkDevice::Config:
    case NetworkDevice::NeedAuth:
    case NetworkDevice::IpConfig:
    case NetworkDevice::IpCheck:
    case NetworkDevice::Secondaries:
    case NetworkDevice::Activated: {
        return true;
    }
    }
}

void WirelessItem::setDeviceEnabled(bool enable)
{
    m_APList->onEnableButtonToggle(enable);
}

WirelessItem::WirelessStatus WirelessItem::getDeviceState()
{
    if (!m_device->enabled()) {
        return Disabled;
    }
    if (m_device->status() == NetworkDevice::Activated
            && NetworkModel::connectivity() != Connectivity::Full) {
        return ConnectNoInternet;
    }
    if (m_device->obtainIpFailed()) {
        return ObtainIpFailed;
    }

    switch (m_device->status()) {
    case NetworkDevice::Unknow:        return Unknow;
    case NetworkDevice::Unmanaged:
    case NetworkDevice::Unavailable:
    case NetworkDevice::Disconnected:  return Disconnected;
    case NetworkDevice::Prepare:
    case NetworkDevice::Config:        return Connecting;
    case NetworkDevice::NeedAuth:      return Authenticating;
    case NetworkDevice::IpConfig:
    case NetworkDevice::IpCheck:
    case NetworkDevice::Secondaries:   return ObtainingIP;
    case NetworkDevice::Activated:     return Connected;
    case NetworkDevice::Deactivation:
    case NetworkDevice::Failed:        return Failed;
    }
}

QJsonObject &WirelessItem::getConnectedApInfo()
{
    return  m_activeApInfo;
}

QJsonObject WirelessItem::getActiveWiredConnectionInfo()
{
    return static_cast<WirelessDevice *>(m_device.data())->activeWirelessConnectionInfo();
}

void WirelessItem::setDeviceInfo(const int index)
{
    m_APList->setDeviceInfo(index);
}

bool WirelessItem::eventFilter(QObject *o, QEvent *e)
{
    if (o == m_APList && e->type() == QEvent::Resize)
        QMetaObject::invokeMethod(this, "adjustHeight", Qt::QueuedConnection);
    if (o == m_APList && e->type() == QEvent::Show)
        Q_EMIT requestWirelessScan();

    return false;
}

void WirelessItem::init()
{
    m_APList = new WirelessList(static_cast<WirelessDevice *>(m_device.data()));
    m_APList->installEventFilter(this);
    m_APList->setObjectName("wireless-" + m_device->path());

    QVBoxLayout *vLayout = new QVBoxLayout;
    vLayout->addWidget(m_APList->controlPanel());
    vLayout->addWidget(m_APList);
    vLayout->setMargin(0);
    vLayout->setSpacing(0);
    m_wirelessApplet->setLayout(vLayout);

    connect(m_APList, &WirelessList::requestSetDeviceEnable, this, &WirelessItem::requestSetDeviceEnable);
    connect(m_APList, &WirelessList::requestActiveAP, this, &WirelessItem::requestActiveAP);
    connect(m_APList, &WirelessList::requestDeactiveAP, this, &WirelessItem::requestDeactiveAP);
    connect(m_APList, &WirelessList::requestWirelessScan, this, &WirelessItem::requestWirelessScan);
    connect(m_APList, &WirelessList::requestUpdatePopup, this, &WirelessItem::deviceStateChanged);

    QTimer::singleShot(0, this, [ = ]() {
        m_refreshTimer->start();
    });
}

void WirelessItem::adjustHeight()
{
    m_wirelessApplet->setFixedHeight(m_APList->height() + m_APList->controlPanel()->height());
}

