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
#include <QApplication>
#include <QSignalSpy>
#include <QTest>

#include <gtest/gtest.h>

#include "imageutil.h"

class Test_ImageUtil : public QObject, public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

};

void Test_ImageUtil::SetUp()
{
}

void Test_ImageUtil::TearDown()
{
}

TEST_F(Test_ImageUtil, coverage_test)
{
    ASSERT_TRUE(ImageUtil::loadSvg("test", QSize(100, 100), 1.5).isNull());
    ASSERT_EQ(ImageUtil::loadSvg("dde-printer", ":/res/dde-calendar.svg", 100, 1.25).size(), QSize(125, 125));
    ASSERT_EQ(ImageUtil::loadSvg("123", "456", 100, 1.25).size(), QSize(125, 125));
}
