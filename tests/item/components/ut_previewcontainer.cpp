/*
 * Copyright (C) 2018 ~ 2028 Uniontech Technology Co., Ltd.
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
#include <QTest>
#include <QEnterEvent>

#include <gtest/gtest.h>

#define private public
#include "previewcontainer.h"
#undef private

class Test_PreviewContainer : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;
};

void Test_PreviewContainer::SetUp()
{
}

void Test_PreviewContainer::TearDown()
{
}

TEST_F(Test_PreviewContainer, coverage_test)
{
    PreviewContainer *container = new PreviewContainer();

    WindowInfoMap map;

    WindowInfo info;
    info.attention = true;
    info.title = "test1";

    map.insert(1, info);
    map.insert(2, info);
    map.insert(3, info);

    container->setWindowInfos(map, map.keys());

    for (const WId id: map.keys()) {
        container->appendSnapWidget(id);
    }

    container->updateSnapshots();
    container->updateLayoutDirection(Dock::Position::Bottom);
    ASSERT_EQ(container->m_windowListLayout->direction(), QBoxLayout::LeftToRight);
    container->updateLayoutDirection(Dock::Position::Top);
    ASSERT_EQ(container->m_windowListLayout->direction(), QBoxLayout::LeftToRight);
    container->updateLayoutDirection(Dock::Position::Left);
    ASSERT_EQ(container->m_windowListLayout->direction(), QBoxLayout::TopToBottom);
    container->updateLayoutDirection(Dock::Position::Right);
    ASSERT_EQ(container->m_windowListLayout->direction(), QBoxLayout::TopToBottom);

    QEnterEvent event(QPoint(10,10), QPoint(100, 100), QPoint(100, 100));
    qApp->sendEvent(container, &event);

    QEvent event2(QEvent::Leave);
    qApp->sendEvent(container, &event);

    QMimeData mimeData;
    mimeData.setText("test");
    QDragEnterEvent dragEnterEvent(QPoint(10, 10), Qt::CopyAction, &mimeData, Qt::LeftButton, Qt::NoModifier);
    qApp->sendEvent(container, &dragEnterEvent);

//    QDragLeaveEvent dragLeaveEvent;
//    qApp->sendEvent(container, &dragLeaveEvent);
}
