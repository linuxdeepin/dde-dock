// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef OVERLAYWARNINGWIDGET_H
#define OVERLAYWARNINGWIDGET_H

#include <QWidget>

class OverlayWarningWidget : public QWidget
{
    Q_OBJECT

public:
    explicit OverlayWarningWidget(QWidget *parent = 0);

protected:
    QSize sizeHint() const;
    void paintEvent(QPaintEvent *e);

private:
    const QPixmap loadSvg(const QString &fileName, const QSize &size) const;
};

#endif // OVERLAYWARNINGWIDGET_H
