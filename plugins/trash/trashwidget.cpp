#include "trashwidget.h"

#include <QPainter>
#include <QIcon>

TrashWidget::TrashWidget(QWidget *parent)
    : QWidget(parent),

      m_popupApplet(new PopupControlWidget(this))
{
    QIcon::setThemeName("deepin");

    m_popupApplet->setVisible(false);
}

QWidget *TrashWidget::popupApplet()
{
    return m_popupApplet;
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
