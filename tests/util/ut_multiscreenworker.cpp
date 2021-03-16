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

#include <DWindowManagerHelper>

#include <gtest/gtest.h>

#include "mainwindow.h"
#include "multiscreenworker.h"

class Test_MultiScreenWorker : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    MainWindow *mainwindow;
    MultiScreenWorker *worker = nullptr;
};

void Test_MultiScreenWorker::SetUp()
{
//    mainwindow = new MainWindow();
//    worker = new MultiScreenWorker(mainwindow, DWindowManagerHelper::instance());
}

void Test_MultiScreenWorker::TearDown()
{
//    delete worker;
//    worker = nullptr;
}

TEST_F(Test_MultiScreenWorker, dockInter_test)
{
//    ASSERT_TRUE(worker->dockInter());
}

