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

#include "imageutil.h"

#include <QIcon>
#include <QPainter>
#include <QCursor>
#include <QGSettings>
#include <QDebug>
#include <QPainterPath>
#include <QRegion>
#include <QBitmap>

#include <X11/Xcursor/Xcursor.h>

const QPixmap ImageUtil::loadSvg(const QString &iconName, const QString &localPath, const int size, const qreal ratio)
{
    QIcon icon = QIcon::fromTheme(iconName);
    if (!icon.isNull()) {
        QPixmap pixmap = icon.pixmap(int(size * ratio), int(size * ratio));
        pixmap.setDevicePixelRatio(ratio);
        return pixmap;
    }

    QPixmap pixmap(int(size * ratio), int(size * ratio));
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

const QPixmap ImageUtil::getShadowPixmap(const QPixmap &pixmap, const QColor &shadowColor, const QSize &backSize)
{
    // 新建一个bitmap作为遮罩，尺寸为外部指定(backSize)
    QBitmap bitmap(backSize);
    QPainter painter(&bitmap);
    painter.setRenderHints(QPainter::Antialiasing | QPainter::SmoothPixmapTransform);
    painter.setPen(Qt::NoPen);
    // 往bitmap上绘制圆形区域
    QRect rectBackground(0, 0, backSize.width(), backSize.height());
    painter.fillRect(rectBackground, Qt::white);
    painter.setBrush(Qt::black);
    painter.drawEllipse(0, 0, backSize.width(), backSize.height());

    // 新建QPixmap，将上述新建的Bitmap作为它的遮罩
    QPixmap newPixmap(backSize);
    newPixmap.setMask(bitmap);
    QPainter pixPainter(&newPixmap);
    pixPainter.setRenderHints(QPainter::Antialiasing | QPainter::SmoothPixmapTransform);
    const QSize iconSize = pixmap.size();
    int startX = (backSize.width() - iconSize.width()) / 2;
    int startY = (backSize.height() - iconSize.height()) / 2;
    // 将传入参数的背景尺寸(backSize)作为背景，通过画笔路径在上面绘制圆形区域
    QRect rectPixmap(startX, startY, iconSize.width(), iconSize.height());
    QPainterPath path;
    path.addEllipse(rectBackground);
    pixPainter.fillPath(path, shadowColor);
    // 将图像绘制到新建的QPixmap的背景区域，并返回新的QPixmap
    pixPainter.drawPixmap(rectPixmap, pixmap);
    pixPainter.setPen(shadowColor);
    pixPainter.drawEllipse(rectBackground);
    return newPixmap;
}
