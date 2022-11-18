/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
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
#ifndef DBUSUTIL_H
#define DBUSUTIL_H

#include "dockinterface.h"
#include "entryinterface.h"

using DockInter = org::deepin::dde::daemon::DdeDock;
using DockEntryInter = org::deepin::dde::daemon::dock::DockEntry;

const QString xEventMonitorService = "org.deepin.dde.XEventMonitor1";
const QString xEventMonitorPath = "/org/deepin/dde/XEventMonitor1";

const QString launcherService = "org.deepin.dde.Launcher1";
const QString launcherPath = "/org/deepin/dde/Launcher1";
const QString launcherInterface = "org.deepin.dde.Launcher1";

const QString controllCenterService = "org.deepin.dde.ControlCenter1";
const QString controllCenterPath = "/org/deepin/dde/ControlCenter1";
const QString controllCenterInterface = "org.deepin.dde.ControlCenter1";

const QString notificationService = "org.deepin.dde.Notification1";
const QString notificationPath = "/org/deepin/dde/Notification1";
const QString notificationInterface = "org.deepin.dde.Notification1";

const QString sessionManagerService = "org.deepin.dde.SessionManager1";
const QString sessionManagerPath = "/org/deepin/dde/SessionManager1";
const QString sessionManagerInterface = "org.deepin.dde.SessionManager1";

inline const QString dockServiceName()
{
    return QString("org.deepin.dde.daemon.Dock1");
}

inline const QString dockServicePath()
{
    return QString("/org/deepin/dde/daemon/Dock1");
}

#endif // DBUSUTIL_H
