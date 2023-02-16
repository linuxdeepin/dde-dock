// Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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

TEST_F(Test_DisplayManager, coverage_test)
{
    DisplayManager::instance()->onGSettingsChanged("onlyShowPrimary");
}
