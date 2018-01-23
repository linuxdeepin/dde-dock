#include "hoverhighlighteffect.h"
#include "util/imagefactory.h"

#include <QPainter>
#include <QDebug>

HoverHighlightEffect::HoverHighlightEffect(QObject *parent)
    : QGraphicsEffect(parent)
{

}

void HoverHighlightEffect::draw(QPainter *painter)
{
    const QPixmap pix = sourcePixmap(Qt::DeviceCoordinates);

    if (isEnabled())
    {
        painter->drawPixmap(0, 0, ImageFactory::lighterEffect(pix));
    } else {
        painter->drawPixmap(0, 0, pix);
    }
}
