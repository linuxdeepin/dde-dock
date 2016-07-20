#ifndef IMAGEUTIL_H
#define IMAGEUTIL_H

#include <QPixmap>
#include <QSvgRenderer>

class ImageUtil
{
public:
    static const QPixmap loadSvg(const QString &path, const int size);
};

#endif // IMAGEUTIL_H
