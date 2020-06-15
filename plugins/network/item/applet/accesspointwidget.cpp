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
#include "util/utils.h"
#include "../frame/util/imageutil.h"
#include "../wireditem.h"
#include "constants.h"
#include "util/statebutton.h"

#include <DGuiApplicationHelper>
#include <DApplication>

#include <QHBoxLayout>
#include <QDebug>
#include <QFontMetrics>
#include <QIcon>

using namespace dde::network;

DWIDGET_USE_NAMESPACE
DGUI_USE_NAMESPACE

extern const QString DarkType;
extern const QString LightType;
extern void initFontColor(QWidget *widget);

AccessPointWidget::AccessPointWidget()
    : QFrame(nullptr)
    , m_activeState(NetworkDevice::Unknow)
    , m_ssidBtn(new SsidButton(this))
    , m_securityLabel(new QLabel)
    , m_strengthLabel(new QLabel)
    , m_stateButton(new StateButton(this))
{
    m_ssidBtn->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Preferred);

    m_ssidBtn->setObjectName("Ssid");
    initFontColor(m_ssidBtn);

    bool isLight = (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType);

    m_stateButton->setFixedSize(PLUGIN_ICON_MAX_SIZE, PLUGIN_ICON_MAX_SIZE);
    m_stateButton->setType(StateButton::Check);
    m_stateButton->setVisible(false);

    auto pixpath = QString(":/wireless/resources/wireless/security");
    pixpath = isLight ? pixpath + DarkType : pixpath + LightType;
    m_securityPixmap = Utils::renderSVG(pixpath, QSize(16, 16), devicePixelRatioF());
    m_securityIconSize = m_securityPixmap.size();
    m_securityLabel->setPixmap(m_securityPixmap);
    m_securityLabel->setFixedSize(m_securityIconSize / devicePixelRatioF());

    QHBoxLayout *infoLayout = new QHBoxLayout;
    infoLayout->addWidget(m_securityLabel);
    infoLayout->setMargin(0);
    infoLayout->setSpacing(0);
    infoLayout->addSpacing(2);
    infoLayout->addWidget(m_strengthLabel);
    infoLayout->addSpacing(10);
    infoLayout->addWidget(m_ssidBtn);
    infoLayout->addWidget(m_stateButton);
    infoLayout->addSpacing(3);
    infoLayout->setSpacing(0);

    QVBoxLayout *centralLayout = new QVBoxLayout;
    centralLayout->addLayout(infoLayout);
    centralLayout->setSpacing(0);
    centralLayout->setMargin(0);

    setLayout(centralLayout);

    connect(m_ssidBtn, &SsidButton::clicked, this, &AccessPointWidget::clicked);
    connect(m_ssidBtn, &SsidButton::clicked, this, &AccessPointWidget::ssidClicked);
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, [ = ] {
        setStrengthIcon(m_ap.strength());
    });

    connect(qApp, &DApplication::iconThemeChanged, this, [ = ] {
        setStrengthIcon(m_ap.strength());
    });

    connect(m_stateButton, &StateButton::click, this, &AccessPointWidget::disconnectBtnClicked);

    setStrengthIcon(m_ap.strength());
}

void AccessPointWidget::updateAP(const AccessPoint &ap)
{
    m_ap = ap;

    QString strSsid = ap.ssid();
    m_ssidBtn->setText(strSsid);

    QFontMetrics fontMetrics(m_ssidBtn->font());
    if(fontMetrics.width(strSsid) > m_ssidBtn->width())
    {
        strSsid = QFontMetrics(m_ssidBtn->font()).elidedText(strSsid, Qt::ElideRight, m_ssidBtn->width());
    }
    m_ssidBtn->setText(strSsid);

    setStrengthIcon(ap.strength());

    if (!ap.secured()) {
        m_securityLabel->clear();
    } else if(!m_securityLabel->pixmap()) {
        m_securityLabel->setPixmap(m_securityPixmap);
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

    const bool isActive = active();

    m_stateButton->setVisible(isActive);
}

void AccessPointWidget::enterEvent(QEvent *e)
{
    QWidget::enterEvent(e);}

void AccessPointWidget::leaveEvent(QEvent *e)
{
    QWidget::leaveEvent(e);}

void AccessPointWidget::setStrengthIcon(const int strength)
{
    QPixmap iconPix;
    const QSize s = QSize(16, 16);

    QString type;
    if (strength == 100)
        type = "80";
    else if (strength < 20)
        type = "0";
    else
        type = QString::number(strength / 10 & ~0x1) + "0";

    QString iconString = QString("wireless-%1-symbolic").arg(type);
    bool isLight = (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType);

    if (isLight) {
        iconString.append("-dark");
    }

    const auto ratio = devicePixelRatioF();
    iconPix = ImageUtil::loadSvg(iconString, ":/wireless/resources/wireless/", s.width(), ratio);

    m_strengthLabel->setPixmap(iconPix);

    m_securityPixmap = QIcon::fromTheme(isLight ? ":/wireless/resources/wireless/security_dark.svg" : ":/wireless/resources/wireless/security.svg").pixmap(s * devicePixelRatioF());
    m_securityPixmap.setDevicePixelRatio(devicePixelRatioF());
    m_securityLabel->setPixmap(m_securityPixmap);
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
