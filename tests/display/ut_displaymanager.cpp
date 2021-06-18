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
#include <QScreen>
#include <QSignalSpy>

#include <gtest/gtest.h>

#define private public
#include "displaymanager.h"
#undef private

using namespace ::testing;

class Test_DisplayManager : public ::testing::Test
{
};

TEST_F(Test_DisplayManager, method_test)
{
    ASSERT_EQ(DisplayManager::instance()->screens().count(), qApp->screens().count());

    for (auto s : qApp->screens()) {
        ASSERT_TRUE(DisplayManager::instance()->screen(s->name()));
    }

    ASSERT_FALSE(DisplayManager::instance()->screen("testname"));

    ASSERT_EQ(DisplayManager::instance()->primary(), qApp->primaryScreen() ? qApp->primaryScreen()->name() : QString());

    // 第一次启动的时候，默认发出一次信号
    QSignalSpy spy(DisplayManager::instance(), &DisplayManager::screenInfoChanged);
    QTest::qWait(100);
    ASSERT_EQ(spy.count(), 1);
}

TEST_F(Test_DisplayManager, coverage_test) // 提高覆盖率,还没想好怎么做这种
{
    DisplayManager::instance()->onGSettingsChanged("onlyShowPrimary");
}
