// Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SCREENRECT_H
#define SCREENRECT_H

#include <QRect>
#include <QDBusArgument>
#include <QDebug>
#include <QDBusMetaType>

struct ScreenRect
{
public:
    ScreenRect();
    operator QRect() const;

    friend QDebug operator<<(QDebug debug, const ScreenRect &rect);
    friend const QDBusArgument &operator>>(const QDBusArgument &arg, ScreenRect &rect);
    friend QDBusArgument &operator<<(QDBusArgument &arg, const ScreenRect &rect);

private:
    qint16 x;
    qint16 y;
    quint16 w;
    quint16 h;
};

Q_DECLARE_METATYPE(ScreenRect)

void registerScreenRectMetaType();

#endif // SCREENRECT_H
