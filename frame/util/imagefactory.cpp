#include "imagefactory.h"

#include <QDebug>

ImageFactory::ImageFactory(QObject *parent)
    : QObject(parent)
{

}

QPixmap ImageFactory::lighter(const QPixmap pixmap, const int delta)
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
            if (qAlpha(rgba) && (qRed(rgba) || qGreen(rgba) || qBlue(rgba)))
                rgba = QColor::fromRgba(rgba).lighter(delta).rgba();
            scanLine += bytesPerPixel;
        }
    }

    return QPixmap::fromImage(image);
}
