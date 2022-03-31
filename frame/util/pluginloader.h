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

#ifndef PLUGINLOADER_H
#define PLUGINLOADER_H

#include <QThread>

class PluginLoader : public QThread
{
    Q_OBJECT

public:
    explicit PluginLoader(const QString &pluginDirPath, QObject *parent);
    static QString libUsedDtkCoreFileName(const QString &fileName);
    /**
     * @brief realFileName 获取软连接的真实文件的路径
     * @param fileName 文件地址
     * @return
     */
    static QString realFileName(QString fileName);

signals:
    void finished() const;
    void pluginFounded(const QString &pluginFile) const;

protected:
    void run();

    QString dtkCoreFileName();

private:
    QString m_pluginDirPath;
};

#endif // PLUGINLOADER_H
