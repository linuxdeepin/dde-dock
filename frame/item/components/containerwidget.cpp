#include "containerwidget.h"

#include <QDebug>

ContainerWidget::ContainerWidget(QWidget *parent)
    : QWidget(parent)
{
}

QSize ContainerWidget::sizeHint() const
{
    return QSize(80, 40);
}
