/*
 * Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
 *
 * Author:     chenjun <chenjun@uniontech.com>
 *
 * Maintainer: chenjun <chenjun@uniontech.com>
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
#include <QSignalSpy>
#include <QTest>

#include <gtest/gtest.h>

#define private public
#include "appsnapshot.h"
#include "floatingpreview.h"
#undef private

class Test_FloatingPreview : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;
};

void Test_FloatingPreview::SetUp()
{
}

void Test_FloatingPreview::TearDown()
{
}

TEST_F(Test_FloatingPreview, view_test)
{
    QWidget *parent = new QWidget;
    FloatingPreview *view = new FloatingPreview(parent);
    AppSnapshot *shot = new AppSnapshot(1000);
    view->trackWindow(shot);

    ASSERT_TRUE(view->m_titleBtn->text() == shot->title());
    ASSERT_EQ(view->trackedWindow(), shot);
//    ASSERT_EQ(view->trackedWid(), shot->wid());

    QSignalSpy spy(shot, &AppSnapshot::clicked);
    QTest::mouseClick(view, Qt::LeftButton, Qt::NoModifier);
    ASSERT_EQ(spy.count(), 1);

    //    view->m_closeBtn3D->click();

    view->hide();
    ASSERT_TRUE(shot->contentsMargins() == QMargins(0, 0, 0, 0));

    view->trackWindow(nullptr);
    ASSERT_TRUE(view->m_titleBtn->text().isEmpty());
    ASSERT_EQ(view->trackedWindow(), shot);

    delete parent;
    parent = nullptr;
}

