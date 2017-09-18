/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
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

#include "imagefactory.h"

#include <QDebug>
#include <QPainter>

ImageFactory::ImageFactory(QObject *parent)
    : QObject(parent)
{

}

QPixmap ImageFactory::lighterEffect(const QPixmap pixmap, const int delta)
{
    QImage image = pixmap.toImage();

    const int width = image.width();
    const int height = image.height();
    const int bytesPerPixel = image.bytesPerLine() / image.width();

    for (int i(0); i != height; ++i)
    {
        uchar *scanLine = image.scanLine(i);
        for (int j(0); j != width; ++j)
        {
            QRgb &rgba = *(QRgb*)scanLine;
            if (qAlpha(rgba) == 0xff && (qRed(rgba) || qGreen(rgba) || qBlue(rgba)))
                rgba = QColor::fromRgba(rgba).lighter(delta).rgba();
            scanLine += bytesPerPixel;
        }
    }

    return QPixmap::fromImage(image);
}
