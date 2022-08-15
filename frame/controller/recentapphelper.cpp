/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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

#include "recentapphelper.h"
#include "appitem.h"

#include <QWidget>

#define ENTRY_NONE      0
#define ENTRY_NORMAL    1
#define ENTRY_RECENT    2

RecentAppHelper::RecentAppHelper(QWidget *appWidget, QWidget *recentWidget, DockInter *dockInter, QObject *parent)
    : QObject(parent)
    , m_appWidget(appWidget)
    , m_recentWidget(recentWidget)
    , m_dockInter(dockInter)
{
    m_appWidget->installEventFilter(this);
    m_recentWidget->installEventFilter(this);
}

void RecentAppHelper::setDisplayMode(Dock::DisplayMode displayMode)
{
    bool lastVisible = dockAppIsVisible();
    m_displayMode = displayMode;
#ifndef USE_AM
    resetDockItems();
#endif
    updateRecentVisible();
    updateDockAppVisible(lastVisible);
}

// 当在应用区域调整位置的时候，需要重新设置索引
void RecentAppHelper::resetAppInfo()
{
    #ifndef USE_AM
    // 获取应用区域和最近打开区域的app图标
    QList<AppItem *> appDockItem = appItems(m_appWidget);

    // 获取应用区域图标在原来列表中的位置
    QList<int> dockIndex;
    for (DockItem *appItem : appDockItem)
        dockIndex << m_sequentDockItems.indexOf(appItem);

    // 按照从小到大排序
    std::sort(dockIndex.begin(), dockIndex.end(), [ = ](int index1, int index2) { return index1 < index2; });
    QMap<DockItem *, int> dockItemIndex;
    for (int i = 0; i < appDockItem.size(); i++) {
        DockItem *item = appDockItem[i];
        dockItemIndex[item] = dockIndex[i];
    }

    // 替换原来的位置
    for (DockItem *appItem : appDockItem) {
        int index = -1;
        if (dockItemIndex.contains(appItem))
            index = dockItemIndex.value(appItem);

        if (index >= 0)
            m_sequentDockItems[index] = appItem;
        else
            m_sequentDockItems << appItem;
    }
#endif
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

#ifdef USE_AM
    connect(appItem, &AppItem::modeChanged, this, &RecentAppHelper::onModeChanged);
#else
    connect(dockItem, &QWidget::destroyed, this, [ this, dockItem ] {
        if (m_sequentDockItems.contains(dockItem))
            m_sequentDockItems.removeOne(dockItem);
    });
    connect(appItem, &AppItem::isDockChanged, this, &RecentAppHelper::onItemChanged);

    // 如果索引值大于0，说明它是插入到固定位置的，否则，则认为它是顺序排列的
    if (index >= 0 && index < m_sequentDockItems.size())
        m_sequentDockItems.insert(index, dockItem);
    else
        m_sequentDockItems << dockItem;
#endif
}

