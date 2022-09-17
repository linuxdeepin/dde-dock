// SPDX-FileCopyrightText: 2019 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SHOWDESKTOPWIDGET_H
#define SHOWDESKTOPWIDGET_H

#include <QWidget>

class ShowDesktopWidget : public QWidget
{
    Q_OBJECT

public:
    explicit ShowDesktopWidget(QWidget *parent = 0);
    void refreshIcon();

signals:
    void requestContextMenu(const QString &itemKey) const;

protected:
    void paintEvent(QPaintEvent *e) override;
};

#endif // SHOWDESKTOPWIDGET_H
