/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef TRAYICON_H
#define TRAYICON_H

#include <QWindow>
#include <QFrame>
#include <QImage>

class QPaintEvent;
class QMouseEvent;
class TrayIcon : public QFrame
{
    Q_OBJECT
public:
    explicit TrayIcon(WId winId, QWidget *parent = 0);
    ~TrayIcon();

    void maskOn();
    void maskOff();

    void hideIcon();
    void updateIcon();

protected:
    void paintEvent(QPaintEvent *);
//    void mousePressEvent(QMouseEvent *);

private:
    WId m_windowId;
    WId m_containerWid;
    bool m_masked;
    bool m_hided;
    QTimer * m_timer;
    QImage m_image;

    void wrapWindow();
    QImage getImageNonComposite();
    void sendClick(uint8_t, int, int);
};

#endif // TRAYICON_H
