// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef DESKTOPINFO_H
#define DESKTOPINFO_H

#include <QSettings>
#include <qscopedpointer.h>
#include <string>
#include <vector>

const QString MainSection        = "Desktop Entry";
const QString KeyType            = "Type";
const QString KeyVersion         = "Version";
const QString KeyName            = "Name";
const QString KeyGenericName     = "GenericName";
const QString KeyNoDisplay       = "NoDisplay";
const QString KeyComment         = "Comment";
const QString KeyIcon            = "Icon";
const QString KeyHidden          = "Hidden";
const QString KeyOnlyShowIn      = "OnlyShowIn";
const QString KeyNotShowIn       = "NotShowIn";
const QString KeyTryExec         = "TryExec";
const QString KeyExec            = "Exec";
const QString KeyPath            = "Path";
const QString KeyTerminal        = "Terminal";
const QString KeyMimeType        = "MimeType";
const QString KeyCategories      = "Categories";
const QString KeyKeywords        = "Keywords";
const QString KeyStartupNotify   = "StartupNotify";
const QString KeyStartupWMClass  = "StartupWMClass";
const QString KeyURL             = "URL";
const QString KeyActions         = "Actions";
const QString KeyDBusActivatable = "DBusActivatable";

const QString TypeApplication    = "Application";
const QString TypeLink           = "Link";
const QString TypeDirectory      = "Directory";

typedef struct DesktopAction
{
    DesktopAction()
        : section("")
        , name("")
        , exec("")
    {
    }
    QString section;
    QString name;
    QString exec;
} DesktopAction;

// 应用Desktop信息类
class DesktopInfo {
public:
    explicit DesktopInfo(const QString &desktopfile);
    ~DesktopInfo();

    static bool isDesktopAction(const QString &name);
    static DesktopInfo getDesktopInfoById(const QString &appId);

    bool shouldShow();
    bool getIsHidden();
    bool isInstalled();
    bool getTerminal();
    bool getNoDisplay();
    bool isExecutableOk();
    bool isValidDesktop();
    bool getShowIn(QStringList desktopEnvs);

    void setDesktopOverrideExec(const QString &execStr);

    QString getId();
    QString getName();
    QString getIcon();
    QString getExecutable();
    QString getGenericName();
    QString getCommandLine();
    QString getDesktopFilePath();

    QStringList getKeywords();
    QStringList getCategories();

    QList<DesktopAction> getActions();

    QSettings *getDesktopFile();

private:
    bool findExecutable(const QString &exec);
    QString getTryExec();
    QString getLocaleStr(const QString &section, const QString &key);

private:
    static QStringList currentDesktops;

    bool m_isValid;
    bool m_isInstalled;

    QString m_id;
    QString m_name;
    QString m_icon;
    QString m_desktopFilePath;

    // Desktopfile ini format
    QScopedPointer<QSettings> m_desktopFile;
    
};
#endif // DESKTOPINFO_H
