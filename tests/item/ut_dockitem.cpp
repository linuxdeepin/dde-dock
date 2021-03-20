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
#include <gmock/gmock.h>

#include "dockitem.h"
class Test_DockItem : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    DockItem *dockItem = nullptr;
};

void Test_DockItem::SetUp()
{
    dockItem = new DockItem();
}

void Test_DockItem::TearDown()
{
    delete dockItem;
    dockItem = nullptr;
}

TEST_F(Test_DockItem, dockitem_test)
{
    ASSERT_NE(dockItem, nullptr);
}

TEST_F(Test_DockItem, dockitem_show_test)
{
    dockItem->show();

    QThread::msleep(450);

    ASSERT_EQ(dockItem->isVisible(), true);
}

TEST_F(Test_DockItem, dockitem_hide_test)
{
    dockItem->hide();

    QThread::msleep(450);

    ASSERT_EQ(dockItem->isVisible(), false);
}
