#include "sounditem.h"

#include <QPainter>

SoundItem::SoundItem(QWidget *parent)
    : QWidget(parent)
{

}

QSize SoundItem::sizeHint() const
{
    return QSize(24, 24);
}

void SoundItem::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    painter.fillRect(rect(), Qt::red);
}
