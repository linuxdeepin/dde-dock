// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef ONBOARDPLUGIN_H
#define ONBOARDPLUGIN_H

#include "pluginsiteminterface.h"
#include "onboarditem.h"

#include <QLabel>
#include <QScopedPointer>

namespace Dock {
class TipsWidget;
}
class OnboardPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "onboard.json")

public:
    explicit OnboardPlugin(QObject *parent = nullptr);

    const QString pluginName() const override;
    const QString pluginDisplayName() const override;
    void init(PluginProxyInterface *proxyInter) override;

    void pluginStateSwitched() override;
    bool pluginIsAllowDisable() override { return true; }
    bool pluginIsDisable() override;

    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;
    const QString itemCommand(const QString &itemKey) override;
    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked) override;
    void displayModeChanged(const Dock::DisplayMode displayMode) override;

    int itemSortKey(const QString &itemKey) override;
    void setSortKey(const QString &itemKey, const int order) override;

    void pluginSettingsChanged() override;
    QIcon icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType) override;
    PluginMode status() const override;
    QString description() const override;

private:
    void loadPlugin();
    void refreshPluginItemsVisible();

private:
    bool m_pluginLoaded;
    bool m_startupState;

    QScopedPointer<OnboardItem> m_onboardItem;
    QScopedPointer<Dock::TipsWidget> m_tipsLabel;
};

#endif // ONBOARDPLUGIN_H
