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
#ifndef QUICKPLUGINWINDOW_H
#define QUICKPLUGINWINDOW_H

#include "constants.h"

#include <QWidget>

class QuickSettingItem;
class PluginsItemInterface;
class QHBoxLayout;
class QuickSettingContainer;
class QStandardItemModel;
class QStandardItem;
class QMouseEvent;
class QBoxLayout;
class QuickDockItem;
enum class DockPart;

namespace Dtk { namespace Gui { class DRegionMonitor; }
                namespace Widget { class DListView; class DStandardItem; } }

using namespace Dtk::Widget;

class QuickPluginWindow : public QWidget
{
    Q_OBJECT

public:
    explicit QuickPluginWindow(QWidget *parent = nullptr);
    ~QuickPluginWindow() override;

    void setPositon(Dock::Position position);
    void dragPlugin(PluginsItemInterface *item);

    QSize suitableSize() const;
    QSize suitableSize(const Dock::Position &position) const;

Q_SIGNALS:
    void itemCountChanged();

protected:
    void mousePressEvent(QMouseEvent *event) override;

private Q_SLOTS:
    void addPlugin(QuickSettingItem *item);
    void removePlugin(PluginsItemInterface *item);
    void onPluginDropItem(QDropEvent *event);
    void onPluginDragMove(QDragMoveEvent *event);
    void onFixedClick();
    void onUpdatePlugin(PluginsItemInterface *itemInter, const DockPart &dockPart);

private:
    void initUi();
    void initConnection();
    void startDrag(PluginsItemInterface *moveItem);
    PluginsItemInterface *findQuickSettingItem(const QPoint &mousePoint, const QList<PluginsItemInterface *> &settingItems);
    int findActiveTargetIndex(QuickDockItem *widget);
    int getDropIndex(QPoint point);
    void resetPluginDisplay();
    QPoint popupPoint() const;
    QuickDockItem *getDockItemByPlugin(PluginsItemInterface *item);

private:
    QBoxLayout *m_mainLayout;
    Dock::Position m_position;
    QList<PluginsItemInterface *> m_activeSettingItems;
    QList<PluginsItemInterface *> m_fixedSettingItems;
};

// 用于在任务栏上显示的插件
class QuickDockItem : public QWidget
{
    Q_OBJECT

public:
    explicit QuickDockItem(PluginsItemInterface *pluginItem, QWidget *parent = nullptr);
    ~QuickDockItem();

    PluginsItemInterface *pluginItem();

Q_SIGNALS:
    void clicked();

protected:
    void paintEvent(QPaintEvent *event);
    void mouseReleaseEvent(QMouseEvent *event);

private:
    PluginsItemInterface *m_pluginItem;
};

#endif // QUICKPLUGINWINDOW_H
