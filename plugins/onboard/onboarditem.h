// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef ONBOARDITEM_H
#define ONBOARDITEM_H

#include "constants.h"

#include <DGuiApplicationHelper>

#include <QWidget>
#include <QIcon>

DGUI_USE_NAMESPACE

class OnboardItem : public QWidget
{
    Q_OBJECT

public:
    explicit OnboardItem(QWidget *parent = nullptr);
    QPixmap iconPixmap(QSize size, DGuiApplicationHelper::ColorType themeType) const;

protected:
    void paintEvent(QPaintEvent *e) override;

private:
    const QPixmap loadSvg(const QString &fileName, const QSize &size) const;
    void mousePressEvent(QMouseEvent *event) override;
    void mouseReleaseEvent(QMouseEvent *event) override;
    void mouseMoveEvent(QMouseEvent *event) override;
    void leaveEvent(QEvent *event) override;
    void resizeEvent(QResizeEvent *event) override;

private:
    Dock::DisplayMode m_displayMode;
    bool m_hover;
    bool m_pressed;
    QIcon m_icon;
};

#endif // ONBOARDITEM_H
