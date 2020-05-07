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
        QObject::connect(&mw, &MainWindow::loaderPlugins, DockItemManager::instance(), &DockItemManager::startLoadPlugins);
    }

    return app.exec();
}
