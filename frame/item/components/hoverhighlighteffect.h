#ifndef HOVERHIGHLIGHTEFFECT_H
#define HOVERHIGHLIGHTEFFECT_H

#include <QGraphicsEffect>

class HoverHighlightEffect : public QGraphicsEffect
{
    Q_OBJECT

public:
    explicit HoverHighlightEffect(QObject *parent = nullptr);

protected:
    void draw(QPainter *painter);
};

#endif // HOVERHIGHLIGHTEFFECT_H
