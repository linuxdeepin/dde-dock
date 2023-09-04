// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "mainwindow.h"
#include "dbusdockadaptors.h"
#include "dockdaemonadaptors.h"
#include "utils.h"
#include "themeappicon.h"
#include "dockapplication.h"
#include "traymainwindow.h"
#include "windowmanager.h"

#include <QDir>
#include <QStandardPaths>
#include <QDateTime>
#include <QDir>

#include <DApplication>
#include <DLog>
#include <DGuiApplicationHelper>

#include <unistd.h>
#include <string>
#include <sys/mman.h>
#include <stdio.h>
#include <time.h>
#include <execinfo.h>
#include <sys/stat.h>
#include <signal.h>

DWIDGET_USE_NAMESPACE
#ifdef DCORE_NAMESPACE
DCORE_USE_NAMESPACE
#else
DUTIL_USE_NAMESPACE
#endif

const QString g_cfgPath = QStandardPaths::standardLocations(QStandardPaths::ConfigLocation)[0] + "/dde-cfg.ini";

/**
 * @brief IsSaveMode
 * @return 判断当前是否应该进入安全模式（安全模式下不加载插件）
 */
bool IsSaveMode()
{
    QSettings settings(g_cfgPath, QSettings::IniFormat);
    settings.beginGroup(qApp->applicationName());
    int collapseNum = settings.value("collapse").toInt();
    /* 崩溃次数达到3次，进入安全模式（不加载插件） */
    if (collapseNum >= 3) {
        settings.remove(""); // 删除记录的数据
        settings.setValue("collapse", 0);
        settings.endGroup();
        settings.sync();
        return true;
    }
    return false;
}

/**
 * @brief sig_crash
 * @return 当应用收到对应的退出信号时，会调用此函数，用于保存一下应用崩溃时间，崩溃次数，用以判断是否应该进入安全模式，见IsSaveMode()
 */
[[noreturn]] void sig_crash(int sig)
{
    QDir dir(QStandardPaths::standardLocations(QStandardPaths::CacheLocation)[0]);
    dir.cdUp();
    QString filePath = dir.path() + "/dde-collapse.log";

    QFile *file = new QFile(filePath);

    // 创建默认配置文件,记录段时间内的崩溃次数
    if (!QFile::exists(g_cfgPath)) {
        QFile cfgFile(g_cfgPath);
        if (!cfgFile.open(QIODevice::WriteOnly))
            exit(0);
        cfgFile.close();
    }

    QSettings settings(g_cfgPath, QSettings::IniFormat);
    settings.beginGroup("dde-dock");

    int collapseNum = settings.value("collapse").toInt();
    /* 第一次崩溃或进入安全模式后的第一次崩溃，将时间重置 */
    if (collapseNum == 0) {
        settings.setValue("first_time", QDateTime::currentDateTime().toString("yyyy-MM-dd hh:mm:ss:zzz"));
    }
    QDateTime lastDate = QDateTime::fromString(settings.value("first_time").toString(), "yyyy-MM-dd hh:mm:ss:zzz");
    /* 将当前崩溃时间与第一次崩溃时间比较，小于9分钟，记录一次崩溃；大于9分钟，覆盖之前的崩溃时间 */
    if (qAbs(lastDate.secsTo(QDateTime::currentDateTime())) < 9 * 60) {
        settings.setValue("collapse", collapseNum + 1);
        switch (collapseNum) {
        case 0:
            settings.setValue("first_time", QDateTime::currentDateTime().toString("yyyy-MM-dd hh:mm:ss:zzz"));
            break;
        case 1:
            settings.setValue("second_time", QDateTime::currentDateTime().toString("yyyy-MM-dd hh:mm:ss:zzz"));
            break;
        case 2:
            settings.setValue("third_time", QDateTime::currentDateTime().toString("yyyy-MM-dd hh:mm:ss:zzz"));
            break;
        default:
            qDebug() << "Error, the collapse is wrong!";
            break;
        }
    } else {
        if (collapseNum == 2){
            settings.setValue("first_time", settings.value("second_time").toString());
            settings.setValue("second_time", QDateTime::currentDateTime().toString("yyyy-MM-dd hh:mm:ss:zzz"));
        } else {
            settings.setValue("first_time", QDateTime::currentDateTime().toString("yyyy-MM-dd hh:mm:ss:zzz"));
        }
    }

    settings.endGroup();
    settings.sync();

    if (!file->open(QIODevice::Text | QIODevice::Append)) {
        qDebug() << file->errorString();
        exit(0);
    }

    if (file->size() >= 10 * 1024 * 1024) {
        // 清空原有内容
        file->close();
        if (file->open(QIODevice::Text | QIODevice::Truncate)) {
            qDebug() << file->errorString();
            exit(0);
        }
    }

    // 捕获异常，打印崩溃日志到配置文件中
    try {
        QString head = "\n#####" + qApp->applicationName() + "#####\n"
                + QDateTime::currentDateTime().toString("[yyyy-MM-dd hh:mm:ss:zzz]")
                + "[crash signal number:" + QString::number(sig) + "]\n";
        file->write(head.toUtf8());
    } catch (...) {
        //
    }
    file->close();
    delete file;
    file = nullptr;
    exit(0);
}

