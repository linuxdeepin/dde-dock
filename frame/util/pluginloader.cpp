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

#include <DSysInfo>

DCORE_USE_NAMESPACE

PluginLoader::PluginLoader(const QString &pluginDirPath, QObject *parent)
    : QThread(parent)
    , m_pluginDirPath(pluginDirPath)
{
}

void PluginLoader::run()
{
    QDir pluginsDir(m_pluginDirPath);
    const QStringList files = pluginsDir.entryList(QDir::Files);

    auto getDisablePluginList = [ = ] {
        if (QGSettings::isSchemaInstalled("com.deepin.dde.dock.disableplugins")) {
            QGSettings gsetting("com.deepin.dde.dock.disableplugins", "/com/deepin/dde/dock/disableplugins/");
            return gsetting.get("disable-plugins-list").toStringList();
        }
        return QStringList();
    };

    const QStringList disable_plugins_list = getDisablePluginList();

    const QString dtkCoreName = dtkCoreFileName();

    QStringList plugins;

    // 查找可用插件
    for (QString file : files) {
        if (!QLibrary::isLibrary(file))
            continue;

        // 社区版需要加载键盘布局，其他不需要
        if (file.contains("libkeyboard-layout") && !DSysInfo::isCommunityEdition())
            continue;

        // TODO: old dock plugins is uncompatible
        if (file.startsWith("libdde-dock-"))
            continue;

        if (disable_plugins_list.contains(file)) {
            qDebug() << "disable loading plugin:" << file;
            continue;
        }
        // 判断当前进程加载的dtkcore库是否和当前库加载的dtkcore的库为同一个，如果不同，则无需加载，
        // 否则在加载dtkcore的时候会报错（dtkcore内部判断如果加载两次会抛出错误）
        QString libUsedDtkCore = libUsedDtkCoreFileName(pluginsDir.absoluteFilePath(file));
        if (!libUsedDtkCore.isEmpty() && libUsedDtkCore != dtkCoreName)
            continue;

        plugins << file;
    }
    for (auto plugin : plugins) {
        emit pluginFounded(pluginsDir.absoluteFilePath(plugin));
    }

    emit finished();
}

/**
 * @brief 获取当前进程使用的dtkcore库的文件名
 * @return 当前进程使用的dtkCore库的文件名
 */
QString PluginLoader::dtkCoreFileName()
{
    QFile f("/proc/self/maps");
    if (!f.open(QIODevice::ReadOnly))
        return QString();

    const QByteArray &data = f.readAll();
    QTextStream ts(data);
    while (Q_UNLIKELY(!ts.atEnd())) {
        const QString line = ts.readLine();
        const QStringList &maps = line.split(' ', QString::SplitBehavior::SkipEmptyParts);
        if (Q_UNLIKELY(maps.size() < 6))
            continue;

        QFileInfo info(maps.value(5));
        QString infoAbPath = info.absoluteFilePath();
        if (info.fileName().contains("dtkcore")) {
            infoAbPath = realFileName(infoAbPath);
            if (infoAbPath.contains("/")) {
                int pathIndex = infoAbPath.lastIndexOf("/");
                infoAbPath = infoAbPath.mid(pathIndex + 1).trimmed();
            }

            f.close();
            return infoAbPath;
        }
    }

    f.close();
    return QString();
}

/**
 * @brief 返回某个so库使用的dtkcore库文件名
 * @param 用于获取dtkcore库的so库文件名
 * @return 返回使用的dtkcore库文件名
 */
QString PluginLoader::libUsedDtkCoreFileName(const QString &fileName)
{
    QString lddCommand = QString("ldd -r %1").arg(fileName);
    FILE *fp = popen(lddCommand.toLocal8Bit().data(), "r");
    if (fp) {
        char buf[2000] = {0};
        while (fgets(buf, sizeof(buf), fp)) {
            QString rowResult = buf;
            rowResult = rowResult.trimmed();
            if (rowResult.contains("dtkcore")) {
                QStringList coreSplits = rowResult.split("=>");
                if (coreSplits.size() < 2)
                    continue;

                pclose(fp);
                QString dtkCorePath = coreSplits[1];

                // 删除后面的括号的内容
                int addrIndex = dtkCorePath.indexOf("(0x");
                dtkCorePath = realFileName(dtkCorePath.left(addrIndex).trimmed());
                if (dtkCorePath.contains("/")) {
                    int pathIndex = dtkCorePath.lastIndexOf("/");
                    dtkCorePath = dtkCorePath.mid(pathIndex + 1).trimmed();
                }

                return dtkCorePath;
            }
        }
        pclose(fp);
    }
    return QString();
}

/**
 * @brief 返回软连接对应的实际的文件名
 * @param 当前软连接的文件名
 * @return 软连接对应的库的实际的文件名
 */
QString PluginLoader::realFileName(QString fileName)
{
    QString command = QString("ls %1 -al").arg(fileName);
    FILE *fp = popen(command.toLocal8Bit().data(), "r");
    if (fp) {
        char buf[2000] = {0};
        if (fgets(buf, sizeof(buf), fp)) {
            QString rowInfo = buf;
            if (rowInfo.contains("->")) {
                pclose(fp);
                QStringList fileInfos = rowInfo.split("->");
                QString srcFileName = fileInfos[1];
                srcFileName = srcFileName.trimmed();
                return srcFileName;
            }
        }
        pclose(fp);
    }
    return fileName;
}
