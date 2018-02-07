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

#include "refreshbutton.h"

#include <QMouseEvent>
#include <QEvent>
#include <QImageReader>

RefreshButton::RefreshButton(QWidget *parent) : QLabel(parent)
{
    setAttribute(Qt::WA_TranslucentBackground);

    m_normalPixmap = loadPixmap(":/wireless/resources/wireless/refresh_normal.svg");
    m_hoverPixmap = loadPixmap(":/wireless/resources/wireless/refresh_hover.svg");
    m_pressPixmap = loadPixmap(":/wireless/resources/wireless/refresh_press.svg");

    setPixmap(m_normalPixmap);
}

void RefreshButton::enterEvent(QEvent *event)
{
    QLabel::enterEvent(event);

    setPixmap(m_hoverPixmap);
}

void RefreshButton::leaveEvent(QEvent *event)
{
    QLabel::leaveEvent(event);

    setPixmap(m_normalPixmap);
}

void RefreshButton::mousePressEvent(QMouseEvent *event)
{
    if (event->button() == Qt::LeftButton)
        setPixmap(m_pressPixmap);
}

void RefreshButton::mouseReleaseEvent(QMouseEvent *event)
{
    if (event->button() == Qt::LeftButton)
        emit clicked();

    setPixmap(m_normalPixmap);
}

QPixmap RefreshButton::loadPixmap(const QString &file)
{
    QPixmap pixmap;

    const qreal ratio = devicePixelRatioF();

    QImageReader reader;
    reader.setFileName(file);
    if (reader.canRead()) {
        reader.setScaledSize(reader.size() * ratio);
        pixmap = QPixmap::fromImage(reader.read());
        pixmap.setDevicePixelRatio(ratio);
    }

    return pixmap;
}