int main(int argc, char *argv[])
{
    QString currentDesktop = QString(getenv("XDG_CURRENT_DESKTOP"));
    if (currentDesktop.compare("DDE", Qt::CaseInsensitive) == 0 ||
        currentDesktop.compare("deepin", Qt::CaseInsensitive) == 0) {
        qDebug() << "Warning: force enable D_DXCB_FORCE_NO_TITLEBAR now!";
        setenv("D_DXCB_FORCE_NO_TITLEBAR", "1", 1);
    }

    DGuiApplicationHelper::setAttribute(DGuiApplicationHelper::UseInactiveColorGroup, false);
    DockApplication::setAttribute(Qt::AA_EnableHighDpiScaling, true);
    DockApplication app(argc, argv);

    //崩溃信号
    signal(SIGSEGV, sig_crash);
    signal(SIGILL,  sig_crash);
    signal(SIGINT,  sig_crash);
    signal(SIGABRT, sig_crash);
    signal(SIGFPE,  sig_crash);

    app.setOrganizationName("deepin");
    app.setApplicationName("dde-dock");
    app.setApplicationDisplayName("DDE Dock");
    app.setApplicationVersion(CVERSION);
    app.loadTranslator();
    app.setAttribute(Qt::AA_UseHighDpiPixmaps, true);

    // 设置日志输出到控制台以及文件
    DLogManager::setLogFormat("%{time}{yyMMdd.HH:mm:ss.zzz}[%{type:1}] [%{category}] [%{function:-25} %{line:-4}] %{message}");
    DLogManager::registerConsoleAppender();
    DLogManager::registerFileAppender();

    // 启动入参 dde-dock --help可以看到一下内容， -x不加载插件 -r 一般用在startdde启动任务栏
    QCommandLineOption disablePlugOption(QStringList() << "x" << "disable-plugins", "do not load plugins.");
    QCommandLineOption runOption(QStringList() << "r" << "run-by-stardde", "run by startdde.");
    QCommandLineParser parser;
    parser.setApplicationDescription("DDE Dock");
    parser.addHelpOption();
    parser.addVersionOption();
    parser.addOption(disablePlugOption);
    parser.addOption(runOption);
    parser.process(app);

    // 任务栏单进程限制
    DGuiApplicationHelper::setSingleInstanceInterval(-1);
    if (!app.setSingleInstance(QString("dde-dock_%1").arg(getuid()))) {
        qDebug() << "set single instance failed!";
        return -1;
    }

#ifndef QT_DEBUG
    QDir::setCurrent(QApplication::applicationDirPath());
#endif

    // 在qApp中记录当前是否为安全模式，如果为安全模式，则无需加载插件，在退出安全模式下，才正常加载插件
    // 此处设置这个属性必须在MainWindow创建之前，因为在mainWindow中会创建加载插件的代理，会在代理中根据这个属性来判断是否需要加载插件
    bool isSafeMode = IsSaveMode();
    bool disablePlugin = parser.isSet(disablePlugOption);
    qApp->setProperty("safeMode", (isSafeMode || disablePlugin));

    MultiScreenWorker multiScreenWorker;

    MainWindow mainWindow(&multiScreenWorker);
    TrayMainWindow trayMainWindow(&multiScreenWorker);

    WindowManager windowManager(&multiScreenWorker);

    // 保证添加窗口的先后顺序，先添加的窗口显示在左边，后添加的窗口显示在右边
    windowManager.addWindow(&mainWindow);
    windowManager.addWindow(&trayMainWindow);

    // 注册任务栏的DBus服务
    DBusDockAdaptors adaptor(&windowManager);
    DockDaemonDBusAdaptor daemonAdaptor(&windowManager);

    QDBusConnection::sessionBus().registerService("org.deepin.dde.Dock1");
    QDBusConnection::sessionBus().registerObject("/org/deepin/dde/Dock1", "org.deepin.dde.Dock1", &windowManager);

    //保证dock daemon接口兼容性,dde-desktop,dde-clipboard,deepin-system-monitor,dde-osd等在调用。
    QDBusConnection::sessionBus().registerService("org.deepin.dde.daemon.Dock1");
    QDBusConnection::sessionBus().registerObject("/org/deepin/dde/daemon/Dock1", "org.deepin.dde.daemon.Dock1", &windowManager);

    // 当任务栏以-r参数启动时，设置CANSHOW未false，之后调用launch不显示任务栏
    qApp->setProperty("CANSHOW", !parser.isSet(runOption));

    windowManager.launch();
    mainWindow.setVisible(true);

    // 判断是否进入安全模式，是否带有入参 -x
    if (!isSafeMode && !disablePlugin) {
        qApp->setProperty("PLUGINSLOADED", true);
    } else {
        windowManager.sendNotifications();
    }

    return app.exec();
}
