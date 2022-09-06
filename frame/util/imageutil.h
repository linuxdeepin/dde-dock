// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef IMAGEUTIL_H
#define IMAGEUTIL_H

#include <QWidget>
#include <QPixmap>
#include <QSvgRenderer>
#include <QApplication>

class QCursor;

class ImageUtil
{
public:
    static const QPixmap loadSvg(const QString &iconName, const QString &localPath, const int size, const qreal ratio);
    static const QPixmap loadSvg(const QString &iconName, const QSize size, const qreal ratio = qApp->devicePixelRatio());
    static QCursor* loadQCursorFromX11Cursor(const char* theme, const char* cursorName, int cursorSize);
};

#endif // IMAGEUTIL_H
