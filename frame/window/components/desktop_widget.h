// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DESKTOP_WIDGET_H
#define DESKTOP_WIDGET_H

#include <QWidget>
#include <QPaintEvent>
#include <QEnterEvent>

class DesktopWidget : public QWidget
{
public:
    explicit DesktopWidget(QWidget *parent = nullptr);

private:
    bool checkNeedShowDesktop();

protected:
    void paintEvent(QPaintEvent *event) override;
    void enterEvent(QEvent *event) override;
    void leaveEvent(QEvent *event) override;
    void mousePressEvent(QMouseEvent *event) override;

private:
    bool m_isHover;         // 判断鼠标是否移到desktop区域
    bool m_needRecoveryWin;
};

#endif // DESKTOP_WIDGET_H
