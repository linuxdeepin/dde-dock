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

#ifndef TOOLAPPHELPER_H
#define TOOLAPPHELPER_H

#include "constants.h"

#include <QObject>

class QWidget;
class DockItem;
class PluginsItem;

using namespace Dock;

class ToolAppHelper : public QObject
{
    Q_OBJECT

public:
    explicit ToolAppHelper(QWidget *pluginAreaWidget, QWidget *toolAreaWidget, QObject *parent = nullptr);

    void setDisplayMode(DisplayMode displayMode);
    void addPluginItem(int index, DockItem *dockItem);
    void removePluginItem(DockItem *dockItem);
    PluginsItem *trashPlugin() const;
    bool toolIsVisible() const;

Q_SIGNALS:
    void requestUpdate();
    void toolVisibleChanged(bool);

private:
    void appendToPluginArea(int index, DockItem *dockItem);
    void appendToToolArea(int index, DockItem *dockItem);
    bool removePluginArea(DockItem *dockItem);
    bool removeToolArea(DockItem *dockItem);

    void resetPluginItems();
    void updateWidgetStatus();
    bool pluginInTool(DockItem *dockItem) const;
    int itemIndex(DockItem *dockItem, bool isTool) const;
    QList<DockItem *> dockItemOnWidget(bool isTool) const;

private:
    QWidget *m_pluginAreaWidget;
    QWidget *m_toolAreaWidget;
    DisplayMode m_displayMode;
    PluginsItem *m_trashItem;
    QList<DockItem *> m_sequentPluginItems;
};

#endif // TOOLAPPHELPER_H
