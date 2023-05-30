// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SOUNDACCESSIBLE_H
#define SOUNDACCESSIBLE_H
#include <QAccessibleInterface>

QAccessibleInterface *soundAccessibleFactory(const QString &classname, QObject *object)
{
    QAccessibleInterface *interface = nullptr;

    return interface;
}

#endif // SOUNDACCESSIBLE_H
