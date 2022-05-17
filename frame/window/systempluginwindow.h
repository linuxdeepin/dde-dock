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
#ifndef SYSTEMPLUGINWINDOW_H
#define SYSTEMPLUGINWINDOW_H

#include "constants.h"
#include "dockpluginscontroller.h"

#include <DBlurEffectWidget>

class DockPluginsController;
class PluginsItem;
class QBoxLayout;

namespace Dtk { namespace Widget { class DListView; } }

DWIDGET_USE_NAMESPACE

class SystemPluginWindow : public DBlurEffectWidget
{
    Q_OBJECT

public:
    explicit SystemPluginWindow(QWidget *parent = nullptr);
    ~SystemPluginWindow() override;
    void setPositon(Dock::Position position);
    QSize suitableSize();

Q_SIGNALS:
    void pluginSizeChanged();

private:
    void initUi();
    int calcIconSize() const;
    void resizeEvent(QResizeEvent *event) override;

private Q_SLOTS:
    void onPluginItemAdded(PluginsItem *pluginItem);
    void onPluginItemRemoved(PluginsItem *pluginItem);
    void onPluginItemUpdated(PluginsItem *pluginItem);

private:
    DockPluginsController *m_pluginController;
    DListView *m_listView;
    Dock::Position m_position;
    QBoxLayout *m_mainLayout;
};

class FixedPluginController : public DockPluginsController
{
    Q_OBJECT

public:
    explicit FixedPluginController(QObject *parent);

protected:
    PluginsItem *createPluginsItem(PluginsItemInterface *const itemInter, const QString &itemKey, const QString &pluginApi) override;
    bool needLoad(PluginsItemInterface *itemInter) override;
};

#endif // SYSTEMPLUGINWINDOW_H
