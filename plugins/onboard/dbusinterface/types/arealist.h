// Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef AREALIST_H
#define AREALIST_H

#include <QDBusMetaType>
#include <QRect>
#include <QList>

struct MonitRect {
    int x1;
    int y1;
    int x2;
    int y2;

    bool operator ==(const MonitRect& rect);
};

typedef QList<MonitRect> AreaList;

Q_DECLARE_METATYPE(MonitRect)
Q_DECLARE_METATYPE(AreaList)

QDBusArgument &operator<<(QDBusArgument &arg, const MonitRect &rect);
const QDBusArgument &operator>>(const QDBusArgument &arg, MonitRect &rect);

void registerAreaListMetaType();

#endif // AREALIST_H
