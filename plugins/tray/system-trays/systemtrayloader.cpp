/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#include "systemtrayloader.h"
#include "systemtrayscontroller.h"

#include <QDebug>

SystemTrayLoader::SystemTrayLoader(QObject *parent)
    : QThread(parent)
{

}

void SystemTrayLoader::run()
{
#ifdef QT_DEBUG
    const QDir pluginsDir("../plugins/system-trays");
#else
    const QDir pluginsDir("../lib/dde-dock/plugins/system-trays");
#endif
    const QStringList plugins = pluginsDir.entryList(QDir::Files);

    for (const QString file : plugins)
    {
        if (!QLibrary::isLibrary(file))
            continue;

        // TODO: old dock plugins is uncompatible
        if (file.startsWith("libdde-dock-"))
            continue;

        emit pluginFounded(pluginsDir.absoluteFilePath(file));

        msleep(500);
    }

    emit finished();
}
