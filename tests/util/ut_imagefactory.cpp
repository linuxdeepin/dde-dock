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

#include "imagefactory.h"

class Test_ImageFactory : public QObject, public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    ImageFactory *factory = nullptr;
};

void Test_ImageFactory::SetUp()
{
    factory = new ImageFactory();
}

void Test_ImageFactory::TearDown()
{
    delete factory;
    factory = nullptr;
}

TEST_F(Test_ImageFactory, factory_test)
{
    QPixmap pix(":/res/all_settings_on.png");
    // 以下是无效值，应该屏蔽才对
    factory->lighterEffect(pix, -1);
    factory->lighterEffect(pix, 256);

    // 传入空的pixmap对象
    QPixmap emptyPix;
    factory->lighterEffect(emptyPix, 150);
}
