// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SYSTEMTRAYSCONTROLLER_H
#define SYSTEMTRAYSCONTROLLER_H

#include "systemtrayitem.h"
#include "pluginproxyinterface.h"
#include "util/abstractpluginscontroller.h"

#include <QPluginLoader>
#include <QList>
#include <QMap>
#include <QDBusConnectionInterface>

class PluginsItemInterface;
class ProxyPluginController;

class SystemTraysController : public AbstractPluginsController
{
    Q_OBJECT

public:
    explicit SystemTraysController(QObject *parent = nullptr);
    ~SystemTraysController() override;

    // implements PluginProxyInterface
    void itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void itemUpdate(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void itemRemoved(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void requestWindowAutoHide(PluginsItemInterface * const itemInter, const QString &itemKey, const bool autoHide) override;
    void requestRefreshWindowVisible(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void requestSetAppletVisible(PluginsItemInterface * const itemInter, const QString &itemKey, const bool visible) override;

    int systemTrayItemSortKey(const QString &itemKey);
    void setSystemTrayItemSortKey(const QString &itemKey, const int order);

    const QVariant getValueSystemTrayItem(const QString &itemKey, const QString &key, const QVariant& fallback = QVariant());
    void saveValueSystemTrayItem(const QString &itemKey, const QString &key, const QVariant &value);

signals:
    void pluginItemAdded(const QString &itemKey, AbstractTrayWidget *pluginItem) const;
    void pluginItemRemoved(const QString &itemKey, AbstractTrayWidget *pluginItem) const;
    void pluginItemUpdated(const QString &itemKey, AbstractTrayWidget *pluginItem) const;

private:
    void loadPlugins(ProxyPluginController *proxyController);
};

#endif // SYSTEMTRAYSCONTROLLER_H
