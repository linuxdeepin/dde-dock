// Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "dockitemmanager.h"
#include "appitem.h"
#include "launcheritem.h"
#include "pluginsitem.h"
#include "taskmanager/entry.h"
#include "taskmanager/taskmanager.h"
#include "taskmanager/windowinfobase.h"
#include "traypluginitem.h"
#include "utils.h"
#include "docksettings.h"
#include "appmultiitem.h"
#include "quicksettingcontroller.h"

#include <QDebug>
#include <QGSettings>

#include <DApplication>

DockItemManager *DockItemManager::INSTANCE = nullptr;
const QGSettings *DockItemManager::m_appSettings = Utils::ModuleSettingsPtr("app");
const QGSettings *DockItemManager::m_activeSettings = Utils::ModuleSettingsPtr("activeapp");
const QGSettings *DockItemManager::m_dockedSettings = Utils::ModuleSettingsPtr("dockapp");

DockItemManager::DockItemManager(QObject *parent)
    : QObject(parent)
    , m_taskmanager(TaskManager::instance())
    , m_loadFinished(false)
{
    //固定区域：启动器
    m_itemList.append(new LauncherItem);

    // 应用区域
    for (auto entry : m_taskmanager->getEntries()) {
        AppItem *it = new AppItem(m_appSettings, m_activeSettings, m_dockedSettings, entry);
        manageItem(it);

        connect(it, &AppItem::requestPreviewWindow, m_taskmanager, &TaskManager::previewWindow);
        connect(it, &AppItem::requestCancelPreview, m_taskmanager, &TaskManager::cancelPreviewWindow);
        connect(it, &AppItem::windowCountChanged, this, &DockItemManager::onAppWindowCountChanged);
        connect(this, &DockItemManager::requestUpdateDockItem, it, &AppItem::requestUpdateEntryGeometries);

        m_itemList.append(it);
        m_appIDist.append(it->appId());
        updateMultiItems(it);
    }

    // 托盘区域和插件区域 由DockPluginsController获取
    QuickSettingController *quickController = QuickSettingController::instance();
    connect(quickController, &QuickSettingController::pluginInserted, this, [ = ](PluginsItemInterface *itemInter, const QuickSettingController::PluginAttribute pluginAttr) {
        if (pluginAttr != QuickSettingController::PluginAttribute::Fixed)
            return;

        m_pluginItems << itemInter;
        pluginItemInserted(quickController->pluginItemWidget(itemInter));
    });

    connect(quickController, &QuickSettingController::pluginRemoved, this, &DockItemManager::onPluginItemRemoved);
    connect(quickController, &QuickSettingController::pluginUpdated, this, &DockItemManager::onPluginUpdate);
    connect(quickController, &QuickSettingController::pluginLoaderFinished, this, &DockItemManager::onPluginLoadFinished, Qt::QueuedConnection);

    // 应用信号
    connect(m_taskmanager, &TaskManager::entryAdded, this, &DockItemManager::appItemAdded, Qt::DirectConnection);
    connect(m_taskmanager, &TaskManager::entryRemoved, this, static_cast<void (DockItemManager::*)(const QString &)>(&DockItemManager::appItemRemoved), Qt::DirectConnection);
    connect(m_taskmanager, &TaskManager::serviceRestarted, this, &DockItemManager::reloadAppItems);
    connect(DockSettings::instance(), &DockSettings::showMultiWindowChanged, this, &DockItemManager::onShowMultiWindowChanged);

    DApplication *app = qobject_cast<DApplication *>(qApp);
    if (app) {
        connect(app, &DApplication::iconThemeChanged, this, &DockItemManager::refreshItemsIcon);
    }

    connect(qApp, &QApplication::aboutToQuit, this, &QObject::deleteLater);

    // 读取已经加载的固定区域插件
    QList<PluginsItemInterface *> plugins = quickController->pluginItems(QuickSettingController::PluginAttribute::Fixed);
    for (PluginsItemInterface *plugin : plugins) {
        m_pluginItems << plugin;
        pluginItemInserted(quickController->pluginItemWidget(plugin));
    }

    // 刷新图标
    QMetaObject::invokeMethod(this, "refreshItemsIcon", Qt::QueuedConnection);
}

