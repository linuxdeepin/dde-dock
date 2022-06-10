/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
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
#ifndef FIXEDPLUGINCONTROLLER_H
#define FIXEDPLUGINCONTROLLER_H

#include "abstractpluginscontroller.h"

class StretchPluginsItem;

class FixedPluginController : public AbstractPluginsController
{
    Q_OBJECT

public:
    explicit FixedPluginController(QObject *parent = nullptr);
    ~FixedPluginController() override;

Q_SIGNALS:
    void pluginItemInserted(StretchPluginsItem *);
    void pluginItemRemoved(StretchPluginsItem *);
    void pluginItemUpdated(StretchPluginsItem *);

protected:
    void itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void itemUpdate(PluginsItemInterface * const itemInter, const QString &) override;
    void itemRemoved(PluginsItemInterface * const itemInter, const QString &) override;

    void requestWindowAutoHide(PluginsItemInterface * const, const QString &, const bool) override {}
    void requestRefreshWindowVisible(PluginsItemInterface * const, const QString &) override {}
    void requestSetAppletVisible(PluginsItemInterface * const, const QString &, const bool) override {}

    bool needLoad(PluginsItemInterface *itemInter) override;

private:
    QList<StretchPluginsItem *> m_pluginItems;
};

#endif // FIXEDPLUGINCONTROLLER_H
