#include "volumeslider.h"

#include <QMouseEvent>
#include <QDebug>
#include <QTimer>

VolumeSlider::VolumeSlider(QWidget *parent)
    : QSlider(Qt::Horizontal, parent),
      m_pressed(false)
{
    setMinimum(0);
    setMaximum(1000);
    setTickInterval(50);
    setPageStep(50);
    setTickPosition(QSlider::NoTicks);
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
        QSlider::setValue(1000.0 * e->x() / rect().width());
    }
}

void VolumeSlider::mouseMoveEvent(QMouseEvent *e)
{
    const int value = minimum() + (double((maximum()) - minimum()) * e->x() / rect().width());

    QSlider::setValue(std::max(std::min(1000, value), 0));
    emit valueChanged(std::max(std::min(1000, value), 0));
}

void VolumeSlider::mouseReleaseEvent(QMouseEvent *e)
{
    if (e->button() == Qt::LeftButton)
        m_pressed = false;
    //        QTimer::singleShot(100, [this] {m_pressed = false;});
}

void VolumeSlider::wheelEvent(QWheelEvent *e)
{
    QSlider::setValue(value() + e->delta() / 5);
}
