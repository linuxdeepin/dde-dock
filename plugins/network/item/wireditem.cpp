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

#include "constants.h"
#include "wireditem.h"
#include "networkplugin.h"
#include "../util/imageutil.h"
#include "../widgets/tipswidget.h"
#include <DGuiApplicationHelper>

#include <QPainter>
#include <QMouseEvent>
#include <QIcon>
#include <QApplication>

using namespace dde::network;
DGUI_USE_NAMESPACE

WiredItem::WiredItem(WiredDevice *device)
    : DeviceItem(device),

      m_itemTips(new TipsWidget(this)),
      m_delayTimer(new QTimer(this))
{
    m_pageName = "wired";

    m_delayTimer->setSingleShot(true);
    m_delayTimer->setInterval(200);

    m_itemTips->setObjectName("wired-" + m_device->path());
    m_itemTips->setVisible(false);
    m_itemTips->setText(tr("Unknown"));

    connect(m_delayTimer, &QTimer::timeout, this, &WiredItem::reloadIcon);
    connect(m_device, static_cast<void (NetworkDevice::*)(NetworkDevice::DeviceStatus) const>(&NetworkDevice::statusChanged), this, &WiredItem::deviceStateChanged);
    connect(static_cast<WiredDevice *>(m_device.data()), &WiredDevice::connectionsChanged, this, &WiredItem::deviceStateChanged);
    connect(static_cast<WiredDevice *>(m_device.data()), &WiredDevice::activeWiredConnectionInfoChanged, this, &WiredItem::deviceStateChanged);
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, [ = ] {
        refreshIcon();
    });

    QTimer::singleShot(0, this, &WiredItem::refreshTips);
    QTimer::singleShot(0, this, &WiredItem::refreshIcon);
}

QWidget *WiredItem::itemTips()
{
    refreshTips();

    return m_itemTips;
}

const QString WiredItem::itemCommand() const
{
    return "dbus-send --print-reply --dest=com.deepin.dde.ControlCenter /com/deepin/dde/ControlCenter com.deepin.dde.ControlCenter.ShowModule \"string:network\"";
}

void WiredItem::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    const auto ratio = devicePixelRatioF();
    const QRectF &rf = rect();
    const QRectF &rfp = QRectF(m_icon.rect());
    const int x = rf.center().x() - rfp.center().x() / ratio;
    const int y = rf.center().y() - rfp.center().y() / ratio;
    painter.drawPixmap(x, y, m_icon);
}

void WiredItem::resizeEvent(QResizeEvent *e)
{
    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    // 保持横纵比
    if (position == Dock::Bottom || position == Dock::Top) {
        setMaximumWidth(height());
        setMaximumHeight(QWIDGETSIZE_MAX);
    } else {
        setMaximumHeight(width());
        setMaximumWidth(QWIDGETSIZE_MAX);
    }

    DeviceItem::resizeEvent(e);

    m_delayTimer->start();
}

void WiredItem::refreshIcon()
{
    m_delayTimer->start();

    refreshTips();
}

void WiredItem::reloadIcon()
{
    Q_ASSERT(sender() == m_delayTimer);

    if (m_device.isNull()) {
        return;
    }

//    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();
    const Dock::DisplayMode displayMode = Dock::DisplayMode::Efficient;
    const auto ratio = devicePixelRatioF();
//    const int iconSize = displayMode == Dock::Efficient ? 16 : std::min(width(), height()) * 0.8;
    const int iconSize = PLUGIN_ICON_MAX_SIZE;

    QString iconName = "network-";
    NetworkDevice::DeviceStatus devState = m_device->status();

    if (m_device->enabled()) {
        switch (devState) {
        case NetworkDevice::DeviceStatus::Unknow:
        case NetworkDevice::DeviceStatus::Unmanaged:
        case NetworkDevice::DeviceStatus::Unavailable: {
            iconName.append("error");
            break;
        }
        case NetworkDevice::DeviceStatus::Disconnected: {
            iconName.append("none");
            break;
        }
        case NetworkDevice::DeviceStatus::Prepare:
        case NetworkDevice::DeviceStatus::Config:
        case NetworkDevice::DeviceStatus::NeedAuth:
        case NetworkDevice::DeviceStatus::IpConfig:
        case NetworkDevice::DeviceStatus::IpCheck:
        case NetworkDevice::DeviceStatus::Secondaries: {
            m_delayTimer->start();
            const quint64 index = QDateTime::currentMSecsSinceEpoch() / 200;
            const int num = (index % 5) + 1;
            m_icon = QIcon(QString(":/wired/resources/wired/network-wired-symbolic-connecting%1.svg").arg(num))
                     .pixmap(iconSize * ratio, iconSize * ratio);
            m_icon.setDevicePixelRatio(ratio);
            update();
            return;
        }
        case NetworkDevice::DeviceStatus::Activated: {
            iconName.append("online");
            break;
        }
        case NetworkDevice::DeviceStatus::Deactivation:
        case NetworkDevice::DeviceStatus::Failed: {
            iconName.append("none");
            break;
        }
        default:;
        }
    } else {
        iconName.append("disabled");
    }

    if (devState == NetworkDevice::DeviceStatus::Activated && !NetworkPlugin::isConnectivity()) {
        iconName = "network-offline";
    }

    if (m_device->obtainIpFailed()) {
        iconName = "network-warning";
    }

    m_delayTimer->stop();

    if (displayMode == Dock::Efficient)
        iconName.append("-symbolic");

    // 最小尺寸时采用深色图标
    if (height() <= PLUGIN_BACKGROUND_MIN_SIZE && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        iconName.append(PLUGIN_MIN_ICON_NAME);

    m_icon = QIcon::fromTheme(iconName).pixmap(iconSize * ratio, iconSize * ratio);
    m_icon.setDevicePixelRatio(ratio);
    update();
}

void WiredItem::deviceStateChanged()
{
    refreshTips();

    m_delayTimer->start();
}

void WiredItem::refreshTips()
{
    if (m_device.isNull()) {
        return;
    }

    m_itemTips->setText(m_device->statusStringDetail());

    if (NetworkPlugin::isConnectivity()) {
        do {
            if (m_device->status() != NetworkDevice::DeviceStatus::Activated) {
                break;
            }
            const QJsonObject info = static_cast<WiredDevice *>(m_device.data())->activeWiredConnectionInfo();
            if (!info.contains("Ip4"))
                break;
            const QJsonObject ipv4 = info.value("Ip4").toObject();
            if (!ipv4.contains("Address"))
                break;
            m_itemTips->setText(tr("Wired connection: %1").arg(ipv4.value("Address").toString()));
        } while (false);
    }
}
