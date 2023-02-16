// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef PLUGINLOADER_H
#define PLUGINLOADER_H

#include <QThread>

class PluginLoader : public QThread
{
    Q_OBJECT

public:
    explicit PluginLoader(const QString &pluginDirPath, QObject *parent);

signals:
    void finished() const;
    void pluginFounded(const QString &pluginFile) const;

protected:
    void run();

private:
    QString m_pluginDirPath;
};

#endif // PLUGINLOADER_H
