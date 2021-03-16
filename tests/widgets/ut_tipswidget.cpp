/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     weizhixiang <weizhixiang@uniontech.com>
 *
 * Maintainer: weizhixiang <weizhixiang@uniontech.com>
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
#include <QDebug>
#include <QApplication>
#include <QTest>

#include <gtest/gtest.h>

#include "../widgets/tipswidget.h"

namespace Dock {
class Test_TipsWidget : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    TipsWidget *tipsWidget;
};

void Test_TipsWidget::SetUp()
{
    tipsWidget = new TipsWidget();
}

void Test_TipsWidget::TearDown()
{
    delete tipsWidget;
    tipsWidget = nullptr;
}

TEST_F(Test_TipsWidget, setText_test)
{
    const QString text = "hello dde dock";
    tipsWidget->setText(text);
    ASSERT_EQ(text, tipsWidget->text());

    tipsWidget->show();
    QTest::qWait(10);

    QEvent event(QEvent::FontChange);
    qApp->sendEvent(tipsWidget, &event);
    QTest::qWait(10);
}

TEST_F(Test_TipsWidget, setTextList_test)
{
    const QStringList textList = {
        "hello",
        "dde",
        "dock"
    };
    tipsWidget->setTextList(textList);
    ASSERT_EQ(textList, tipsWidget->textList());

    tipsWidget->show();
    QTest::qWait(10);

    QEvent event(QEvent::FontChange);
    qApp->sendEvent(tipsWidget, &event);
    QTest::qWait(10);
}
}
