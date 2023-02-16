// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SHUTDOWNPLUGIN_H
#define SHUTDOWNPLUGIN_H

#include "pluginsiteminterface.h"
#include "shutdownwidget.h"

#include <QLabel>

class DBusPowerManager;

namespace Dock {
class TipsWidget;
}
class QGSettings;
class ShutdownPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "shutdown.json")

public:
    explicit ShutdownPlugin(QObject *parent = nullptr);

    const QString pluginName() const override;
    const QString pluginDisplayName() const override;
    void init(PluginProxyInterface *proxyInter) override;

    void pluginStateSwitched() override;
    bool pluginIsAllowDisable() override { return true; }
    bool pluginIsDisable() override;

    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;
    const QString itemCommand(const QString &itemKey) override;
    const QString itemContextMenu(const QString &itemKey) override;
    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked) override;
    void displayModeChanged(const Dock::DisplayMode displayMode) override;

    int itemSortKey(const QString &itemKey) override;
    void setSortKey(const QString &itemKey, const int order) override;

    QIcon icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType) override;
    PluginFlags flags() const override;

    // 休眠待机配置，保持和sessionshell一致
    const QStringList session_ui_configs {
        "/etc/lightdm/lightdm-deepin-greeter.conf",
        "/etc/deepin/dde-session-ui.conf",
        "/usr/share/dde-session-ui/dde-session-ui.conf"
    };
    template <typename T>
    T findValueByQSettings(const QStringList &configFiles,
                           const QString &group,
                           const QString &key,
                           const QVariant &failback)
    {
        for (const QString &path : configFiles) {
            QSettings settings(path, QSettings::IniFormat);
            if (!group.isEmpty()) {
                settings.beginGroup(group);
            }

            const QVariant& v = settings.value(key);
            if (v.isValid()) {
                T t = v.value<T>();
                return t;
            }
        }

        return failback.value<T>();
    }

    template <typename T>
    T valueByQSettings(const QString & group,
                       const QString & key,
                       const QVariant &failback) {
        return findValueByQSettings<T>(session_ui_configs,
                                       group,
                                       key,
                                       failback);
    }

    std::pair<bool, qint64> checkIsPartitionType(const QStringList &list);
    qint64 get_power_image_size();

private:
    void loadPlugin();
    bool checkSwap();

private:
    bool m_pluginLoaded;

    QScopedPointer<ShutdownWidget> m_shutdownWidget;
    QScopedPointer<Dock::TipsWidget> m_tipsLabel;
    DBusPowerManager* m_powerManagerInter;
    const QGSettings *m_gsettings;
    const QGSettings *m_sessionShellGsettings;
};

#endif // SHUTDOWNPLUGIN_H
