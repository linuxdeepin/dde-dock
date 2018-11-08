/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *             listenerri <listenerri@gmail.com>
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

#include "wirelesstraywidget.h"
#include "../util/imageutil.h"
#include "../widgets/tipswidget.h"

#include <QPainter>
#include <QMouseEvent>
#include <QApplication>
#include <QIcon>

using namespace dde::network;

WirelessTrayWidget::WirelessTrayWidget(WirelessDevice *device, QWidget *parent)
    : AbstractNetworkTrayWidget(device, parent),
      m_refershTimer(new QTimer(this)),
      m_wirelessApplet(new QWidget),
      m_wirelessPopup(new TipsWidget),
      m_APList(nullptr),
      m_reloadIcon(false)
{
    m_refershTimer->setSingleShot(false);
    m_refershTimer->setInterval(100);

    m_wirelessApplet->setVisible(false);
    m_wirelessPopup->setObjectName("wireless-" + m_device->path());
    m_wirelessPopup->setVisible(false);

    connect(m_refershTimer, &QTimer::timeout, [=] {
        WirelessDevice *dev = static_cast<WirelessDevice *>(m_device);
        // the status is Activated and activeApInfo is empty if the hotspot status of this wireless device is enabled
        if (m_device->status() == NetworkDevice::Activated && dev->activeApInfo().isEmpty() && !dev->hotspotEnabled()) {
            Q_EMIT queryActiveConnInfo();
            return;
        }
        updateIcon();
    });
    connect(m_device, static_cast<void (NetworkDevice::*) (const QString &statStr) const>(&NetworkDevice::statusChanged), m_refershTimer, static_cast<void (QTimer::*) ()>(&QTimer::start));
    connect(static_cast<WirelessDevice *>(m_device), &WirelessDevice::activeApInfoChanged, m_refershTimer, static_cast<void (QTimer::*) ()>(&QTimer::start));
    connect(static_cast<WirelessDevice *>(m_device), &WirelessDevice::activeConnectionChanged, m_refershTimer, static_cast<void (QTimer::*) ()>(&QTimer::start));

//    QMetaObject::invokeMethod(this, "init", Qt::QueuedConnection);
    init();
}

WirelessTrayWidget::~WirelessTrayWidget()
{
    m_APList->deleteLater();
    m_APList->controlPanel()->deleteLater();
}

void WirelessTrayWidget::setActive(const bool active)
{

}

void WirelessTrayWidget::updateIcon()
{
    QString type;

    const auto state = m_device->status();
    if (state <= NetworkDevice::Disconnected) {
        type = "disconnect";
        m_refershTimer->stop();
    } else {
        int strength = 0;
        if (state == NetworkDevice::Activated) {
            const auto &activeApInfo = static_cast<WirelessDevice *>(m_device)->activeApInfo();
            if (!activeApInfo.isEmpty()) {
                strength = activeApInfo.value("Strength").toInt();
                m_refershTimer->stop();
            }
        } else {
            strength = QTime::currentTime().msec() / 10 % 100;
            if (!m_refershTimer->isActive())
                m_refershTimer->start();
        }

        if (strength == 100)
            type = "80";
        else if (strength < 20)
            type = "0";
        else
            type = QString::number(strength / 10 & ~0x1) + "0";
    }

    const QString key = QString("wireless-%1%2")
                                .arg(type)
                                .arg("-symbolic");

    const auto ratio = qApp->devicePixelRatio();
    m_reloadIcon = true;
    m_pixmap = cachedPix(key, 16 * ratio);
    m_pixmap.setDevicePixelRatio(ratio);

    update();
}

const QImage WirelessTrayWidget::trayImage()
{
    return m_pixmap.toImage();
}

QWidget *WirelessTrayWidget::trayTipsWidget()
{
    const NetworkDevice::DeviceStatus stat = m_device->status();

    m_wirelessPopup->setText(tr("No Network"));

    if (stat == NetworkDevice::Activated)
    {
        const QJsonObject obj = static_cast<WirelessDevice *>(m_device)->activeConnectionInfo();
        if (obj.contains("Ip4"))
        {
            const QJsonObject ip4 = obj["Ip4"].toObject();
            if (ip4.contains("Address"))
            {
                m_wirelessPopup->setText(tr("Wireless Connection: %1").arg(ip4["Address"].toString()));
            }
        }
    }

    return m_wirelessPopup;
}

QWidget *WirelessTrayWidget::trayPopupApplet()
{
    return m_wirelessApplet;
}

void WirelessTrayWidget::setDeviceInfo(const int index)
{
    m_APList->setDeviceInfo(index);
}

bool WirelessTrayWidget::eventFilter(QObject *o, QEvent *e)
{
    if (o == m_APList && e->type() == QEvent::Resize)
        QMetaObject::invokeMethod(this, "adjustHeight", Qt::QueuedConnection);
    if (o == m_APList && e->type() == QEvent::Show)
        Q_EMIT requestWirelessScan();

    return false;
}

void WirelessTrayWidget::paintEvent(QPaintEvent *e)
{
    AbstractNetworkTrayWidget::paintEvent(e);

    QPainter painter(this);

    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(m_pixmap.rect());
    const QPointF &p = rf.center() - rfp.center() / m_pixmap.devicePixelRatioF();
    painter.drawPixmap(p, m_pixmap);

    if (m_reloadIcon)
        m_reloadIcon = false;
}

void WirelessTrayWidget::resizeEvent(QResizeEvent *e)
{
    AbstractNetworkTrayWidget::resizeEvent(e);

    m_icons.clear();
}

const QPixmap WirelessTrayWidget::cachedPix(const QString &key, const int size)
{
    if (m_reloadIcon || !m_icons.contains(key)) {
        m_icons.insert(key, QIcon::fromTheme(key, QIcon(":/icons/system-trays/network/resources/wireless/" + key + ".svg")).pixmap(size));
    }

    return m_icons.value(key);
}

void WirelessTrayWidget::init()
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

    connect(m_APList, &WirelessList::requestSetDeviceEnable, this, &WirelessTrayWidget::requestSetDeviceEnable);
    connect(m_APList, &WirelessList::requestActiveAP, this, &WirelessTrayWidget::requestActiveAP);
    connect(m_APList, &WirelessList::requestDeactiveAP, this, &WirelessTrayWidget::requestDeactiveAP);
    connect(m_APList, &WirelessList::requestWirelessScan, this, &WirelessTrayWidget::requestWirelessScan);

    QTimer::singleShot(0, this, [=]() {
        m_refershTimer->start();
    });
}

// called in eventFilter method
void WirelessTrayWidget::adjustHeight()
{
    m_wirelessApplet->setFixedHeight(m_APList->height() + m_APList->controlPanel()->height());
}

