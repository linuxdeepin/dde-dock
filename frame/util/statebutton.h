// Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef STATEBUTTON_H
#define STATEBUTTON_H

#include <QWidget>

class StateButton : public QWidget
{
    Q_OBJECT

public:
    enum Type {
        Check,
        Fork
    };

public:
    explicit StateButton(QWidget *parent = nullptr);
    void setType(Type type);
    void setSwitchFork(bool switchFork);

signals:
    void click();

protected:
    void paintEvent(QPaintEvent *event) override;
    void mousePressEvent(QMouseEvent *event) override;
    void enterEvent(QEvent *event) override;
    void leaveEvent(QEvent *event) override;

private:
    void drawCheck(QPainter &painter, QPen &pen, int radius);
    void drawFork(QPainter &painter, QPen &pen, int radius);

private:
    Type m_type;
    bool m_switchFork;
};

#endif // STATEBUTTON_H
