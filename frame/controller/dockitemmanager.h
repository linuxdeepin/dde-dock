// Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DOCKITEMMANAGER_H
#define DOCKITEMMANAGER_H

#include "pluginsiteminterface.h"
#include "dockitem.h"
#include "appitem.h"
#include "placeholderitem.h"
#include "dbusutil.h"
#include "taskmanager/taskmanager.h"
#include "taskmanager/windowinfobase.h"

#include <QObject>

class AppMultiItem;
class PluginsItem;

/**
 * @brief The DockItemManager class
 * 管理类，管理所有的应用数据，插件数据
 */
class DockItemManager : public QObject
{
    Q_OBJECT

public:
    static DockItemManager *instance(QObject *parent = nullptr);

    const QList<QPointer<DockItem> > itemList() const;
    bool appIsOnDock(const QString &appDesktop) const;

signals:
    void itemInserted(const int index, DockItem *item) const;
    void itemRemoved(DockItem *item) const;
    void itemUpdated(DockItem *item) const;
    void trayVisableCountChanged(const int &count) const;
    void requestWindowAutoHide(const bool autoHide) const;
    void requestRefershWindowVisible() const;

    void requestUpdateDockItem() const;

public slots:
    void refreshItemsIcon();
    void itemMoved(DockItem *const sourceItem, DockItem *const targetItem);
    void itemAdded(const QString &appDesktop, int idx);

private Q_SLOTS:
    void onPluginLoadFinished();
    void onPluginItemRemoved(PluginsItemInterface *itemInter);
    void onPluginUpdate(PluginsItemInterface *itemInter);

    void onAppWindowCountChanged();
    void onShowMultiWindowChanged();

private:
    explicit DockItemManager(QObject *parent = nullptr);
    void appItemAdded(const Entry *entry, const int index);
    void appItemRemoved(const QString &appId);
    void appItemRemoved(AppItem *appItem);
    void updatePluginsItemOrderKey();
    void reloadAppItems();
    void manageItem(DockItem *item);
    void pluginItemInserted(PluginsItem *item);

    void updateMultiItems(AppItem *appItem, bool emitSignal = false);
    bool multiWindowExist(quint32 winId) const;
    bool needRemoveMultiWindow(AppMultiItem *multiItem) const;

private:
    TaskManager *m_taskmanager;

    static DockItemManager *INSTANCE;

    QList<QPointer<DockItem>> m_itemList;
    QList<QString> m_appIDist;
    QList<PluginsItemInterface *> m_pluginItems;

    bool m_loadFinished; // 记录所有插件是否加载完成

    static const QGSettings *m_appSettings;
    static const QGSettings *m_activeSettings;
    static const QGSettings *m_dockedSettings;
};

#endif // DOCKITEMMANAGER_H
