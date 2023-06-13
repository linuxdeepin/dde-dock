// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef WINDOWINFOBASE_H
#define WINDOWINFOBASE_H

#include "processinfo.h"
#include "xcbutils.h"

#include <QString>
#include <QVector>
#include <qobject.h>
#include <qobjectdefs.h>
#include <qscopedpointer.h>

class Entry;
class AppInfo;

class WindowInfoBase : public QObject
{
    Q_OBJECT
public:
    WindowInfoBase(QObject *parent = nullptr) : QObject(parent), entry(nullptr), app(nullptr), m_processInfo(nullptr) {}
    virtual ~WindowInfoBase() {};

    virtual bool shouldSkip() = 0;
    virtual QString getIcon() = 0;
    virtual QString getTitle() = 0;
    virtual bool isDemandingAttention() = 0;
    virtual void close(uint32_t timestamp) = 0;
    virtual void activate() = 0;
    virtual void minimize() = 0;
    virtual bool isMinimized() = 0;
    virtual int64_t getCreatedTime() = 0;
    virtual QString getWindowType() = 0;
    virtual QString getDisplayName() = 0;
    virtual bool allowClose() = 0;
    virtual void update() = 0;
    virtual void killClient() = 0;
    virtual QString uuid() = 0;
    virtual QString getInnerId() { return innerId; }

    XWindow getXid() {return xid;}
    void setEntry(Entry *value) { entry = value; }
    Entry *getEntry() { return entry; }
    QString getEntryInnerId() { return entryInnerId; }
    void setEntryInnerId(QString value) { entryInnerId = value; }
    AppInfo *getAppInfo() { return app; }
    void setAppInfo(AppInfo *value) { app = value; }
    int getPid() { return pid; }
    ProcessInfo *getProcess() { return m_processInfo.data(); }
    bool containAtom(QVector<XCBAtom> supports, XCBAtom ty) {return supports.indexOf(ty) != -1;}

protected:
    XWindow xid;            // 窗口id
    QString title;          // 窗口标题
    QString icon;           // 窗口图标
    int pid;                // 窗口所属应用进程
    QString entryInnerId;   // 窗口所属应用对应的innerId
    QString innerId;        // 窗口对应的innerId
    Entry *entry;           // 窗口所属应用
    AppInfo *app;           // 窗口所属应用对应的desktopFile信息
    int64_t m_createdTime;    // 创建时间
    QScopedPointer<ProcessInfo> m_processInfo; // 窗口所属应用的进程信息
};

#endif // WINDOWINFOBASE_H
