// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "dragwidget.h"
#include "utils.h"
#include "constants.h"

#include <QCoreApplication>
#include <QMouseEvent>

DragWidget::DragWidget(QWidget *parent)
    : QWidget(parent)
    , m_dragStatus(false)
{
    setObjectName("DragWidget");
}

bool DragWidget::isDraging() const
{
    return m_dragStatus;
}

void DragWidget::onTouchMove(double scaleX, double scaleY)
{
    Q_UNUSED(scaleX);
    Q_UNUSED(scaleY);

    static QPoint lastPos;
    QPoint curPos = QCursor::pos();
    if (lastPos == curPos) {
        return;
    }
    lastPos = curPos;
    qApp->postEvent(this, new QMouseEvent(QEvent::MouseMove, mapFromGlobal(curPos)
                                                  , QPoint(), curPos, Qt::LeftButton, Qt::LeftButton
                                          , Qt::NoModifier, Qt::MouseEventSynthesizedByApplication));
}

void DragWidget::mousePressEvent(QMouseEvent *event)
{
    // qt转发的触屏按下信号不进行响应
    if (event->source() == Qt::MouseEventSynthesizedByQt)
        return;

    if (event->button() == Qt::LeftButton) {
        m_resizePoint = event->globalPos();
        m_dragStatus = true;
        this->grabMouse();
    }
}

void DragWidget::mouseMoveEvent(QMouseEvent *)
{
    if (m_dragStatus) {
        QPoint offset = QPoint(QCursor::pos() - m_resizePoint);
        emit dragPointOffset(offset);
    }
}

void DragWidget::mouseReleaseEvent(QMouseEvent *)
{
    if (!m_dragStatus)
        return;

    m_dragStatus =  false;
    releaseMouse();
    emit dragFinished();
}

void DragWidget::enterEvent(QEvent *)
{
    QApplication::setOverrideCursor(cursor());
}

void DragWidget::leaveEvent(QEvent *)
{
    QApplication::setOverrideCursor(Qt::ArrowCursor);
}

