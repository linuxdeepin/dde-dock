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
#include "applet/horizontalseperator.h"
#include "../widgets/tipswidget.h"
#include "util/utils.h"

#include <DGuiApplicationHelper>
#include <NetworkModel>

DGUI_USE_NAMESPACE

const int ItemHeight = 30;

WiredItem::WiredItem(WiredDevice *device, QWidget *parent)
    : DeviceItem(device, parent)
    , m_connectedName(new QLabel(this))
    , m_wiredIcon(new QLabel(this))
    , m_stateButton(new StateLabel(this))
    , m_loadingStat(new DSpinner(this))
//    , m_line(new HorizontalSeperator(this))
{
    setFixedHeight(ItemHeight);

//    m_connectedName->setVisible(false);
    auto iconPix = Utils::renderSVG(":/wired/resources/wired/network-wired-symbolic-dark.svg",
                                    QSize(PLUGIN_ICON_MAX_SIZE, PLUGIN_ICON_MAX_SIZE), devicePixelRatioF());
    m_wiredIcon->setPixmap(iconPix);
    m_wiredIcon->setVisible(false);
    iconPix = Utils::renderSVG(":/common/resources/common/list_select@2x.png",
                               QSize(PLUGIN_ICON_MAX_SIZE, PLUGIN_ICON_MAX_SIZE), devicePixelRatioF());
    m_stateButton->setSizePolicy(QSizePolicy::Preferred, QSizePolicy::Fixed);
    m_stateButton->setPixmap(iconPix);
    m_stateButton->setVisible(false);
    m_loadingStat->setFixedSize(PLUGIN_ICON_MAX_SIZE, PLUGIN_ICON_MAX_SIZE);
    m_loadingStat->setVisible(false);
//    m_line->setVisible(false);

    auto connectionLayout = new QVBoxLayout(this);
    connectionLayout->setMargin(0);
    connectionLayout->setSpacing(0);
//    connectionLayout->addWidget(m_line);
    auto itemLayout = new QHBoxLayout;
    itemLayout->setMargin(0);
    itemLayout->setSpacing(0);
    itemLayout->addSpacing(10);
    itemLayout->addWidget(m_wiredIcon);
    itemLayout->addWidget(m_connectedName);
    itemLayout->addStretch();
    itemLayout->addWidget(m_stateButton);
    itemLayout->addWidget(m_loadingStat);
    itemLayout->addSpacing(10);
    connectionLayout->addLayout(itemLayout);
    setLayout(connectionLayout);

    connect(m_device, static_cast<void (NetworkDevice::*)(NetworkDevice::DeviceStatus) const>(&NetworkDevice::statusChanged),
            this, &WiredItem::deviceStateChanged);

    connect(static_cast<WiredDevice*>(m_device.data()), &WiredDevice::activeWiredConnectionInfoChanged,
            this, &WiredItem::changedActiveWiredConnectionInfo);

    connect(m_stateButton, &StateLabel::click, this, [&]{
        auto enableState = m_device->enabled();
        emit requestSetDeviceEnable(path(), !enableState);
    });
    connect(m_stateButton, &StateLabel::enter, this , &WiredItem::buttonEnter);
    connect(m_stateButton, &StateLabel::leave, this , &WiredItem::buttonLeave);

    deviceStateChanged(m_device->status());
}

bool WiredItem::deviceActivated()
{
    switch (m_device->status()) {
    case NetworkDevice::Unknow:
    case NetworkDevice::Unmanaged:
    case NetworkDevice::Unavailable: {
        return false;
    }
    case NetworkDevice::Disconnected:
    case NetworkDevice::Deactivation:
    case NetworkDevice::Failed:
    case NetworkDevice::Prepare:
    case NetworkDevice::Config:
    case NetworkDevice::NeedAuth:
    case NetworkDevice::IpConfig:
    case NetworkDevice::IpCheck:
    case NetworkDevice::Secondaries:
    case NetworkDevice::Activated: {
        if (m_device->enabled())
            return true;
        else
            return false;
    }
    }
}

void WiredItem::setDeviceEnabled(bool enabled)
{
    emit requestSetDeviceEnable(path(), enabled);
}

WiredItem::WiredStatus WiredItem::getDeviceState()
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
    case NetworkDevice::Unavailable:   return Nocable;
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

QJsonObject WiredItem::getActiveWiredConnectionInfo()
{
    return static_cast<WiredDevice *>(m_device.data())->activeWiredConnectionInfo();
}

void WiredItem::deviceStateChanged(NetworkDevice::DeviceStatus state)
{
    emit wiredState(state);
    QPixmap iconPix;
    switch (state) {
    case NetworkDevice::Unknow:
    case NetworkDevice::Unmanaged:
    case NetworkDevice::Unavailable:
    case NetworkDevice::Disconnected:
    case NetworkDevice::Deactivation:
    case NetworkDevice::Failed: {
        m_loadingStat->stop();
        m_loadingStat->hide();
        m_loadingStat->setVisible(false);
    }
        break;
    case NetworkDevice::Prepare:
    case NetworkDevice::Config:
    case NetworkDevice::NeedAuth:
    case NetworkDevice::IpConfig:
    case NetworkDevice::IpCheck:
    case NetworkDevice::Secondaries: {
        m_stateButton->setVisible(false);
        m_loadingStat->setVisible(true);
        m_loadingStat->start();
        m_loadingStat->show();
    }
        break;
    case NetworkDevice::Activated: {
        m_loadingStat->stop();
        m_loadingStat->hide();
        m_loadingStat->setVisible(false);
        m_stateButton->setVisible(true);
    }
        break;
    }

    emit wiredStateChanged();
}

void WiredItem::changedActiveWiredConnectionInfo(const QJsonObject &connInfo)
{
    auto strTitle = connInfo.value("ConnectionName").toString();
    m_connectedName->setText(strTitle);
    QFontMetrics fontMetrics(m_connectedName->font());
    if(fontMetrics.width(strTitle) > m_connectedName->width()) {
        strTitle = QFontMetrics(m_connectedName->font()).elidedText(strTitle, Qt::ElideRight, m_connectedName->width());
    }

    if (strTitle.isEmpty())
        m_connectedName->setText(m_device->usingHwAdr());
    else
        m_connectedName->setText(strTitle);
}

void WiredItem::buttonEnter()
{
    if (m_device) {
        if (NetworkDevice::Activated == m_device->status()) {
            auto iconPix = Utils::renderSVG(":/common/resources/common/notify_close_press@2x.png",
                                            QSize(PLUGIN_ICON_MAX_SIZE, PLUGIN_ICON_MAX_SIZE), devicePixelRatioF());
            m_stateButton->setPixmap(iconPix);
        }
    }
}

void WiredItem::buttonLeave()
{
    if (m_device) {
        if (NetworkDevice::Activated == m_device->status()) {
            auto iconPix = Utils::renderSVG(":/common/resources/common/list_select@2x.png",
                                            QSize(PLUGIN_ICON_MAX_SIZE, PLUGIN_ICON_MAX_SIZE), devicePixelRatioF());
            m_stateButton->setPixmap(iconPix);
        }
    }
}
