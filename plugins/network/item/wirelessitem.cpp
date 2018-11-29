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
#include "../util/imageutil.h"
#include "../widgets/tipswidget.h"

#include <QPainter>
#include <QMouseEvent>
#include <QApplication>
#include <QIcon>

using namespace dde::network;

WirelessItem::WirelessItem(WirelessDevice *device)
    : DeviceItem(device),

      m_reloadIcon(false),
      m_refreshTimer(new QTimer(this)),
      m_wirelessApplet(new QWidget),
      m_wirelessTips(new TipsWidget),
      m_APList(nullptr)
{
    m_refreshTimer->setSingleShot(true);
    m_refreshTimer->setInterval(100);

    m_wirelessApplet->setVisible(false);
    m_wirelessTips->setObjectName("wireless-" + m_device->path());
    m_wirelessTips->setVisible(false);
    m_wirelessTips->setText(tr("No Network"));

    connect(m_refreshTimer, &QTimer::timeout, this, &WirelessItem::onRefreshTimeout);
    connect(m_device, static_cast<void (NetworkDevice::*) (const QString &statStr) const>(&NetworkDevice::statusChanged), this, &WirelessItem::deviceStateChanged);
    connect(static_cast<WirelessDevice *>(m_device), &WirelessDevice::activeApInfoChanged, m_refreshTimer, static_cast<void (QTimer::*) ()>(&QTimer::start));
    connect(static_cast<WirelessDevice *>(m_device), &WirelessDevice::activeConnectionChanged, m_refreshTimer, static_cast<void (QTimer::*) ()>(&QTimer::start));

    //QMetaObject::invokeMethod(this, "init", Qt::QueuedConnection);
    QMetaObject::invokeMethod(this, &WirelessItem::refreshTips, Qt::QueuedConnection);
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

QWidget *WirelessItem::itemTips()
{
    refreshTips();

    return m_wirelessTips;
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

void WirelessItem::paintEvent(QPaintEvent *e)
{
    DeviceItem::paintEvent(e);

//    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();
    const Dock::DisplayMode displayMode = Dock::DisplayMode::Efficient;

    const auto ratio = qApp->devicePixelRatio();
    const int iconSize = displayMode == Dock::Fashion ? std::min(width(), height()) * 0.8 : 16;
    QPixmap pixmap = iconPix(displayMode, iconSize * ratio);
    pixmap.setDevicePixelRatio(ratio);

    QPainter painter(this);
    if (displayMode == Dock::Fashion)
    {
        QPixmap pixmap = backgroundPix(iconSize * ratio);
        pixmap.setDevicePixelRatio(ratio);
        painter.drawPixmap(rect().center() - pixmap.rect().center() / ratio, pixmap);
    }
    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(pixmap.rect());
    painter.drawPixmap(rf.center() - rfp.center() / ratio, pixmap);

    if (m_reloadIcon)
        m_reloadIcon = false;
}

void WirelessItem::resizeEvent(QResizeEvent *e)
{
    DeviceItem::resizeEvent(e);

    m_icons.clear();
}

void WirelessItem::mousePressEvent(QMouseEvent *e)
{
    if (e->button() != Qt::RightButton)
        return e->ignore();

    const QPoint p(e->pos() - rect().center());
    if (p.manhattanLength() < std::min(width(), height()) * 0.8 * 0.5)
    {
        emit requestContextMenu();
        return;
    }

    return QWidget::mousePressEvent(e);
}

const QPixmap WirelessItem::iconPix(const Dock::DisplayMode displayMode, const int size)
{
    QString type;
    const auto state = m_device->status();

    if (m_device->enabled()) {
        // get strength in switch-case
        int strength = 0;
        switch (state) {
        case NetworkDevice::DeviceStatus::Unknow:
        case NetworkDevice::DeviceStatus::Unmanaged:
        case NetworkDevice::DeviceStatus::Unavailable:
        case NetworkDevice::DeviceStatus::Disconnected: {
            strength = 0;
            break;
        }
        case NetworkDevice::DeviceStatus::Prepare:
        case NetworkDevice::DeviceStatus::Config:
        case NetworkDevice::DeviceStatus::NeedAuth:
        case NetworkDevice::DeviceStatus::IpConfig:
        case NetworkDevice::DeviceStatus::IpCheck:
        case NetworkDevice::DeviceStatus::Secondaries: {
            strength = QTime::currentTime().msec() / 10 % 100;
            if (!m_refreshTimer->isActive()) {
                m_refreshTimer->start();
            }
            break;
        }
        case NetworkDevice::DeviceStatus::Activated: {
            const auto &activeApInfo = static_cast<WirelessDevice *>(m_device)->activeApInfo();
            if (activeApInfo.isEmpty()) {
                strength = 100;
                m_refreshTimer->start();
            } else {
                strength = activeApInfo.value("Strength").toInt();
            }
            break;
        }
        case NetworkDevice::DeviceStatus::Deactivation:
        case NetworkDevice::DeviceStatus::Failed: {
            strength = 0;
            break;
        }
        default:;
        }

        // set wireless icon by strength
        if (strength == 100) {
            type = "80";
        } else if (strength < 20) {
            type = "0";
        } else {
            type = QString::number(strength / 10 & ~0x1) + "0";
        }
    } else {
        type = "disabled";
    }

    QString key = QString("wireless-%1%2")
                                .arg(type)
                                .arg(displayMode == Dock::Fashion ? "" : "-symbolic");

    if (state == NetworkDevice::DeviceStatus::Activated && !NetworkPlugin::isConnectivity()) {
        key = "network-wireless-offline-symbolic";
    }

    if (m_device->obtainIpFailed()) {
        key = "network-wireless-warning-symbolic";
    }

    return cachedPix(key, size);
}

const QPixmap WirelessItem::backgroundPix(const int size)
{
    return cachedPix("wireless-background", size);
}

const QPixmap WirelessItem::cachedPix(const QString &key, const int size)
{
    if (m_reloadIcon || !m_icons.contains(key)) {
        m_icons.insert(key, QIcon::fromTheme(key, QIcon(":/wireless/resources/wireless/" + key + ".svg")).pixmap(size));
    }

    return m_icons.value(key);
}

void WirelessItem::init()
{
    m_APList = new WirelessList(static_cast<WirelessDevice *>(m_device));
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

    QTimer::singleShot(0, this, [=]() {
        m_refreshTimer->start();
    });
}

void WirelessItem::adjustHeight()
{
    m_wirelessApplet->setFixedHeight(m_APList->height() + m_APList->controlPanel()->height());
}

void WirelessItem::refreshIcon()
{
    m_reloadIcon = true;
    m_refreshTimer->start();

    refreshTips();
}

void WirelessItem::refreshTips()
{
    m_wirelessTips->setText(m_device->statusStringDetail());

    if (NetworkPlugin::isConnectivity()) {
        do {
            if (m_device->status() != NetworkDevice::DeviceStatus::Activated) {
                break;
            }
            const QJsonObject info = static_cast<WirelessDevice *>(m_device)->activeConnectionInfo();
            if (!info.contains("Ip4"))
                break;
            const QJsonObject ipv4 = info.value("Ip4").toObject();
            if (!ipv4.contains("Address"))
                break;
            m_wirelessTips->setText(tr("Wireless Connection: %1").arg(ipv4.value("Address").toString()));
        } while (false);
    }
}

void WirelessItem::deviceStateChanged()
{
    refreshTips();

    m_refreshTimer->start();
}

void WirelessItem::onRefreshTimeout()
{
    Q_ASSERT(sender() == m_refreshTimer);

    WirelessDevice *dev = static_cast<WirelessDevice *>(m_device);
    // the status is Activated and activeApInfo is empty if the hotspot status of this wireless device is enabled
    if (m_device->status() == NetworkDevice::Activated && dev->activeApInfo().isEmpty() && !dev->hotspotEnabled()) {
        Q_EMIT queryActiveConnInfo();
        return;
    }
    update();
}

