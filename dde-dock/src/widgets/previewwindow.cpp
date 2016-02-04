/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include "previewwindow.h"

PreviewWindow::PreviewWindow(ArrowDirection direction, QWidget *parent) : DArrowRectangle(direction, parent)
{
    setWindowFlags(Qt::X11BypassWindowManagerHint | Qt::Tool);
    setArrowWidth(ARROW_WIDTH);
    setArrowHeight(ARROW_HEIGHT);

    //QWidget calls this function after it has been fully constructed but before it is shown the very first time.
    //so call it to make sure the style sheet won't repolish and cover the value set by seter-function
    ensurePolished();
    setRadius(4);
    setBorderWidth(1);
    setBorderColor(QColor(0, 0, 0, 0.2 * 255));

    m_showTimer = new QTimer(this);
    m_showTimer->setSingleShot(true);
    connect(m_showTimer, &QTimer::timeout, this, &PreviewWindow::onShowTimerTriggered);

    m_hideTimer = new QTimer(this);
    m_hideTimer->setSingleShot(true);
    connect(m_hideTimer, &QTimer::timeout, this, &PreviewWindow::hide);

    m_animation = new QPropertyAnimation(this, "arrowPos");
    m_animation->setDuration(MOVE_ANIMATION_DURATION);
    m_animation->setEasingCurve(MOVE_ANIMATION_CURVE);
}

PreviewWindow::~PreviewWindow()
{

}

void PreviewWindow::showPreview(int x, int y, int interval)
{
    m_hideTimer->stop();

    if (m_showTimer->isActive())
        return;

    m_lastPos = QPoint(m_x, m_y);
    m_x = x;
    m_y = y;

    m_showTimer->start(interval);
}

void PreviewWindow::hidePreview(int interval)
{
    m_showTimer->stop();

    if (interval <= 0) {
        m_animation->stop();
        hide();
    }
    else
        m_hideTimer->start(interval);
}

void PreviewWindow::setContent(QWidget *content)
{
    m_currentContent = content;
}

void PreviewWindow::setArrowPos(const QPoint &pos)
{
    show(pos.x(), pos.y());
}

void PreviewWindow::hide()
{
    if (m_lastContent != m_currentContent)
        emit hideFinish(m_lastContent);

    DArrowRectangle::hide();
}

void PreviewWindow::enterEvent(QEvent *)
{
    m_hideTimer->stop();
}

void PreviewWindow::leaveEvent(QEvent *)
{
    m_hideTimer->start();
}

void PreviewWindow::onShowTimerTriggered()
{
    if (!m_lastContent.isNull()) {
        m_lastContent.data()->setParent(NULL);
        if (m_lastContent != m_currentContent)
            emit showFinish(m_lastContent);
    }

    DArrowRectangle::setContent(m_currentContent);
    m_lastContent = m_currentContent;

    if (isHidden())
        show(m_x, m_y);
    else{
        m_animation->setStartValue(m_lastPos);
        m_animation->setEndValue(QPoint(m_x, m_y));
        m_animation->start();
    }
}

