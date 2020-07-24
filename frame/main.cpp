/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#include "window/mainwindow.h"
#include "window/accessible.h"
#include "util/themeappicon.h"
#include "controller/dockitemmanager.h"

#include <QAccessible>
#include <QDir>
#include <QStandardPaths>
#include <QDateTime>

#include <DApplication>
#include <DLog>
#include <DDBusSender>
#include <DGuiApplicationHelper>

#include <unistd.h>
#include "dbus/dbusdockadaptors.h"
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

const int MAX_STACK_FRAMES = 128;
const QString strPath = QStandardPaths::standardLocations(QStandardPaths::ConfigLocation)[0] + "/dde-collapse.log";
const QString cfgPath = QStandardPaths::standardLocations(QStandardPaths::ConfigLocation)[0] + "/dde-cfg.ini";

using namespace std;

// let startdde know that we've already started.
void RegisterDdeSession()
{
    QString envName("DDE_SESSION_PROCESS_COOKIE_ID");

    QByteArray cookie = qgetenv(envName.toUtf8().data());
    qunsetenv(envName.toUtf8().data());

    if (!cookie.isEmpty()) {
        QDBusPendingReply<bool> r = DDBusSender()
                                    .interface("com.deepin.SessionManager")
                                    .path("/com/deepin/SessionManager")
                                    .service("com.deepin.SessionManager")
                                    .method("Register")
                                    .arg(QString(cookie))
                                    .call();

        qDebug() << Q_FUNC_INFO << r.value();
    }
}

bool IsSaveMode()
{
    QSettings settings(cfgPath, QSettings::IniFormat);
    settings.beginGroup("dde-dock");
    int collapseNum = settings.value("collapse").toInt();

    // 自动进入安全模式
    if (collapseNum >= 3) {
        settings.setValue("collapse", 0);
        settings.endGroup();
        settings.sync();
        return true;
    }

    return false;
}
void sig_crash(int sig)
{
    FILE *fd;
    struct stat buf;
    char path[100];
    memset(path, 0, 100);

    // 创建默认配置文件,记录段时间内的崩溃次数
    if (!QFile::exists(cfgPath)) {
        QFile file(cfgPath);
        if (!file.open(QIODevice::WriteOnly))
            exit(0);
        file.close();
    }

    QSettings settings(cfgPath, QSettings::IniFormat);
    settings.beginGroup("dde-dock");

    QDateTime lastDate = QDateTime::fromString(settings.value("lastDate").toString(), "yyyy-MM-dd hh:mm:ss:zzz");
    int collapseNum = settings.value("collapse").toInt();

    //3分钟以内发生崩溃则累加,记录到文件中
    if (qAbs(lastDate.secsTo(QDateTime::currentDateTime())) < 60 * 3) {
        settings.setValue("collapse", collapseNum + 1);
    }
    settings.setValue("lastDate", QDateTime::currentDateTime().toString("yyyy-MM-dd hh:mm:ss:zzz"));
    settings.endGroup();
    settings.sync();

    memcpy(path, strPath.toStdString().data(), strPath.length());

    stat(path, &buf);
    if (buf.st_size > 10 * 1024 * 1024) {
        // 超过10兆则清空内容
        fd = fopen(path, "w");
    } else {
        fd = fopen(path, "at");
    }

    if (nullptr == fd) {
        exit(0);
    }
    //捕获异常，打印崩溃日志到配置文件中
    try {
        char szLine[512] = {0};
        time_t t = time(nullptr);
        tm *now = localtime(&t);
        QString log = "#####" + qApp->applicationName() + "#####\n[%04d-%02d-%02d %02d:%02d:%02d][crash signal number:%d]\n";
        int nLen1 = sprintf(szLine, log.toStdString().c_str(),
                            now->tm_year + 1900,
                            now->tm_mon + 1,
                            now->tm_mday,
                            now->tm_hour,
                            now->tm_min,
                            now->tm_sec,
                            sig);
        fwrite(szLine, 1, strlen(szLine), fd);

#ifdef __linux
        void *array[MAX_STACK_FRAMES];
        size_t size = 0;
        char **strings = nullptr;
        size_t i, j;
        signal(sig, SIG_DFL);
        size = backtrace(array, MAX_STACK_FRAMES);
        strings = (char **)backtrace_symbols(array, size);
        for (i = 0; i < size; ++i) {
            char szLine[512] = {0};
            sprintf(szLine, "%d %s\n", i, strings[i]);
            fwrite(szLine, 1, strlen(szLine), fd);

            std::string symbol(strings[i]);

            size_t pos1 = symbol.find_first_of("[");
            size_t pos2 = symbol.find_last_of("]");
            std::string address = symbol.substr(pos1 + 1, pos2 - pos1 - 1);
            char cmd[128] = {0};
            sprintf(cmd, "addr2line -C -f -e dde-dock %s", address.c_str()); // 打印当前进程的id和地址
            FILE *fPipe = popen(cmd, "r");
            if (fPipe != nullptr) {
                char buff[1024];
                memset(buff, 0, sizeof(buff));
                char *ret = fgets(buff, sizeof(buff), fPipe);
                pclose(fPipe);
                fwrite(ret, 1, strlen(ret), fd);
            }
        }
        free(strings);
#endif // __linux
    } catch (...) {
        //
    }
    fflush(fd);
    fclose(fd);
    fd = nullptr;
    exit(0);
}

