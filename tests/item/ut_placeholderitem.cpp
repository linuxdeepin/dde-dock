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
#include <QTest>

#include <gtest/gtest.h>

#include "placeholderitem.h"

class Test_PlaceholderItem : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    PlaceholderItem *placeholderitem = nullptr;
};

void Test_PlaceholderItem::SetUp()
{
    placeholderitem = new PlaceholderItem();
}

void Test_PlaceholderItem::TearDown()
{
    delete placeholderitem;
    placeholderitem = nullptr;
}

TEST_F(Test_PlaceholderItem, placeholder_test)
{
    QCOMPARE(placeholderitem->itemType(), PlaceholderItem::Placeholder);
}
