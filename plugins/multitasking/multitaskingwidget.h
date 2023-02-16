// Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef MultitaskingWidget_H
#define MultitaskingWidget_H

#include <QWidget>
#include <QIcon>

#define PLUGIN_KEY "multitasking"

class MultitaskingWidget : public QWidget
{
    Q_OBJECT

public:
    explicit MultitaskingWidget(QWidget *parent = nullptr);
    void refreshIcon();

signals:
    void requestContextMenu(const QString &itemKey) const;

protected:
    void paintEvent(QPaintEvent *e) override;
    QIcon m_icon;
};

#endif // MULTITASKINGWIDGET_H
