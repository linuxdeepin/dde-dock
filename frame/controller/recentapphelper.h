// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef RECENTAPPHELPER_H
#define RECENTAPPHELPER_H

#include "constants.h"
#include "dbusutil.h"

#include <QObject>

class DockItem;
class AppItem;
class QWidget;

/** 用来管理最近打开区域和APP应用区域交互的类
 * @brief The RecentAppManager class
 */

class RecentAppHelper : public QObject
{
    Q_OBJECT

public:
    explicit RecentAppHelper(QWidget *appWidget, QWidget *recentWidget, QObject *parent = nullptr);
    void setDisplayMode(Dock::DisplayMode displayMode);
    void resetAppInfo();
    void addAppItem(int index, DockItem *appItem);
    void removeAppItem(DockItem *dockItem);
    bool recentIsVisible() const;
    bool dockAppIsVisible() const;

Q_SIGNALS:
    void requestUpdate();
    void recentVisibleChanged(bool);                    // 最近区域是否可见发生变化的信号
    void dockAppVisibleChanged(bool);                       // 驻留应用区域是否可见发生变化的信号
    void requestUpdateRecentVisible();

protected:
    bool eventFilter(QObject *watched, QEvent *event) override;

private:
    bool appInRecent(DockItem *item) const;
    void addAppAreaItem(int index, DockItem *wdg);
    void addRecentAreaItem(int index, DockItem *wdg);
    void updateDockAppVisible(bool lastVisible);

    void removeRecentAreaItem(DockItem *wdg);
    void removeAppAreaItem(DockItem *wdg);

    int getEntryIndex(DockItem *dockItem, QWidget *widget) const;

    QList<AppItem *> appItems(QWidget *widget) const;

private Q_SLOTS:
    void onModeChanged(int mode);
    void updateRecentVisible();

private:
    QWidget *m_appWidget;
    QWidget *m_recentWidget;
    Dock::DisplayMode m_displayMode;
};

#endif // RECENTAPPHELPER_H
