/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
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

#include "resolution.h"

#include <QDebug>

void registerResolutionMetaType()
{
    qRegisterMetaType<Resolution>("Resolution");
    qDBusRegisterMetaType<Resolution>();
}

Resolution::Resolution()
{

}

bool Resolution::operator!=(const Resolution &other) const
{
    return m_width != other.m_width || m_height != other.m_height || m_rate != other.m_rate;
}

bool Resolution::operator==(const Resolution &other) const
{
    return !(other != *this);
}

QDBusArgument &operator<<(QDBusArgument &arg, const Resolution &value)
{
    arg.beginStructure();
    arg << quint32(value.id()) << quint16(value.width()) << quint16(value.height()) << value.rate();
    arg.endStructure();

    return arg;
}

const QDBusArgument &operator>>(const QDBusArgument &arg, Resolution &value)
{
    quint32 id;
    quint16 w, h;
    double rate;

    arg.beginStructure();
    arg >> id >> w >> h >> rate;
    arg.endStructure();

    value.setId(id);
    value.setWidth(w);
    value.setHeight(h);
    value.setRate(rate);

    return arg;
}
