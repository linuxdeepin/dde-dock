/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
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
#ifndef DRAGWIDGET_H
#define DRAGWIDGET_H

#include <QWidget>

class DragWidget : public QWidget
{
    Q_OBJECT

public:
    explicit DragWidget(QWidget *parent = nullptr);

    bool isDraging() const;

public Q_SLOTS:
    void onTouchMove(double scaleX, double scaleY);

Q_SIGNALS:
    void dragPointOffset(QPoint);
    void dragFinished();

protected:
    void mousePressEvent(QMouseEvent *event) override;
    void mouseMoveEvent(QMouseEvent *) override;
    void mouseReleaseEvent(QMouseEvent *) override;
    void enterEvent(QEvent *) override;
    void leaveEvent(QEvent *) override;

private:
    void updateCursor();

private:
    bool m_dragStatus;
    QPoint m_resizePoint;
};

#endif // DRAGWIDGET_H
