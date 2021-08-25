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
#include <QTest>

#include <DWindowManagerHelper>

#include <gtest/gtest.h>
#define private public
#include "mainwindow.h"
#include "multiscreenworker.h"
#undef private

class Test_MultiScreenWorker : public ::testing::Test
{};

TEST_F(Test_MultiScreenWorker, coverage_test)
{
    MainWindow window;
    MultiScreenWorker *worker = new MultiScreenWorker(&window, DWindowManagerHelper::instance());

    qDebug() << worker->dockRect("test screen");

    worker->m_displayMode = DisplayMode::Fashion;
    worker->updateDaemonDockSize(40);

    worker->m_displayMode = DisplayMode::Efficient;
    worker->updateDaemonDockSize(20);

    QDBusMessage msg;
    worker->handleDbusSignal(msg);

    worker->onRegionMonitorChanged(0, 0, worker->m_registerKey);

    Dock::Position pos = Dock::Position::Bottom;
    Dock::DisplayMode dis = Dock::DisplayMode::Fashion;
    worker->getDockShowGeometry("", pos, dis);
    worker->getDockHideGeometry("", pos, dis);

    worker->checkXEventMonitorService();
    worker->showAniFinished();
    worker->hideAniFinished();
    worker->primaryScreenChanged();
    worker->onRequestUpdateFrontendGeometry();
    worker->isCopyMode();
    worker->onRequestUpdatePosition(Dock::Position::Top, Dock::Position::Bottom);
    worker->onAutoHideChanged(false);
    worker->onOpacityChanged(0.5);
    worker->onRequestDelayShowDock();

    delete worker;
    ASSERT_TRUE(true);
}

TEST_F(Test_MultiScreenWorker, onDisplayModeChanged)
{
    MainWindow window;
    MultiScreenWorker *worker = new MultiScreenWorker(&window, DWindowManagerHelper::instance());

    worker->onDisplayModeChanged(static_cast<DisplayMode>(0));
    worker->m_hideMode = HideMode::KeepShowing;
    worker->onDisplayModeChanged(static_cast<DisplayMode>(1));

    worker->onHideModeChanged(static_cast<HideMode>(0));
    worker->m_hideMode = HideMode::KeepShowing;
    worker->onHideModeChanged(static_cast<HideMode>(1));
    worker->m_hideMode = HideMode::KeepHidden;
    worker->onHideModeChanged(static_cast<HideMode>(3));

    worker->onHideStateChanged(static_cast<HideState>(0));
    worker->m_hideMode = HideMode::KeepShowing;
    worker->onHideStateChanged(static_cast<HideState>(1));
    worker->m_hideMode = HideMode::KeepHidden;
    worker->onHideStateChanged(static_cast<HideState>(2));

    delete worker;
    ASSERT_TRUE(true);
}

TEST_F(Test_MultiScreenWorker, displayAnimation_onRequestUpdateRegionMonitor)
{
    MainWindow window;
    MultiScreenWorker *worker = new MultiScreenWorker(&window, DWindowManagerHelper::instance());

    worker->m_position = Dock::Position::Left;
    worker->displayAnimation("primary", MultiScreenWorker::AniAction::Show);
    QTest::qWait(300);

    worker->m_position = Dock::Position::Top;
    worker->displayAnimation("primary", MultiScreenWorker::AniAction::Show);
    QTest::qWait(300);

    worker->m_position = Dock::Position::Bottom;
    worker->displayAnimation("primary", MultiScreenWorker::AniAction::Hide);
    QTest::qWait(300);

    worker->m_position = Dock::Position::Right;
    worker->displayAnimation("primary", MultiScreenWorker::AniAction::Hide);
    QTest::qWait(300);

    worker->m_position = Dock::Position::Top;
    worker->onRequestUpdateRegionMonitor();

    worker->m_position = Dock::Position::Bottom;
    worker->onRequestUpdateRegionMonitor();

    worker->m_position = Dock::Position::Left;
    worker->onRequestUpdateRegionMonitor();

    worker->m_position = Dock::Position::Right;
    worker->onRequestUpdateRegionMonitor();

    ASSERT_EQ(worker->parent(), &window);

    delete worker;
    ASSERT_TRUE(true);
}

