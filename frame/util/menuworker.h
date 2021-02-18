/*
 * Copyright (C) 2018 ~ 2028 Deepin Technology Co., Ltd.
 *
 * Author:     fanpengcheng <fanpengcheng_cm@deepin.com>
 *
 * Maintainer: fanpengcheng <fanpengcheng_cm@deepin.com>
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
#ifndef MENUWORKER_H
#define MENUWORKER_H
#include <QObject>

#include "constants.h"

#include <com_deepin_dde_daemon_dock.h>

using DBusDock = com::deepin::dde::daemon::Dock;
class QMenu;
class QAction;
class DockItemManager;
/**
 * @brief The MenuWorker class  此类用于处理任务栏右键菜单的逻辑
 */
class MenuWorker : public QObject
{
    Q_OBJECT
public:
    explicit MenuWorker(DBusDock *dockInter,QWidget *parent = nullptr);
    ~ MenuWorker();

    void initMember();
    void initUI();
    void initConnection();

    void showDockSettingsMenu();

    inline bool menuEnable() const { return m_menuEnable; }
    inline quint8 Opacity() const { return quint8(m_opacity * 255); }

    void onGSettingsChanged(const QString &key);
    // TODO 是否还有其他的插件未处理其gsettings配置,这里只是移植之前的代码
    void onTrashGSettingsChanged(const QString &key);

private:
    void setSettingsMenu();

signals:
    void autoHideChanged(const bool autoHide) const;
    void trayCountChanged();

private slots:
    void menuActionClicked(QAction *action);
    void trayVisableCountChanged(const int &count);
    void gtkIconThemeChanged();

public slots:
    void setAutoHide(const bool autoHide);

private:
    DockItemManager *m_itemManager;
    DBusDock *m_dockInter;

    QMenu *m_settingsMenu;
    QMenu *m_hideSubMenu;
    QAction *m_fashionModeAct;
    QAction *m_efficientModeAct;
    QAction *m_topPosAct;
    QAction *m_bottomPosAct;
    QAction *m_leftPosAct;
    QAction *m_rightPosAct;
    QAction *m_keepShownAct;
    QAction *m_keepHiddenAct;
    QAction *m_smartHideAct;
    QAction *m_modeSubMenuAct;
    QAction *m_locationSubMenuAct;
    QAction *m_statusSubMenuAct;
    QAction *m_hideSubMenuAct;

    bool m_menuEnable;
    bool m_autoHide;
    bool m_trashPluginShow;
    double m_opacity;
};

#endif // MENUWORKER_H
