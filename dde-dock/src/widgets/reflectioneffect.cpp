/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include "reflectioneffect.h"

ReflectionEffect::ReflectionEffect(QWidget * source, QWidget *parent) :
    QWidget(parent),
    m_source(source),
    m_opacity(0.1)
{
    this->setFixedWidth(m_source->width());
    setAttribute(Qt::WA_TransparentForMouseEvents);
}

qreal ReflectionEffect::opacity() const
{
    return m_opacity;
}

void ReflectionEffect::setOpacity(const qreal &opacity)
{
    m_opacity = opacity;
}

void ReflectionEffect::paintEvent(QPaintEvent *)
{
    if (m_source) {
        QPixmap pixmap = m_source->grab();

        // flip the pixmap
        pixmap = pixmap.transformed(QTransform().scale(1, -1));

        if (!pixmap.isNull()) {
            QPainter painter;
            painter.begin(this);

            painter.setClipRect(rect());
            painter.setOpacity(m_opacity);
            painter.drawPixmap(0, 0, pixmap);

            painter.end();
        }
    }
}

void ReflectionEffect::updateReflection()
{
    update();
}
