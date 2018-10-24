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
      m_refershTimer(new QTimer(this)),
      m_wirelessApplet(new QWidget),
      m_wirelessPopup(new TipsWidget),
      m_APList(nullptr)
{
    m_refershTimer->setSingleShot(false);
    m_refershTimer->setInterval(100);

    m_wirelessApplet->setVisible(false);
    m_wirelessPopup->setObjectName("wireless-" + m_device->path());
    m_wirelessPopup->setVisible(false);

    connect(m_device, static_cast<void (NetworkDevice::*) (const QString &statStr) const>(&NetworkDevice::statusChanged), m_refershTimer, static_cast<void (QTimer::*) ()>(&QTimer::start));
    connect(static_cast<WirelessDevice *>(m_device), &WirelessDevice::activeApInfoChanged, m_refershTimer, static_cast<void (QTimer::*) ()>(&QTimer::start));
    connect(static_cast<WirelessDevice *>(m_device), &WirelessDevice::activeConnectionChanged, m_refershTimer, static_cast<void (QTimer::*) ()>(&QTimer::start));
    connect(m_refershTimer, &QTimer::timeout, [=] {
        WirelessDevice *dev = static_cast<WirelessDevice *>(m_device);
        // the status is Activated and activeApInfo is empty if the hotspot status of this wireless device is enabled
        if (m_device->status() == NetworkDevice::Activated && dev->activeApInfo().isEmpty() && !dev->hotspotEnabled()) {
            Q_EMIT queryActiveConnInfo();
            return;
        }
        update();
    });

    //QMetaObject::invokeMethod(this, "init", Qt::QueuedConnection);
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

void WirelessItem::onNeedSecrets(const QString &info)
{
    m_APList->onNeedSecrets(info);
}

void WirelessItem::onNeedSecretsFinished(const QString &info0, const QString &info1)
{
    m_APList->onNeedSecretsFinished(info0, info1);
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

    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();

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
    painter.drawPixmap(rect().center() - pixmap.rect().center() / ratio, pixmap);

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
    if (state <= NetworkDevice::Disconnected)
    {
        type = "disconnect";
        m_refershTimer->stop();
    }
    else
    {
        int strength = 0;
        if (state == NetworkDevice::Activated)
        {
            const auto &activeApInfo = static_cast<WirelessDevice *>(m_device)->activeApInfo();
            if (!activeApInfo.isEmpty()) {
                strength = activeApInfo.value("Strength").toInt();
                m_refershTimer->stop();
            }
        }
        else
        {
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
                                .arg(displayMode == Dock::Fashion ? "" : "-symbolic");

    return cachedPix(key, size);
}

const QPixmap WirelessItem::backgroundPix(const int size)
{
    return cachedPix("wireless-background", size);
}

const QPixmap WirelessItem::cachedPix(const QString &key, const int size)
{
    if (m_reloadIcon || !m_icons.contains(key)) {
        m_icons.insert(key, QIcon::fromTheme(key, QIcon(":/icons/system-trays/network/resources/wireless/" + key + ".svg")).pixmap(size));
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
    connect(m_APList, &WirelessList::feedSecret, this, &WirelessItem::feedSecret);
    connect(m_APList, &WirelessList::cancelSecret, this, &WirelessItem::cancelSecret);
    connect(m_APList, &WirelessList::requestWirelessScan, this, &WirelessItem::requestWirelessScan);
    connect(m_APList, &WirelessList::queryConnectionSession, this, &WirelessItem::queryConnectionSession);
    connect(m_APList, &WirelessList::createApConfig, this, &WirelessItem::createApConfig);

    QTimer::singleShot(0, this, [=]() {
        m_refershTimer->start();
    });
}

void WirelessItem::adjustHeight()
{
    m_wirelessApplet->setFixedHeight(m_APList->height() + m_APList->controlPanel()->height());
}

void WirelessItem::refreshIcon()
{
    m_reloadIcon = true;
    update();
}

