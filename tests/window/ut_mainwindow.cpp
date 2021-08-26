/*
 * Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
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
#include <QThread>
#include <QTest>
#include <QSignalSpy>

#include <gtest/gtest.h>

#define private public
#include "mainwindow.h"
#undef private

using namespace ::testing;

class Test_MainWindow : public ::testing::Test
{
};

TEST_F(Test_MainWindow, onDbusNameOwnerChanged)
{
    MainWindow window;
    window.onDbusNameOwnerChanged("org.kde.StatusNotifierWatcher", "old", "new");
}

TEST_F(Test_MainWindow, compositeChanged)
{
    MainWindow window;
    window.compositeChanged();

    ASSERT_TRUE(true);
}

TEST_F(Test_MainWindow, moveEvent)
{
    MainWindow window;
    QMouseEvent moveEvent(QEvent::MouseMove, QPointF(0, 0), Qt::LeftButton, Qt::LeftButton, Qt::ControlModifier);
    window.mouseMoveEvent(&moveEvent);

    ASSERT_TRUE(true);
}

TEST_F(Test_MainWindow, mousePressEvent)
{
    // 显示菜单会阻塞住
//    MainWindow window;
//    QMouseEvent mousePressEvent(QEvent::MouseButtonPress, QPointF(0, 0), Qt::RightButton, Qt::RightButton, Qt::ControlModifier);
//    window.mousePressEvent(&mousePressEvent);

    ASSERT_TRUE(true);
}

TEST_F(Test_MainWindow, launch)
{
    MainWindow window;

    qApp->setProperty("CANSHOW", false);
    window.launch();

    qApp->setProperty("CANSHOW", true);
    window.launch();

    window.callShow();

    ASSERT_TRUE(true);
}

TEST_F(Test_MainWindow, RegisterDdeSession)
{
    MainWindow window;
    window.RegisterDdeSession();

    qputenv("DDE_SESSION_PROCESS_COOKIE_ID", "111");
    window.RegisterDdeSession();

    ASSERT_TRUE(true);
}

TEST_F(Test_MainWindow, resetDragWindow_test)
{
    MainWindow window;
    window.m_multiScreenWorker->m_position = Position::Top;
    window.resetDragWindow();
    window.onMainWindowSizeChanged(QPoint(10, 10));
    window.touchRequestResizeDock();

    window.m_multiScreenWorker->m_position = Position::Bottom;
    window.resetDragWindow();
    window.onMainWindowSizeChanged(QPoint(10, 10));
    window.touchRequestResizeDock();

    window.m_multiScreenWorker->m_position = Position::Left;
    window.resetDragWindow();
    window.onMainWindowSizeChanged(QPoint(10, 10));
    window.touchRequestResizeDock();

    window.m_multiScreenWorker->m_position = Position::Right;
    window.resetDragWindow();
    window.onMainWindowSizeChanged(QPoint(10, 10));
    window.touchRequestResizeDock();
}

TEST_F(Test_MainWindow, adjustShadowMask)
{
    MainWindow *window = new MainWindow;

    window->RegisterDdeSession();

    window->m_launched = true;
    window->adjustShadowMask();

//        window->onDbusNameOwnerChanged(SNI_WATCHER_SERVICE, "old", "new");
    delete window;
}

TEST_F(Test_MainWindow, event_test)
{
    MainWindow *window = new MainWindow;

    QMouseEvent event4(QEvent::MouseMove, QPointF(0, 0), Qt::RightButton, Qt::RightButton, Qt::ControlModifier);
    window->mouseMoveEvent(&event4);

    QResizeEvent event5((QSize()), QSize());
    window->resizeEvent(&event5);

    QMimeData data;
    data.setText("test");

    QDragEnterEvent event9(QPoint(), Qt::DropAction::CopyAction, &data, Qt::LeftButton, Qt::NoModifier);
    window->dragEnterEvent(&event9);

    QKeyEvent event11(QEvent::Type::KeyPress, Qt::Key_Escape, Qt::ControlModifier);
    window->keyPressEvent(&event11);

    QEnterEvent event12(QPointF(0.0, 0.0), QPointF(0.0, 0.0), QPointF(0.0, 0.0));
    window->enterEvent(&event12);

    delete window;
}

TEST_F(Test_MainWindow, coverage_test)
{
    MainWindow *window = new MainWindow;

    window->getTrayVisableItemCount();
    window->adjustShadowMask();
    window->resetDragWindow();
    window->onMainWindowSizeChanged(QPoint(10, 10));
    window->touchRequestResizeDock();
    window->sendNotifications();

    window->m_multiScreenWorker->m_hideMode = HideMode::SmartHide;
    window->m_multiScreenWorker->m_hideState = HideState::Hide;

//    window->callShow();

//    window->m_multiScreenWorker->m_hideState = HideState::Show;
//    window->callShow();

//    window->m_multiScreenWorker->m_hideMode = HideMode::KeepShowing;
//    window->callShow();

//    window->m_multiScreenWorker->m_hideMode = HideMode::KeepHidden;
//    window->callShow();

    delete window;
}

TEST_F(Test_MainWindow, dragWidget_test)
{
//    DragWidget w;

//    ASSERT_FALSE(w.m_dragStatus);
//    QTest::mousePress(&w, Qt::LeftButton);
//    ASSERT_TRUE(w.m_dragStatus);

//    QMouseEvent event(QEvent::MouseMove, QPointF(0, 0), Qt::LeftButton, Qt::LeftButton, Qt::ControlModifier);
//    w.mouseMoveEvent(&event);

//    QTest::mouseRelease(&w, Qt::LeftButton);
//    ASSERT_FALSE(w.m_dragStatus);

//    ASSERT_EQ(w.objectName(), "DragWidget");

//    w.onTouchMove(1.0, 1.0);

//    QEvent enterEvent(QEvent::Enter);
//    w.enterEvent(&enterEvent);

//    QEvent leaveEvent(QEvent::Leave);
//    w.leaveEvent(&leaveEvent);
}

TEST_F(Test_MainWindow, test4)
{
    MainWindow *window = new MainWindow;
    MultiScreenWorker *worker = window->m_multiScreenWorker;

    worker->reInitDisplayData();

    worker->setStates(MultiScreenWorker::AutoHide, true);
    worker->onAutoHideChanged(true);
    QTest::qWait(510);

    worker->dockRectWithoutScale("", Position::Top, HideMode::KeepShowing, DisplayMode::Fashion);
    worker->dockRectWithoutScale("", Position::Top, HideMode::KeepHidden, DisplayMode::Fashion);
    worker->dockRectWithoutScale("", Position::Top, HideMode::SmartHide, DisplayMode::Fashion);

    worker->onWindowSizeChanged(1);
    worker->onRequestUpdateMonitorInfo();

    worker->setStates(MultiScreenWorker::LauncherDisplay, false);
    worker->onRequestDelayShowDock();

    worker->updateParentGeometry(QRect(), Position::Top);
    worker->updateParentGeometry(QRect(), Position::Bottom);
    worker->updateParentGeometry(QRect(), Position::Left);
    worker->updateParentGeometry(QRect(), Position::Right);

    worker->m_position = Position::Top;
    worker->onRequestNotifyWindowManager();
    worker->m_position = Position::Bottom;
    worker->onRequestNotifyWindowManager();
    worker->m_position = Position::Left;
    worker->onRequestNotifyWindowManager();
    worker->m_position = Position::Right;
    worker->onRequestNotifyWindowManager();

    worker->m_hideMode = HideMode::SmartHide;
    worker->onExtralRegionMonitorChanged(0, 0, worker->m_registerKey);

    worker->m_hideMode = HideMode::KeepHidden;
    worker->onExtralRegionMonitorChanged(0, 0, worker->m_registerKey);

    worker->m_hideMode = HideMode::KeepShowing;
    worker->onExtralRegionMonitorChanged(0, 0, worker->m_registerKey);

    ASSERT_TRUE(true);
    delete window;
}
