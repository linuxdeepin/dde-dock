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

#include <DHiDPIHelper>
#include <QTimer>
#include <QHBoxLayout>
#include <QDebug>
#include <QEvent>

DWIDGET_USE_NAMESPACE

DeviceControlWidget::DeviceControlWidget(QWidget *parent)
    : QWidget(parent)
{
    m_deviceName = new QLabel;
    m_deviceName->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Preferred);
    m_deviceName->setStyleSheet("color:white;");

    m_switchBtn = new DSwitchButton;

    const QPixmap pixmap = DHiDPIHelper::loadNxPixmap(":/wireless/resources/wireless/refresh_normal.svg");

    m_loadingIndicator = new DLoadingIndicator;
    m_loadingIndicator->setImageSource(pixmap);
    m_loadingIndicator->setLoading(false);
    m_loadingIndicator->setSmooth(true);
    m_loadingIndicator->setAniDuration(1000);
    m_loadingIndicator->setAniEasingCurve(QEasingCurve::InOutCirc);
    m_loadingIndicator->installEventFilter(this);
    m_loadingIndicator->setFixedSize(pixmap.size() / devicePixelRatioF());

    QHBoxLayout *infoLayout = new QHBoxLayout;
    infoLayout->addWidget(m_deviceName);
    infoLayout->addWidget(m_loadingIndicator);
    infoLayout->addSpacing(10);
    infoLayout->addWidget(m_switchBtn);
    infoLayout->setSpacing(0);
    infoLayout->setContentsMargins(15, 0, 5, 0);

//    m_seperator = new HorizontalSeperator;
//    m_seperator->setFixedHeight(1);
//    m_seperator->setColor(Qt::black);

    QVBoxLayout *centralLayout = new QVBoxLayout;
    centralLayout->addStretch();
    centralLayout->addLayout(infoLayout);
    centralLayout->addStretch();
//    centralLayout->addWidget(m_seperator);
    centralLayout->setMargin(0);
    centralLayout->setSpacing(0);

    setLayout(centralLayout);
    setFixedHeight(30);

    connect(m_switchBtn, &DSwitchButton::checkedChanged, this, &DeviceControlWidget::deviceEnableChanged);
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

    QTimer::singleShot(1000, this, [=] {
        m_loadingIndicator->setLoading(false);
    });
}

//void DeviceControlWidget::setSeperatorVisible(const bool visible)
//{
//    m_seperator->setVisible(visible);
//}
