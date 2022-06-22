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

#ifdef USE_AM
#include "dockinterface.h"
#include "entryinterface.h"
#else
#include <com_deepin_dde_daemon_dock.h>
#include <com_deepin_dde_daemon_dock_entry.h>
#endif

#ifdef USE_AM
using DockInter = org::deepin::dde::daemon::DdeDock;
using DockEntryInter = org::deepin::dde::daemon::dock::DockEntry;
#else
using DockInter = com::deepin::dde::daemon::Dock;
using DockEntryInter = com::deepin::dde::daemon::dock::Entry;
#endif

inline const QString dockServiceName()
{
#ifdef USE_AM
    return QString("org.deepin.dde.daemon.Dock1");
#else
    return QString("com.deepin.dde.daemon.Dock");
#endif
}

inline const QString dockServicePath()
{
#ifdef USE_AM
    return QString("/org/deepin/dde/daemon/Dock1");
#else
    return QString("/com/deepin/dde/daemon/Dock");
#endif
}

#endif // DBUSUTIL_H
