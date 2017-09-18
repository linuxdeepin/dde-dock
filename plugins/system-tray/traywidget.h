/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#ifndef TRAYWIDGET_H
#define TRAYWIDGET_H

#include <QWidget>
#include <QTimer>

class TrayWidget : public QWidget
{
    Q_OBJECT

public:
    explicit TrayWidget(quint32 winId, QWidget *parent = 0);
    ~TrayWidget();

    void updateIcon();
    const QImage trayImage();
    void sendClick(uint8_t mouseButton, int x, int y);
    void setActive(const bool active);

signals:
    void iconChanged();

private:
    QSize sizeHint() const;
    void showEvent(QShowEvent *e);
    void hideEvent(QHideEvent *e);
    void paintEvent(QPaintEvent *e);
    void mousePressEvent(QMouseEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);
    void moveEvent(QMoveEvent *e);
    void enterEvent(QEvent *e);
    void configContainerPosition();

    void wrapWindow();
//    void hideIcon();
    void refershIconImage();

private slots:
    void setX11PassMouseEvent(const bool pass);
    void setWindowOnTop(const bool top);
    bool isBadWindow();

private:
    bool m_active = false;
    WId m_windowId;
    WId m_containerWid;
    QImage m_image;

    QTimer *m_updateTimer;
    QTimer *m_ignoreRepeat;
    QPoint m_pressPoint;
};

#endif // TRAYWIDGET_H
