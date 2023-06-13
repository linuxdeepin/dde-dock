// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DOCKRECT_H
#define DOCKRECT_H

#include <QRect>
#include <QDBusMetaType>

struct DockRect
{
public:
    DockRect();
    operator QRect() const;

    friend QDebug operator<<(QDebug debug, const DockRect &rect);
    friend const QDBusArgument &operator>>(const QDBusArgument &arg, DockRect &rect);
    friend QDBusArgument &operator<<(QDBusArgument &arg, const DockRect &rect);

    int x;
    int y;
    uint w;
    uint h;
};

Q_DECLARE_METATYPE(DockRect)

void registerDockRectMetaType();

#endif // DOCKRECT_H
