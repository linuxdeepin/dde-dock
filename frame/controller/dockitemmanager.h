// SPDX-FileCopyrightText: 2019 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DOCKITEMMANAGER_H
#define DOCKITEMMANAGER_H

#include "dockpluginscontroller.h"
#include "pluginsiteminterface.h"
#include "dockitem.h"
#include "appitem.h"
#include "placeholderitem.h"

#include <com_deepin_dde_daemon_dock.h>

#include <QObject>

using DBusDock = com::deepin::dde::daemon::Dock;

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
    const QList<PluginsItemInterface *> pluginList() const;
    bool appIsOnDock(const QString &appDesktop) const;
    void startLoadPlugins() const;

signals:
    void itemInserted(const int index, DockItem *item) const;
    void itemRemoved(DockItem *item) const;
    void itemUpdated(DockItem *item) const;
    void trayVisableCountChanged(const int &count) const;
    void requestWindowAutoHide(const bool autoHide) const;
    void requestRefershWindowVisible() const;
    void requestUpdateDockItem() const;
    void requestUpdateItemMinimizedGeometry(AppItem *item, const QRect) const;

public slots:
    void refreshItemsIcon();
    void itemMoved(DockItem *const sourceItem, DockItem *const targetItem);
    void itemAdded(const QString &appDesktop, int idx);

private Q_SLOTS:
    void onPluginLoadFinished();

private:
    explicit DockItemManager(QObject *parent = nullptr);
    void appItemAdded(const QDBusObjectPath &path, const int index);
    void appItemRemoved(const QString &appId);
    void appItemRemoved(AppItem *appItem);
    void pluginItemInserted(PluginsItem *item);
    void pluginItemRemoved(PluginsItem *item);
    void updatePluginsItemOrderKey();
    void reloadAppItems();
    void manageItem(DockItem *item);

private:
    DBusDock *m_appInter;
    DockPluginsController *m_pluginsInter;

    static DockItemManager *INSTANCE;

    QList<QPointer<DockItem>> m_itemList;
    QList<QString> m_appIDist;

    static const QGSettings *m_appSettings;
    static const QGSettings *m_activeSettings;
    static const QGSettings *m_dockedSettings;
};

#endif // DOCKITEMMANAGER_H