void RecentAppHelper::removeAppItem(DockItem *dockItem)
{
    if (m_recentWidget->children().contains(dockItem))
        removeRecentAreaItem(dockItem);
    else
        removeAppAreaItem(dockItem);
#ifndef USE_AM
    AppItem *appItem = qobject_cast<AppItem *>(dockItem);
    disconnect(appItem, &AppItem::isDockChanged, this, &RecentAppHelper::onItemChanged);
#endif
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

#ifdef USE_AM
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
#else
void RecentAppHelper::onItemChanged()
{
    bool lastVisible = dockAppIsVisible();
    resetDockItems();
    updateRecentVisible();
    updateDockAppVisible(lastVisible);
}
#endif

bool RecentAppHelper::appInRecent(DockItem *item) const
{
#ifdef USE_AM
    AppItem *appItem = qobject_cast<AppItem *>(item);
    if (!appItem)
        return false;

    return (appItem->mode() == ENTRY_RECENT);
#else
    // 先判断当前是否为时尚模式，只有时尚模式下才支持最近打开的应用
    if (m_displayMode != Dock::DisplayMode::Fashion)
        return false;
    // 只有当应用没有固定到任务栏上才认为它是最新打开的应用
    AppItem *appItem = qobject_cast<AppItem *>(item);
    return (appItem && !appItem->isDocked());
#endif
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
    updateRecentVisible();
}

void RecentAppHelper::removeAppAreaItem(DockItem *wdg)
{
    QBoxLayout *boxLayout = static_cast<QBoxLayout *>(m_appWidget->layout());
    bool lastVisible = dockAppIsVisible();
    boxLayout->removeWidget(wdg);
    updateDockAppVisible(lastVisible);
}

#ifndef USE_AM
QList<DockItem *> RecentAppHelper::dockItemToAppArea() const
{
    QList<DockItem *> dockItems;
    if (m_displayMode == Dock::DisplayMode::Efficient) {
        // 由特效模式变成高效模式，将所有的最近打开的应用移动到左侧的应用区域
        for (int i = 0; i < m_recentWidget->layout()->count(); i++) {
            DockItem *appItem = qobject_cast<DockItem *>(m_recentWidget->layout()->itemAt(i)->widget());
            if (!appItem)
                continue;

            dockItems << appItem;
        }
    } else {
        // 如果是特效模式下，则查找所有已驻留的应用，将其移动到应用区域
        for (int i = 0; i < m_recentWidget->layout()->count(); i++) {
            DockItem *appItem = qobject_cast<DockItem *>(m_recentWidget->layout()->itemAt(i)->widget());
            if (!appItem || appInRecent(appItem))
                continue;

            dockItems << appItem;
        }
    }

    return dockItems;
}

void RecentAppHelper::resetDockItems()
{
    // 先将所有的最近打开的区域移动到应用区域
    QList<DockItem *> recentAppItems = dockItemToAppArea();

    // 从最近使用应用中移除
    for (DockItem *appItem : recentAppItems)
        m_recentWidget->layout()->removeWidget(appItem);

    // 将这些图标添加到应用区域
    QBoxLayout *boxLayout = static_cast<QBoxLayout *>(m_appWidget->layout());
    for (DockItem *appItem : recentAppItems) {
        int index = getDockItemIndex(appItem, false);
        if (index >= 0)
            boxLayout->insertWidget(index, appItem);
        else
            boxLayout->addWidget(appItem);
    }

    if (m_displayMode == Dock::DisplayMode::Fashion) {
        // 由高效模式变成特效模式后，将左侧未驻留的应用移动到右侧的最近打开应用中
        QList<DockItem *> unDockItems;
        for (int i = 0; i < m_appWidget->layout()->count() ; i++) {
            DockItem *appItem = qobject_cast<DockItem *>(m_appWidget->layout()->itemAt(i)->widget());
            if (!appInRecent(appItem))
                continue;

            unDockItems << appItem;
        }

        // 从应用区域中删除未驻留的应用
        for (DockItem *appItem : unDockItems)
            m_appWidget->layout()->removeWidget(appItem);

        // 将这些图标添加到最近打开区域
        QBoxLayout *recentLayout = static_cast<QBoxLayout *>(m_recentWidget->layout());
        for (DockItem *appItem : unDockItems) {
            int index = getDockItemIndex(appItem, true);
            if (index >= 0)
                recentLayout->insertWidget(index, appItem);
            else
                recentLayout->addWidget(appItem);
        }
    }
}

int RecentAppHelper::getDockItemIndex(DockItem *dockItem, bool isRecent) const
{
    // 当从最近区域移动到应用区域的时候，重新计算插入索引值
    if (!m_sequentDockItems.contains(dockItem))
        return -1;

    QList<DockItem *> sequeDockItems = m_sequentDockItems;
    if (isRecent) {
        // 如果是最近打开区域，需要按照时间从小到大排列(先打开的排在前面)
        std::sort(sequeDockItems.begin(), sequeDockItems.end(), [](DockItem *item1, DockItem *item2) {
            AppItem *appItem1 = qobject_cast<AppItem *>(item1);
            AppItem *appItem2 = qobject_cast<AppItem *>(item2);
            if (!appItem1 || !appItem2)
                return false;

            return appItem1->appOpenMSecs() < appItem2->appOpenMSecs();
        });
    }

    int index = sequeDockItems.indexOf(dockItem);
    // 查找所有在应用区域的图标
    QList<AppItem *> dockApps = appItems(isRecent ? m_recentWidget : m_appWidget);
    for (int i = index + 1; i < sequeDockItems.size(); i++) {
        DockItem *item = sequeDockItems[i];
        int itemIndex = dockApps.indexOf(static_cast<AppItem *>(item));
        if (itemIndex >= 0)
            return itemIndex;
    }

    return -1;
}
#endif

int RecentAppHelper::getEntryIndex(DockItem *dockItem, QWidget *widget) const
{
    AppItem *appItem = qobject_cast<AppItem *>(dockItem);
    if (!appItem)
        return -1;

    // 查找当前的应用在所有的应用中的排序
    QStringList entryIds = m_dockInter->GetEntryIDs();
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
