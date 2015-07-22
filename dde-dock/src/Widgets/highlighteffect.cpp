#include <QColor>
#include <QPainter>
#include <QBitmap>

#include "highlighteffect.h"

HighlightEffect::HighlightEffect(QWidget * source, QWidget *parent) :
    QWidget(parent),
    m_source(source)
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
}

int HighlightEffect::darker() const
{
    return m_darker;
}

void HighlightEffect::setDarker(int darker)
{
    m_darker = darker;
}

void HighlightEffect::showDarker()
{
    m_effectState = ESDarker;

    this->repaint();
}

void HighlightEffect::showLighter()
{
    m_effectState = ESLighter;

    this->repaint();
}

void HighlightEffect::showNormal()
{
    m_effectState = ESNormal;

    this->repaint();
}

void HighlightEffect::paintEvent(QPaintEvent *)
{
    if (m_source)
    {
        QPixmap pixmap = m_source->grab();

        switch (m_effectState)
        {
        case ESDarker:
            pixmapDarker(&pixmap);
            break;
        case ESLighter:
            pixmapLigher(&pixmap);
        default:
            break;
        }

        QPainter painter;
        painter.begin(this);

        painter.setClipRect(rect());

        painter.drawPixmap(0, 0, pixmap);

        painter.end();
    }
}

void HighlightEffect::pixmapLigher(QPixmap *pixmap)
{
    QImage img = pixmap->toImage();  // slow

    for (int y=0; y < img.height(); y++)
    {
        for (int x = 0; x < img.width(); x++)
        {
            QRgb pix = img.pixel(x,y);
            QColor col(pix);
            col = col.lighter(m_lighter);
            img.setPixel(x, y, qRgba(col.red(), col.green(), col.blue(), qAlpha(pix)));
        }
    }
    pixmap->convertFromImage(img); // slow
}

void HighlightEffect::pixmapDarker(QPixmap *pixmap)
{
    QImage img = pixmap->toImage();  // slow

    for (int y=0; y < img.height(); y++)
    {
        for (int x = 0; x < img.width(); x++)
        {
            QRgb pix = img.pixel(x,y);
            QColor col(pix);
            col = col.darker(m_darker);
            img.setPixel(x, y, qRgba(col.red(), col.green(), col.blue(), qAlpha(pix)));
        }
    }
    pixmap->convertFromImage(img); // slow
}
