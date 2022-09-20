// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef PLUGINLOADER_H
#define PLUGINLOADER_H

#include <QThread>

class PluginLoader : public QThread
{
    Q_OBJECT

public:
    explicit PluginLoader(const QStringList &pluginDirPaths, QObject *parent);
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
    QStringList m_pluginDirPaths;
};

#endif // PLUGINLOADER_H
