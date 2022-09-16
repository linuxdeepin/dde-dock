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

#ifndef RESOLUTION_H
#define RESOLUTION_H

#include <QDBusMetaType>

class Resolution
{
public:
    friend QDBusArgument &operator<<(QDBusArgument &arg, const Resolution &value);
    friend const QDBusArgument &operator>>(const QDBusArgument &arg, Resolution &value);

    explicit Resolution();

    bool operator!=(const Resolution &other) const;
    bool operator==(const Resolution &other) const;

    int id() const { return m_id; }
    int width() const { return m_width; }
    int height() const { return m_height; }
    double rate() const { return m_rate; }

private:
    void setId(const int id) { m_id = id; }
    void setWidth(const int w) { m_width = w; }
    void setHeight(const int h) { m_height = h; }
    void setRate(const double rate) { m_rate = rate; }

private:
    int m_id;
    int m_width;
    int m_height;
    double m_rate;
};


Q_DECLARE_METATYPE(Resolution)

void registerResolutionMetaType();

#endif // RESOLUTION_H
