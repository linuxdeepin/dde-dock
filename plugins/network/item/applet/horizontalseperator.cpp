#include "horizontalseperator.h"

#include <QPainter>

HorizontalSeperator::HorizontalSeperator(QWidget *parent)
    : QWidget(parent)
{
    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
}

void HorizontalSeperator::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    QPainter painter(this);
    painter.fillRect(rect(), QColor(255, 255, 255, 255 * 0.1));
}
