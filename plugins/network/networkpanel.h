/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
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

#ifndef NETWORKPANEL_H
#define NETWORKPANEL_H

#include "item/devicestatushandler.h"

#include <DGuiApplicationHelper>
#include <DSwitchButton>
#include <DListView>
#include <DStyledItemDelegate>
#include <dloadingindicator.h>

#include <QWidget>
#include <QScopedPointer>

DGUI_USE_NAMESPACE
DWIDGET_USE_NAMESPACE

namespace Dock {
  class TipsWidget;
}

namespace dde {
  namespace network {
    enum class UDeviceType;
    class UNetworkDeviceBase;
  }
}

class QTimer;
class NetItem;

using namespace dde::network;

class NetworkPanel : public QWidget
{
    Q_OBJECT

public:
    explicit NetworkPanel(QWidget *parent = Q_NULLPTR);
    ~NetworkPanel();

    void invokeMenuItem(const QString &menuId);
    bool needShowControlCenter();
    const QString contextMenu() const;
    QWidget *itemTips();
    QWidget *itemApplet();
    bool hasDevice();

    void refreshIcon();

private:
    void setControlBackground();
    void initUi();
    void initConnection();
    void getPluginState();
    void updateView();                                                  // 更新网络列表内容大小
    void updateTooltips();                                              // 更新提示的内容
    void updateItems(QList<NetItem *> &removeItems);
    bool deviceEnabled(const UDeviceType &deviceType) const;
    void setDeviceEnabled(const UDeviceType &deviceType, bool enabeld);

    int getStrongestAp();
    int deviceCount(const UDeviceType &devType);
    QStringList ipTipsMessage(const UDeviceType &devType);

    void paintEvent(QPaintEvent *e);
    void resizeEvent(QResizeEvent *e);

private Q_SLOTS:
    void onDeviceAdded(QList<UNetworkDeviceBase *> devices);

    void updatePlugView();
    void wirelessChanged();

    void onClickListView(const QModelIndex &index);

private:
    PluginState m_pluginState;

    QTimer *m_refreshIconTimer;
    QTimer *m_switchWireTimer;
    QTimer *m_wirelessScanTimer;
    int m_wirelessScanInterval;

    Dock::TipsWidget *m_tipsWidget;
    bool m_switchWire;
    QPixmap m_iconPixmap;

    QStandardItemModel *m_model;

    QScrollArea *m_applet;
    QWidget *m_centerWidget;
    DListView *m_netListView;
    // 判断定时的时间是否到,否则不重置计时器
    bool m_timeOut;

    QList<NetItem *> m_items;
};

class NetworkDelegate : public DStyledItemDelegate
{
    Q_OBJECT

public:
    NetworkDelegate(QAbstractItemView *parent = Q_NULLPTR);
    ~NetworkDelegate() Q_DECL_OVERRIDE;

private:
    void paint(QPainter *painter, const QStyleOptionViewItem &option, const QModelIndex &index) const Q_DECL_OVERRIDE;

    bool needDrawLine(const QModelIndex &index) const;
};

#endif // NETWORKPANEL_H
