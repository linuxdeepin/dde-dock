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

#include "refreshbutton.h"

#include <QMouseEvent>
#include <QEvent>

RefreshButton::RefreshButton(QWidget *parent) : QLabel(parent)
{
    setAttribute(Qt::WA_TranslucentBackground);

    setPixmap(QPixmap(":/wireless/resources/wireless/refresh_normal.svg"));
}

void RefreshButton::enterEvent(QEvent *event)
{
    QLabel::enterEvent(event);

    setPixmap(QPixmap(":/wireless/resources/wireless/refresh_hover.svg"));
}

void RefreshButton::leaveEvent(QEvent *event)
{
    QLabel::leaveEvent(event);

    setPixmap(QPixmap(":/wireless/resources/wireless/refresh_normal.svg"));
}

void RefreshButton::mousePressEvent(QMouseEvent *event)
{
    if (event->button() == Qt::LeftButton)
        setPixmap(QPixmap(":/wireless/resources/wireless/refresh_press.svg"));
}

void RefreshButton::mouseReleaseEvent(QMouseEvent *event)
{
    if (event->button() == Qt::LeftButton)
        emit clicked();

    setPixmap(QPixmap(":/wireless/resources/wireless/refresh_normal.svg"));
}
