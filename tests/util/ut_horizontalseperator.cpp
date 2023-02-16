// Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
