// Copyright (C) 2018 ~ 2020 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include <QObject>
#include <QApplication>
#include <QSignalSpy>
#include <QTest>

#include <gtest/gtest.h>

#include "dockpopupwindow.h"

#include <DRegionMonitor>

DWIDGET_USE_NAMESPACE

class Test_DockPopupWindow : public QObject, public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;
};

void Test_DockPopupWindow::SetUp()
{
}

void Test_DockPopupWindow::TearDown()
{
}

TEST_F(Test_DockPopupWindow, coverage_test)
{
    DockPopupWindow *window = new DockPopupWindow;
    QWidget *w = new QWidget;
    w->setObjectName("test widget");
    window->setContent(w);

    window->show(QCursor::pos(), false);
    ASSERT_FALSE(window->model());

    window->hide();

    window->show(QCursor::pos(), true);
    ASSERT_TRUE(window->model());

    delete window;
    window = nullptr;

    ASSERT_TRUE(true);
}

TEST_F(Test_DockPopupWindow, onGlobMouseRelease)
{
    DockPopupWindow *window = new DockPopupWindow;
    QWidget *w = new QWidget;
    w->setObjectName("test widget");
    window->setContent(w);

    window->show(QCursor::pos(), true);

    ASSERT_TRUE(window->model());

    window->onGlobMouseRelease(QPoint(0, 0), DTK_GUI_NAMESPACE::DRegionMonitor::WatchedFlags::Button_Middle);
    window->onGlobMouseRelease(QPoint(0, 0), DTK_GUI_NAMESPACE::DRegionMonitor::WatchedFlags::Button_Left);

    qApp->processEvents();
    QTest::qWait(10);
    window->ensureRaised();

    QResizeEvent event(QSize(10, 10), QSize(20, 20));
    qApp->sendEvent(w, &event);
    QTest::qWait(15);

    delete window;
}
