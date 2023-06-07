// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef PLUGINWIDGET_H
#define PLUGINWIDGET_H

#include "constants.h"

#include <QWidget>
#include <QIcon>

class ShutdownWidget : public QWidget
{
    Q_OBJECT

public:
    explicit ShutdownWidget(QWidget *parent = Q_NULLPTR);
    QPixmap loadPixmap() const;

protected:
    void paintEvent(QPaintEvent *e) override;
    void mousePressEvent(QMouseEvent *event) override;
    void mouseReleaseEvent(QMouseEvent *event) override;
    void mouseMoveEvent(QMouseEvent *event) override;
    void leaveEvent(QEvent *event) override;

private:
    Dock::DisplayMode m_displayMode;
    bool m_hover;
    bool m_pressed;
};

#endif // PLUGINWIDGET_H
