// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef WINDOWINFOLIST_H
#define WINDOWINFOLIST_H

#include <QDebug>
#include <QList>
#include <QDBusArgument>

class WindowInfo
{
public:
    friend QDebug operator<<(QDebug argument, const WindowInfo &info);
    friend QDBusArgument &operator<<(QDBusArgument &argument, const WindowInfo &info);
    friend const QDBusArgument &operator>>(const QDBusArgument &argument, WindowInfo &info);

    bool operator==(const WindowInfo &rhs) const;

public:
    bool attention;
    QString title;
    QString uuid;
};
typedef QMap<quint32, WindowInfo> WindowInfoMap;

#endif // WINDOWINFOLIST_H
