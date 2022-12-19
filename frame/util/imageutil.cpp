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
#include <QDebug>
#include <QPainterPath>
#include <QRegion>
#include <QBitmap>
#include <QDBusInterface>
#include <QDBusReply>
#include <QFile>
#include <QDBusUnixFileDescriptor>
#include <QDir>

#include <X11/Xcursor/Xcursor.h>

#include <fcntl.h>
#include <unistd.h>
#include <iosfwd>

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
        QPixmap pixmap = icon.pixmap(QCoreApplication::testAttribute(Qt::AA_UseHighDpiPixmaps) ? size : QSize(size * ratio));
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

QPixmap ImageUtil::loadWindowThumb(const QString &winInfoId)
{
    // 在tmp下创建临时目录，用来存放缩略图
    QString thumbPath(imagePath());
    QDir dir(thumbPath);
    if (!dir.exists())
        dir.mkpath(thumbPath);

    QString fileName = QString("%1/%2").arg(thumbPath).arg(winInfoId);
    int fileId = open(fileName.toLocal8Bit().data(), O_CREAT | O_RDWR, S_IWUSR | S_IRUSR);
    if (fileId < 0) {
        //打开文件失败
        return QPixmap();
    }

    QDBusInterface interface(QStringLiteral("org.kde.KWin"), QStringLiteral("/org/kde/KWin/ScreenShot2"), QStringLiteral("org.kde.KWin.ScreenShot2"));
    // 第一个参数，winID或者UUID
    QList<QVariant> args;
    args << QVariant::fromValue(winInfoId);
    // 第二个参数，需要截图的选项
    QVariantMap option;
    option["include-decoration"] = true;
    option["include-cursor"] = false;
    option["native-resolution"] = true;
    args << QVariant::fromValue(option);
    // 第三个参数，文件描述符
    args << QVariant::fromValue(QDBusUnixFileDescriptor(fileId));

    QDBusReply<QVariantMap> reply = interface.callWithArgumentList(QDBus::Block, QStringLiteral("CaptureWindow"), args);
    if(!reply.isValid()) {
        close(fileId);
        qDebug() << "get current workspace background error: "<< reply.error().message();
        return QPixmap();
    }

    QVariantMap imageInfo = reply.value();
    int imageWidth = imageInfo.value("width").toUInt();
    int imageHeight = imageInfo.value("height").toUInt();
    int imageStride = imageInfo.value("stride").toUInt();
    int imageFormat = imageInfo.value("format").toUInt();

    QFile file;
    if (!file.open(fileId, QIODevice::ReadOnly)) {
        close(fileId);
        return QPixmap();
    }

    if (file.size() == 0) {
        file.close();
        return QPixmap();
    }

    QByteArray fileContent = file.readAll();
    QImage image(reinterpret_cast<uchar *>(fileContent.data()), imageWidth, imageHeight, imageStride, static_cast<QImage::Format>(imageFormat));
    QPixmap pixmap = QPixmap::fromImage(image);
    close(fileId);
    return pixmap;
}

QString ImageUtil::imagePath()
{
    return QString("%1/dde-dock/windowthumb").arg(QDir::tempPath());
}
