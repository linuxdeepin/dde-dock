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

class QuickSettingController : public AbstractPluginsController
{
    Q_OBJECT

public:
    static QuickSettingController *instance();
    const QList<QuickSettingItem *> &settingItems() const { return m_quickSettingItems; }

Q_SIGNALS:
    void pluginInserted(QuickSettingItem *);
    void pluginRemoved(QuickSettingItem *);
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
    void sortPlugins();

private:
    QList<QuickSettingItem *> m_quickSettingItems;
};

#endif // CONTAINERPLUGINSCONTROLLER_H
