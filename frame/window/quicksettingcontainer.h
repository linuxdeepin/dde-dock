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
class MediaWidget;
class VolumeModel;
class BrightnessModel;
class VolumeWidget;
class BrightnessWidget;
class QuickSettingItem;
class DockPopupWindow;
class QStackedLayout;
class VolumeDevicesWidget;
class QLabel;
class PluginChildPage;
class QGridLayout;
class DisplaySettingWidget;

DWIDGET_USE_NAMESPACE

class QuickSettingContainer : public QWidget
{
    Q_OBJECT

public:
    static DockPopupWindow *popWindow();
    static void setPosition(Dock::Position position);

protected:
    void mousePressEvent(QMouseEvent *event) override;
    void mouseReleaseEvent(QMouseEvent *event) override;
    void mouseMoveEvent(QMouseEvent *event) override;

    explicit QuickSettingContainer(QWidget *parent = nullptr);
    ~QuickSettingContainer() override;
    void showHomePage();

private Q_SLOTS:
    void onPluginInsert(QuickSettingItem *quickItem);
    void onPluginRemove(QuickSettingItem *quickItem);
    void onItemDetailClick(PluginsItemInterface *pluginInter);
    bool eventFilter(QObject *watched, QEvent *event) override;

private:
    // 加载UI
    void initUi();
    // 初始化槽函数
    void initConnection();
    // 调整尺寸
    void resizeView();
    // 调整控件位置
    void updateItemLayout();
    // 初始化控件项目
    void initQuickItem(QuickSettingItem *quickItem);
    // 显示具体的窗体
    void showWidget(QWidget *widget, const QString &title);
    // 清除移动轨迹
    void clearDragPoint();

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
    MediaWidget *m_playerWidget;
    VolumeModel *m_volumeModel;
    BrightnessModel *m_brightnessModel;
    VolumeWidget *m_volumnWidget;
    BrightnessWidget *m_brihtnessWidget;

    VolumeDevicesWidget *m_volumeSettingWidget;
    DisplaySettingWidget *m_displaySettingWidget;
    PluginChildPage *m_childPage;
    QPoint m_dragPluginPosition;
};

class QuickPluginMimeData : public QMimeData
{
    Q_OBJECT

public:
    explicit QuickPluginMimeData(PluginsItemInterface *item) : QMimeData(), m_item(item) {}
    ~QuickPluginMimeData() {}
    PluginsItemInterface *pluginItemInterface() const { return m_item; }

private:
     PluginsItemInterface *m_item;
};

#endif // PLUGINCONTAINER_H
