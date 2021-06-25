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
#include <QThread>
#include <QTest>
#include <QSignalSpy>
#include <QThread>

#include <gtest/gtest.h>

#define private public
#include "appsnapshot.h"
#undef private

class Test_AppSnapshot : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    AppSnapshot *shot = nullptr;
};

void Test_AppSnapshot::SetUp()
{
    shot = new AppSnapshot(1000000);
}

void Test_AppSnapshot::TearDown()
{
    delete shot;
    shot = nullptr;
}
