// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef PROCESSINFO_H
#define PROCESSINFO_H

#include <QMap>
#include <QVector>
#include <QStringList>

typedef QMap<QString, QString> Status;

// 进程信息
class ProcessInfo
{
public:
    explicit ProcessInfo(int pid);
    explicit ProcessInfo(QStringList cmd);
    virtual ~ProcessInfo();

    bool isValid();
    bool initWithPid();

    int getPid();
    int getPpid();

    QString getExe();
    QString getCwd();
    Status getStatus();
    QString getEnv(const QString &key);

    QStringList getArgs();
    QStringList getCmdLine();
    QMap<QString, QString> getEnviron();

private:
    bool isExist();
    QString getJoinedExeArgs();
    QString getFile(const QString &file);
    QStringList readFile(const QString &filePath);

private:

    int m_pid;
    int m_ppid;
    bool m_hasPid;
    bool m_isValid;

    Status m_status;
    QString m_exe;
    QString m_cwd;
    QStringList m_args;
    QStringList m_cmdLine;
    QVector<int> m_uids;
    QMap<QString, QString> m_environ;
};

#endif // PROCESSINFO_H
