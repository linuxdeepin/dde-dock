// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef APPINFO_H
#define APPINFO_H

#include "desktopinfo.h"

#include <QVector>

// 应用信息类
class AppInfo
{
public:
    explicit AppInfo(DesktopInfo &info);
    explicit AppInfo(const QString &_fileName);

    void init(DesktopInfo &info);
    void setIdentifyMethod(QString method) {m_identifyMethod = method;}

    bool isValidApp() {return m_isValid;}
    bool isInstalled() {return m_installed;}
    

    QString getId() {return m_id;}
    QString getIcon() {return m_icon;}
    QString getName() {return m_name;}
    QString getInnerId() {return m_innerId;}
    QString getFileName() {return m_fileName;}
    QString getIdentifyMethod() {return m_identifyMethod;}

    QVector<DesktopAction> getActions() {return m_actions;}

private:
    QString genInnerIdWithDesktopInfo(DesktopInfo &info);

private:

    bool m_installed;
    bool m_isValid;
    
    QString m_id;
    QString m_name;
    QString m_icon;
    QString m_innerId;
    QString m_fileName;
    QString m_identifyMethod;
    QVector<DesktopAction> m_actions;

};

#endif // APPINFO_H
