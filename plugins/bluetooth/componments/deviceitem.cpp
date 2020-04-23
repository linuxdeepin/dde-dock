/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     zhaolong <zhaolong@uniontech.com>
 *
 * Maintainer: zhaolong <zhaolong@uniontech.com>
 *
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

#include "deviceitem.h"

#include <DStyle>

#include <QHBoxLayout>
#include <QPainter>

extern const int ItemHeight = 30;

DeviceItem::DeviceItem(const QString &title, QWidget *parent)
    : QWidget(parent)
    , m_title(new QLabel(this))
    , m_state(new QLabel(this))
    , m_loadingStat(new DSpinner)
    , m_line(new HorizontalSeparator(this))
{
    setFixedHeight(ItemHeight);
//    m_state->setPixmap(QPixmap(":/list_select@2x.png"));
    m_state->setPixmap(QPixmap(":/select_dark.svg"));

    auto strTitle = title;
    m_title->setText(strTitle);
    QFontMetrics fontMetrics(m_title->font());
    if(fontMetrics.width(strTitle) > m_title->width())
    {
        strTitle = QFontMetrics(m_title->font()).elidedText(strTitle, Qt::ElideRight, m_title->width());
    }
    m_title->setText(strTitle);

    m_line->setVisible(true);
    m_state->setVisible(false);
    m_loadingStat->setFixedSize(20, 20);
    m_loadingStat->setVisible(false);

    auto deviceLayout = new QVBoxLayout(this);
    deviceLayout->setMargin(0);
    deviceLayout->setSpacing(0);
    deviceLayout->addWidget(m_line);
    auto itemLayout = new QHBoxLayout(this);
    itemLayout->setMargin(0);
    itemLayout->setSpacing(0);
    itemLayout->addSpacing(12);
    itemLayout->addWidget(m_title);
    itemLayout->addStretch();
    itemLayout->addWidget(m_state);
    itemLayout->addWidget(m_loadingStat);
    itemLayout->addSpacing(12);
    deviceLayout->addLayout(itemLayout);
    setLayout(deviceLayout);
}

void DeviceItem::mousePressEvent(QMouseEvent *event)
{
    emit clicked(m_device);
    QWidget::mousePressEvent(event);
}

void DeviceItem::enterEvent(QEvent *event)
{
    QWidget::enterEvent(event);
    if (m_device) {
        if (Device::StateConnected == m_device->state()) {
//            m_state->setPixmap(QPixmap(":/notify_close_press@2x.png"));
            m_state->setPixmap(QPixmap(":/disconnect_dark.svg"));
        }
    }
}

void DeviceItem::leaveEvent(QEvent *event)
{
    QWidget::enterEvent(event);
    if (m_device) {
        if (Device::StateConnected == m_device->state()) {
//            m_state->setPixmap(QPixmap(":/list_select@2x.png"));
            m_state->setPixmap(QPixmap(":/select_dark.svg"));
        }
    }
}

void DeviceItem::changeState(const Device::State state)
{
    switch (state) {
    case Device::StateUnavailable: {
        m_state->setVisible(false);
        m_loadingStat->stop();
        m_loadingStat->hide();
        m_loadingStat->setVisible(false);
    }
        break;
    case Device::StateAvailable: {
        m_state->setVisible(false);
        m_loadingStat->start();
        m_loadingStat->show();
        m_loadingStat->setVisible(true);
    }
        break;
    case Device::StateConnected: {
        m_loadingStat->stop();
        m_loadingStat->hide();
        m_loadingStat->setVisible(false);
        m_state->setVisible(true);
    }
        break;
    }
}

HorizontalSeparator::HorizontalSeparator(QWidget *parent)
    : QWidget(parent)
{
    setFixedHeight(1);
    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
}

void HorizontalSeparator::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    painter.fillRect(rect(), QColor(0, 0, 0, 0));
}

MenueItem::MenueItem(QWidget *parent)
    : QLabel(parent)
{
}

void MenueItem::mousePressEvent(QMouseEvent *event)
{
    QLabel::mousePressEvent(event);
    emit clicked();
}
