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

#ifndef MONITOR_H
#define MONITOR_H
#include "constants.h"

#include <QObject>
#include <QDebug>

#include <com_deepin_daemon_display_monitor.h>

using MonitorInter = com::deepin::daemon::display::Monitor;
using namespace Dock;
class Monitor : public QObject
{
    Q_OBJECT
public:
    struct DockPosition {
        // 左、上、右、下
        bool leftDock = true;
        bool topDock = true;
        bool rightDock = true;
        bool bottomDock = true;
        DockPosition(bool l = true, bool t = true, bool r = true, bool b = true)
        {
            qDebug() << leftDock << topDock << rightDock << bottomDock;
            leftDock = l;
            topDock = t;
            rightDock = r;
            bottomDock = b;
        }

        bool docked(const Position &pos)
        {
            switch (pos) {
            case Position::Top:
                return topDock;
            case Position::Bottom:
                return bottomDock;
            case Position::Left:
                return leftDock;
            case Position::Right:
                return rightDock;
            }
            Q_UNREACHABLE();
        }

        void reset()
        {
            leftDock = true;
            topDock = true;
            rightDock = true;
            bottomDock = true;
        }
    };

public:
    explicit Monitor(QObject *parent = nullptr);

    inline int x() const { return m_x; }
    inline int y() const { return m_y; }
    inline int w() const { return m_w; }
    inline int h() const { return m_h; }
    inline int left() const { return  m_x; }
    inline int right() const { return  m_x + m_w; }
    inline int top() const { return  m_y; }
    inline int bottom() const { return  m_y + m_h; }
    inline QPoint topLeft() const { return QPoint(m_x, m_y); }
    inline QPoint topRight() const { return QPoint(m_x + m_w, m_y); }
    inline QPoint bottomLeft() const { return QPoint(m_x, m_y + m_h); }
    inline QPoint bottomRight() const { return QPoint(m_x + m_w, m_y + m_h); }
    inline const QRect rect() const { return QRect(m_x, m_y, m_w, m_h); }

    inline const QString name() const { Q_ASSERT(!m_name.isEmpty()); return m_name; }
    inline const QString path() const { return m_path; }
    inline bool enable() const { return m_enable; }

    inline void setDockPosition(const DockPosition &position) { m_dockPosition = position; }
    inline DockPosition &dockPosition() { return m_dockPosition; }

Q_SIGNALS:
    void geometryChanged() const;
    void enableChanged(bool enable) const;

public Q_SLOTS:
    void setX(const int x);
    void setY(const int y);
    void setW(const int w);
    void setH(const int h);
    void setName(const QString &name);
    void setPath(const QString &path);
    void setMonitorEnable(bool enable);

private:
    int m_x;
    int m_y;
    int m_w;
    int m_h;

    QString m_name;
    QString m_path;

    bool m_enable;
    DockPosition m_dockPosition;
};

#endif // MONITOR_H
