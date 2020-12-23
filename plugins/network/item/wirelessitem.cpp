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
#include <QLayout>

using namespace dde::network;
DGUI_USE_NAMESPACE

WirelessItem::WirelessItem(WirelessDevice *device, NetworkModel *model)
    : DeviceItem(device),
      m_model(model),
      m_wirelessApplet(new QWidget),
      m_APList(nullptr)
{
    m_wirelessApplet->setVisible(false);
    initConnect();
    init();
}

void WirelessItem::initConnect()
{
    //获取状态提示语
    connect(m_device, static_cast<void (NetworkDevice::*)(const QString &statStr) const>(&NetworkDevice::statusChanged), this, &WirelessItem::deviceStateChanged);
    //获取wifi的连接，获取wifi的信号强度
    connect(static_cast<WirelessDevice *>(m_device.data()), &WirelessDevice::apInfoChanged, this, [ = ]() {
        const auto &activeApInfo = static_cast<WirelessDevice *>(m_device.data())->activeApInfo();
        if (activeApInfo != m_activeApInfo) {
            m_activeApInfo = activeApInfo;
        }
        update();
    });
    connect(m_model, &NetworkModel::activeConnInfoChanged, this , [ = ]() {
        m_activeApInfo = static_cast<WirelessDevice *>(m_device.data())->activeApInfo();
        update();
    });
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, [ = ] {
        update();
    });

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

bool WirelessItem::deviceEanbled()
{
    return m_device->enabled();
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
    case NetworkDevice::Unknown:       return Unknown;
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
    Q_UNREACHABLE();
}

QJsonObject &WirelessItem::getConnectedApInfo()
{
    return  m_activeApInfo;
}

QJsonObject WirelessItem::getActiveWirelessConnectionInfo()
{
    return static_cast<WirelessDevice *>(m_device.data())->activeWirelessConnectionInfo();
}

void WirelessItem::setControlPanelVisible(bool visible)
{
    auto layout = m_wirelessApplet->layout();
    auto controlPanel = m_APList->controlPanel();
    if (layout && controlPanel) {
        if (visible) {
            layout->removeWidget(controlPanel);
            layout->removeWidget(m_APList);

            layout->addWidget(controlPanel);
            layout->addWidget(m_APList);
        } else {
            layout->removeWidget(controlPanel);
        }
        controlPanel->setVisible(visible);
        adjustHeight(visible);
    }
}

void WirelessItem::setDeviceInfo(const int index)
{
    m_APList->setDeviceInfo(index);
    m_index = index;
}

//bool WirelessItem::eventFilter(QObject *o, QEvent *e)
//{
//    if (o == m_APList && e->type() == QEvent::Resize)
//        QMetaObject::invokeMethod(this, "adjustHeight", Qt::QueuedConnection,Q_ARG(bool, m_APList->controlPanel()->isVisible()));
//    if (o == m_APList && e->type() == QEvent::Show)
//        Q_EMIT requestWirelessScan();

//    return false;
//}

void WirelessItem::init()
{
    m_APList = new WirelessList(static_cast<WirelessDevice *>(m_device.data()), m_model);
    m_APList->installEventFilter(this);
    m_APList->setObjectName("wireless-" + m_device->path());

    QVBoxLayout *vLayout = new QVBoxLayout;
    vLayout->addWidget(m_APList->controlPanel());
    vLayout->addWidget(m_APList);
    vLayout->setMargin(0);
    vLayout->setSpacing(0);
    m_wirelessApplet->setLayout(vLayout);

//    connect(m_APList, &WirelessList::requestDeactiveAP, this, &WirelessItem::requestDeactiveAP);
    connect(m_APList, &WirelessList::requestUpdatePopup, this, &WirelessItem::deviceStateChanged);
}

void WirelessItem::adjustHeight(bool visibel)
{
    auto controlPanel = m_APList->controlPanel();
    if (!controlPanel)
        return;

    auto height = visibel ? (m_APList->height() + controlPanel->height())
                  : m_APList->height();
    m_wirelessApplet->setFixedHeight(height);
}