TEST_F(Test_MultiScreenWorker, onTouchPress_onTouchRelease)
{
    MainWindow window;
    MultiScreenWorker *worker = new MultiScreenWorker(&window, DWindowManagerHelper::instance());

    QPoint p(0, 0);
    worker->rawXPosition(p);

    worker->onTouchPress(0, 0, 0, worker->m_touchRegisterKey);
    ASSERT_TRUE(worker->testState(MultiScreenWorker::RunState::TouchPress));
    worker->m_position = Dock::Position::Top;
    worker->onTouchRelease(0, 0, 100, worker->m_touchRegisterKey);
    ASSERT_FALSE(worker->testState(MultiScreenWorker::RunState::TouchPress));

    worker->onTouchPress(0, 0, 0, worker->m_touchRegisterKey);
    ASSERT_TRUE(worker->testState(MultiScreenWorker::RunState::TouchPress));
    worker->m_position = Dock::Position::Bottom;
    worker->onTouchRelease(0, 0, 100, worker->m_touchRegisterKey);
    ASSERT_FALSE(worker->testState(MultiScreenWorker::RunState::TouchPress));

    worker->onTouchPress(0, 0, 0, worker->m_touchRegisterKey);
    ASSERT_TRUE(worker->testState(MultiScreenWorker::RunState::TouchPress));
    worker->m_position = Dock::Position::Left;
    worker->onTouchRelease(0, 0, 100, worker->m_touchRegisterKey);
    ASSERT_FALSE(worker->testState(MultiScreenWorker::RunState::TouchPress));

    worker->onTouchPress(0, 0, 0, worker->m_touchRegisterKey);
    ASSERT_TRUE(worker->testState(MultiScreenWorker::RunState::TouchPress));
    worker->m_position = Dock::Position::Right;
    worker->onTouchRelease(0, 0, 100, worker->m_touchRegisterKey);
    ASSERT_FALSE(worker->testState(MultiScreenWorker::RunState::TouchPress));

    delete worker;
    ASSERT_TRUE(true);
}

TEST_F(Test_MultiScreenWorker, onDelayAutoHideChanged)
{
    MainWindow window;
    MultiScreenWorker *worker = new MultiScreenWorker(&window, DWindowManagerHelper::instance());

    worker->m_hideMode = HideMode::SmartHide;
    worker->m_hideState = HideState::Show;
    worker->onDelayAutoHideChanged();
    worker->m_hideState = HideState::Hide;
    worker->onDelayAutoHideChanged();

    worker->m_hideMode = HideMode::KeepShowing;
    worker->onDelayAutoHideChanged();

    worker->m_hideMode = HideMode::KeepHidden;
    worker->onDelayAutoHideChanged();

    delete worker;
    ASSERT_TRUE(true);
}

TEST_F(Test_MultiScreenWorker, onPositionChanged)
{
    MainWindow window;
    MultiScreenWorker *worker = new MultiScreenWorker(&window, DWindowManagerHelper::instance());

    worker->m_hideMode = HideMode::KeepHidden;
    worker->onPositionChanged(static_cast<Position>(0));
    QTest::qWait(400);

    worker->m_hideMode = HideMode::KeepShowing;
    worker->onPositionChanged(static_cast<Position>(1));
    QTest::qWait(400);

    worker->m_hideMode = HideMode::KeepHidden;
    worker->onPositionChanged(static_cast<Position>(2));
    QTest::qWait(400);

    worker->m_hideMode = HideMode::KeepShowing;
    worker->onPositionChanged(static_cast<Position>(3));
    QTest::qWait(400);

    delete worker;
    ASSERT_TRUE(true);
}