DockItemManager *DockItemManager::instance(QObject *parent)
{
    if (!INSTANCE)
        INSTANCE = new DockItemManager(parent);

    return INSTANCE;
}

const QList<QPointer<DockItem>> DockItemManager::itemList() const
{
    return m_itemList;
}

bool DockItemManager::appIsOnDock(const QString &appDesktop) const
{
    return m_taskmanager->isOnDock(appDesktop);
}

void DockItemManager::refreshItemsIcon()
{
    for (auto item : m_itemList) {
        if (item.isNull())
            continue;

        item->refreshIcon();
        item->update();
    }
}

/**
 * @brief 将插件的参数(Order, Visible, etc)写入gsettings
 * 自动化测试需要通过dbus(GetPluginSettings)获取这些参数
 */
void DockItemManager::updatePluginsItemOrderKey()
{
    int index = 0;
    for (auto item : m_itemList) {
        if (item.isNull() || item->itemType() != DockItem::Plugins)
            continue;
        static_cast<PluginsItem *>(item.data())->setItemSortKey(++index);
    }

    // 固定区域插件排序
    index = 0;
    for (auto item : m_itemList) {
        if (item.isNull() || item->itemType() != DockItem::FixedPlugin)
            continue;
        static_cast<PluginsItem *>(item.data())->setItemSortKey(++index);
    }
}

void DockItemManager::itemMoved(DockItem *const sourceItem, DockItem *const targetItem)
{
    Q_ASSERT(sourceItem != targetItem);

    const DockItem::ItemType moveType = sourceItem->itemType();
    const DockItem::ItemType replaceType = targetItem->itemType();

    // app move
    if (moveType == DockItem::App || moveType == DockItem::Placeholder)
        if (replaceType != DockItem::App)
            return;

    // plugins move
    if (moveType == DockItem::Plugins || moveType == DockItem::TrayPlugin)
        if (replaceType != DockItem::Plugins && replaceType != DockItem::TrayPlugin)
            return;

    const int moveIndex = m_itemList.indexOf(sourceItem);
    const int replaceIndex = m_itemList.indexOf(targetItem);

    m_itemList.removeAt(moveIndex);
    m_itemList.insert(replaceIndex, sourceItem);

    // update plugins sort key if order changed
    if (moveType == DockItem::Plugins || replaceType == DockItem::Plugins
            || moveType == DockItem::TrayPlugin || replaceType == DockItem::TrayPlugin
            || moveType == DockItem::FixedPlugin || replaceType == DockItem::FixedPlugin) {
        updatePluginsItemOrderKey();
    }

    // for app move, index 0 is launcher item, need to pass it.
    if (moveType == DockItem::App && replaceType == DockItem::App)
        m_taskmanager->moveEntry(moveIndex - 1, replaceIndex - 1);
}

void DockItemManager::itemAdded(const QString &appDesktop, int idx)
{
    m_taskmanager->requestDock(appDesktop, idx);
}

void DockItemManager::appItemAdded(const Entry *entry, const int index)
{
    // 第一个是启动器
    int insertIndex = 1;

    // -1 for append to app list end
    if (index != -1) {
        insertIndex += index;
    } else {
        for (auto item : m_itemList)
            if (!item.isNull() && item->itemType() == DockItem::App)
                ++insertIndex;
    }

    AppItem *item = new AppItem(m_appSettings, m_activeSettings, m_dockedSettings, entry);

    if (m_appIDist.contains(item->appId())) {
        item->deleteLater();
        return;
    }

    manageItem(item);

    connect(item, &AppItem::requestPreviewWindow, m_taskmanager, &TaskManager::previewWindow);
    connect(item, &AppItem::requestCancelPreview, m_taskmanager, &TaskManager::cancelPreviewWindow);
    connect(item, &AppItem::windowCountChanged, this, &DockItemManager::onAppWindowCountChanged);
    connect(this, &DockItemManager::requestUpdateDockItem, item, &AppItem::requestUpdateEntryGeometries);

    m_itemList.insert(insertIndex, item);
    m_appIDist.append(item->appId());

    int itemIndex = insertIndex;
    if (index != -1)
        itemIndex = insertIndex - 1;

    // 插入dockItem
    emit itemInserted(itemIndex, item);
    // 向后插入多开窗口
    updateMultiItems(item, true);
}

