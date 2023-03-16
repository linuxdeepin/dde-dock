// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
        if (ratio == 1)
            return pixmap;
        return pixmap.scaled(size * ratio, size * ratio);
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

    if (ratio == 1)
        return pixmap;

    return pixmap.scaled(size * ratio, size * ratio);
}

const QPixmap ImageUtil::loadSvg(const QString &iconName, const QSize size, const qreal ratio)
{
    QIcon icon = QIcon::fromTheme(iconName);
    if (!icon.isNull()) {
        QPixmap pixmap = icon.pixmap(QCoreApplication::testAttribute(Qt::AA_UseHighDpiPixmaps) ? size : QSize(size * ratio));
        pixmap.setDevicePixelRatio(ratio);
        if (ratio == 1)
            return pixmap;

        if (pixmap.size().width() > size.width() * ratio)
            pixmap = pixmap.scaledToWidth(size.width() * ratio);
        if (pixmap.size().height() > size.height() * ratio)
            pixmap = pixmap.scaledToHeight(size.height() * ratio);

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

    // pipe read write fd
    int fd[2];

    if (pipe(fd) < 0) {
        qDebug() << "failed to create pipe";
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
    args << QVariant::fromValue(QDBusUnixFileDescriptor(fd[1]));

    QDBusReply<QVariantMap> reply = interface.callWithArgumentList(QDBus::Block, QStringLiteral("CaptureWindow"), args);
    if(!reply.isValid()) {
        close(fd[1]);
        close(fd[0]);
        qDebug() << "get current workspace background error: "<< reply.error().message();
        return QPixmap();
    }

    // close write
    close(fd[1]);

    QVariantMap imageInfo = reply.value();
    int imageWidth = imageInfo.value("width").toUInt();
    int imageHeight = imageInfo.value("height").toUInt();
    int imageStride = imageInfo.value("stride").toUInt();
    int imageFormat = imageInfo.value("format").toUInt();

    QFile file;
    if (!file.open(fd[0], QIODevice::ReadOnly)) {
        file.close();
        close(fd[0]);
        return QPixmap();
    }

    QImage::Format qimageFormat = static_cast<QImage::Format>(imageFormat);
    int bitsCountPerPixel = QImage::toPixelFormat(qimageFormat).bitsPerPixel();

    QByteArray fileContent = file.read(imageHeight * imageWidth * bitsCountPerPixel / 8);
    QImage image(reinterpret_cast<uchar *>(fileContent.data()), imageWidth, imageHeight, imageStride, qimageFormat);
    QPixmap pixmap = QPixmap::fromImage(image);

    // close read
    close(fd[0]);

    return pixmap;
}
