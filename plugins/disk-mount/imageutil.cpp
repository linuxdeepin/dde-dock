#include "imageutil.h"

#include <QPainter>

const QPixmap ImageUtil::loadSvg(const QString &path, const int size)
{
    QPixmap pixmap(size, size);
    QSvgRenderer renderer(path);
    pixmap.fill(Qt::transparent);

    QPainter painter;
    painter.begin(&pixmap);
    renderer.render(&painter);
    painter.end();

    return pixmap;
}