void DockItemManager::appItemRemoved(const QString &appId)
{
    for (int i(0); i != m_itemList.size(); ++i) {
        AppItem *app = static_cast<AppItem *>(m_itemList[i].data());
        if (!app) {
            continue;
        }

        if (m_itemList[i]->itemType() != DockItem::App)
            continue;

        if (!app->isValid() || app->appId() == appId) {
            appItemRemoved(app);
            break;
        }
    }

    m_appIDist.removeAll(appId);
}

void DockItemManager::appItemRemoved(AppItem *appItem)
{
    emit itemRemoved(appItem);
    m_itemList.removeOne(appItem);

    if (appItem->isDragging()) {
        QDrag::cancel();
    }
    appItem->deleteLater();
}

void DockItemManager::reloadAppItems()
{
    // remove old item
    for (auto item : m_itemList)
        if (item->itemType() == DockItem::App)
            appItemRemoved(static_cast<AppItem *>(item.data()));

    // append new item
    for (Entry* entry : m_taskmanager->getEntries())
        appItemAdded(entry, -1);
}

void DockItemManager::manageItem(DockItem *item)
{
    connect(item, &DockItem::requestRefreshWindowVisible, this, &DockItemManager::requestRefershWindowVisible, Qt::UniqueConnection);
    connect(item, &DockItem::requestWindowAutoHide, this, &DockItemManager::requestWindowAutoHide, Qt::UniqueConnection);
}

void DockItemManager::pluginItemInserted(PluginsItem *item)
{
    manageItem(item);

    DockItem::ItemType pluginType = item->itemType();

    // find first plugins item position
    int firstPluginPosition = -1;
    for (int i(0); i != m_itemList.size(); ++i) {
        DockItem::ItemType type = m_itemList[i]->itemType();
        if (type != pluginType)
            continue;

        firstPluginPosition = i;
        break;
    }

    if (firstPluginPosition == -1)
        firstPluginPosition = m_itemList.size();

    // find insert position
    int insertIndex = 0;
    const int itemSortKey = item->itemSortKey();
    if (itemSortKey == -1 || firstPluginPosition == -1) {
        insertIndex = m_itemList.size();
    } else if (itemSortKey == 0) {
        insertIndex = firstPluginPosition;
    } else {
        insertIndex = m_itemList.size();
        for (int i(firstPluginPosition + 1); i != m_itemList.size() + 1; ++i) {
            PluginsItem *pItem = static_cast<PluginsItem *>(m_itemList[i - 1].data());
            Q_ASSERT(pItem);

            const int sortKey = pItem->itemSortKey();
            if (pluginType == DockItem::FixedPlugin) {
                if (sortKey != -1 && itemSortKey > sortKey)
                    continue;
                insertIndex = i - 1;
                break;
            }
            if (sortKey != -1 && itemSortKey > sortKey && pItem->itemType() != DockItem::FixedPlugin)
                continue;
            insertIndex = i - 1;
            break;
        }
    }

    m_itemList.insert(insertIndex, item);
    if(pluginType == DockItem::FixedPlugin)
        insertIndex ++;

    if (!Utils::SettingValue(QString("com.deepin.dde.dock.module.") + item->pluginName(), QByteArray(), "enable", true).toBool())
        item->setVisible(false);
    else
        item->setVisible(true);

    emit itemInserted(insertIndex - firstPluginPosition, item);
}

void DockItemManager::onPluginItemRemoved(PluginsItemInterface *itemInter)
{
    if (!m_pluginItems.contains(itemInter))
        return;

    PluginsItem *item = QuickSettingController::instance()->pluginItemWidget(itemInter);
    item->hidePopup();
    item->hide();

    emit itemRemoved(item);

    m_itemList.removeOne(item);

    if (m_loadFinished) {
        updatePluginsItemOrderKey();
    }
}

void DockItemManager::onPluginUpdate(PluginsItemInterface *itemInter)
{
    if (!m_pluginItems.contains(itemInter))
        return;

    Q_EMIT itemUpdated(QuickSettingController::instance()->pluginItemWidget(itemInter));
}

void DockItemManager::onPluginLoadFinished()
{
    updatePluginsItemOrderKey();
    m_loadFinished = true;
}

void DockItemManager::onAppWindowCountChanged()
{
    AppItem *appItem = static_cast<AppItem *>(sender());
    updateMultiItems(appItem, true);
}

