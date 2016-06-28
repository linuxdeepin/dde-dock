#include "imagefactory.h"

#include <QDebug>
#include <QPainter>

ImageFactory::ImageFactory(QObject *parent)
    : QObject(parent)
{

}

QPixmap ImageFactory::lighterEffect(const QPixmap pixmap, const int delta)
{
    QPixmap result(pixmap);
    QPainter painter(&result);
    painter.setCompositionMode(QPainter::CompositionMode_SourceIn);
    painter.fillRect(result.rect(), QColor::fromRgb(255, 255, 255, delta));

    return result;
}
