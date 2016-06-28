#ifndef IMAGEFACTORY_H
#define IMAGEFACTORY_H

#include <QObject>
#include <QPixmap>
#include <QImage>

class ImageFactory : public QObject
{
    Q_OBJECT

public:
    explicit ImageFactory(QObject *parent = 0);

    static QPixmap lighterEffect(const QPixmap pixmap, const int delta = 50);
};

#endif // IMAGEFACTORY_H
