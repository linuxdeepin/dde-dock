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

#include "accesspointwidget.h"
#include "horizontalseperator.h"

#include <QHBoxLayout>
#include <QDebug>

#include <DSvgRenderer>

using namespace dde::network;

DWIDGET_USE_NAMESPACE

AccessPointWidget::AccessPointWidget()
    : QFrame(nullptr),

      m_activeState(NetworkDevice::Unknow),
      m_ssidBtn(new QPushButton(this)),
      m_disconnectBtn(new DImageButton(this)),
      m_securityIcon(new QLabel),
      m_strengthIcon(new QLabel)
{
    const auto ratio = devicePixelRatioF();
    m_ssidBtn->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Preferred);

    m_ssidBtn->setObjectName("Ssid");

    m_disconnectBtn->setVisible(false);
    m_disconnectBtn->setNormalPic(":/wireless/resources/wireless/select.svg");
    m_disconnectBtn->setHoverPic(":/wireless/resources/wireless/disconnect_hover.svg");
    m_disconnectBtn->setPressPic(":/wireless/resources/wireless/disconnect_press.svg");

    QPixmap iconPix = DSvgRenderer::render(":/wireless/resources/wireless/security.svg", QSize(16, 16) * ratio);
    iconPix.setDevicePixelRatio(ratio);
    m_securityIconSize = iconPix.size();
    m_securityIcon->setPixmap(iconPix);

    QHBoxLayout *infoLayout = new QHBoxLayout;
    infoLayout->addWidget(m_securityIcon);
    infoLayout->addSpacing(5);
    infoLayout->addWidget(m_strengthIcon);
    infoLayout->addSpacing(10);
    infoLayout->addWidget(m_ssidBtn);
    infoLayout->addWidget(m_disconnectBtn);
    infoLayout->addSpacing(20);
    infoLayout->setSpacing(0);
    infoLayout->setContentsMargins(15, 0, 0, 0);

    QVBoxLayout *centralLayout = new QVBoxLayout;
    centralLayout->addLayout(infoLayout);
    centralLayout->setSpacing(0);
    centralLayout->setMargin(0);

    setLayout(centralLayout);
    setStyleSheet("AccessPointWidget #Ssid {"
                  "color:white;"
                  "background-color:transparent;"
                  "border:none;"
                  "text-align:left;"
                  "}"
                  "AccessPointWidget {"
                  "border-radius:4px;"
                  "margin:0 2px;"
                  "border-top:1px solid rgba(255, 255, 255, .05);"
                  "}"
                  "AccessPointWidget:hover {"
                  "border:none;"
                  "margin:0;"
                  "background-color:rgba(255, 255, 255, .1);"
                  "}"
                  "AccessPointWidget[active=true] #Ssid {"
//                  "color:#2ca7f8;"
                  "}");

    connect(m_ssidBtn, &QPushButton::clicked, this, &AccessPointWidget::clicked);
    connect(m_ssidBtn, &QPushButton::clicked, this, &AccessPointWidget::ssidClicked);
    connect(m_disconnectBtn, &DImageButton::clicked, this, &AccessPointWidget::disconnectBtnClicked);
}

void AccessPointWidget::updateAP(const AccessPoint &ap)
{
    m_ap = ap;

    m_ssidBtn->setText(ap.ssid());

    setStrengthIcon(ap.strength());

    if (!ap.secured()) {
        m_securityIcon->clear();
        m_securityIcon->setFixedSize(m_securityIconSize);
    }

    // reset state
    setActiveState(NetworkDevice::Unknow);
}

bool AccessPointWidget::active() const
{
    return m_activeState == NetworkDevice::Activated;
}

void AccessPointWidget::setActiveState(const NetworkDevice::DeviceStatus state)
{
    if (m_activeState == state)
        return;

    m_activeState = state;
    setStyleSheet(styleSheet());

    const bool isActive = active();
    m_disconnectBtn->setVisible(isActive);
}

void AccessPointWidget::enterEvent(QEvent *e)
{
    QWidget::enterEvent(e);
    m_disconnectBtn->setNormalPic(":/wireless/resources/wireless/disconnect.svg");
}

void AccessPointWidget::leaveEvent(QEvent *e)
{
    QWidget::leaveEvent(e);
    m_disconnectBtn->setNormalPic(":/wireless/resources/wireless/select.svg");
}

void AccessPointWidget::setStrengthIcon(const int strength)
{
    QPixmap iconPix;
    const auto ratio = devicePixelRatioF();
    const QSize s = QSize(16, 16) * ratio;

    QString type;
    if (strength == 100)
        type = "80";
    else if (strength < 20)
        type = "0";
    else
        type = QString::number(strength / 10 & ~0x1) + "0";

    iconPix = DSvgRenderer::render(QString(":/wireless/resources/wireless/wireless-%1-symbolic.svg").arg(type), s);
    iconPix.setDevicePixelRatio(ratio);

    m_strengthIcon->setPixmap(iconPix);
}

void AccessPointWidget::ssidClicked()
{
    if (m_activeState == NetworkDevice::Activated)
        return;

    setActiveState(NetworkDevice::Prepare);
    emit requestActiveAP(m_ap.path(), m_ap.ssid());
}

void AccessPointWidget::disconnectBtnClicked()
{
    setActiveState(NetworkDevice::Unknow);
    emit requestDeactiveAP(m_ap);
}
