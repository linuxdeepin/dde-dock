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
#include <QPaintEvent>

#include <gtest/gtest.h>

#include "horizontalseperator.h"

class Test_HorizontalSeperator : public QObject, public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

};

void Test_HorizontalSeperator::SetUp()
{
}

void Test_HorizontalSeperator::TearDown()
{
}

TEST_F(Test_HorizontalSeperator, coverage_test)
{
    HorizontalSeperator seperator;
    ASSERT_EQ(seperator.sizeHint().height(), 2);

    seperator.show();

    ASSERT_TRUE(true);
}

TEST_F(Test_HorizontalSeperator, paintEvent)
{
    HorizontalSeperator seperator;

    QRect rect(0, 0, 10, 10);
    QPaintEvent e(rect);
    seperator.paintEvent(&e);

    ASSERT_TRUE(true);

}
