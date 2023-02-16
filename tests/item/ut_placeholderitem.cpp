// Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
