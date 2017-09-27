/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
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
#include "refreshbutton.h"

#include <QHBoxLayout>
#include <QDebug>

DWIDGET_USE_NAMESPACE

DeviceControlWidget::DeviceControlWidget(QWidget *parent)
    : QWidget(parent)
{
    m_deviceName = new QLabel;
    m_deviceName->setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Preferred);
    m_deviceName->setStyleSheet("color:white;");

    m_switchBtn = new DSwitchButton;

    m_refreshBtn = new RefreshButton;

    QHBoxLayout *infoLayout = new QHBoxLayout;
    infoLayout->addWidget(m_deviceName);
    infoLayout->addWidget(m_refreshBtn);
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
    connect(m_refreshBtn, &RefreshButton::clicked, this, &DeviceControlWidget::requestRefresh);
}

void DeviceControlWidget::setDeviceName(const QString &name)
{
    m_deviceName->setText(name);
}

void DeviceControlWidget::setDeviceEnabled(const bool enable)
{
    m_switchBtn->blockSignals(true);
    m_switchBtn->setChecked(enable);
    m_switchBtn->blockSignals(false);

    m_refreshBtn->setVisible(enable);
}

//void DeviceControlWidget::setSeperatorVisible(const bool visible)
//{
//    m_seperator->setVisible(visible);
//}
