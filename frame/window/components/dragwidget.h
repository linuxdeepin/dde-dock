// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
    bool m_dragStatus;
    QPoint m_resizePoint;
};

#endif // DRAGWIDGET_H
