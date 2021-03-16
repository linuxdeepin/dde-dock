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
#include "themeappicon.h"

#include <QPixmap>
#include <QDebug>
#include <QApplication>

#include <gtest/gtest.h>

class Ut_ThemeAppIcon : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;
};

void Ut_ThemeAppIcon::SetUp()
{
}

void Ut_ThemeAppIcon::TearDown()
{
}

TEST_F(Ut_ThemeAppIcon, getIcon_test)
{
    ThemeAppIcon appIcon;
    const QPixmap &pix1 = appIcon.getIcon("", 50, 1.0);
    ASSERT_FALSE(pix1.isNull());
    appIcon.getIcon("dde-calendar", 50, 1.0);
    const QPixmap &pix2 = appIcon.getIcon("data:image/test", 50, 1.0);
    ASSERT_FALSE(pix2.isNull());
    const QPixmap &pix3 = appIcon.getIcon(":/res/all_settings_on.png", 50, 1.0);
    ASSERT_FALSE(pix3.isNull());
}
