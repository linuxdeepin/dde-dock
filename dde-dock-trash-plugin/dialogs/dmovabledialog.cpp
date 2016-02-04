/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include "dmovabledialog.h"
#include <QMouseEvent>
#include <QApplication>
#include <QDesktopWidget>
#include <QPushButton>
#include <QResizeEvent>


DMovabelDialog::DMovabelDialog(QWidget *parent):QDialog(parent)
{
    setWindowFlags(Qt::FramelessWindowHint | Qt::Dialog);
    setAttribute(Qt::WA_TranslucentBackground);
    m_closeButton = new QPushButton(this);
    m_closeButton->setObjectName("CloseButton");
    m_closeButton->setFixedSize(25, 25);
    m_closeButton->setAttribute(Qt::WA_NoMousePropagation);
    connect(m_closeButton, SIGNAL(clicked()), this, SLOT(close()));
}

void DMovabelDialog::setMovableHeight(int height){
    m_movableHeight = height;
}


QPushButton* DMovabelDialog::getCloseButton(){
    return m_closeButton;
}

void DMovabelDialog::moveCenter(){
    QRect qr = frameGeometry();
    QPoint cp = qApp->desktop()->availableGeometry().center();
    qr.moveCenter(cp);
    move(qr.topLeft());
}

void DMovabelDialog::mousePressEvent(QMouseEvent *event)
{
    if(event->button() & Qt::LeftButton)
    {
        m_dragPosition = event->globalPos() - frameGeometry().topLeft();
    }
    QDialog::mousePressEvent(event);
}

void DMovabelDialog::mouseReleaseEvent(QMouseEvent *event)
{
    QDialog::mouseReleaseEvent(event);
}

void DMovabelDialog::mouseMoveEvent(QMouseEvent *event)
{
    move(event->globalPos() - m_dragPosition);
    QDialog::mouseMoveEvent(event);
}

void DMovabelDialog::resizeEvent(QResizeEvent *event){
    m_closeButton->move(width() - m_closeButton->width() - 4, 4);
    m_closeButton->raise();
    moveCenter();
    QDialog::resizeEvent(event);
}

DMovabelDialog::~DMovabelDialog()
{

}

