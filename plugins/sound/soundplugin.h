/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#ifndef SOUNDPLUGIN_H
#define SOUNDPLUGIN_H

#include "pluginsiteminterface.h"
#include "sounditem.h"

class SoundPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "sound.json")

public:
    explicit SoundPlugin(QObject *parent = 0);

    const QString pluginName() const Q_DECL_OVERRIDE;
    const QString pluginDisplayName() const Q_DECL_OVERRIDE;
    void init(PluginProxyInterface *proxyInter) Q_DECL_OVERRIDE;
    void pluginStateSwitched() Q_DECL_OVERRIDE;
    bool pluginIsAllowDisable() Q_DECL_OVERRIDE { return true; }
    bool pluginIsDisable() Q_DECL_OVERRIDE;
    QWidget *itemWidget(const QString &itemKey) Q_DECL_OVERRIDE;
    QWidget *itemTipsWidget(const QString &itemKey) Q_DECL_OVERRIDE;
    QWidget *itemPopupApplet(const QString &itemKey) Q_DECL_OVERRIDE;
    const QString itemContextMenu(const QString &itemKey) Q_DECL_OVERRIDE;
    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked) Q_DECL_OVERRIDE;
    int itemSortKey(const QString &itemKey) Q_DECL_OVERRIDE;
    void setSortKey(const QString &itemKey, const int order) Q_DECL_OVERRIDE;
    void refreshIcon(const QString &itemKey) Q_DECL_OVERRIDE;
    void pluginSettingsChanged() override;

private:
    void refreshPluginItemsVisible();

private:
    SoundItem *m_soundItem;
};

#endif // SOUNDPLUGIN_H
