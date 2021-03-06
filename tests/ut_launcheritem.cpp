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

#include "launcheritem.h"
#include "qgsettingsinterfacemock.h"

class Test_LauncherItem : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    LauncherItem *launcherItem = nullptr;
};

void Test_LauncherItem::SetUp()
{
    launcherItem = new LauncherItem(new QGSettingsInterfaceMock("com.deepin.dde.dock.module.launcher"));
}

void Test_LauncherItem::TearDown()
{
    delete launcherItem;
    launcherItem = nullptr;
}

TEST_F(Test_LauncherItem, dockitem_test)
{
    ASSERT_NE(launcherItem, nullptr);
}
