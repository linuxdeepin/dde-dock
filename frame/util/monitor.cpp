/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *             kirigaya <kirigaya@mkacg.com>
 *             Hualet <mr.asianwang@gmail.com>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *             kirigaya <kirigaya@mkacg.com>
 *             Hualet <mr.asianwang@gmail.com>
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

#include "monitor.h"

Monitor::Monitor(QObject *parent)
    : QObject(parent)
    , m_x(0)
    , m_y(0)
    , m_w(0)
    , m_h(0)
    , m_enable(false)
{

}

void Monitor::setX(const int x)
{
    if (m_x == x)
        return;

    m_x = x;

    Q_EMIT geometryChanged();
}

void Monitor::setY(const int y)
{
    if (m_y == y)
        return;

    m_y = y;

    Q_EMIT geometryChanged();
}

void Monitor::setW(const int w)
{
    if (m_w == w)
        return;

    m_w = w;

    Q_EMIT geometryChanged();
}

void Monitor::setH(const int h)
{
    if (m_h == h)
        return;

    m_h = h;

    Q_EMIT geometryChanged();
}

void Monitor::setName(const QString &name)
{
    qDebug() << "screen name change from :" << m_name << " to: " << name;
    m_name = name;
}

void Monitor::setPath(const QString &path)
{
    m_path = path;
}

bool compareResolution(const Resolution &first, const Resolution &second)
{
    long firstSum = long(first.width()) * first.height();
    long secondSum = long(second.width()) * second.height();
    if (firstSum > secondSum)
        return true;
    else if (firstSum == secondSum) {
        if (first.rate() - second.rate() > 0.000001)
            return true;
        else
            return false;
    } else
        return false;

}

void Monitor::setMonitorEnable(bool enable)
{
    if (m_enable == enable)
        return;

    m_enable = enable;
    Q_EMIT enableChanged(enable);
}
