// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "multiwindowhelper.h"
#include "appmultiitem.h"
#include "appitem.h"

MultiWindowHelper::MultiWindowHelper(QWidget *appWidget, QWidget *multiWindowWidget, QObject *parent)
    : QObject(parent)
    , m_appWidget(appWidget)
    , m_multiWindowWidget(multiWindowWidget)
    , m_displayMode(Dock::DisplayMode::Efficient)
{
    m_appWidget->installEventFilter(this);
    m_multiWindowWidget->installEventFilter(this);
}

void MultiWindowHelper::setDisplayMode(Dock::DisplayMode displayMode)
{
    if (m_displayMode == displayMode)
        return;

    m_displayMode = displayMode;
    resetMultiItemPosition();
}

void MultiWindowHelper::addMultiWindow(int, AppMultiItem *item)
{
    int index = itemIndex(item);
    if (m_displayMode == Dock::DisplayMode::Efficient) {
        // 将多开窗口项目插入到对应的APP的后面
        insertChildWidget(m_appWidget, index, item);
    } else {
        // 将多开窗口插入到工具区域的前面
        insertChildWidget(m_multiWindowWidget, index, item);
    }
}

void MultiWindowHelper::removeMultiWindow(AppMultiItem *item)
{
    if (m_appWidget->children().contains(item))
        m_appWidget->layout()->removeWidget(item);
    else
        m_multiWindowWidget->layout()->removeWidget(item);
}

bool MultiWindowHelper::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == m_appWidget || watched == m_multiWindowWidget) {
        switch(event->type()) {
        case QEvent::ChildAdded:
        case QEvent::ChildRemoved: {
            /* 这里用异步的方式，因为收到QEvent::ChildAdded信号的时候，
             此时应用还没有插入到Widget中，收到QEvent::ChildRemoved信号的时候，
             此时应用还未从任务栏上移除，通过异步的方式保证同步新增或移除成功后才执行，这样更新的界面才是最准确的
            */
            QMetaObject::invokeMethod(this, &MultiWindowHelper::requestUpdate, Qt::QueuedConnection);
            break;
        }
        default:
            break;
        }
    }

    return QObject::eventFilter(watched, event);
}

int MultiWindowHelper::itemIndex(AppMultiItem *item)
{
    if (m_displayMode != Dock::DisplayMode::Efficient)
        return -1;

    // 高效模式，查找对应的应用或者这个应用所有的子窗口所在的位置，然后插入到最大的值的后面
    int lastIndex = -1;
    for (int i = 0; i < m_appWidget->layout()->count(); i++) {
        DockItem *dockItem = qobject_cast<DockItem *>(m_appWidget->layout()->itemAt(i)->widget());
        if (!dockItem)
            continue;

        if (dockItem != item->appItem()) {
            AppMultiItem *multiItem = qobject_cast<AppMultiItem *>(dockItem);
            if (!multiItem || multiItem->appItem() != item->appItem())
                continue;
        }

        lastIndex = i;
    }

    if (lastIndex >= 0)
        return ++lastIndex;

    return -1;
}

void MultiWindowHelper::insertChildWidget(QWidget *parentWidget, int index, AppMultiItem *item)
{
    QBoxLayout *layout = static_cast<QBoxLayout *>(parentWidget->layout());
    if (index >= 0)
        layout->insertWidget(index, item);
    else
        layout->addWidget(item);
}

void MultiWindowHelper::resetMultiItemPosition()
{
    QWidget *fromWidget = nullptr;
    QWidget *toWidget = nullptr;

    if (m_displayMode == Dock::DisplayMode::Efficient) {
        // 从时尚模式变换为高效模式
        fromWidget = m_multiWindowWidget;
        toWidget = m_appWidget;
    } else {
        // 从高效模式变换到时尚模式
        fromWidget = m_appWidget;
        toWidget = m_multiWindowWidget;
    }

    QList<AppMultiItem *> moveWidgetItem;
    for (int i = 0; i < fromWidget->layout()->count(); i++) {
        AppMultiItem *multiItem = qobject_cast<AppMultiItem *>(fromWidget->layout()->itemAt(i)->widget());
        if (!multiItem)
            continue;

        moveWidgetItem << multiItem;
    }

    QBoxLayout *toLayout = static_cast<QBoxLayout *>(toWidget->layout());
    for (AppMultiItem *item : moveWidgetItem) {
        fromWidget->layout()->removeWidget(item);
        int index = itemIndex(item);
        if (index >= 0)
            toLayout->insertWidget(index, item);
        else
            toLayout->addWidget(item);
    }
}
