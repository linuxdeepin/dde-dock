#include <QColor>
#include <QPainter>
#include <QBitmap>

#include "highlighteffect.h"

HighlightEffect::HighlightEffect(QWidget * source, QWidget *parent) :
    QWidget(parent),
    m_source(source),
    m_lighter(150)
{
    setFixedSize(m_source->size());
}


int HighlightEffect::lighter() const
{
    return m_lighter;
}

void HighlightEffect::setLighter(int lighter)
{
    m_lighter = lighter;

    this->repaint();
}

void HighlightEffect::paintEvent(QPaintEvent *)
{
    if (m_source) {
        QPixmap pixmap = m_source->grab();
        this->pixmapLigher(&pixmap, 150);

        QPainter painter;
        painter.begin(this);

        painter.setClipRect(rect());

        painter.drawPixmap(0, 0, pixmap);

        painter.end();
    }
}

void HighlightEffect::pixmapLigher(QPixmap *pixmap, int lighter)
{
    QImage img = pixmap->toImage();  // slow

    for (int y=0; y < img.height(); y++)
    {
        for (int x = 0; x < img.width(); x++)
        {
            QRgb pix = img.pixel(x,y);
            QColor col(pix);
            col = col.lighter(lighter);
            img.setPixel(x, y, qRgba(col.red(), col.green(), col.blue(), qAlpha(pix)));
        }
    }
    pixmap->convertFromImage(img); // slow
}
