/*
 * Copyright (C) 2018 ~ 2028 Uniontech Technology Co., Ltd.
 *
 * Author:     weizhixiang <weizhixiang@uniontech.com>
 *
 * Maintainer: weizhixiang <weizhixiang@uniontech.com>
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
#include "mainpanelcontrol.h"
#undef private

using namespace ::testing;

class Test_MainPanelControl : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    MainPanelControl *mainPanel = nullptr;
};

void Test_MainPanelControl::SetUp()
{
    mainPanel = new MainPanelControl();
}

void Test_MainPanelControl::TearDown()
{
    delete mainPanel;
    mainPanel = nullptr;
}

TEST_F(Test_MainPanelControl, coverage_test)
{
    ASSERT_TRUE(mainPanel);

    mainPanel->setPositonValue(Dock::Position::Top);
    mainPanel->updateMainPanelLayout();
    QTest::qWait(10);

    mainPanel->setPositonValue(Dock::Position::Bottom);
    mainPanel->updateMainPanelLayout();
    QTest::qWait(10);

    mainPanel->setPositonValue(Dock::Position::Left);
    mainPanel->updateMainPanelLayout();
    QTest::qWait(10);

    mainPanel->setPositonValue(Dock::Position::Right);
    mainPanel->updateMainPanelLayout();
    QTest::qWait(10);
}
