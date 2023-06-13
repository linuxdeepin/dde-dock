// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "processinfo.h"

#include <string>
#include <fstream>
#include <unistd.h>

#include <QDir>
#include <QDebug>
#include <QFileInfo>

ProcessInfo::ProcessInfo(int pid)
    : m_pid(pid)
    , m_ppid(0)

{
    if (pid == 0)
        return;
    m_exe = getExe();
    m_cwd = getCwd();
    m_cmdLine = getCmdLine();
    getStatus();
    // 部分root进程在/proc文件系统查找不到exe、cwd、cmdline信息
    if (m_exe.isEmpty() || m_cwd.isEmpty() || m_cmdLine.size() == 0) {
        m_isValid = false;
        return;
    }

    // args
    qInfo() << "ProcessInfo: exe=" << m_exe << " cwd=" << m_cwd << " cmdLine=" << (m_cmdLine[0].isEmpty() ? " " : m_cmdLine[0]);
    auto verifyExe =  [](QString exe, QString cwd, QString firstArg){
        if (firstArg.size() == 0) return false;

        QFileInfo info(firstArg);
        if (info.completeBaseName() == firstArg) return true;

        if (!QDir::isAbsolutePath(firstArg))
            firstArg = cwd + firstArg;

        return exe == firstArg;
    };

    if (!m_cmdLine[0].isEmpty()) {
        if (!verifyExe(m_exe, m_cwd, m_cmdLine[0])) {
            auto parts = m_cmdLine[0].split(' ');
            // try again
            if (verifyExe(m_exe, m_cwd, parts[0])) {
                m_args.append(parts.mid(1, parts.size() - 1));
                m_args.append(m_cmdLine.mid(1, m_cmdLine.size() - 1));
            }
        } else {
            m_args.append(m_cmdLine.mid(1, m_cmdLine.size() - 1));
        }
    }
}

ProcessInfo::ProcessInfo(QStringList cmd)
    : m_hasPid(false)
    , m_isValid(true)
{
    if (cmd.size() == 0) {
        m_isValid = false;
        return;
    }

    m_cmdLine = cmd;
    m_exe = cmd[0];
    m_args.append(cmd.mid(1, cmd.size() - 1));
}

ProcessInfo::~ProcessInfo()
{
}

QString ProcessInfo::getEnv(const QString &key)
{
    if (m_environ.size() == 0) getEnviron();
    return m_environ[key];
}

Status ProcessInfo::getStatus()
{
    if (!m_status.empty()){
        return m_status;
    }

    QString statusFile = getFile("status");

    std::ifstream fs(statusFile.toStdString());
    if (!fs.is_open()) {
        return m_status;
    }

    std::string tmp = "";
    while (std::getline(fs, tmp)) {
        auto pos = tmp.find_first_of(':');
        if (pos == std::string::npos) {
            continue;
        }

        QString value;
        if (pos + 1 < tmp.length()) {
            value = QString::fromStdString(tmp.substr(pos + 1));
        }

        m_status[QString::fromStdString(tmp.substr(0, pos))] = value;
    }

    return m_status;
}

QStringList ProcessInfo::getCmdLine()
{
    if (m_cmdLine.size() == 0) {
        QString cmdlineFile = getFile("cmdline");
        m_cmdLine = readFile(cmdlineFile);
    }

    return m_cmdLine;
}

QStringList ProcessInfo::getArgs()
{
    return m_args;
}

int ProcessInfo::getPid()
{

    return m_pid;
}

int ProcessInfo::getPpid()
{
    if (m_ppid == 0) {
        if (m_status.find("PPid") != m_status.end()) {
            m_ppid = std::stoi(m_status["PPid"].toStdString());
        }
    }
    return m_ppid;
}

bool ProcessInfo::initWithPid()
{
    return m_hasPid;
}

bool ProcessInfo::isValid()
{
    return m_isValid;
}

QString ProcessInfo::getExe()
{
    if (m_exe.isEmpty()) {
        QString cmdLineFile = getFile("exe");
        QFileInfo path(cmdLineFile);
        m_exe = path.canonicalFilePath();
    }

    return m_exe;
}

bool ProcessInfo::isExist()
{
    QString procDir = "/proc/" + QString(m_pid);
    return QFile::exists(procDir);
}

QStringList ProcessInfo::readFile(const QString &filePath)
{
    QStringList ret;
    std::ifstream fs(filePath.toStdString());
    if (!fs.is_open()) {
            return ret;
    }

    std::string tmp;
    while (std::getline(fs, tmp, '\0')) {
        ret.append(QString::fromStdString(tmp));
    }
    return ret;
}

QString ProcessInfo::getFile(const QString &file)
{
    return QString("/proc/").append(QString::number(m_pid).append('/').append(file));
}

QString ProcessInfo::getCwd()
{
    if (m_cwd.isEmpty()) {
        QString cwdFile = getFile("cwd");
        QFileInfo path(cwdFile);
        m_cwd = path.canonicalFilePath();
    }
    return m_cwd;
}

QMap<QString, QString> ProcessInfo::getEnviron()
{
    if (m_environ.size() == 0) {
        QString envFile = getFile("environ");
        QStringList contents = readFile(envFile);
        for (auto line : contents){
            int index = line.indexOf('=');
            m_environ.insert(line.left(index), line.right(line.size() - index - 1));
        }
    }
    return m_environ;
}