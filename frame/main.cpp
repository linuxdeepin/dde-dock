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

#include <DApplication>
#include <DLog>
#include <DDBusSender>

#include <QDir>
#include <DGuiApplicationHelper>

#include <unistd.h>
#include "dbus/dbusdockadaptors.h"

#include <sys/mman.h>

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

QAccessibleInterface *accessibleFactory(const QString &classname, QObject *object)
{
    QAccessibleInterface *interface = nullptr;

    USE_ACCESSIBLE(classname,MainPanelControl);
    USE_ACCESSIBLE(classname,LauncherItem);
    USE_ACCESSIBLE(classname,AppItem);
    USE_ACCESSIBLE(classname,PreviewContainer);
    USE_ACCESSIBLE(classname,PluginsItem);
    USE_ACCESSIBLE(classname,TrayPluginItem);
    USE_ACCESSIBLE(classname,PlaceholderItem);
    USE_ACCESSIBLE(classname,AppDragWidget);
    USE_ACCESSIBLE(classname,AppSnapshot);
    USE_ACCESSIBLE(classname,FloatingPreview);
    USE_ACCESSIBLE(classname,SNITrayWidget);
    USE_ACCESSIBLE(classname,SystemTrayItem);
    USE_ACCESSIBLE(classname,FashionTrayItem);
    USE_ACCESSIBLE(classname,FashionTrayWidgetWrapper);
    USE_ACCESSIBLE(classname,FashionTrayControlWidget);
    USE_ACCESSIBLE(classname,AttentionContainer);
    USE_ACCESSIBLE(classname,HoldContainer);
    USE_ACCESSIBLE(classname,NormalContainer);
    USE_ACCESSIBLE(classname,SpliterAnimated);
    USE_ACCESSIBLE(classname,IndicatorTrayWidget);
    USE_ACCESSIBLE(classname,XEmbedTrayWidget);
    USE_ACCESSIBLE(classname,ShowDesktopWidget);
    USE_ACCESSIBLE(classname,SoundItem);
    USE_ACCESSIBLE(classname,SoundApplet);
    USE_ACCESSIBLE(classname,SinkInputWidget);
    USE_ACCESSIBLE(classname,VolumeSlider);
    USE_ACCESSIBLE(classname,HorizontalSeparator);
    USE_ACCESSIBLE(classname,TipsWidget);
    USE_ACCESSIBLE(classname,DatetimeWidget);
    USE_ACCESSIBLE(classname,OnboardItem);
    USE_ACCESSIBLE(classname,TrashWidget);
    USE_ACCESSIBLE(classname,PopupControlWidget);
    USE_ACCESSIBLE(classname,ShutdownWidget);
    USE_ACCESSIBLE(classname,MultitaskingWidget);
//    USE_ACCESSIBLE(classname,OverlayWarningWidget);
    return interface;
}

int main(int argc, char *argv[])
{
    DGuiApplicationHelper::setUseInactiveColorGroup(false);
    DApplication::loadDXcbPlugin();
    DApplication app(argc, argv);

    // 锁定物理内存，用于国测测试
    qDebug() << "lock memory result:" << mlockall(MCL_CURRENT | MCL_FUTURE);

    app.setOrganizationName("deepin");
    app.setApplicationName("dde-dock");
    app.setApplicationDisplayName("DDE Dock");
    app.setApplicationVersion("2.0");
    app.loadTranslator();
    app.setAttribute(Qt::AA_EnableHighDpiScaling, true);
    app.setAttribute(Qt::AA_UseHighDpiPixmaps, false);

    QAccessible::installFactory(accessibleFactory);

    DLogManager::registerConsoleAppender();
    DLogManager::registerFileAppender();

    QCommandLineOption disablePlugOption(QStringList() << "x" << "disable-plugins", "do not load plugins.");
    QCommandLineParser parser;
    parser.setApplicationDescription("DDE Dock");
    parser.addHelpOption();
    parser.addVersionOption();
    parser.addOption(disablePlugOption);
    parser.process(app);

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

    if (!parser.isSet(disablePlugOption)) {
        DockItemManager::instance()->startLoadPlugins();
    }

    return app.exec();
}
