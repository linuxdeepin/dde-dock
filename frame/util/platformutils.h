// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef PLATFORMUTILS_H
#define PLATFORMUTILS_H

#include <QObject>

class PlatformUtils
{
public:
    static QString getAppNameForWindow(quint32 winId);

private:
    static QString getWindowProperty(quint32 winId, QString propName);
};

#endif // PLATFORMUTILS_H
