#ifndef HOVERHIGHLIGHTEFFECT_H
#define HOVERHIGHLIGHTEFFECT_H

#include <QGraphicsEffect>

class HoverHighlightEffect : public QGraphicsEffect
{
    Q_OBJECT

public:
    explicit HoverHighlightEffect(QObject *parent = nullptr);

    void setHighlighting(const bool highlighting) { m_highlighting = highlighting; }

protected:
    void draw(QPainter *painter);

private:
    bool m_highlighting;
};

#endif // HOVERHIGHLIGHTEFFECT_H
