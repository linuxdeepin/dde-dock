// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "imageutil.h"

#include <QIcon>
#include <QPainter>
#include <QCursor>
#include <QDebug>

#include <X11/Xcursor/Xcursor.h>

const QPixmap ImageUtil::loadSvg(const QString &iconName, const QString &localPath, const int size, const qreal ratio)
{
    QIcon icon = QIcon::fromTheme(iconName);
    int pixmapSize = QCoreApplication::testAttribute(Qt::AA_UseHighDpiPixmaps) ? size : int(size * ratio);
    if (!icon.isNull()) {
        QPixmap pixmap = icon.pixmap(pixmapSize);
        pixmap.setDevicePixelRatio(ratio);
        return pixmap;
    }

    QPixmap pixmap(pixmapSize, pixmapSize);
    QString localIcon = QString("%1%2%3").arg(localPath).arg(iconName).arg(iconName.contains(".svg") ? "" : ".svg");
    QSvgRenderer renderer(localIcon);
    pixmap.fill(Qt::transparent);

    QPainter painter;
    painter.begin(&pixmap);
    renderer.render(&painter);
    painter.end();
    pixmap.setDevicePixelRatio(ratio);

    return pixmap;
}

const QPixmap ImageUtil::loadSvg(const QString &iconName, const QSize size, const qreal ratio)
{
    QIcon icon = QIcon::fromTheme(iconName);
    if (!icon.isNull()) {
        QPixmap pixmap = icon.pixmap(size*ratio);
        pixmap.setDevicePixelRatio(ratio);
        return pixmap;
    }
    return QPixmap();
}

QCursor* ImageUtil::loadQCursorFromX11Cursor(const char* theme, const char* cursorName, int cursorSize)
{
    if (!theme || !cursorName || cursorSize <= 0)
        return nullptr;

    XcursorImages *images = XcursorLibraryLoadImages(cursorName, theme, cursorSize);
    if (!images || !(images->images[0])) {
        qWarning() << "loadCursorFalied, theme =" << theme << ", cursorName=" << cursorName;
        return nullptr;
    }
    const int imgW = images->images[0]->width;
    const int imgH = images->images[0]->height;
    QImage img((const uchar*)images->images[0]->pixels, imgW, imgH, QImage::Format_ARGB32);
    QPixmap pixmap = QPixmap::fromImage(img);
    QCursor *cursor = new QCursor(pixmap, images->images[0]->xhot, images->images[0]->yhot);
    XcursorImagesDestroy(images);
    return cursor;
}
