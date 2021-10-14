/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     chenwei <chenwei@uniontech.com>
 *
 * Maintainer: chenwei <chenwei@uniontech.com>
 *
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#include "statebutton.h"

#include <QIcon>
#include <QPainter>
#include <QtMath>

StateButton::StateButton(QWidget *parent)
    : QWidget(parent)
    , m_type(Check)
    , m_switchFork(true)
{
    setAttribute(Qt::WA_TranslucentBackground, true);
}

void StateButton::setType(StateButton::Type type)
{
    m_type =  type;
    update();
}

void StateButton::setSwitchFork(bool switchFork)
{
    m_switchFork = switchFork;
}

void StateButton::paintEvent(QPaintEvent *event)
{
    QWidget::paintEvent(event);
    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing, true);

    int radius = qMin(width(), height());
    painter.setPen(QPen(Qt::NoPen));
    painter.setBrush(palette().color(QPalette::Highlight));
    painter.drawPie(rect(), 0, 360 * 16);

    QPen pen(Qt::white, radius / 100.0 * 6.20, Qt::SolidLine, Qt::SquareCap, Qt::MiterJoin);
    switch (m_type) {
    case Check: drawCheck(painter, pen, radius);    break;
    case Fork:  drawFork(painter, pen, radius);     break;
    }
}

void StateButton::mousePressEvent(QMouseEvent *event)
{
    QWidget::mousePressEvent(event);
    if (m_switchFork)
        emit click();
}

void StateButton::enterEvent(QEvent *event)
{
    QWidget::enterEvent(event);
    if (m_switchFork)
        setType(Fork);
}

void StateButton::leaveEvent(QEvent *event)
{
    QWidget::leaveEvent(event);
    if (m_switchFork)
        setType(Check);
}

void StateButton::drawCheck(QPainter &painter, QPen &pen, int radius)
{
    painter.setPen(pen);

    QPointF points[3] = {
        QPointF(radius / 100.0 * 32,  radius / 100.0 * 57),
        QPointF(radius / 100.0 * 45,  radius / 100.0 * 70),
        QPointF(radius / 100.0 * 75,  radius / 100.0 * 35)
    };

    painter.drawPolyline(points, 3);
}

void StateButton::drawFork(QPainter &painter, QPen &pen, int radius)
{
    pen.setCapStyle(Qt::RoundCap);
    painter.setPen(pen);

    QPointF pointsl[2] = {
        QPointF(radius / 100.0 * 35,  radius / 100.0 * 35),
        QPointF(radius / 100.0 * 65,  radius / 100.0 * 65)
    };

    painter.drawPolyline(pointsl, 2);

    QPointF pointsr[2] = {
        QPointF(radius / 100.0 * 65,  radius / 100.0 * 35),
        QPointF(radius / 100.0 * 35,  radius / 100.0 * 65)
    };

    painter.drawPolyline(pointsr, 2);
}
