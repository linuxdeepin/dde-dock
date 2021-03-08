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

#define private public
#include "appitem.h"
#include "qgsettingsinterfacemock.h"
#undef private

class Test_AppItem : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    AppItem *appItem = nullptr;
};

void Test_AppItem::SetUp()
{
    appItem = new AppItem(QDBusObjectPath("/com/deepin/dde/daemon/Dock/entries/e0T6045b766"), QGSettingsInterface::Type::MockType);
}

void Test_AppItem::TearDown()
{
    delete appItem;
    appItem = nullptr;
}

TEST_F(Test_AppItem, coverage_test)
{
    appItem->checkEntry();
//    ASSERT_FALSE(appItem->isValid());

    ASSERT_TRUE(appItem->itemType() == AppItem::App);

    appItem->setDockInfo(Dock::Position::Top, QRect(QPoint(0,0), QPoint(1920, 40)));

    ASSERT_EQ(appItem->accessibleName(), appItem->m_itemEntryInter->name());
}

TEST_F(Test_AppItem, AppItem_show_test)
{
    appItem->show();

    QThread::msleep(450);

    ASSERT_EQ(appItem->isVisible(), true);
}

TEST_F(Test_AppItem, AppItem_hide_test)
{
    appItem->hide();

    QThread::msleep(450);

    ASSERT_EQ(appItem->isVisible(), false);
}
