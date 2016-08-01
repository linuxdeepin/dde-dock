#include "trashwidget.h"

#include <QPainter>
#include <QIcon>

TrashWidget::TrashWidget(QWidget *parent)
    : QWidget(parent)
{
    QIcon::setThemeName("deepin");
}

void TrashWidget::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);

    const int size = std::min(width(), height()) * 0.8;
    QIcon icon = QIcon::fromTheme("user-trash");
    QPixmap pixmap = icon.pixmap(size, size);

    QPainter painter(this);
    painter.drawPixmap(rect().center() - pixmap.rect().center(), pixmap);
}
