// Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include <QObject>
#include <QTest>
#include <QEnterEvent>

#include <gtest/gtest.h>

#include "previewcontainer.h"

class Test_PreviewContainer : public ::testing::Test
{};

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

    WId id(1000);
    AppSnapshot *snap = new AppSnapshot(id);
    container->m_snapshots.insert(id, snap);
    snap->requestCloseAppSnapshot();
    container->setWindowInfos(map, QList<quint32> () << 1 << 2 << 3 << 4);

    for (const WId snapId: map.keys()) {
        container->appendSnapWidget(snapId);
    }

    container->previewEntered(id);
    container->m_waitForShowPreviewTimer->start();

    container->updateLayoutDirection(Dock::Position::Bottom);
    ASSERT_EQ(container->m_windowListLayout->direction(), container->m_wmHelper->hasComposite() ? QBoxLayout::LeftToRight : QBoxLayout::TopToBottom);
    container->updateLayoutDirection(Dock::Position::Top);
    ASSERT_EQ(container->m_windowListLayout->direction(), container->m_wmHelper->hasComposite() ? QBoxLayout::LeftToRight : QBoxLayout::TopToBottom);
    container->updateLayoutDirection(Dock::Position::Left);
    ASSERT_EQ(container->m_windowListLayout->direction(), QBoxLayout::TopToBottom);
    container->updateLayoutDirection(Dock::Position::Right);
    ASSERT_EQ(container->m_windowListLayout->direction(), QBoxLayout::TopToBottom);

    QEnterEvent event(QPoint(10,10), QPoint(100, 100), QPoint(100, 100));
    qApp->sendEvent(container, &event);

    QEvent event2(QEvent::Leave);
    qApp->sendEvent(container, &event2);

    QMimeData mimeData;
    mimeData.setText("test");
    QDragEnterEvent dragEnterEvent(QPoint(10, 10), Qt::CopyAction, &mimeData, Qt::LeftButton, Qt::NoModifier);
    qApp->sendEvent(container, &dragEnterEvent);

    container->prepareHide();
    container->adjustSize(true);
    container->adjustSize(false);

    delete snap;
    delete container;
    ASSERT_TRUE(true);
}

TEST_F(Test_PreviewContainer, checkMouseLeave)
{
    PreviewContainer container;
    container.checkMouseLeave();
    ASSERT_TRUE(true);
}

TEST_F(Test_PreviewContainer, dragLeaveEvent)
{
    PreviewContainer container;
    QDragLeaveEvent dragLeaveEvent_;
    container.dragLeaveEvent(&dragLeaveEvent_);
    ASSERT_TRUE(true);
}
TEST_F(Test_PreviewContainer, previewFloating)
{
    PreviewContainer container;
    container.previewFloating();

    ASSERT_TRUE(true);
}

TEST_F(Test_PreviewContainer, event_test)
{
    PreviewContainer *container = new PreviewContainer();

    QMouseEvent event1(QEvent::MouseButtonPress, QPointF(0, 0), Qt::LeftButton, Qt::RightButton, Qt::ControlModifier);
    container->mousePressEvent(&event1);

    QMouseEvent event2(QEvent::MouseButtonRelease, QPointF(0, 0), Qt::RightButton, Qt::RightButton, Qt::ControlModifier);
    container->mouseReleaseEvent(&event2);

    QMouseEvent event3(QEvent::MouseMove, QPointF(0, 0), Qt::LeftButton, Qt::LeftButton, Qt::ControlModifier);
    container->mouseMoveEvent(&event3);

    QMouseEvent event4(QEvent::MouseMove, QPointF(0, 0), Qt::RightButton, Qt::RightButton, Qt::ControlModifier);
    container->mouseMoveEvent(&event4);

    QResizeEvent event5((QSize()), QSize());
    container->resizeEvent(&event5);

    QEvent event6(QEvent::Leave);
    container->leaveEvent(&event6);

    QShowEvent event7;
    container->showEvent(&event7);

    QMimeData *data = new QMimeData;
    data->setText("test");
    QDropEvent event8(QPointF(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::ControlModifier);
    container->dropEvent(&event8);

    QDragEnterEvent event9(QPoint(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::NoModifier);
    container->dragEnterEvent(&event9);

    QDragMoveEvent event10(QPoint(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::NoModifier);
    container->dragMoveEvent(&event10);

    delete container;
    data->deleteLater();
    ASSERT_TRUE(true);
}
