// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "recentapphelper.h"
#include "dockitemmanager.h"
#include "appitem.h"

#include <QWidget>

#define ENTRY_NONE      0
#define ENTRY_NORMAL    1
#define ENTRY_RECENT    2

RecentAppHelper::RecentAppHelper(QWidget *appWidget, QWidget *recentWidget, QObject *parent)
    : QObject(parent)
    , m_appWidget(appWidget)
    , m_recentWidget(recentWidget)
{
    m_appWidget->installEventFilter(this);
    m_recentWidget->installEventFilter(this);
    connect(this, &RecentAppHelper::requestUpdateRecentVisible, this, &RecentAppHelper::updateRecentVisible, Qt::QueuedConnection);
}

void RecentAppHelper::setDisplayMode(Dock::DisplayMode displayMode)
{
    bool lastVisible = dockAppIsVisible();
    m_displayMode = displayMode;
    updateRecentVisible();
    updateDockAppVisible(lastVisible);
}

// 当在应用区域调整位置的时候，需要重新设置索引
void RecentAppHelper::resetAppInfo()
{

}

void RecentAppHelper::addAppItem(int index, DockItem *dockItem)
{
    if (appInRecent(dockItem)) {
        addRecentAreaItem(index, dockItem);
        updateRecentVisible();
    } else {
        bool lastVisible = dockAppIsVisible();
        addAppAreaItem(index, dockItem);
        updateDockAppVisible(lastVisible);
    }

    AppItem *appItem = qobject_cast<AppItem *>(dockItem);

    connect(appItem, &AppItem::modeChanged, this, &RecentAppHelper::onModeChanged);
}

void RecentAppHelper::removeAppItem(DockItem *dockItem)
{
    if (m_recentWidget->children().contains(dockItem))
        removeRecentAreaItem(dockItem);
    else
        removeAppAreaItem(dockItem);
}

bool RecentAppHelper::recentIsVisible() const
{
    return m_recentWidget->isVisible();
}

bool RecentAppHelper::dockAppIsVisible() const
{
    return (m_displayMode == Dock::DisplayMode::Efficient
            || m_appWidget->layout()->count() > 0);
}

bool RecentAppHelper::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == m_appWidget || watched == m_recentWidget) {
        switch(event->type()) {
        case QEvent::ChildAdded:
        case QEvent::ChildRemoved: {
            QMetaObject::invokeMethod(this, [ this ] {
                /* 这里用异步的方式，因为收到QEvent::ChildAdded信号的时候，
                 此时应用还没有插入到Widget中，收到QEvent::ChildRemoved信号的时候，
                 此时应用还未从任务栏上移除，通过异步的方式保证同步新增或移除成功后才执行，这样更新的界面才是最准确的
                */
                Q_EMIT requestUpdate();
            }, Qt::QueuedConnection);
        }
            break;
        default:
            break;
        }
    }

    return QObject::eventFilter(watched, event);
}

void RecentAppHelper::onModeChanged(int mode)
{
    AppItem *appItem = qobject_cast<AppItem *>(sender());
    if (!appItem)
        return;

    auto moveItemToWidget = [ = ](QWidget *widget) {
        int index = getEntryIndex(appItem, widget);
        removeAppItem(appItem);
        QBoxLayout *layout = static_cast<QBoxLayout *>(widget->layout());
        layout->insertWidget(index, appItem);
    };

    if (mode == ENTRY_NORMAL) {
        // 添加到应用区域
        moveItemToWidget(m_appWidget);
    } else if (mode == ENTRY_RECENT) {
        // 添加到最近打开应用区域
        moveItemToWidget(m_recentWidget);
    }
    updateRecentVisible();
}

bool RecentAppHelper::appInRecent(DockItem *item) const
{
    AppItem *appItem = qobject_cast<AppItem *>(item);
    if (!appItem)
        return false;

    return (appItem->mode() == ENTRY_RECENT);
}

