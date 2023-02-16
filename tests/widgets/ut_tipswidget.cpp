// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include <QObject>
#include <QDebug>
#include <QApplication>
#include <QPaintEvent>
#include <QTest>

#include <gtest/gtest.h>

#define protected public
#include "../widgets/tipswidget.h"
#undef protected

using namespace ::testing;
using namespace Dock;

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

TEST_F(Test_TipsWidget, setText)
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

TEST_F(Test_TipsWidget, setTextList)
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

TEST_F(Test_TipsWidget, paintEvent)
{
    QPaintEvent paintEvent((QRect()));
    tipsWidget->paintEvent(&paintEvent);

    QTest::qWait(10);
}
