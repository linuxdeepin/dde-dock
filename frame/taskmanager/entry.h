// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef ENTRY_H
#define ENTRY_H

#include "appinfo.h"
#include "appmenu.h"
#include "windowinfomap.h"
#include "windowinfobase.h"

#include <QMap>
#include <QVector>
#include <QObject>
#include <qscopedpointer.h>

#define ENTRY_NONE      0
#define ENTRY_NORMAL    1
#define ENTRY_RECENT    2

// 单个应用类
class TaskManager;
class DBusAdaptorEntry;
class WindowInfo;

typedef QMap<quint32, WindowInfo> WindowInfoMap;

class Entry: public QObject
{
    Q_OBJECT
public:
    Entry(TaskManager *_taskmanager, AppInfo *_app, QString _innerId, QObject *parent = nullptr);
    ~Entry();

    void updateName();
    void updateMenu();
    void updateIcon();
    void updateMode();
    void updateIsActive();
    void forceUpdateIcon();
    void updateExportWindowInfos();
    void launchApp(uint32_t timestamp);

    void setIsDocked(bool value);
    void setMenu(AppMenu *_menu);
    void setPropIcon(QString value);
    void setPropName(QString value);
    void setPropIsActive(bool active);
    void setInnerId(QString _innerId);
    void setAppInfo(AppInfo *appinfo);
    void setPropCurrentWindow(XWindow value);
    void setCurrentWindowInfo(WindowInfoBase *windowInfo);

    void check();
    void forceQuit();
    void presentWindows();
    void active(uint32_t timestamp);
    void activeWindow(quint32 winId);
    void newInstance(uint32_t timestamp);
    void requestDock(bool dockToEnd = false);
    void requestUndock(bool dockToEnd = false);
    void handleMenuItem(uint32_t timestamp, QString itemId);
    void handleDragDrop(uint32_t timestamp, QStringList files);

    bool containsWindow(XWindow xid);
    bool detachWindow(WindowInfoBase *info);
    bool attachWindow(WindowInfoBase *info);

    bool getIsDocked() const;
    bool getIsActive() const;

    QString getId() const;
    QString getMenu() const;

    bool isValid();
    bool hasWindow();

    int mode();

    QString getName();
    QString getIcon();
    QString getInnerId();
    QString getFileName();
    QString getDesktopFile();
    QString getExec();
    QString getCmdLine();

    XWindow getCurrentWindow();

    AppInfo *getAppInfo();

    WindowInfoBase *findNextLeader();
    WindowInfoBase *getCurrentWindowInfo();
    WindowInfoBase *getWindowInfoByPid(int pid);
    WindowInfoBase *getWindowInfoByWinId(XWindow windowId);

    WindowInfoMap getExportWindowInfos();
    QVector<XWindow> getAllowedClosedWindowIds();

public Q_SLOTS:
    QVector<WindowInfoBase *> getAllowedCloseWindows();

Q_SIGNALS:
    void modeChanged(int);
    void isActiveChanged(bool);
    void isDockedChanged(bool);
    void menuChanged(QString);
    void iconChanged(QString);
    void nameChanged(QString);
    void desktopFileChanged(QString);
    void currentWindowChanged(uint32_t);
    void windowInfosChanged(const WindowInfoMap&);

private:
    // 右键菜单项
    bool killProcess(int pid);
    bool setPropDesktopFile(QString value);
    bool isShowOnDock() const;
    int getCurrentMode();

    AppMenuItem getMenuItemLaunch();
    AppMenuItem getMenuItemCloseAll();
    AppMenuItem getMenuItemForceQuit();
    
    AppMenuItem getMenuItemDock();
    AppMenuItem getMenuItemUndock();
    AppMenuItem getMenuItemAllWindows();
    AppMenuItem getMenuItemForceQuitAndroid();
    QVector<AppMenuItem> getMenuItemDesktopActions();

private:
    bool m_isActive;
    bool m_isValid;
    bool m_isDocked;
    bool m_winIconPreferred;
    int m_mode;

    QString m_id;
    QString m_name;
    QString m_icon;
    QString m_innerId;
    QString m_desktopFile;

    DBusAdaptorEntry *m_adapterEntry;
    TaskManager *m_taskmanager;
    WindowInfoMap m_exportWindowInfos;      // 该应用导出的窗口属性
    WindowInfoBase *m_current; // 当前窗口
    XWindow m_currentWindow; //当前窗口Id

    QScopedPointer<AppInfo> m_appInfo;
    QScopedPointer<AppMenu> m_appMenu;
    QMap<XWindow, WindowInfoBase *> m_windowInfoMap; // 该应用所有窗口
};

#endif // ENTRY_H
