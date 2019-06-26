/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#include "volumeslider.h"

#include <QMouseEvent>
#include <QDebug>
#include <QTimer>

VolumeSlider::VolumeSlider(QWidget *parent)
    : QSlider(Qt::Horizontal, parent),
      m_pressed(false),
      m_timer(new QTimer(this))
{
    setTickInterval(50);
    setPageStep(50);
    setTickPosition(QSlider::NoTicks);
    setFixedHeight(22);
    setStyleSheet("QSlider::groove {"
                  "margin-left:11px;"
                  "margin-right:11px;"
                  "border:none;"
                  "height:2px;"
//                  "border-width:0 0px 0 0px;"
//                  "background:url(://slider_bg.png) 0 2 0 2 stretch;"
                  "}"
                  "QSlider::handle{"
                  "background:url(://slider_handle.svg) no-repeat;"
                  "width:22px;"
                  "height:22px;"
                  "margin:-9px -14px -11px -14px;"
                  "}"
                  "QSlider::add-page {"
                  "background-color:rgba(255, 255, 255, .1);"
                  "}"
                  "QSlider::sub-page {"
                  "background-color:rgba(255, 255, 255, .8);"
                  "}");

    m_timer->setInterval(100);

    connect(m_timer, &QTimer::timeout, this, &VolumeSlider::onTimeout);
}

void VolumeSlider::setValue(const int value)
{
    if (m_pressed)
        return;

    blockSignals(true);
    QSlider::setValue(value);
    blockSignals(false);
}

void VolumeSlider::mousePressEvent(QMouseEvent *e)
{
    if (e->button() == Qt::LeftButton)
    {
        if (!rect().contains(e->pos()))
            return;
        m_pressed = true;
        QSlider::setValue(maximum() * e->x() / rect().width());
    }
}

void VolumeSlider::mouseMoveEvent(QMouseEvent *e)
{
    const int value = minimum() + (double((maximum()) - minimum()) * e->x() / rect().width());
    const int normalized = std::max(std::min(maximum(), value), 0);

    QSlider::setValue(normalized);

    emit valueChanged(normalized);
}

void VolumeSlider::mouseReleaseEvent(QMouseEvent *e)
{
    if (e->button() == Qt::LeftButton)
    {
        m_pressed = false;
        emit requestPlaySoundEffect();
    }
}

void VolumeSlider::wheelEvent(QWheelEvent *e)
{
    e->accept();

    m_timer->start();

    QSlider::setValue(value() + (e->delta() > 0 ? 10 : -10));
}

void VolumeSlider::onTimeout()
{
    m_timer->stop();
    emit requestPlaySoundEffect();
}
