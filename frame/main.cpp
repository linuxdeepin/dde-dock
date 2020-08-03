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
    QFile *file = new QFile(strPath);

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

    // 3分钟以内发生崩溃则累加,记录到文件中
    if (qAbs(lastDate.secsTo(QDateTime::currentDateTime())) < 60 * 3) {
        settings.setValue("collapse", collapseNum + 1);
    }
    settings.setValue("lastDate", QDateTime::currentDateTime().toString("yyyy-MM-dd hh:mm:ss:zzz"));
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

#ifdef Q_OS_LINUX
        void *array[MAX_STACK_FRAMES];
        size_t size = 0;
        char **strings = nullptr;
        size_t i;
        signal(sig, SIG_DFL);
        size = static_cast<size_t>(backtrace(array, MAX_STACK_FRAMES));
        strings = backtrace_symbols(array, int(size));
        for (i = 0; i < size; ++i) {
            QString line = QString::number(i) + " " + QString::fromStdString(strings[i]) + "\n";
            file->write(line.toUtf8());

            std::string symbol(strings[i]);
            QString strSymbol = QString::fromStdString(symbol);
            int pos1 = strSymbol.indexOf("[");
            int pos2 = strSymbol.lastIndexOf("]");
            QString address = strSymbol.mid(pos1 + 1,pos2 - pos1 - 1);

            // 按照内存地址找到对应代码的行号
            QString cmd = "addr2line -C -f -e " + qApp->applicationName() + " " + address;
            QProcess *p = new QProcess;
            p->setReadChannel(QProcess::StandardOutput);
            p->start(cmd);
            p->waitForFinished();
            p->waitForReadyRead();
            file->write(p->readAllStandardOutput());
            delete p;
            p = nullptr;
        }
        free(strings);
#endif // __linux
    } catch (...) {
        //
    }
    file->close();
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
    signal(SIGILL,  sig_crash);
    signal(SIGINT,  sig_crash);
    signal(SIGABRT, sig_crash);
    signal(SIGFPE,  sig_crash);

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

    DGuiApplicationHelper::setSingleInstanceInterval(-1);
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

    mw.launch();

    if (!IsSaveMode() && !parser.isSet(disablePlugOption)) {
        DockItemManager::instance()->startLoadPlugins();
    }

    return app.exec();
}
