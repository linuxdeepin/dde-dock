// Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include <QObject>

#include <gtest/gtest.h>

#include "appdrag.h"

class Test_AppDrag : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;
};

void Test_AppDrag::SetUp()
{
}

void Test_AppDrag::TearDown()
{
}

TEST_F(Test_AppDrag, coverage_test)
{
    QWidget w;
    AppDrag drag(&w);
    QPixmap pix(":/res/all_settings_on.png");
    drag.setPixmap(pix);

    ASSERT_TRUE(drag.appDragWidget());

//    drag->exec();
//    drag->exec(Qt::MoveAction, Qt::IgnoreAction);
}
