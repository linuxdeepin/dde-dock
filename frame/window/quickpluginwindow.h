// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
    explicit QuickPluginWindow(Dock::DisplayMode displayMode, QWidget *parent = nullptr);
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
    Dock::DisplayMode m_displayMode;
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
    int shadowRadius() const;
    int iconSize() const;

private Q_SLOTS:
    void onMenuActionClicked(QAction *action);

private:
    PluginsItemInterface *m_pluginItem;
    QString m_itemKey;
    Dock::Position m_position;
    DockPopupWindow *m_popupWindow;
    QMenu *m_contextMenu;
    QWidget *m_tipParent;
    QHBoxLayout *m_topLayout;
    QWidget *m_mainWidget;
    QHBoxLayout *m_mainLayout;
    QWidget *m_dockItemParent;
    bool m_isEnter;
};

#endif // QUICKPLUGINWINDOW_H
