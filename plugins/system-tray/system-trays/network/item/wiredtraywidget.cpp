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

#include "wiredtraywidget.h"
#include "constants.h"
#include "../util/imageutil.h"
#include "../widgets/tipswidget.h"

#include <QPainter>
#include <QMouseEvent>
#include <QIcon>
#include <QApplication>

using namespace dde::network;

WiredTrayWidget::WiredTrayWidget(WiredDevice *device, QWidget *parent)
    : AbstractNetworkTrayWidget(device, parent),

      m_itemTips(new TipsWidget(this)),
      m_delayTimer(new QTimer(this))
{
    m_delayTimer->setSingleShot(false);
    m_delayTimer->setInterval(200);

    m_itemTips->setObjectName("wired-" + m_device->path());
    m_itemTips->setVisible(false);

    connect(m_delayTimer, &QTimer::timeout, this, &WiredTrayWidget::updateIcon);
    connect(m_device, static_cast<void (NetworkDevice::*)(NetworkDevice::DeviceStatus) const>(&NetworkDevice::statusChanged),
            m_delayTimer, static_cast<void (QTimer::*)()>(&QTimer::start));
}

void WiredTrayWidget::setActive(const bool active)
{

}

void WiredTrayWidget::updateIcon()
{
    QString iconName = "network-";
    NetworkDevice::DeviceStatus devState = m_device->status();
    if (devState != NetworkDevice::Activated)
    {
        if (devState < NetworkDevice::Disconnected)
            iconName.append("error");
        else
            iconName.append("offline");
    } else {
        if (devState >= NetworkDevice::Prepare && devState <= NetworkDevice::Secondaries) {
            m_delayTimer->start();
            const quint64 index = QDateTime::currentMSecsSinceEpoch() / 200;
            const int num = (index % 5) + 1;
            m_icon = QPixmap(QString(":/icons/system-trays/network/resources/wired/resources/wired/network-wired-symbolic-connecting%1.svg").arg(num));
            update();
            return;
        }

        if (devState == NetworkDevice::Activated)
            iconName.append("online");
        else
            iconName.append("idle");
    }

    m_delayTimer->stop();

    iconName.append("-symbolic");

    const auto ratio = qApp->devicePixelRatio();
    const int size = 16;
    m_icon = QIcon::fromTheme(iconName).pixmap(size * ratio, size * ratio);
    m_icon.setDevicePixelRatio(ratio);
    update();
}

const QImage WiredTrayWidget::trayImage()
{
    return m_icon.toImage();
}

QWidget *WiredTrayWidget::trayTipsWidget()
{
    m_itemTips->setText(tr("Unknown"));

    do {
        if (m_device->status() != NetworkDevice::Activated)
        {
            m_itemTips->setText(tr("No Network"));
            break;
        }

        const QJsonObject info = static_cast<WiredDevice *>(m_device)->activeConnection();
        if (!info.contains("Ip4"))
            break;
        const QJsonObject ipv4 = info.value("Ip4").toObject();
        if (!ipv4.contains("Address"))
            break;
        m_itemTips->setText(tr("Wired connection: %1").arg(ipv4.value("Address").toString()));
    } while (false);

    return m_itemTips;
}

const QString WiredTrayWidget::trayClickCommand()
{
    return "dbus-send --print-reply --dest=com.deepin.dde.ControlCenter /com/deepin/dde/ControlCenter com.deepin.dde.ControlCenter.ShowModule \"string:network\"";
}

void WiredTrayWidget::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);

    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(m_icon.rect());
    const QPointF &p = rf.center() - rfp.center() / m_icon.devicePixelRatioF();
    painter.drawPixmap(p, m_icon);
}

void WiredTrayWidget::resizeEvent(QResizeEvent *e)
{
    AbstractNetworkTrayWidget::resizeEvent(e);

    m_delayTimer->start();
}
