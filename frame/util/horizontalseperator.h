// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef HORIZONTALSEPERATOR_H
#define HORIZONTALSEPERATOR_H

#include <QWidget>

class HorizontalSeperator : public QWidget
{
    Q_OBJECT

public:
    explicit HorizontalSeperator(QWidget *parent = nullptr);

    QSize sizeHint() const override;

protected:
    void paintEvent(QPaintEvent *e) override;
};

#endif // HORIZONTALSEPERATOR_H
