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

#include "horizontalseperator.h"

#include <DApplicationHelper>

#include <QPainter>

/**
 * @brief HorizontalSeperator::HorizontalSeperator 分割线控件,高度值初始化为2个像素
 * @param parent
 */
HorizontalSeperator::HorizontalSeperator(QWidget *parent)
    : QWidget(parent)
{
    setFixedHeight(2);
    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
}

void HorizontalSeperator::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e)

    QPainter painter(this);
    QColor c = palette().color(QPalette::BrightText);
    c.setAlpha(int(0.1 * 255));
    painter.fillRect(rect(), c);
}
