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
#include "themeappicon.h"

#include <QPixmap>
#include <QDebug>
#include <QApplication>
#include <QFile>

#include <gtest/gtest.h>

class Ut_ThemeAppIcon : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;
};

void Ut_ThemeAppIcon::SetUp()
{
    ThemeAppIcon icon;
}

void Ut_ThemeAppIcon::TearDown()
{
}

TEST_F(Ut_ThemeAppIcon, getIcon_test1)
{
    QPixmap pix;

    // 无效图标
    QString name = "123";
    ThemeAppIcon::getIcon(name);

    // 有效图标
    ThemeAppIcon::getIcon("dde-calendar");

    ThemeAppIcon::getIcon(pix, "", 50);
    ThemeAppIcon::getIcon(pix, "", 50, true);

    // 获取base64编码的png图片数据
    const QString &iconName = "data:image/png;base64,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";
    ThemeAppIcon::getIcon(pix, iconName, 60);
}

TEST_F(Ut_ThemeAppIcon, getIcon_test2)
{
    QPixmap pix;
    ThemeAppIcon::getIcon(pix, "dde-calendar", 50);
    ASSERT_TRUE(true);
}

TEST_F(Ut_ThemeAppIcon, getIcon_test3)
{
    QPixmap pix;
    ThemeAppIcon::getIcon(pix, "data:image/test", 50);
    ASSERT_FALSE(pix.isNull());
}

TEST_F(Ut_ThemeAppIcon, getIcon_test4)
{
    QPixmap pix;
    ThemeAppIcon::getIcon(pix, ":/res/all_settings_on.png", 50);
    ASSERT_FALSE(pix.isNull());
}

TEST_F(Ut_ThemeAppIcon, createCalendarIcon_test)
{
    const QString &filePath = "/tmp/calendar.svg";
    ASSERT_TRUE(ThemeAppIcon::createCalendarIcon(filePath));
    QFile::remove(filePath);
}
