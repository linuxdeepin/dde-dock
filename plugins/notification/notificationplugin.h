// SPDX-FileCopyrightText: 2024 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later
#ifndef NOTIFICATIONPLUGIN_H
#define NOTIFICATIONPLUGIN_H
#include "pluginsiteminterface.h"
#include "notification.h"
#include "tipswidget.h"

class NotificationPlugin : public QObject, public PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "notification.json")

public:
    explicit NotificationPlugin(QObject *parent = nullptr);
    const QString pluginName() const override { return "notification"; }
    const QString pluginDisplayName() const override { return tr("Notification"); }
    PluginMode status() const override { return PluginMode::Active; }
    PluginType type() override { return PluginType::Normal; }
    PluginFlags flags() const override { return PluginFlag::Type_Common | PluginFlag::Attribute_CanSetting; }
    QString description() const override { return pluginDisplayName(); }
    bool pluginIsAllowDisable() override { return true; }

    void init(PluginProxyInterface *proxyInter) override;
    void pluginStateSwitched() override;
    bool pluginIsDisable() override;

    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;
    const QString itemCommand(const QString &itemKey) override;
    const QString itemContextMenu(const QString &itemKey) override;
    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked) override;

    int itemSortKey(const QString &itemKey) override;
    void setSortKey(const QString &itemKey, const int order) override;

    void pluginSettingsChanged() override;
    QIcon icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType) override;
    void refreshIcon(const QString &itemKey) override;

private:
    void loadPlugin();
    void refreshPluginItemsVisible();
    void updateTipsText(uint notificationCount);
    QString toggleDndText() const;

private:
    bool m_pluginLoaded;
    QScopedPointer<Notification> m_notification;
    QScopedPointer<Dock::TipsWidget> m_tipsLabel;
};

#endif  // NOTIFICATIONPLUGIN_H
