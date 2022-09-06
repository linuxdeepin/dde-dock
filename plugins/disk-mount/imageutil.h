// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef IMAGEUTIL_H
#define IMAGEUTIL_H

#include <QPixmap>
#include <QSvgRenderer>

class ImageUtil
{
public:
    static const QPixmap loadSvg(const QString &path, const int size);
};

#endif // IMAGEUTIL_H