void DockItemManager::updateMultiItems(AppItem *appItem, bool emitSignal)
{
    // 如果系统设置不开启应用多窗口拆分，则无需之后的操作
    if (!m_taskmanager->showMultiWindow())
        return;

    // 如果开启了多窗口拆分，则同步窗口和多窗口应用的信息
    const WindowInfoMap &windowInfoMap = appItem->windowsInfos();
    QList<AppMultiItem *> removeItems;
    // 同步当前已经存在的多开窗口的列表，删除不存在的多开窗口
    for (int i = 0; i < m_itemList.size(); i++) {
        QPointer<DockItem> dockItem = m_itemList[i];
        AppMultiItem *multiItem = qobject_cast<AppMultiItem *>(dockItem.data());
        if (!multiItem || multiItem->appItem() != appItem)
            continue;

        // 如果查找到的当前的应用的窗口不需要移除，则继续下一个循环
        if (!needRemoveMultiWindow(multiItem))
            continue;

        removeItems << multiItem;
    }
    // 从itemList中移除多开窗口
    for (AppMultiItem *dockItem : removeItems)
        m_itemList.removeOne(dockItem);
    if (emitSignal) {
        // 移除发送每个多开窗口的移除信号
        for (AppMultiItem *dockItem : removeItems)
            Q_EMIT itemRemoved(dockItem);
    }
    qDeleteAll(removeItems);

    // 遍历当前APP打开的所有窗口的列表，如果不存在多开窗口的应用，则新增，同时发送信号
    for (auto it = windowInfoMap.begin(); it != windowInfoMap.end(); it++) {
        if (multiWindowExist(it.key()))
            continue;

        const WindowInfo &windowInfo = it.value();
        // 如果不存在这个窗口对应的多开窗口，则新建一个窗口，同时发送窗口新增的信号
        AppMultiItem *multiItem = new AppMultiItem(appItem, it.key(), windowInfo);
        m_itemList << multiItem;
        if (emitSignal)
            Q_EMIT itemInserted(-1, multiItem);
    }
}

// 检查对应的窗口是否存在多开窗口
bool DockItemManager::multiWindowExist(quint32 winId) const
{
    for (QPointer<DockItem> dockItem : m_itemList) {
        AppMultiItem *multiItem = qobject_cast<AppMultiItem *>(dockItem.data());
        if (!multiItem)
            continue;

        if (multiItem->winId() == winId)
            return true;
    }

    return false;
}

// 检查当前多开窗口是否需要移除
// 如果当前多开窗口图标对应的窗口在这个窗口所属的APP中所有打开窗口中不存在，那么则认为该多窗口已经被关闭
bool DockItemManager::needRemoveMultiWindow(AppMultiItem *multiItem) const
{
    // 查找多分窗口对应的窗口在应用所有的打开的窗口中是否存在，只要它对应的窗口存在，就无需删除
    // 只要不存在，就需要删除
    AppItem *appItem = multiItem->appItem();
    const WindowInfoMap &windowInfoMap = appItem->windowsInfos();
    for (auto it = windowInfoMap.begin(); it != windowInfoMap.end(); it++) {
        if (it.key() == multiItem->winId())
            return false;
    }

    return true;
}

void DockItemManager::onShowMultiWindowChanged()
{
    if (m_taskmanager->showMultiWindow()) {
        // 如果当前设置支持窗口多开，那么就依次对每个APPItem加载多开窗口
        for (int i = 0; i < m_itemList.size(); i++) {
            const QPointer<DockItem> &dockItem = m_itemList[i];
            if (dockItem->itemType() != DockItem::ItemType::App)
                continue;

            updateMultiItems(static_cast<AppItem *>(dockItem.data()), true);
        }
    } else {
        // 如果当前设置不支持窗口多开，则删除所有的多开窗口
        QList<DockItem *> multiWindows;
        for (const QPointer<DockItem> &dockItem : m_itemList) {
            if (dockItem->itemType() != DockItem::AppMultiWindow)
                continue;

            multiWindows << dockItem.data();
        }
        for (DockItem *multiItem : multiWindows) {
            m_itemList.removeOne(multiItem);
            Q_EMIT itemRemoved(multiItem);
            multiItem->deleteLater();
        }
    }
}
