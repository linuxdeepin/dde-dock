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

#include "horizontalseperator.h"

#include <QPainter>

HorizontalSeperator::HorizontalSeperator(QWidget *parent)
    : QWidget(parent),
      m_color(255, 255, 255, 255 * 0.1)
{
    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
}

void HorizontalSeperator::setColor(const QColor color)
{
    m_color = color;
}

void HorizontalSeperator::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    painter.fillRect(rect(), m_color);
}
