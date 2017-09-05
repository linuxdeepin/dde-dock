#include "refreshbutton.h"

#include <QMouseEvent>
#include <QEvent>

RefreshButton::RefreshButton(QWidget *parent) : QLabel(parent)
{
    setAttribute(Qt::WA_TranslucentBackground);

    setPixmap(QPixmap(":/wireless/resources/wireless/refresh_normal.svg"));
}

void RefreshButton::enterEvent(QEvent *event)
{
    QLabel::enterEvent(event);

    setPixmap(QPixmap(":/wireless/resources/wireless/refresh_hover.svg"));
}

void RefreshButton::leaveEvent(QEvent *event)
{
    QLabel::leaveEvent(event);

    setPixmap(QPixmap(":/wireless/resources/wireless/refresh_normal.svg"));
}

void RefreshButton::mousePressEvent(QMouseEvent *event)
{
    if (event->button() == Qt::LeftButton)
        setPixmap(QPixmap(":/wireless/resources/wireless/refresh_press.svg"));
}

void RefreshButton::mouseReleaseEvent(QMouseEvent *event)
{
    if (event->button() == Qt::LeftButton)
        emit clicked();

    setPixmap(QPixmap(":/wireless/resources/wireless/refresh_normal.svg"));
}
