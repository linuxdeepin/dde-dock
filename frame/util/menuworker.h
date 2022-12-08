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

#include "constants.h"
#include "dbusutil.h"

#include <QObject>

class QMenu;
class QGSettings;
/**
 * @brief The MenuWorker class  此类用于处理任务栏右键菜单的逻辑
 */
class MenuWorker : public QObject
{
    Q_OBJECT
public:
    explicit MenuWorker(DockInter *dockInter, QObject *parent = nullptr);

    void exec();

private:
    void createMenu(QMenu *settingsMenu);

private slots:
    void onDockSettingsTriggered();

private:
    DockInter *m_dockInter;
};

#endif // MENUWORKER_H
