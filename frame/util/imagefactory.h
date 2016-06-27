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

    static QPixmap lighter(const QPixmap pixmap, const int delta = 120);
};

#endif // IMAGEFACTORY_H
