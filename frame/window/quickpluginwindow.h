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
    void dragPlugin(QuickSettingItem *item);

    QSize suitableSize();

Q_SIGNALS:
    void itemCountChanged();

protected:
    void mouseReleaseEvent(QMouseEvent *event) override;
    void mousePressEvent(QMouseEvent *event) override;

private Q_SLOTS:
    void addPlugin(QuickSettingItem *item);
    void removePlugin(QuickSettingItem *item);
    void onPluginDragMove(QDragMoveEvent *event);

private:
    void initUi();
    void initConnection();
    void startDrag(QuickSettingItem *moveItem);
    QList<QuickSettingItem *> settingItems();
    QuickSettingItem *findQuickSettingItem(const QPoint &mousePoint, const QList<QuickSettingItem *> &settingItems);
    int findActiveTargetIndex(QWidget *widget);
    void resetPluginDisplay();
    QPoint popupPoint() const;

private:
    QBoxLayout *m_mainLayout;
    Dock::Position m_position;
    QList<QuickSettingItem *> m_activeSettingItems;
    QList<QuickSettingItem *> m_fixedSettingItems;
};

#endif // QUICKPLUGINWINDOW_H
