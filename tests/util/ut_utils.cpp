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
#include "utils.h"

#include <QTest>

#include <gtest/gtest.h>

class Ut_Utils : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;
};

void Ut_Utils::SetUp()
{
}

void Ut_Utils::TearDown()
{
}

TEST_F(Ut_Utils, comparePluginApi_test)
{
    QString v1("1.0.0");
    QString v2("1.0.1");
    QString v3("1.0.0.0");

    QCOMPARE(Utils::comparePluginApi(v1, v1), 0);
    QCOMPARE(Utils::comparePluginApi(v1, v2), -1);
    QCOMPARE(Utils::comparePluginApi(v2, v1), 1);
    QCOMPARE(Utils::comparePluginApi(v1, v3), -1);
    QCOMPARE(Utils::comparePluginApi(v3, v1), 1);
}

TEST_F(Ut_Utils, screenAt_test)
{
    Utils::screenAt(QPoint(0, 0));
    QCOMPARE(Utils::screenAt(QPoint(-1, -1)), nullptr);
}

TEST_F(Ut_Utils, screenAtByScaled_test)
{
    Utils::screenAtByScaled(QPoint(0, 0));
    QCOMPARE(Utils::screenAtByScaled(QPoint(-1, -1)), nullptr);
}

TEST_F(Ut_Utils, renderSVG_test)
{
    QPixmap pix(":/res/all_settings_on.png");
    const QSize &size = pix.size();

    ASSERT_TRUE(Utils::renderSVG("", size, 1.0).isNull());
    QCOMPARE(Utils::renderSVG(":/res/all_settings_on.png", size, 1.0).size(), size);
    QCOMPARE(Utils::renderSVG(":/res/all_settings_on.png", QSize(50, 50), 1.0).size(), QSize(50, 50));
    QCOMPARE(Utils::renderSVG(":/res/all_settings_on.png", QSize(50, 50), 0.5).size(), QSize(25, 25));
}

TEST_F(Ut_Utils, gsettings_test)
{
    ASSERT_FALSE(Utils::SettingValue("", "").isValid());
}