int main(int argc, char *argv[])
{
    DGuiApplicationHelper::setUseInactiveColorGroup(false);
    DApplication::loadDXcbPlugin();
    DApplication app(argc, argv);
    //崩溃信号
    signal(SIGTERM, sig_crash);
    signal(SIGSEGV, sig_crash);
    signal(SIGILL, sig_crash);
    signal(SIGINT, sig_crash);
    signal(SIGABRT, sig_crash);
    signal(SIGFPE, sig_crash);

    // 锁定物理内存，用于国测测试[会显著增加内存占用]
    //    qDebug() << "lock memory result:" << mlockall(MCL_CURRENT | MCL_FUTURE);

    app.setOrganizationName("deepin");
    app.setApplicationName("dde-dock");
    app.setApplicationDisplayName("DDE Dock");
    app.setApplicationVersion("2.0");
    app.loadTranslator();
    app.setAttribute(Qt::AA_EnableHighDpiScaling, true);
    app.setAttribute(Qt::AA_UseHighDpiPixmaps, false);

    QAccessible::installFactory(accessibleFactory);

    // load dde-network-utils translator
    QTranslator translator;
    translator.load("/usr/share/dde-network-utils/translations/dde-network-utils_" + QLocale::system().name());
    app.installTranslator(&translator);

    DLogManager::registerConsoleAppender();
    DLogManager::registerFileAppender();

    QCommandLineOption disablePlugOption(QStringList() << "x" << "disable-plugins", "do not load plugins.");
    QCommandLineParser parser;
    parser.setApplicationDescription("DDE Dock");
    parser.addHelpOption();
    parser.addVersionOption();
    parser.addOption(disablePlugOption);
    parser.process(app);

    DGuiApplicationHelper::setSingelInstanceInterval(-1);
    if (!app.setSingleInstance(QString("dde-dock_%1").arg(getuid()))) {
        qDebug() << "set single instance failed!";
        return -1;
    }

    qDebug() << "\n\ndde-dock startup";
    RegisterDdeSession();

#ifndef QT_DEBUG
    QDir::setCurrent(QApplication::applicationDirPath());
#endif

    MainWindow mw;
    DBusDockAdaptors adaptor(&mw);
    QDBusConnection::sessionBus().registerService("com.deepin.dde.Dock");
    QDBusConnection::sessionBus().registerObject("/com/deepin/dde/Dock", "com.deepin.dde.Dock", &mw);

    QTimer::singleShot(1, &mw, &MainWindow::launch);

    if (!IsSaveMode() && !parser.isSet(disablePlugOption)) {
        DockItemManager::instance()->startLoadPlugins();
    }

    return app.exec();
}
