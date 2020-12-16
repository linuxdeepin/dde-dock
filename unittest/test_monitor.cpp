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

TEST_F(Test_Monitor, dockitem_test)
{
    ASSERT_NE(monitor, nullptr);

    int x = 10;
    int y = 10;
    int w = 100;
    int h = 100;

    monitor->setX(x);
    ASSERT_EQ(monitor->x(), x);

    monitor->setY(y);
    ASSERT_EQ(monitor->y(), y);

    monitor->setW(w);
    ASSERT_EQ(monitor->w(), w);

    monitor->setH(h);
    ASSERT_EQ(monitor->h(), h);

    ASSERT_EQ(monitor->topLeft(), QPoint(x, y));
    ASSERT_EQ(monitor->topRight(), QPoint(x + w, y));
    ASSERT_EQ(monitor->bottomLeft(), QPoint(x, y + h));
    ASSERT_EQ(monitor->bottomRight(), QPoint(x + w, y + h));
    ASSERT_EQ(monitor->rect(), QRect(x, y, w, h));

    QString name = "MonitorTestName";
    monitor->setName(name);
    ASSERT_EQ(monitor->name(), name);

    QString path = "testPath";
    monitor->setPath(path);
    ASSERT_EQ(monitor->path(), path);

    bool monitorEnable = true;
    monitor->setMonitorEnable(monitorEnable);
    ASSERT_EQ(monitor->enable(), monitorEnable);

    Monitor::DockPosition dockPosition;
    dockPosition.leftDock = true;
    dockPosition.rightDock = true;
    dockPosition.topDock = true;
    dockPosition.bottomDock = true;
    monitor->setDockPosition(dockPosition);
    ASSERT_EQ(monitor->dockPosition(), dockPosition);
}
