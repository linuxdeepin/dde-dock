/*
 * Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
 *
 * Author:     wangshaojun <wangshaojun_cm@deepin.com>
 *
 * Maintainer: wangshaojun <wangshaojun_cm@deepin.com>
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

#ifndef MULTITASKINGPLUGIN_H
#define MULTITASKINGPLUGIN_H

#include "pluginsiteminterface.h"
#include "multitaskingwidget.h"

#include <QLabel>
namespace Dock {
class TipsWidget;
}
class MultitaskingPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "multitasking.json")

public:
    explicit MultitaskingPlugin(QObject *parent = nullptr);

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
    void refreshIcon(const QString &itemKey) override;
    int itemSortKey(const QString &itemKey) override;
    void setSortKey(const QString &itemKey, const int order) override;
    void pluginSettingsChanged() override;
    PluginType type() override;

private:
    void updateVisible();
    void loadPlugin();
    void refreshPluginItemsVisible();

private:
    bool m_pluginLoaded;

    MultitaskingWidget *m_multitaskingWidget;
    Dock::TipsWidget *m_tipsLabel;
};

#endif // MULTITASKINGPLUGIN_H
