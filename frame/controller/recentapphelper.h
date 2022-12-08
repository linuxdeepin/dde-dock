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
    explicit RecentAppHelper(QWidget *appWidget, QWidget *recentWidget, DockInter *dockInter, QObject *parent = nullptr);
    void setDisplayMode(Dock::DisplayMode displayMode);
    void resetAppInfo();
    void addAppItem(int index, DockItem *appItem);
    void removeAppItem(DockItem *dockItem);
    bool recentIsVisible() const;
    bool dockAppIsVisible() const;
    void updateDockInter(DockInter *dockInter);

Q_SIGNALS:
    void requestUpdate();
    void recentVisibleChanged(bool);                    // 最近区域是否可见发生变化的信号
    void dockAppVisibleChanged(bool);                       // 驻留应用区域是否可见发生变化的信号

protected:
    bool eventFilter(QObject *watched, QEvent *event) override;

private:
    bool appInRecent(DockItem *item) const;
    void addAppAreaItem(int index, DockItem *wdg);
    void addRecentAreaItem(int index, DockItem *wdg);
    void updateRecentVisible();
    void updateDockAppVisible(bool lastVisible);

    void removeRecentAreaItem(DockItem *wdg);
    void removeAppAreaItem(DockItem *wdg);

    int getEntryIndex(DockItem *dockItem, QWidget *widget) const;

    QList<AppItem *> appItems(QWidget *widget) const;

private Q_SLOTS:
    void onModeChanged(int mode);

private:
    QWidget *m_appWidget;
    QWidget *m_recentWidget;
    Dock::DisplayMode m_displayMode;
    DockInter *m_dockInter;
};

#endif // RECENTAPPHELPER_H
