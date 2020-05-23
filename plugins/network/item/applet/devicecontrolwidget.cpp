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

#include "devicecontrolwidget.h"
#include "horizontalseperator.h"
#include "../widgets/tipswidget.h"

#include <DHiDPIHelper>
#include <DGuiApplicationHelper>

#include <QTimer>
#include <QHBoxLayout>
#include <QDebug>
#include <QEvent>
<<<<<<< HEAD
#include <DGuiApplicationHelper>
#include <QDBusConnection>
=======
#include <QLabel>
>>>>>>> 7cac6001... feat:format code

DWIDGET_USE_NAMESPACE
DGUI_USE_NAMESPACE

extern const int ItemHeight = 30;
extern void initFontColor(QWidget *widget);

DeviceControlWidget::DeviceControlWidget(QWidget *parent)
    : QWidget(parent)
    , m_airplaninter(new AirplanInter("com.deepin.daemon.AirplaneMode","/com/deepin/daemon/AirplaneMode",QDBusConnection::systemBus(),this))
{
    m_deviceName = new QLabel(this);
    m_deviceName->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Preferred);
    initFontColor(m_deviceName);

    m_switchBtn = new DSwitchButton;

    const QPixmap pixmap = DHiDPIHelper::loadNxPixmap(":/wireless/resources/wireless/refresh.svg");

    m_loadingIndicator = new DLoadingIndicator;
    m_loadingIndicator->setLoading(false);
    m_loadingIndicator->setSmooth(true);
    m_loadingIndicator->setAniDuration(1000);
    m_loadingIndicator->setAniEasingCurve(QEasingCurve::InOutCirc);
    m_loadingIndicator->installEventFilter(this);
    m_loadingIndicator->setFixedSize(pixmap.size() / devicePixelRatioF());
    m_loadingIndicator->viewport()->setAutoFillBackground(false);
    m_loadingIndicator->setFrameShape(QFrame::NoFrame);
    refreshIcon();

    QHBoxLayout *infoLayout = new QHBoxLayout;
    infoLayout->setMargin(0);
    infoLayout->setSpacing(0);
    infoLayout->addSpacing(3);
    infoLayout->addWidget(m_deviceName);
    infoLayout->addStretch();
    infoLayout->addWidget(m_loadingIndicator);
    infoLayout->addSpacing(10);
    infoLayout->addWidget(m_switchBtn);
    infoLayout->addSpacing(3);

    QVBoxLayout *centralLayout = new QVBoxLayout;
    centralLayout->addStretch();
    centralLayout->addLayout(infoLayout);
    centralLayout->addStretch();
    centralLayout->setMargin(0);
    centralLayout->setSpacing(0);

    setLayout(centralLayout);
    setFixedHeight(ItemHeight);

    connect(m_switchBtn, &DSwitchButton::clicked, this, &DeviceControlWidget::enableButtonToggled);
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &DeviceControlWidget::refreshIcon);
}

void DeviceControlWidget::setDeviceName(const QString &name)
{
    m_deviceName->setText(name);
}

void DeviceControlWidget::setDeviceEnabled(const bool enable)
{
    m_switchBtn->blockSignals(true);
    m_switchBtn->setChecked(enable);
    m_loadingIndicator->setVisible(enable);
    m_switchBtn->blockSignals(false);
}

bool DeviceControlWidget::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == m_loadingIndicator) {
        if (event->type() == QEvent::MouseButtonPress) {
            if (!m_loadingIndicator->loading()) {
                refreshNetwork();
            }
        }
    }

    return QWidget::eventFilter(watched, event);
}

void DeviceControlWidget::refreshNetwork()
{
    emit requestRefresh();

    m_loadingIndicator->setLoading(true);

    QTimer::singleShot(1000, this, [ = ] {
        m_loadingIndicator->setLoading(false);
    });
}

void DeviceControlWidget::refreshIcon()
{
    QPixmap pixmap;
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        pixmap = DHiDPIHelper::loadNxPixmap(":/wireless/resources/wireless/refresh_dark.svg");
    else
        pixmap = DHiDPIHelper::loadNxPixmap(":/wireless/resources/wireless/refresh.svg");

    m_loadingIndicator->setImageSource(pixmap);
}
