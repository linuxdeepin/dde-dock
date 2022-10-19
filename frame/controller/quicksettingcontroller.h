/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#ifndef QUICKSETTINGCONTROLLER_H
#define QUICKSETTINGCONTROLLER_H

#include "abstractpluginscontroller.h"

class QuickSettingItem;
class PluginsItem;

class QuickSettingController : public AbstractPluginsController
{
    Q_OBJECT

public:
    enum class PluginAttribute {
        Quick = 0,
        Tool,
        Fixed
    };

public:
    static QuickSettingController *instance();
    QList<PluginsItemInterface *> pluginItems(const PluginAttribute &pluginClass) const;
    QString itemKey(PluginsItemInterface *pluginItem) const;
    QJsonObject metaData(PluginsItemInterface *pluginItem) const;
    PluginsItem *pluginItemWidget(PluginsItemInterface *pluginItem);

Q_SIGNALS:
    void pluginInserted(PluginsItemInterface *itemInter, const PluginAttribute &);
    void pluginRemoved(PluginsItemInterface *itemInter);
    void pluginUpdated(PluginsItemInterface *, const DockPart &);

protected:
    explicit QuickSettingController(QObject *parent = Q_NULLPTR);
    ~QuickSettingController() override;

protected:
    void itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void itemUpdate(PluginsItemInterface * const itemInter, const QString &) override;
    void itemRemoved(PluginsItemInterface * const itemInter, const QString &) override;
    void requestWindowAutoHide(PluginsItemInterface * const, const QString &, const bool) override {}
    void requestRefreshWindowVisible(PluginsItemInterface * const, const QString &) override {}
    void requestSetAppletVisible(PluginsItemInterface * const, const QString &, const bool) override {}
    void updateDockInfo(PluginsItemInterface * const itemInter, const DockPart &part) override;

private:
    QMap<PluginAttribute, QList<PluginsItemInterface *>> m_quickPlugins;
    QMap<PluginsItemInterface *, QString> m_quickPluginsMap;
    QMap<PluginsItemInterface *, PluginsItem *> m_pluginItemWidgetMap;
};

#endif // CONTAINERPLUGINSCONTROLLER_H