TEST_F(Test_MultiScreenWorker, reInitDisplayData)
{
    MainWindow window;
    MultiScreenWorker *worker = new MultiScreenWorker(&window, DWindowManagerHelper::instance());

    worker->reInitDisplayData();

    delete worker;
    ASSERT_TRUE(true);
}

TEST_F(Test_MultiScreenWorker, onRequestUpdateMonitorInfo)
{
    MainWindow window;
    MultiScreenWorker *worker = new MultiScreenWorker(&window, DWindowManagerHelper::instance());

    worker->onRequestUpdateMonitorInfo();

    delete worker;
    ASSERT_TRUE(true);
}

TEST_F(Test_MultiScreenWorker, updateParentGeometry)
{
    MainWindow window;
    MultiScreenWorker *worker = new MultiScreenWorker(&window, DWindowManagerHelper::instance());

    worker->updateParentGeometry(QRect(0, 0, 10, 10), Position::Top);
    worker->updateParentGeometry(QRect(0, 0, 10, 10), Position::Bottom);
    worker->updateParentGeometry(QRect(0, 0, 10, 10), Position::Left);
    worker->updateParentGeometry(QRect(0, 0, 10, 10), Position::Right);

    delete worker;
    ASSERT_TRUE(true);
}

TEST_F(Test_MultiScreenWorker, onWindowSizeChanged)
{
    MainWindow window;
    MultiScreenWorker *worker = new MultiScreenWorker(&window, DWindowManagerHelper::instance());

    worker->onWindowSizeChanged(-1);

    delete worker;
    ASSERT_TRUE(true);
}

TEST_F(Test_MultiScreenWorker, updateDisplay)
{
    MainWindow window;
    MultiScreenWorker *worker = new MultiScreenWorker(&window, DWindowManagerHelper::instance());

    worker->updateDisplay();

    delete worker;
    ASSERT_TRUE(true);
}

TEST_F(Test_MultiScreenWorker, onExtralRegionMonitorChanged)
{
    MainWindow window;
    MultiScreenWorker *worker = new MultiScreenWorker(&window, DWindowManagerHelper::instance());

    worker->onExtralRegionMonitorChanged(0, 0, "test");

    worker->m_hideMode = HideMode::KeepShowing;
    worker->m_hideState = HideState::Show;
    worker->onExtralRegionMonitorChanged(0, 0, worker->m_registerKey);

    worker->m_hideMode = HideMode::SmartHide;
    worker->onExtralRegionMonitorChanged(0, 0, worker->m_registerKey);

    worker->m_hideMode = HideMode::KeepHidden;
    worker->m_hideState = HideState::Hide;
    worker->onExtralRegionMonitorChanged(0, 0, worker->m_registerKey);

    delete worker;
    ASSERT_TRUE(true);
}

TEST_F(Test_MultiScreenWorker, onRegionMonitorChanged)
{
    MainWindow window;
    MultiScreenWorker *worker = new MultiScreenWorker(&window, DWindowManagerHelper::instance());

    worker->onRegionMonitorChanged(0, 0, "test");

    delete worker;
    ASSERT_TRUE(true);
}

TEST_F(Test_MultiScreenWorker, dockScreen)
{
    DockScreen ds("primary");

    ds.updateDockedScreen("screen1");

    ASSERT_EQ(ds.current(), "screen1");
    ASSERT_EQ(ds.last(), "primary");
    ASSERT_EQ(ds.primary(), "primary");

    ds.updatePrimary("screen2");
    ASSERT_EQ(ds.primary(), "screen2");
}

TEST_F(Test_MultiScreenWorker, screenworker_test3)
{
    MainWindow window;
    MultiScreenWorker *worker = new MultiScreenWorker(&window, DWindowManagerHelper::instance());

    worker->resetDockScreen();

    delete worker;
    ASSERT_TRUE(true);
}
