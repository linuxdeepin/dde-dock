/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef REFLECTIONEFFECT_H
#define REFLECTIONEFFECT_H

#include <QWidget>
#include <QPainter>

class QPaintEvent;
class ReflectionEffect : public QWidget
{
    Q_OBJECT
public:
    ReflectionEffect(QWidget * source, QWidget *parent = 0);

    qreal opacity() const;
    void setOpacity(const qreal &opacity);
    void updateReflection();

protected:
    void paintEvent(QPaintEvent * event) Q_DECL_OVERRIDE;

private:
    QWidget * m_source;
    qreal m_opacity;
};

#endif // REFLECTIONEFFECT_H
