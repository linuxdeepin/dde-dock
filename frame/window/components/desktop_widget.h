// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DESKTOP_WIDGET_H
#define DESKTOP_WIDGET_H

#include <QWidget>
#include <QPaintEvent>
#include <QEnterEvent>
#include <QTimer>

class DesktopWidget : public QWidget
{
public:
    explicit DesktopWidget(QWidget *parent = nullptr);

    void setToggleDesktopInterval(int ms);
private:
    bool checkNeedShowDesktop();

protected:
    void paintEvent(QPaintEvent *event) override;
    void enterEvent(QEvent *event) override;
    void leaveEvent(QEvent *event) override;
    void mousePressEvent(QMouseEvent *event) override;
    void toggleDesktop();

private:
    bool m_isHover;         // 判断鼠标是否移到desktop区域
    bool m_needRecoveryWin;
    QTimer *m_timer;
};

#endif // DESKTOP_WIDGET_H
