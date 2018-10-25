/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     listenerri <listenerri@gmail.com>
 *
 * Maintainer: listenerri <listenerri@gmail.com>
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

#include "fashiontraywidgetwrapper.h"

#include <QPainter>
#include <QDebug>

FashionTrayWidgetWrapper::FashionTrayWidgetWrapper(AbstractTrayWidget *absTrayWidget, QWidget *parent)
    : QWidget(parent),
      m_absTrayWidget(absTrayWidget),
      m_layout(new QVBoxLayout(this)),
      m_attention(false)

{
    m_layout->setSpacing(0);
    m_layout->setMargin(0);
    m_layout->setContentsMargins(0, 0, 0, 0);

    m_layout->addWidget(m_absTrayWidget);

    setLayout(m_layout);

    connect(m_absTrayWidget, &AbstractTrayWidget::iconChanged, this, &FashionTrayWidgetWrapper::onTrayWidgetIconChanged);
    connect(m_absTrayWidget, &AbstractTrayWidget::clicked, this, &FashionTrayWidgetWrapper::onTrayWidgetClicked);
}

AbstractTrayWidget *FashionTrayWidgetWrapper::absTrayWidget() const
{
    return m_absTrayWidget;
}

void FashionTrayWidgetWrapper::paintEvent(QPaintEvent *event)
{

    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing, true);
    painter.setOpacity(0.5);

    QPainterPath path;
    path.addRoundedRect(rect(), 10, 10);
    painter.fillPath(path, QColor::fromRgb(40, 40, 40));
}

void FashionTrayWidgetWrapper::onTrayWidgetIconChanged()
{
    setAttention(true);
}

void FashionTrayWidgetWrapper::onTrayWidgetClicked()
{
    setAttention(false);
}

bool FashionTrayWidgetWrapper::attention() const
{
    return m_attention;
}

void FashionTrayWidgetWrapper::setAttention(bool attention)
{
    m_attention = attention;

    Q_EMIT attentionChanged(m_attention);
}
