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
#ifndef QUICKSETTINGCONTAINER_H
#define QUICKSETTINGCONTAINER_H

#include "pluginproxyinterface.h"

#include "dtkwidget_global.h"

#include <DListView>

#include <QWidget>

class DockItem;
class QVBoxLayout;
class QuickSettingController;
class BrightnessModel;
class BrightnessWidget;
class QuickSettingItem;
class DockPopupWindow;
class QStackedLayout;
class VolumeDevicesWidget;
class QLabel;
class PluginChildPage;
class QGridLayout;
class DisplaySettingWidget;
struct QuickDragInfo;

DWIDGET_USE_NAMESPACE

class QuickSettingContainer : public QWidget
{
    Q_OBJECT

public:
    static DockPopupWindow *popWindow();
    static void setPosition(Dock::Position position);
    void showPage(QWidget *widget, PluginsItemInterface *pluginInter = nullptr, bool canBack = false);

protected:
    void mouseMoveEvent(QMouseEvent *event) override;
    bool eventFilter(QObject *watched, QEvent *event) override;

    explicit QuickSettingContainer(QWidget *parent = nullptr);
    ~QuickSettingContainer() override;

private Q_SLOTS:
    void onPluginRemove(PluginsItemInterface *itemInter);
    void onShowChildWidget(QWidget *childWidget);
    void onResizeView();
    void onPluginUpdated(PluginsItemInterface *itemInter, const DockPart dockPart);

private:
    // 加载UI
    void initUi();
    // 初始化槽函数
    void initConnection();
    // 调整控件位置
    void updateItemLayout();
    // 调整全列插件的位置
    void updateFullItemLayout();
    // 插入插件
    void appendPlugin(PluginsItemInterface *itemInter, bool needLayout = true);

private:
    static DockPopupWindow *m_popWindow;
    static Dock::Position m_position;
    QStackedLayout *m_switchLayout;
    QWidget *m_mainWidget;
    QWidget *m_pluginWidget;
    QGridLayout *m_pluginLayout;
    QWidget *m_componentWidget;
    QVBoxLayout *m_mainlayout;
    QuickSettingController *m_pluginLoader;
    PluginChildPage *m_childPage;
    QuickDragInfo *m_dragInfo;
    QList<QuickSettingItem *> m_quickSettings;
    PluginsItemInterface *m_childShowPlugin;
};

#endif // PLUGINCONTAINER_H
