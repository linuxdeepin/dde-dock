#include "closebutton.h"

CloseButton::CloseButton(QWidget *parent) : QLabel(parent)
{
    QPixmap iconPixmap(28,28);
    iconPixmap.load(ICON_NORMAL_PATH);
    this->setPixmap(iconPixmap);
}

void CloseButton::mousePressEvent(QMouseEvent *ev)
{
    QPixmap iconPixmap;
    iconPixmap.load(ICON_PRESS_PATH);
    this->setPixmap(iconPixmap);
    emit pressed();
    isPressed = true;
}

void CloseButton::mouseReleaseEvent(QMouseEvent *ev)
{
    QPixmap iconPixmap;
    iconPixmap.load(ICON_NORMAL_PATH);
    this->setPixmap(iconPixmap);
    emit released();
    if (isPressed)
    {
        emit clicked();
        isPressed = false;
    }
}

void CloseButton::enterEvent(QEvent *)
{
    QPixmap iconPixmap;
    iconPixmap.load(ICON_HOVER_PATH);
    this->setPixmap(iconPixmap);
    emit hovered();
}

void CloseButton::leaveEvent(QEvent *)
{
    QPixmap iconPixmap;
    iconPixmap.load(ICON_NORMAL_PATH);
    this->setPixmap(iconPixmap);
    emit exited();
}

