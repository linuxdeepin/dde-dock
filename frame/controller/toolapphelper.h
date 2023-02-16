// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef TOOLAPPHELPER_H
#define TOOLAPPHELPER_H

#include "constants.h"

#include <QObject>

class QWidget;
class DockItem;
class PluginsItem;
class PluginsItemInterface;

using namespace Dock;

class ToolAppHelper : public QObject
{
    Q_OBJECT

public:
    explicit ToolAppHelper(QWidget *toolAreaWidget, QObject *parent = nullptr);

    void setDisplayMode(DisplayMode displayMode);
    void setPosition(Dock::Position position);
    bool toolIsVisible() const;

Q_SIGNALS:
    void requestUpdate();
    void toolVisibleChanged(bool);

private:
    void appendToToolArea(int index, DockItem *dockItem);
    bool removeToolArea(PluginsItemInterface *itemInter);
    void moveToolWidget();
    void updateToolArea();

    void updateWidgetStatus();
    bool pluginInTool(PluginsItemInterface *itemInter) const;
    void pluginItemAdded(PluginsItemInterface *itemInter);
    void pluginItemRemoved(PluginsItemInterface *itemInter);
    bool pluginExists(PluginsItemInterface *itemInter) const;

private:
    QWidget *m_toolAreaWidget;
    DisplayMode m_displayMode;
    Dock::Position m_position;
    QList<DockItem *> m_sequentPluginItems;
};

#endif // TOOLAPPHELPER_H
