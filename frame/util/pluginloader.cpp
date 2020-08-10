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

#include "pluginloader.h"

#include <QDir>
#include <QDebug>
#include <QLibrary>
#include <QGSettings>

PluginLoader::PluginLoader(const QString &pluginDirPath, QObject *parent)
    : QThread(parent)
    , m_pluginDirPath(pluginDirPath)
{
}

void PluginLoader::run()
{
    QDir pluginsDir(m_pluginDirPath);
    const QStringList plugins = pluginsDir.entryList(QDir::Files);
    static const QGSettings gsetting("com.deepin.dde.dock.disableplugins", "/com/deepin/dde/dock/disableplugins/");
    static const auto disable_plugins_list = gsetting.get("disable-plugins-list").toStringList();

    for (const QString& file : plugins)
    {
        if (!QLibrary::isLibrary(file)){
#ifdef QT_DEBUG
            qDebug() << "--------not a library: " << pluginsDir.absoluteFilePath(file);
#endif
            continue;
        }

        // TODO: old dock plugins is uncompatible
        if (file.startsWith("libdde-dock-")) {
#ifdef QT_DEBUG
            qDebug() << "--------uncompatible library: " << pluginsDir.absoluteFilePath(file);
#endif
            continue;
        }

        if (disable_plugins_list.contains(file)) {
            qDebug() << "disable loading plugin:" << pluginsDir.absoluteFilePath(file);
            continue;
        }

#ifdef QT_DEBUG
        qDebug() << "----------loaded: " << pluginsDir.absoluteFilePath(file);
#endif
        emit pluginFounded(pluginsDir.absoluteFilePath(file));
    }

    emit finished();
}
