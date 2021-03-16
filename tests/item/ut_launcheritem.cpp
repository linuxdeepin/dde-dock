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
#include <gmock/gmock.h>

using namespace ::testing;

#define private public
#include "launcheritem.h"
#undef private

#include "mock/qgsettingsmock.h"

class Test_LauncherItem : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;
};

void Test_LauncherItem::SetUp()
{
}

void Test_LauncherItem::TearDown()
{
}

TEST_F(Test_LauncherItem, launcher_test)
{
    LauncherItem *launcherItem = new LauncherItem;

    ASSERT_EQ(launcherItem->itemType(), LauncherItem::Launcher);
    launcherItem->refreshIcon();
    launcherItem->show();
    QThread::msleep(10);

    launcherItem->hide();
    QThread::msleep(10);

    launcherItem->resize(100,100);
    launcherItem->popupTips();

    QTest::mouseClick(launcherItem, Qt::LeftButton, Qt::NoModifier, launcherItem->geometry().center());

    delete launcherItem;
    launcherItem = nullptr;
}