void RecentAppHelper::addAppAreaItem(int index, DockItem *wdg)
{
    QBoxLayout *boxLayout = static_cast<QBoxLayout *>(m_appWidget->layout());
    boxLayout->insertWidget(index, wdg);
}

void RecentAppHelper::addRecentAreaItem(int index, DockItem *wdg)
{
    QBoxLayout *recentLayout = static_cast<QBoxLayout *>(m_recentWidget->layout());
    recentLayout->insertWidget(index, wdg);
}

void RecentAppHelper::updateRecentVisible()
{
    bool lastRecentVisible = m_recentWidget->isVisible();
    bool recentVisible = lastRecentVisible;

    if (m_displayMode == Dock::DisplayMode::Efficient) {
        // 如果是高效模式，不显示最近打开应用区域
        m_recentWidget->setVisible(false);
        recentVisible = false;
    } else {
        QBoxLayout *recentLayout = static_cast<QBoxLayout *>(m_recentWidget->layout());
        qInfo() << "recent Widget count:" << recentLayout->count() << ", app Widget count" << m_appWidget->layout()->count();
        // 如果是时尚模式，则判断当前打开应用数量是否为0，为0则不显示，否则显示
        recentVisible = (recentLayout->count() > 0);
        m_recentWidget->setVisible(recentVisible);
    }

    if (lastRecentVisible != recentVisible)
        Q_EMIT recentVisibleChanged(recentVisible);
}

void RecentAppHelper::updateDockAppVisible(bool lastVisible)
{
    bool visible = dockAppIsVisible();
    if (lastVisible != visible)
        Q_EMIT dockAppVisibleChanged(visible);
}

void RecentAppHelper::removeRecentAreaItem(DockItem *wdg)
{
    QBoxLayout *recentLayout = static_cast<QBoxLayout *>(m_recentWidget->layout());
    recentLayout->removeWidget(wdg);
    Q_EMIT requestUpdateRecentVisible();
}

void RecentAppHelper::removeAppAreaItem(DockItem *wdg)
{
    QBoxLayout *boxLayout = static_cast<QBoxLayout *>(m_appWidget->layout());
    bool lastVisible = dockAppIsVisible();
    boxLayout->removeWidget(wdg);
    updateDockAppVisible(lastVisible);
}

int RecentAppHelper::getEntryIndex(DockItem *dockItem, QWidget *widget) const
{
    AppItem *appItem = qobject_cast<AppItem *>(dockItem);
    if (!appItem)
        return -1;

    // 查找当前的应用在所有的应用中的排序
    QStringList entryIds = TaskManager::instance()->getEntryIDs();
    int index = entryIds.indexOf(appItem->appId());
    if (index < 0)
        return -1;

    QList<AppItem *> filterAppItems = appItems(widget);
    // 获取当前在最近应用中的所有的APP，并计算它的位置
    int lastIndex = -1;
    // 从后面向前面遍历，找到对应的位置，插入
    for (int i = filterAppItems.size() - 1; i >= 0; i--) {
        AppItem *item = filterAppItems[i];
        // 如果所在的索引在要查找的APP索引的后面，说明当前的索引在要查找的索引之后，跳过即可
        // 如果所在索引不在列表中（一般情况下不存在，这里是容错处理），也跳过
        int curIndex = entryIds.indexOf(item->appId());
        if (item == appItem || curIndex < 0 || curIndex >= index)
            continue;

        if (lastIndex < curIndex)
            lastIndex = curIndex;
    }

    return ++lastIndex;
}

QList<AppItem *> RecentAppHelper::appItems(QWidget *widget) const
{
    QLayout *layout = widget->layout();

    QList<AppItem *> dockItems;
    for (int i = 0; i < layout->count(); i++) {
        AppItem *dockItem = qobject_cast<AppItem *>(layout->itemAt(i)->widget());
        if (!dockItem)
            continue;

        dockItems << dockItem;
    }

    return dockItems;
}
