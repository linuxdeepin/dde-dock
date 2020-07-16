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

#include <QAccessible>
#include <QDir>
#include <QStandardPaths>

#include <DApplication>
#include <DLog>
#include <DDBusSender>

#include <QDir>
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
const int MAX_STACK_FRAMES = 128;

using namespace std;

void sig_crash(int sig)
{
    FILE *fd;
    struct stat buf;
    char path[100];
    memset(path, 0, 100);
    //崩溃日志路径
    QString strPath = QStandardPaths::standardLocations(QStandardPaths::ConfigLocation)[0] + "/dde-collapse.log";
    memcpy(path, strPath.toStdString().data(), strPath.length());
    qDebug() << path;

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

QAccessibleInterface *accessibleFactory(const QString &classname, QObject *object)
{
    QAccessibleInterface *interface = nullptr;

    GET_ACCESSIBLE(classname, MainWindow);
    GET_ACCESSIBLE(classname, MainPanelControl);
    GET_ACCESSIBLE(classname, TipsWidget);
    GET_ACCESSIBLE(classname, DockPopupWindow);
    GET_ACCESSIBLE(classname, LauncherItem);
    GET_ACCESSIBLE(classname, AppItem);
    GET_ACCESSIBLE(classname, PreviewContainer);
    GET_ACCESSIBLE(classname, PluginsItem);
    GET_ACCESSIBLE(classname, TrayPluginItem);
    GET_ACCESSIBLE(classname, PlaceholderItem);
    GET_ACCESSIBLE(classname, AppDragWidget);
    GET_ACCESSIBLE(classname, AppSnapshot);
    GET_ACCESSIBLE(classname, FloatingPreview);
    GET_ACCESSIBLE(classname, SNITrayWidget);
    GET_ACCESSIBLE(classname, SystemTrayItem);
    GET_ACCESSIBLE(classname, FashionTrayItem);
    GET_ACCESSIBLE(classname, FashionTrayWidgetWrapper);
    GET_ACCESSIBLE(classname, FashionTrayControlWidget);
    GET_ACCESSIBLE(classname, AttentionContainer);
    GET_ACCESSIBLE(classname, HoldContainer);
    GET_ACCESSIBLE(classname, NormalContainer);
    GET_ACCESSIBLE(classname, SpliterAnimated);
    GET_ACCESSIBLE(classname, IndicatorTrayWidget);
    GET_ACCESSIBLE(classname, XEmbedTrayWidget);
    GET_ACCESSIBLE(classname, SoundItem);
    GET_ACCESSIBLE(classname, SoundApplet);
    GET_ACCESSIBLE(classname, SinkInputWidget);
    GET_ACCESSIBLE(classname, VolumeSlider);
    GET_ACCESSIBLE(classname, HorizontalSeparator);
    GET_ACCESSIBLE(classname, DatetimeWidget);
    GET_ACCESSIBLE(classname, OnboardItem);
    GET_ACCESSIBLE(classname, TrashWidget);
    GET_ACCESSIBLE(classname, PopupControlWidget);
    GET_ACCESSIBLE(classname, ShutdownWidget);
    GET_ACCESSIBLE(classname, MultitaskingWidget);
    GET_ACCESSIBLE(classname, ShowDesktopWidget);
    //    USE_ACCESSIBLE(classname,OverlayWarningWidget);
    GET_ACCESSIBLE_BY_OBJECTNAME(classname, QWidget, "Btn_showdesktoparea");//TODO 点击坐标有偏差
    GET_ACCESSIBLE_BY_OBJECTNAME(QString(classname).replace("Dtk::Widget::", ""), DImageButton, "closebutton-2d");
    GET_ACCESSIBLE_BY_OBJECTNAME(QString(classname).replace("Dtk::Widget::", ""), DImageButton, "closebutton-3d");
    GET_ACCESSIBLE_BY_OBJECTNAME(QString(classname).replace("Dtk::Widget::", ""), DSwitchButton, "");

    return interface;
}

int main(int argc, char *argv[])
{
    qputenv("QT_WAYLAND_SHELL_INTEGRATION", "kwayland-shell");
    DGuiApplicationHelper::setUseInactiveColorGroup(false);
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

    // load dde-network-utils translator
    QTranslator translator;
    translator.load("/usr/share/dde-network-utils/translations/dde-network-utils_" + QLocale::system().name());
    app.installTranslator(&translator);

    DLogManager::registerConsoleAppender();
    DLogManager::registerFileAppender();

    QAccessible::installFactory(accessibleFactory);

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
    mw.setAttribute(Qt::WA_NativeWindow);
    mw.windowHandle()->setProperty("_d_dwayland_window-type" , "dock");
    QDBusConnection::sessionBus().registerService("com.deepin.dde.Dock");
    QDBusConnection::sessionBus().registerObject("/com/deepin/dde/Dock", "com.deepin.dde.Dock", &mw);

    QTimer::singleShot(1, &mw, &MainWindow::launch);

    if (!parser.isSet(disablePlugOption)) {
        DockItemManager::instance()->startLoadPlugins();
    }

    return app.exec();
}
