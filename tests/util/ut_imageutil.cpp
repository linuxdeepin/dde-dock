// Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include <QObject>
#include <QApplication>
#include <QSignalSpy>
#include <QTest>

#include <gtest/gtest.h>

#include "imageutil.h"

class Test_ImageUtil : public QObject, public ::testing::Test
{};

TEST_F(Test_ImageUtil, coverage_test)
{
    ASSERT_TRUE(ImageUtil::loadSvg("test", QSize(100, 100), 1.5).isNull());
    ASSERT_FALSE(ImageUtil::loadSvg(":/res/dde-calendar.svg", QSize(100, 100), 1.5).isNull());
    ASSERT_EQ(ImageUtil::loadSvg(":/res/dde-calendar.svg", "dde-printer", 100, 1.25).size(), QSize(125, 125));
    ASSERT_EQ(ImageUtil::loadSvg("123", "456", 100, 1.25).size(), QSize(125, 125));
}
