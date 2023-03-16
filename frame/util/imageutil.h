// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
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
    // 加载窗口的预览图
    static QPixmap loadWindowThumb(const QString &winInfoId);                      // 加载图片，参数为windowId或者窗口的UUID

};

#endif // IMAGEUTIL_H
