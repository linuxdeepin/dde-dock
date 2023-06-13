// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#ifndef WINDOWIDENTIFY_H
#define WINDOWIDENTIFY_H

#include "taskmanager/entry.h"
#include "windowpatterns.h"
#include "windowinfok.h"
#include "windowinfox.h"

#include <QObject>
#include <QVector>
#include <QMap>

class AppInfo;
class TaskManager;

typedef AppInfo *(*IdentifyFunc)(TaskManager *, WindowInfoX*, QString &innerId);

// 应用窗口识别类
class WindowIdentify : public QObject
{
    Q_OBJECT

public:
    explicit WindowIdentify(TaskManager *_taskmanager, QObject *parent = nullptr);

    AppInfo *identifyWindow(WindowInfoBase *winInfo, QString &innerId);
    AppInfo *identifyWindowX11(WindowInfoX *winInfo, QString &innerId);
    AppInfo *identifyWindowWayland(WindowInfoK *winInfo, QString &innerId);

    static AppInfo *identifyWindowAndroid(TaskManager *_dock, WindowInfoX *winInfo, QString &innerId);
    static AppInfo *identifyWindowByPidEnv(TaskManager *_dock, WindowInfoX *winInfo, QString &innerId);
    static AppInfo *identifyWindowByCmdlineTurboBooster(TaskManager *_dock, WindowInfoX *winInfo, QString &innerId);
    static AppInfo *identifyWindowByCmdlineXWalk(TaskManager *_dock, WindowInfoX *winInfo, QString &innerId);
    static AppInfo *identifyWindowByFlatpakAppID(TaskManager *_dock, WindowInfoX *winInfo, QString &innerId);
    static AppInfo *identifyWindowByCrxId(TaskManager *_dock, WindowInfoX *winInfo, QString &innerId);
    static AppInfo *identifyWindowByRule(TaskManager *_dock, WindowInfoX *winInfo, QString &innerId);
    static AppInfo *identifyWindowByBamf(TaskManager *_dock, WindowInfoX *winInfo, QString &innerId);
    static AppInfo *identifyWindowByPid(TaskManager *_dock, WindowInfoX *winInfo, QString &innerId);
    static AppInfo *identifyWindowByScratch(TaskManager *_dock, WindowInfoX *winInfo, QString &innerId);
    static AppInfo *identifyWindowByGtkAppId(TaskManager *_dock, WindowInfoX *winInfo, QString &innerId);
    static AppInfo *identifyWindowByWmClass(TaskManager *_dock, WindowInfoX *winInfo, QString &innerId);

private:
    AppInfo *fixAutostartAppInfo(QString fileName);
    static int32_t getAndroidUengineId(XWindow winId);
    static QString getAndroidUengineName(XWindow winId);

private:
    TaskManager *m_taskmanager;
    QList<QPair<QString, IdentifyFunc>> m_identifyWindowFuns;
};

#endif // IDENTIFYWINDOW_H
