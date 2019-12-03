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

#include <QJsonDocument>
#include <QScreen>
#include <QDebug>
#include <QGuiApplication>

#include <dinputdialog.h>
#include <QScrollBar>
#include <DDBusSender>

DWIDGET_USE_NAMESPACE

using namespace dde::network;

#define WIDTH           300
#define MAX_HEIGHT      300
#define ITEM_HEIGHT     30

WirelessList::WirelessList(WirelessDevice *deviceIter, QWidget *parent)
    : QScrollArea(parent),

      m_device(deviceIter),
      m_activeAP(),

      m_updateAPTimer(new QTimer(this)),

      m_centralLayout(new QVBoxLayout),
      m_centralWidget(new QWidget),
      m_controlPanel(new DeviceControlWidget)
{
    setFixedHeight(WIDTH);

    const auto ratio = devicePixelRatioF();

    m_updateAPTimer->setSingleShot(true);
    m_updateAPTimer->setInterval(100);

    m_centralWidget->setFixedWidth(WIDTH);
    m_centralWidget->setLayout(m_centralLayout);

    m_centralLayout->addWidget(m_controlPanel);
    m_centralLayout->setSpacing(0);
    m_centralLayout->setMargin(0);

    setWidget(m_centralWidget);
    setFrameShape(QFrame::NoFrame);
    setFixedWidth(300);
    setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    m_centralWidget->setAutoFillBackground(false);
    viewport()->setAutoFillBackground(false);

    m_indicator = new DPictureSequenceView(this);
    m_indicator->setPictureSequence(":/wireless/indicator/resources/wireless/spinner14/Spinner%1.png", QPair<int, int>(1, 91), 2);
    m_indicator->setFixedSize(QSize(14, 14) * ratio);
    m_indicator->setVisible(false);
    isHotposActive = false;

    connect(m_device, &WirelessDevice::apAdded, this, &WirelessList::APAdded);
    connect(m_device, &WirelessDevice::apRemoved, this, &WirelessList::APRemoved);
    connect(m_device, &WirelessDevice::apInfoChanged, this, &WirelessList::APPropertiesChanged);
    connect(m_device, &WirelessDevice::enableChanged, this, &WirelessList::onDeviceEnableChanged);
    connect(m_device, &WirelessDevice::activateAccessPointFailed, this, &WirelessList::onActivateApFailed);
    connect(m_device, &WirelessDevice::hotspotEnabledChanged, this, &WirelessList::onHotspotEnabledChanged);

    connect(m_controlPanel, &DeviceControlWidget::enableButtonToggled, this, &WirelessList::onEnableButtonToggle);
    connect(m_controlPanel, &DeviceControlWidget::requestRefresh, this, &WirelessList::requestWirelessScan);

    connect(m_updateAPTimer, &QTimer::timeout, this, &WirelessList::updateAPList);

    connect(m_device, &WirelessDevice::activeWirelessConnectionInfoChanged, this, &WirelessList::onActiveConnectionInfoChanged);
    connect(m_device, static_cast<void (WirelessDevice:: *) (NetworkDevice::DeviceStatus stat) const>(&WirelessDevice::statusChanged), m_updateAPTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
    connect(m_device, &WirelessDevice::activeConnectionsChanged, this, &WirelessList::updateIndicatorPos, Qt::QueuedConnection);

    connect(this->verticalScrollBar(), &QScrollBar::valueChanged, this, [=] {
        auto apw = accessPointWidgetByAp(m_activatingAP);
        if (!apw) return;

        const int h = -(apw->height() - m_indicator->height()) / 2;
        m_indicator->move(apw->mapTo(this, apw->rect().topRight()) - QPoint(35, h));
    });

    QMetaObject::invokeMethod(this, "loadAPList", Qt::QueuedConnection);
}

WirelessList::~WirelessList()
{
}

QWidget *WirelessList::controlPanel()
{
    return m_controlPanel;
}


void WirelessList::APAdded(const QJsonObject &apInfo)
{
    AccessPoint ap(apInfo);
    const auto mIndex = m_apList.indexOf(ap);
    if (mIndex != -1) {
        if (ap > m_apList.at(mIndex)) {
            m_apList.replace(mIndex, ap);
        } else {
            return;
        }
    } else {
        m_apList.append(ap);
    }

    m_updateAPTimer->start();
}

void WirelessList::APRemoved(const QJsonObject &apInfo)
{
    AccessPoint ap(apInfo);
    const auto mIndex = m_apList.indexOf(ap);
    if (mIndex != -1) {
        if (ap.path() == m_apList.at(mIndex).path()) {
            m_apList.removeAt(mIndex);
            m_updateAPTimer->start();
        }
    }
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

void WirelessList::loadAPList()
{
    if (m_device.isNull()) {
        return;
    }

    for (auto item : m_device->apList()) {
        AccessPoint ap(item.toObject());
        const auto mIndex = m_apList.indexOf(ap);
        if (mIndex != -1) {
            // indexOf() will use AccessPoint reimplemented function "operator==" as comparison condition
            // this means that the ssid of the AP is a comparison condition
            if (ap > m_apList.at(mIndex)) {
                m_apList.replace(mIndex, ap);
            }
        } else {
            m_apList.append(ap);
        }
    }

    m_updateAPTimer->start();
}

void WirelessList::APPropertiesChanged(const QJsonObject &apInfo)
{
    AccessPoint ap(apInfo);
    const auto mIndex = m_apList.indexOf(ap);
    if (mIndex != -1) {
        if (ap > m_apList.at(mIndex)) {
            m_apList.replace(mIndex, ap);
            m_updateAPTimer->start();
        }
    }
}

void WirelessList::updateAPList()
{
    Q_ASSERT(sender() == m_updateAPTimer);

    if (m_device.isNull()) {
        return;
    }

    int avaliableAPCount = 0;

    //if (m_networkInter->IsDeviceEnabled(m_device.dbusPath()))
    if (m_device->enabled())
    {
        if (m_device->hotspotEnabled()) {
            m_apList.clear();
            m_apList.append(m_activeHotspotAP);
        }

        // sort ap list by strength
        // std::sort(m_apList.begin(), m_apList.end(), std::greater<AccessPoint>());
        //        const bool wirelessActived = m_device.state() == NetworkDevice::Activated;

        // NOTE: Keep the amount consistent
        if(m_apList.size() > m_apwList.size()) {
            int i = m_apList.size() - m_apwList.size();
            for (int index = 0; index != i; index++) {
                AccessPointWidget *apw = new AccessPointWidget;
                apw->setFixedHeight(ITEM_HEIGHT);
                m_apwList << apw;
                m_centralLayout->addWidget(apw);

                connect(apw, &AccessPointWidget::requestActiveAP, this, &WirelessList::activateAP);
                connect(apw, &AccessPointWidget::requestDeactiveAP, this, &WirelessList::deactiveAP);
                connect(apw, &AccessPointWidget::requestActiveAP, this, [=] {
                    m_clickedAPW = apw;
                }, Qt::UniqueConnection);
            }
        } else if (m_apList.size() < m_apwList.size()) {
            if (!m_apwList.isEmpty()) {
                int i = m_apwList.size() - m_apList.size();
                for (int index = 0; index != i; index++) {
                    AccessPointWidget *apw = m_apwList.last();
                    m_apwList.removeLast();
                    m_centralLayout->removeWidget(apw);
                    apw->deleteLater();
                }
            }
        }

        std::sort(m_apList.begin(), m_apList.end(), [&] (const AccessPoint &ap1, const AccessPoint &ap2) {
            if (ap1 == m_activeAP)
                return true;

            if (ap2 == m_activeAP)
                return false;

            return ap1.strength() > ap2.strength();
        });

        // update the content of AccessPointWidget by the order of ApList
        for (int i = 0; i != m_apList.size(); i++) {
            m_apwList[i]->updateAP(m_apList[i]);
            ++avaliableAPCount;
        }

        // update active AP state
        NetworkDevice::DeviceStatus deviceStatus = m_device->status();
        if (!m_apwList.isEmpty()) {
            AccessPointWidget *apw = m_apwList.first();

            apw->setActiveState(deviceStatus);
        }

        // If the order of item changes
        if (m_indicator->isVisible() && !m_activatingAP.isEmpty() && m_apList.contains(m_activatingAP)) {
            AccessPointWidget *apw = accessPointWidgetByAp(m_activatingAP);
            if (apw) {
                const int h = -(apw->height() - m_indicator->height()) / 2;
                m_indicator->move(apw->mapTo(this, apw->rect().topRight()) - QPoint(35, h));
            }
        }

        if (deviceStatus <= NetworkDevice::Disconnected || deviceStatus >= NetworkDevice::Activated) {
            m_indicator->stop();
            m_indicator->hide();
        }
    }

    const int contentHeight = avaliableAPCount * ITEM_HEIGHT;
    m_centralWidget->setFixedHeight(contentHeight);
    setFixedHeight(std::min(contentHeight, MAX_HEIGHT));

    QTimer::singleShot(100, this, &WirelessList::updateIndicatorPos);
}

void WirelessList::onEnableButtonToggle(const bool enable)
{
    if (m_device.isNull()) {
        return;
    }

    Q_EMIT requestSetDeviceEnable(m_device->path(), enable);
    m_updateAPTimer->start();
}

void WirelessList::onDeviceEnableChanged(const bool enable)
{
    m_controlPanel->setDeviceEnabled(enable);
    m_updateAPTimer->start();
}

void WirelessList::activateAP(const QString &apPath, const QString &ssid)
{
    if (m_device.isNull()) {
        return;
    }

    QString uuid;

    QList<QJsonObject> connections = m_device->connections();
    for (auto item : connections) {
        if (item.value("Ssid").toString() != ssid)
            continue;

        uuid = item.value("Uuid").toString();
        if (!uuid.isEmpty())
            break;
    }

    Q_EMIT requestActiveAP(m_device->path(), apPath, uuid);
}

void WirelessList::deactiveAP()
{
    if (m_device.isNull()) {
        return;
    }

    Q_EMIT requestDeactiveAP(m_device->path());
}

void WirelessList::updateIndicatorPos()
{
    QString activeSsid;
    for (auto activeConnObj : m_device->activeConnections()) {
        if (activeConnObj.value("Vpn").toBool(false)) {
            continue;
        }
        // the State of Active Connection
        // 0:Unknow, 1:Activating, 2:Activated, 3:Deactivating, 4:Deactivated
        if (activeConnObj.value("State").toInt(0) != 1) {
            break;
        }
        activeSsid = activeConnObj.value("Id").toString();
        break;
    }

    // if current is not activating wireless this will make m_activatingAP empty
    m_activatingAP = accessPointBySsid(activeSsid);
    AccessPointWidget *apw = accessPointWidgetByAp(m_activatingAP);

    if (activeSsid.isEmpty() || m_activatingAP.isEmpty() || !apw) {
        m_indicator->hide();
        return;
    }

    const int h = -(apw->height() - m_indicator->height()) / 2;
    m_indicator->move(apw->mapTo(this, apw->rect().topRight()) - QPoint(35, h));
    m_indicator->show();
    m_indicator->play();
}

void WirelessList::onActiveConnectionInfoChanged()
{
    if (m_device.isNull()) {
        return;
    }

    // 在这个方法中需要通过m_device->activeApSsid()的信息设置m_activeAP的值
    // m_activeAP的值应该从m_apList中拿到，但在程序第一次启动后，当后端扫描无线网的数据还没有发过来，
    // 这时m_device中的ap list为空，导致本类初始化时调用loadAPList()后m_apList也是空的，
    // 那么也就无法给m_activeAP正确的值，所以在这里使用timer等待一下后端的数据，再执行遍历m_apList给m_activeAP赋值的操作
    if (m_device->enabled() && m_device->status() == NetworkDevice::Activated
            && m_apList.size() == 0) {
        QTimer::singleShot(1000, [=]{onActiveConnectionInfoChanged();});
        return;
    }

    for (int i = 0; i < m_apList.size(); ++i) {
        if (m_apList.at(i).ssid() == m_device->activeApSsid()) {
            m_activeAP = m_apList.at(i);
            m_updateAPTimer->start();
            break;
        }
    }
}

void WirelessList::onActivateApFailed(const QString &apPath, const QString &uuid)
{
    if (m_device.isNull() && !m_clickedAPW) {
        return;
    }

    AccessPoint clickedAP = m_clickedAPW->ap();

    if (clickedAP.isEmpty()) {
        return;
    }

    if (clickedAP.path() == apPath) {
        qDebug() << "wireless connect failed and may require more configuration,"
            << "path:" << clickedAP.path() << "ssid" << clickedAP.ssid()
            << "secret:" << clickedAP.secured() << "strength" << clickedAP.strength();
        m_updateAPTimer->start();

        DDBusSender()
                .service("com.deepin.dde.ControlCenter")
                .interface("com.deepin.dde.ControlCenter")
                .path("/com/deepin/dde/ControlCenter")
                .method("ShowPage")
                .arg(QString("network"))
                .arg(apPath)
                .call();

        Q_EMIT requestSetAppletVisible(false);
    }
}

void WirelessList::onHotspotEnabledChanged(const bool enabled)
{
    // Note: the obtained hotspot info is not complete
    m_activeHotspotAP = enabled ? AccessPoint(m_device->activeHotspotInfo().value("Hotspot").toObject())
                                : AccessPoint();
    isHotposActive = enabled;
    m_updateAPTimer->start();
}

AccessPoint WirelessList::accessPointBySsid(const QString &ssid)
{
    if (ssid.isEmpty()) {
        return AccessPoint();
    }

    for (auto ap : m_apList) {
        if (ap.ssid() == ssid) {
            return ap;
        }
    }

    return AccessPoint();
}

AccessPointWidget *WirelessList::accessPointWidgetByAp(const AccessPoint ap)
{
    if (ap.isEmpty()) {
        return nullptr;
    }

    for (auto apw : m_apwList) {
        if (apw->ap().path() == ap.path()) {
            return apw;
        }
    }

    return nullptr;
}
