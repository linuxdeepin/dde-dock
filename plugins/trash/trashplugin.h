/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *               2016 ~ 2018 dragondjf
 *
 * Author:     sbw <sbw@sbw.so>
 *             dragondjf<dingjiangfeng@deepin.com>
 *             zccrs<zhangjide@deepin.com>
 *             Tangtong<tangtong@deepin.com>
 *
 * Maintainer: dragondjf<dingjiangfeng@deepin.com>
 *             zccrs<zhangjide@deepin.com>
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

#ifndef TRASHPLUGIN_H
#define TRASHPLUGIN_H

#include "pluginsiteminterface.h"
#include "trashwidget.h"

#include <QLabel>
#include <QSettings>
namespace Dock{
class TipsWidget;
}
class TrashPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "trash.json")

public:
    explicit TrashPlugin(QObject *parent = 0);

    const QString pluginName() const override;
    const QString pluginDisplayName() const override;
    void init(PluginProxyInterface *proxyInter) override;

    QWidget *itemWidget(const QString &itemKey) override;
    QWidget *itemTipsWidget(const QString &itemKey) override;
    QWidget *itemPopupApplet(const QString &itemKey) override;
    const QString itemCommand(const QString &itemKey) override;
    const QString itemContextMenu(const QString &itemKey) override;
    void refreshIcon(const QString &itemKey) override;
    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked) override;
    bool pluginIsAllowDisable() override { return true; }
    bool pluginIsDisable() override;
    void pluginStateSwitched() override;
    int itemSortKey(const QString &itemKey) override;
    void setSortKey(const QString &itemKey, const int order) override;
    void displayModeChanged(const Dock::DisplayMode displayMode) override;
    void pluginSettingsChanged() override;

private:
    void refreshPluginItemsVisible();

    QScopedPointer<TrashWidget> m_trashWidget;
    QScopedPointer<Dock::TipsWidget> m_tipsLabel;
};

#endif // TRASHPLUGIN_H
