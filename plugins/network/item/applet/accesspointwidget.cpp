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
#include "constants.h"
#include "util/statebutton.h"

#include <DGuiApplicationHelper>
#include <DApplication>

#include <QHBoxLayout>
#include <QDebug>
#include <QFontMetrics>
#include <QIcon>

using namespace Dock;
using namespace dde::network;

DWIDGET_USE_NAMESPACE
DGUI_USE_NAMESPACE

extern const QString DarkType;
extern const QString LightType;
extern void initFontColor(QWidget *widget);

AccessPointWidget::AccessPointWidget(const QJsonObject &apInfo)
    : QFrame(nullptr)
    , m_activeState(AccessPointWidget::ApState::Unknown)
    , m_ssidBtn(new SsidButton(this))
    , m_securityLabel(new QLabel)
    , m_strengthLabel(new QLabel)
    , m_ap(AccessPoint(apInfo))
    , m_stateButton(new StateButton(this))
    , m_loadingStat(new DSpinner(this))

{
    initUI();
    initConnect();
    //初始化数据
    updateApInfo(apInfo);
}

AccessPointWidget::~AccessPointWidget()
{
    Q_EMIT apChange();
}

void AccessPointWidget::initUI()
{
    m_ssidBtn->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Preferred);

    m_ssidBtn->setObjectName("Ssid");
    initFontColor(m_ssidBtn);

    bool isLight = (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType);

    m_stateButton->setFixedSize(PLUGIN_ICON_MAX_SIZE, PLUGIN_ICON_MAX_SIZE);
    m_stateButton->setType(StateButton::Check);
    m_stateButton->setVisible(false);

    m_loadingStat->setFixedSize(PLUGIN_ICON_MAX_SIZE, PLUGIN_ICON_MAX_SIZE);
    m_loadingStat->setVisible(false);
    m_loadingStat->start();
    m_loadingStat->move(QPoint(200,5));

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

    setStrengthIcon(m_ap.strength());
}

void AccessPointWidget::initConnect()
{
    //连接wifi
    connect(m_ssidBtn, &SsidButton::clicked, this, &AccessPointWidget::ssidClicked);
    //断开连接
    connect(m_stateButton, &StateButton::click, this, &AccessPointWidget::disconnectBtnClicked);

    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, [ = ] {
        setStrengthIcon(m_ap.strength());
    });
    connect(qApp, &DApplication::iconThemeChanged, this, [ = ] {
        setStrengthIcon(m_ap.strength());
    });

}

void AccessPointWidget::updateApInfo(const QJsonObject &apInfo)
{
    //刷新网络数据
    m_ap.updateApInfo(apInfo);
    //设置SSID  对长度太长的Ssid进行截断
    QString strSsid = ssid();
    QFontMetrics fontMetrics(m_ssidBtn->font());
    if(fontMetrics.width(strSsid) > m_ssidBtn->width())
    {
        strSsid = QFontMetrics(m_ssidBtn->font()).elidedText(strSsid, Qt::ElideRight, m_ssidBtn->width());
    }
    m_ssidBtn->setText(strSsid);
    //设置信号强度
    setStrengthIcon(strength());

    //设置是否是密码状态
    if (!secured()) {
        m_securityLabel->clear();
    } else if(!m_securityLabel->pixmap()) {
        m_securityLabel->setPixmap(m_securityPixmap);
    }
    Q_EMIT apChange();
}

void AccessPointWidget::setActiveState(ApState state)
{
    if (m_activeState == state)
        return;

    m_activeState = state;
    if (state == ApState::Activating) {
        m_loadingStat->show();
        m_stateButton->hide();
    } else if (state == ApState::Activated) {
        m_loadingStat->hide();
        m_stateButton->show();
    } else {
        m_loadingStat->hide();
        m_stateButton->hide();
    }
    Q_EMIT apChange();
}

bool AccessPointWidget::operator==(const AccessPointWidget *ap) const
{
    return this->ssid() == ap->ssid();
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
    //这个区间是需求文档中规定的
    if (strength > 65)
        type = "80";
    else if (strength > 55)
        type = "60";
    else if (strength > 30)
        type = "40";
    else if (strength > 5)
        type = "20";
    else
        type = "0";

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
    if (m_activeState == ApState::Activated)
        return;
    emit requestConnectAP(m_ap.path(), m_ap.uuid());
}

void AccessPointWidget::disconnectBtnClicked()
{
    setActiveState(ApState::Unknown);
    emit requestDisconnectAP(uuid());
}
