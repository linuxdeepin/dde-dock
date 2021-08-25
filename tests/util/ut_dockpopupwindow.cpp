/*
 * Copyright (C) 2018 ~ 2020 Deepin Technology Co., Ltd.
 *
 * Author:     fanpengcheng <fanpengcheng@uniontech.com>
 *
 * Maintainer: fanpengcheng <fanpengcheng@uniontech.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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

    window->onGlobMouseRelease(QPoint(0, 0), DRegionMonitor::WatchedFlags::Button_Middle);
    window->onGlobMouseRelease(QPoint(0, 0), DRegionMonitor::WatchedFlags::Button_Left);

    qApp->processEvents();
    QTest::qWait(10);
    window->ensureRaised();

    QResizeEvent event(QSize(10, 10), QSize(20, 20));
    qApp->sendEvent(w, &event);
    QTest::qWait(15);

    delete window;
}
