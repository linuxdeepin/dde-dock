// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "appinfo.h"
#include "common.h"

#include <QDebug>
#include <QString>
#include <QCryptographicHash>

AppInfo::AppInfo(DesktopInfo &info)
 : m_installed(false)
 , m_isValid(true)
{
    init(info);
}

AppInfo::AppInfo(const QString &_fileName)
 : m_isValid(true)
{
    DesktopInfo info(_fileName);
    init(info);
}

void AppInfo::init(DesktopInfo &info)
{
    if (!info.isValidDesktop()) {
        m_isValid = false;
        return;
    }

    QString xDeepinVendor = info.getDesktopFile()->value(MainSection + "/X-Deepin-Vendor").toString();
    if (xDeepinVendor == "deepin") {
        m_name = info.getGenericName();
        if (m_name.isEmpty()) {
            m_name = info.getName();
        }
    } else {
        m_name = info.getName();
    }

    m_innerId = genInnerIdWithDesktopInfo(info);
    m_fileName = info.getDesktopFilePath();
    m_id = info.getId();
    m_icon = info.getIcon();
    m_installed = info.isInstalled();
    auto actions = info.getActions();
    std::copy(actions.begin(), actions.end(), std::back_inserter(m_actions));

}

QString AppInfo::genInnerIdWithDesktopInfo(DesktopInfo &info)
{
    QString cmdline = info.getCommandLine();
    QByteArray encryText = QCryptographicHash::hash(QString(cmdline).toLatin1(), QCryptographicHash::Md5);
    QString innerId = desktopHashPrefix + encryText.toHex();
    qInfo() << "app: " << info.getId() << " generate innerId: " << innerId;
    return innerId;
}
