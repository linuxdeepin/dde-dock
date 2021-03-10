/*
 * Copyright (C) 2018 ~ 2028 Uniontech Technology Co., Ltd.
 *
 * Author:     chenjun <chenjun@uniontech.com>
 *
 * Maintainer: chenjun <chenjun@uniontech.com>
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

#include "monitor.h"

//因为GTest不能断言自定义结构数据，需要重载<<和==操作符
std::ostream & operator<<(std::ostream & os, const Monitor::DockPosition & dockPosition) {
    return os << "leftDock = "
              << dockPosition.leftDock
              << " rightDock = "
              << dockPosition.rightDock
              << "topDock = "
              << dockPosition.topDock
              << " bottomDock = "
              << dockPosition.bottomDock;
}

bool operator==(const Monitor::DockPosition & p1, const Monitor::DockPosition & p2) {
    return p1.leftDock == p2.leftDock && p1.rightDock == p2.rightDock && p1.topDock ==  p2.topDock && p1.bottomDock == p2.bottomDock;
}

class Test_Monitor : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    Monitor *monitor = nullptr;
};

void Test_Monitor::SetUp()
{
    monitor = new Monitor();
}

void Test_Monitor::TearDown()
{
    delete monitor;
    monitor = nullptr;
}

TEST_F(Test_Monitor, monitor_test)
{
    ASSERT_NE(monitor, nullptr);

    int x = 10;
    int y = 10;
    int w = 100;
    int h = 100;

    monitor->setX(x);
    QCOMPARE(monitor->x(), x);
    monitor->setX(x);

    monitor->setY(y);
    QCOMPARE(monitor->y(), y);
    monitor->setY(y);

    monitor->setW(w);
    QCOMPARE(monitor->w(), w);
    monitor->setW(w);

    monitor->setH(h);
    QCOMPARE(monitor->h(), h);
    monitor->setH(h);

    QCOMPARE(monitor->left(), x);
    QCOMPARE(monitor->right(), x + w);
    QCOMPARE(monitor->top(), y);
    QCOMPARE(monitor->bottom(), y + h);

    QCOMPARE(monitor->topLeft(), QPoint(x, y));
    QCOMPARE(monitor->topRight(), QPoint(x + w, y));
    QCOMPARE(monitor->bottomLeft(), QPoint(x, y + h));
    QCOMPARE(monitor->bottomRight(), QPoint(x + w, y + h));
    QCOMPARE(monitor->rect(), QRect(x, y, w, h));

    QString name = "MonitorTestName";
    monitor->setName(name);
    QCOMPARE(monitor->name(), name);

    QString path = "testPath";
    monitor->setPath(path);
    QCOMPARE(monitor->path(), path);

    bool monitorEnable = true;
    monitor->setMonitorEnable(monitorEnable);
    QCOMPARE(monitor->enable(), monitorEnable);
    monitor->setMonitorEnable(monitorEnable);

    Monitor::DockPosition dockPosition;
    dockPosition.leftDock = true;
    dockPosition.rightDock = true;
    dockPosition.topDock = true;
    dockPosition.bottomDock = true;
    monitor->setDockPosition(dockPosition);
    QCOMPARE(monitor->dockPosition(), dockPosition);
}

TEST_F(Test_Monitor, dockPosition_test)
{
    monitor->setDockPosition(Monitor::DockPosition(false, false, false, false));
    ASSERT_FALSE(monitor->dockPosition().docked(Position::Top));
    ASSERT_FALSE(monitor->dockPosition().docked(Position::Bottom));
    ASSERT_FALSE(monitor->dockPosition().docked(Position::Left));
    ASSERT_FALSE(monitor->dockPosition().docked(Position::Right));

    monitor->dockPosition().reset();
    ASSERT_TRUE(monitor->dockPosition().docked(Position::Top));
    ASSERT_TRUE(monitor->dockPosition().docked(Position::Bottom));
    ASSERT_TRUE(monitor->dockPosition().docked(Position::Left));
    ASSERT_TRUE(monitor->dockPosition().docked(Position::Right));
}
