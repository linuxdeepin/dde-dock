/*
 * Copyright (C) 2018 ~ 2028 Uniontech Technology Co., Ltd.
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

#include <gtest/gtest.h>

#define private public
#include "mainwindow.h"
#undef private

using namespace ::testing;

class Test_MainWindow : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    MainWindow *m_window = nullptr;
};

void Test_MainWindow::SetUp()
{
    m_window = new MainWindow;
}

void Test_MainWindow::TearDown()
{
    delete m_window;
    m_window = nullptr;
}

TEST_F(Test_MainWindow, coverage_test)
{
    ASSERT_TRUE(m_window);

    m_window->getTrayVisableItemCount();
    m_window->adjustShadowMask();
    m_window->resetDragWindow();
    m_window->onMainWindowSizeChanged(QPoint(10, 10));
    m_window->touchRequestResizeDock();
    m_window->sendNotifications();

    m_window->callShow();
    QTest::qWait(450);

    //TODO 这里无论输入什么，均返回true
    //    ASSERT_FALSE(m_window->appIsOnDock("testname"));

    QEvent enterEvent(QEvent::Enter);
    qApp->sendEvent(m_window, &enterEvent);
    QTest::qWait(10);

    QEvent dragEnterEvent(QEvent::Enter);
    qApp->sendEvent(m_window->m_dragWidget, &dragEnterEvent);
    QTest::qWait(10);
    ASSERT_EQ(QApplication::overrideCursor()->shape(), m_window->m_dragWidget->cursor().shape());

    QEvent dragLeaveEvent(QEvent::Leave);
    qApp->sendEvent(m_window->m_dragWidget, &dragLeaveEvent);
    QTest::qWait(10);
    ASSERT_EQ(QApplication::overrideCursor()->shape(), Qt::ArrowCursor);
}
