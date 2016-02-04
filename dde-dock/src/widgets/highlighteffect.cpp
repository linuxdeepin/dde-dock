/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include "highlighteffect.h"

HighlightEffect::HighlightEffect(QWidget * source, QWidget *parent) :
    QWidget(parent),
    m_source(source)
{
    setFixedSize(m_source->size());
    setAttribute(Qt::WA_TransparentForMouseEvents);
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
    setVisible(true);

    m_effectState = ESDarker;

    update();
}

void HighlightEffect::showLighter()
{
    setVisible(true);

    m_effectState = ESLighter;

    update();
}

void HighlightEffect::showNormal()
{
    setVisible(true);

    m_effectState = ESNormal;

    update();
}

void HighlightEffect::resizeEvent(QResizeEvent * event)
{
    QWidget::resizeEvent(event);

    update();
}

void HighlightEffect::paintEvent(QPaintEvent *)
{
    if (m_source && m_source->isVisible())
    {
        QPixmap pixmap = m_source->grab();

        QPainter painter;
        painter.begin(&pixmap);

        painter.setCompositionMode(QPainter::CompositionMode_SourceIn);

        if (m_effectState == ESLighter) {
            painter.fillRect(pixmap.rect(), QColor::fromRgbF(1, 1, 1, 0.2));
        } else if (m_effectState == ESDarker) {
            painter.fillRect(pixmap.rect(), QColor::fromRgbF(0, 0, 0, 0.2));
        }

        painter.end();

        painter.begin(this);
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
