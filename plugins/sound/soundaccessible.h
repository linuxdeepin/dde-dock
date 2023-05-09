// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SOUNDACCESSIBLE_H
#define SOUNDACCESSIBLE_H
#include "accessibledefine.h"
#include <qglobal.h>

QAccessibleInterface *soundAccessibleFactory(const QString &classname, QObject *object)
{
    Q_UNUSED(classname)
    Q_UNUSED(object)
    QAccessibleInterface *interface = nullptr;

    return interface;
}

#endif // SOUNDACCESSIBLE_H
