/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef PREVIEWWINDOW_H
#define PREVIEWWINDOW_H

#include <QTimer>
#include <QPointer>
#include <QPropertyAnimation>

#include <libdui/darrowrectangle.h>

DUI_USE_NAMESPACE

class PreviewWindow : public DArrowRectangle
{
    Q_OBJECT
    Q_PROPERTY(QPoint arrowPos READ pos WRITE setArrowPos)
public:
    explicit PreviewWindow(DArrowRectangle::ArrowDirection direction, QWidget *parent = 0);
    ~PreviewWindow();

    void showPreview(int x, int y, int interval);
    void hidePreview(int interval = 0);
    void setContent(QWidget *content);
    void setArrowPos(const QPoint &pos);
    void hide();

signals:
    void showFinish(QWidget *lastContent);
    void hideFinish(QWidget *lastContent);

protected:
    void enterEvent(QEvent *);
    void leaveEvent(QEvent *);

private:
    void onShowTimerTriggered();

private:
    QTimer *m_showTimer = NULL;
    QTimer *m_hideTimer = NULL;
    QPointer<QWidget> m_currentContent;
    QPointer<QWidget> m_lastContent;
    QPropertyAnimation *m_animation = NULL;
    QPoint m_lastPos = QPoint(0, 0);
    int m_x = 0;
    int m_y = 0;
    const int ARROW_WIDTH = 18;
    const int ARROW_HEIGHT = 11;
    const int MOVE_ANIMATION_DURATION = 300;
    const QEasingCurve MOVE_ANIMATION_CURVE = QEasingCurve::OutCirc;
};

#endif // PREVIEWWINDOW_H
