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
class DockPopupWindow;
class QMenu;
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
    void requestDrop(QDropEvent *dropEvent);

protected:
    bool eventFilter(QObject *watched, QEvent *event) override;

private Q_SLOTS:
    void onRequestUpdate();
    void onPluginDropItem(QDropEvent *event);
    void onPluginDragMove(QDragMoveEvent *event);
    void onUpdatePlugin(PluginsItemInterface *itemInter, const DockPart &dockPart);
    void onRequestAppletShow(PluginsItemInterface * itemInter, const QString &itemKey);

private:
    void initUi();
    void initConnection();
    void startDrag();
    PluginsItemInterface *findQuickSettingItem(const QPoint &mousePoint, const QList<PluginsItemInterface *> &settingItems);
    int getDropIndex(QPoint point);
    QPoint popupPoint(QWidget *widget) const;
    QuickDockItem *getDockItemByPlugin(PluginsItemInterface *item);
    QuickDockItem *getActiveDockItem(QPoint point) const;
    void showPopup(QuickDockItem *item, QWidget *childPage = nullptr);

private:
    QBoxLayout *m_mainLayout;
    Dock::Position m_position;
    struct DragInfo *m_dragInfo;
};

// 用于在任务栏上显示的插件
class QuickDockItem : public QWidget
{
    Q_OBJECT

public:
    explicit QuickDockItem(PluginsItemInterface *pluginItem, const QJsonObject &metaData, const QString itemKey, QWidget *parent = nullptr);
    ~QuickDockItem();

    void setPositon(Dock::Position position);
    PluginsItemInterface *pluginItem();
    bool isPrimary() const;
    void hideToolTip();

protected:
    void paintEvent(QPaintEvent *event) override;
    void mousePressEvent(QMouseEvent *event) override;
    void enterEvent(QEvent *event) override;
    void leaveEvent(QEvent *event) override;

private:
    QPoint topleftPoint() const;
    QPoint popupMarkPoint() const;

    QPixmap iconPixmap() const;

private Q_SLOTS:
    void onMenuActionClicked(QAction *action);

private:
    PluginsItemInterface *m_pluginItem;
    QJsonObject m_metaData;
    QString m_itemKey;
    Dock::Position m_position;
    DockPopupWindow *m_popupWindow;
    QMenu *m_contextMenu;
    QWidget *m_tipParent;
};

#endif // QUICKPLUGINWINDOW_H
