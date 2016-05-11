/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#include <QApplication>
#include <QFile>
#include <QDebug>
#include <QTranslator>
#include <QDBusConnection>

#include <DLog>
#include <dapplication.h>

#include "mainwidget.h"
#include "Logger.h"
#include "controller/stylemanager.h"
#include "controller/signalmanager.h"

#include <unistd.h>

#undef signals
extern "C" {
  #include <gtk/gtk.h>
}
#define signals public

//#define ARCH_MIPSEL

static void requrestUpdateIcons()
{
    //can not passing QObject to the callback function,so use signal
    emit SignalManager::instance()->requestAppIconUpdate();
}

void initGtkThemeWatcher()
{
    GtkSettings* gs = gtk_settings_get_default();
    g_signal_connect(gs, "notify::gtk-icon-theme-name",
                     G_CALLBACK(requrestUpdateIcons), NULL);
}

// let startdde know that we've already started.
void RegisterDdeSession()
{
    QString envName("DDE_SESSION_PROCESS_COOKIE_ID");

    QByteArray cookie = qgetenv(envName.toUtf8().data());
    qunsetenv(envName.toUtf8().data());

    if (!cookie.isEmpty()) {
        QDBusInterface iface("com.deepin.SessionManager",
                             "/com/deepin/SessionManager",
                             "com.deepin.SessionManager",
                             QDBusConnection::sessionBus());
        iface.asyncCall("Register", QString(cookie));
    }
}

DUTIL_USE_NAMESPACE

int main(int argc, char *argv[])
{
    DApplication a(argc, argv);
    if (!a.setSingleInstance(QString("dde-dock_%1").arg(getuid()))) {
        qDebug() << "set single instance failed!";
        return -1;
    }
    a.setOrganizationName("deepin");
    a.setApplicationName("dde-dock");
    a.setApplicationDisplayName("Dock");

    // install translators
    QTranslator translator;
    translator.load("/usr/share/dde-dock/translations/dde-dock_" + QLocale::system().name());
    a.installTranslator(&translator);

	// translations from dde-control-center, used by those plugins provided by dde-control-center,
	// but below lines should be moved to individual plugins in the future.
    QTranslator translator1;
    translator1.load("/usr/share/dde-control-center/translations/dde-control-center_" + QLocale::system().name());
    a.installTranslator(&translator1);

    DLogManager::registerConsoleAppender();
    DLogManager::registerFileAppender();
    dInfo()<< "LogFile:" << DLogManager::getlogFilePath();

    QDBusConnection::sessionBus().registerService(DBUS_NAME);
    RegisterDdeSession();

    StyleManager::instance()->initStyleSheet();

    MainWidget w;
    qWarning() << "Start Dock, The main window has been shown.";
    w.loadResources();
//#ifndef ARCH_MIPSEL
//    w.show();
//#endif

    initGtkThemeWatcher();

    return a.exec();
}
