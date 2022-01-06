/*
 * Copyright (C) 2018 ~ 2020 Deepin Technology Co., Ltd.
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
class QGSettings;
/**
 * @brief The MenuWorker class  此类用于处理任务栏右键菜单的逻辑
 */
class MenuWorker : public QObject
{
    Q_OBJECT
public:
    explicit MenuWorker(DBusDock *dockInter,QWidget *parent = nullptr);

    void showDockSettingsMenu(QMenu *menu);

signals:
    void autoHideChanged(const bool autoHide) const;

public slots:
    void setAutoHide(const bool autoHide);

private:
    QMenu *createMenu(QMenu *settingsMenu);

private slots:
    void onDockSettingsTriggered();

private:
    DBusDock *m_dockInter;
    bool m_autoHide;
};

#endif // MENUWORKER_H
