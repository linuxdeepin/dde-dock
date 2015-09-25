#include <QApplication>
#include <QFile>
#include <QDebug>
#include <QDBusConnection>

#include "mainwidget.h"
#include "logmanager.h"
#include "Logger.h"

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

int main(int argc, char *argv[])
{
    QApplication a(argc, argv);
    a.setOrganizationName("deepin");
    a.setApplicationName("dde-dock");
    a.setApplicationDisplayName("Dock");

    LogManager::instance()->debug_log_console_on();
    LOG_INFO() << LogManager::instance()->getlogFilePath();

    if (QDBusConnection::sessionBus().registerService(DBUS_NAME)) {
        QFile file("://qss/resources/dark/qss/dde-dock.qss");
        if (file.open(QFile::ReadOnly)) {
            QString styleSheet = QLatin1String(file.readAll());
            qApp->setStyleSheet(styleSheet);
            file.close();
        } else {
            qWarning() << "[Error:] Open  style file errr!";
        }


        MainWidget w;
        w.show();

        RegisterDdeSession();

        return a.exec();
    } else {
        qWarning() << "dde dock is running...";
        return 0;
    }
}
