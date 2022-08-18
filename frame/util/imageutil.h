/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
    static QImage loadWindowThumb(const WId &windowId, int width, int height);     // x11下加载，用windowID
    static QImage loadWindowThumb(const QString &uuid, int width, int height);     // wayland下加载，用uuid
};

#endif // IMAGEUTIL_H
