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
#define private public
#include "touchsignalmanager.h"
#undef private

#include <QTest>
#include <QSignalSpy>

#include <gtest/gtest.h>

class Ut_TouchSignalManager : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;
};

void Ut_TouchSignalManager::SetUp()
{
}

void Ut_TouchSignalManager::TearDown()
{
}

TEST_F(Ut_TouchSignalManager, isDragIconPress_test)
{
    QCOMPARE(TouchSignalManager::instance()->isDragIconPress(), false);
    TouchSignalManager::instance()->m_dragIconPressed = true;
    QCOMPARE(TouchSignalManager::instance()->isDragIconPress(), true);
    TouchSignalManager::instance()->m_dragIconPressed = false;
    QCOMPARE(TouchSignalManager::instance()->isDragIconPress(), false);
}

TEST_F(Ut_TouchSignalManager, dealShortTouchPress_test)
{
    QSignalSpy spy(TouchSignalManager::instance(), SIGNAL(shortTouchPress(int, double, double)));
    TouchSignalManager::instance()->dealShortTouchPress(1, 0, 0);
    QCOMPARE(spy.count(), 1);
    QCOMPARE(TouchSignalManager::instance()->isDragIconPress(), true);

    const QList<QVariant> &arguments = spy.takeFirst();
    QCOMPARE(arguments.size(), 3);
    QCOMPARE(arguments.at(0), 1);
    ASSERT_TRUE(qAbs(arguments.at(1).toDouble()) < 0.00001);
    ASSERT_TRUE(qAbs(arguments.at(2).toDouble()) < 0.00001);
}

TEST_F(Ut_TouchSignalManager, dealTouchRelease_test)
{
    QSignalSpy spy(TouchSignalManager::instance(), SIGNAL(touchRelease(double, double)));
    TouchSignalManager::instance()->dealTouchRelease(0, 0);
    QCOMPARE(spy.count(), 1);
    QCOMPARE(TouchSignalManager::instance()->isDragIconPress(), false);

    const QList<QVariant> &arguments = spy.takeFirst();
    QCOMPARE(arguments.size(), 2);
    ASSERT_TRUE(qAbs(arguments.at(0).toDouble()) < 0.00001);
    ASSERT_TRUE(qAbs(arguments.at(1).toDouble()) < 0.00001);
}

TEST_F(Ut_TouchSignalManager, dealTouchPress_test)
{
    QSignalSpy spy(TouchSignalManager::instance(), SIGNAL(middleTouchPress(double, double)));
    TouchSignalManager::instance()->dealTouchPress(1, 1000, 0, 0);
    QCOMPARE(spy.count(), 1);
    const QList<QVariant> &arguments = spy.takeFirst();
    QCOMPARE(arguments.size(), 2);
    ASSERT_TRUE(qAbs(arguments.at(0).toDouble()) < 0.00001);
    ASSERT_TRUE(qAbs(arguments.at(1).toDouble()) < 0.00001);

    TouchSignalManager::instance()->dealTouchPress(1, 2000, 0, 0);
    QCOMPARE(spy.count(), 0);
    TouchSignalManager::instance()->dealTouchPress(1, 500, 0, 0);
    QCOMPARE(spy.count(), 0);
    TouchSignalManager::instance()->dealTouchPress(2, 0000, 0, 0);
    QCOMPARE(spy.count(), 0);
}
