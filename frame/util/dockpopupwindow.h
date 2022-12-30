// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DOCKPOPUPWINDOW_H
#define DOCKPOPUPWINDOW_H

#include <darrowrectangle.h>
#include <DRegionMonitor>
#include <DWindowManagerHelper>

DGUI_USE_NAMESPACE

class DockPopupWindow : public Dtk::Widget::DArrowRectangle
{
    Q_OBJECT

public:
    explicit DockPopupWindow(QWidget *parent = 0);
    ~DockPopupWindow();

    bool model() const;

    void setContent(QWidget *content);

public slots:
    void show(const QPoint &pos, const bool model = false);
    void show(const int x, const int y);
    void hide();

signals:
    void accept() const;
    // 在把专业版的仓库降级到debian的stable时, dock出现了一个奇怪的问题:
    // 在plugins/tray/system-trays/systemtrayitem.cpp中的showPopupWindow函数中
    // 无法连接到上面这个信号: "accept", qt给出一个运行时警告提示找不到信号
    // 目前的解决方案就是在下面增加了这个信号
    void unusedSignal();

protected:
    void showEvent(QShowEvent *e);
    void enterEvent(QEvent *e);
    bool eventFilter(QObject *o, QEvent *e);
    void blockButtonRelease();

private slots:
    void onGlobMouseRelease(const QPoint &mousePos, const int flag);
    void compositeChanged();
    void ensureRaised();

private:
    bool m_model;
    QPoint m_lastPoint;

    DRegionMonitor *m_regionInter;
    DWindowManagerHelper *m_wmHelper;
    bool m_enableMouseRelease;
};

#endif // DOCKPOPUPWINDOW_H
