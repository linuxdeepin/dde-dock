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
#include <QProcess>

PluginLoader::PluginLoader(const QString &pluginDirPath, QObject *parent)
    : QThread(parent)
    , m_pluginDirPath(pluginDirPath)
    , m_isPanguV(false)
{
    //判断当前是否为盘古v机器，如果为盘古v则退出初始化飞行模式插件
    QProcess process;
    process.start("gdbus introspect -y -d com.deepin.system.SystemInfo -o /com/deepin/system/SystemInfo -p");
    process.waitForFinished();
    QString pcType = process.readAllStandardOutput();
    process.close();
    if (pcType.contains("panguV"))
        m_isPanguV = true;
}

void PluginLoader::run()
{
    QDir pluginsDir(m_pluginDirPath);
    const QStringList plugins = pluginsDir.entryList(QDir::Files);

    for (const QString file : plugins)
    {
        if (!QLibrary::isLibrary(file))
            continue;

        // TODO: old dock plugins is uncompatible
        if (file.startsWith("libdde-dock-"))
            continue;

        if (m_isPanguV && file.contains("airplane-mode"))
            continue;

        emit pluginFounded(pluginsDir.absoluteFilePath(file));
    }

    emit finished();
}
