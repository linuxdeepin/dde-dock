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
class QStandardItemModel;
class QStandardItem;
class QMouseEvent;
class QBoxLayout;
class QuickDockItem;
class DockPopupWindow;
class QMenu;
class QuickPluginMimeData;
enum class DockPart;

namespace Dtk { namespace Widget { class DListView; class DStandardItem; } }

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

    bool isQuickWindow(QObject *object) const;

Q_SIGNALS:
    void itemCountChanged();

protected:
    bool eventFilter(QObject *watched, QEvent *event) override;
    void dragEnterEvent(QDragEnterEvent *event) override;
    void dragLeaveEvent(QDragLeaveEvent *event) override;
    void dragMoveEvent(QDragMoveEvent *event) override;
    void resizeEvent(QResizeEvent *event) override;

private Q_SLOTS:
    void onRequestUpdate();
    void onUpdatePlugin(PluginsItemInterface *itemInter, const DockPart &dockPart);
    void onRequestAppletVisible(PluginsItemInterface * itemInter, const QString &itemKey, bool visible);

private:
    void initUi();
    void initConnection();
    void startDrag();
    PluginsItemInterface *findQuickSettingItem(const QPoint &mousePoint, const QList<PluginsItemInterface *> &settingItems);
    int getDropIndex(QPoint point);
    QPoint popupPoint(QWidget *widget) const;
    QuickDockItem *getDockItemByPlugin(PluginsItemInterface *item);
    QuickDockItem *getActiveDockItem(QPoint point) const;
    void showPopup(QuickDockItem *item, PluginsItemInterface *itemInter = nullptr, QWidget *childPage = nullptr, bool isClicked = true);
    QList<QuickDockItem *> quickDockItems();
    DockPopupWindow *getPopWindow() const;
    void updateDockItemSize(QuickDockItem *dockItem);
    void resizeDockItem();

private:
    QBoxLayout *m_mainLayout;
    Dock::Position m_position;
    struct DragInfo *m_dragInfo;
    QuickPluginMimeData *m_dragEnterMimeData;
};

// 用于在任务栏上显示的插件
class QuickDockItem : public QWidget
{
    Q_OBJECT

public:
    explicit QuickDockItem(PluginsItemInterface *pluginItem, const QString &itemKey, QWidget *parent = nullptr);
    ~QuickDockItem();

    void setPosition(Dock::Position position);
    PluginsItemInterface *pluginItem();
    bool canInsert() const;
    bool canMove() const;
    void hideToolTip();

    QSize suitableSize() const;

protected:
    void paintEvent(QPaintEvent *event) override;
    void mousePressEvent(QMouseEvent *event) override;
    void enterEvent(QEvent *event) override;
    void leaveEvent(QEvent *event) override;
    void showEvent(QShowEvent *event) override;
    void hideEvent(QHideEvent *event) override;
    bool eventFilter(QObject *watched, QEvent *event) override;
    void resizeEvent(QResizeEvent *event) override;

private:
    QPoint topleftPoint() const;
    QPoint popupMarkPoint() const;

    QPixmap iconPixmap() const;

    void initUi();
    void initAttribute();
    void initConnection();

    void updateWidgetSize();

private Q_SLOTS:
    void onMenuActionClicked(QAction *action);