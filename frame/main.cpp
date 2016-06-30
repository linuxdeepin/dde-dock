
#include "window/mainwindow.h"

#include <dapplication.h>
#include <QDir>

#include <unistd.h>

DWIDGET_USE_NAMESPACE

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
    DApplication app(argc, argv);
    if (!app.setSingleInstance(QString("dde-dock_%1").arg(getuid()))) {
        qDebug() << "set single instance failed!";
        return -1;
    }
    app.setOrganizationName("deepin");
    app.setApplicationName("dde-dock");
    app.setApplicationDisplayName("DDE Dock");
    app.setApplicationVersion("2.0");

#ifndef QT_DEBUG
    QDir::setCurrent(QApplication::applicationDirPath());
#endif

    MainWindow mw;
    QDBusConnection::sessionBus().registerService("com.deepin.dde.dock");
    QDBusConnection::sessionBus().registerObject("/com/deepin/dde/dock", "com.deepin.dde.dock", &mw);
    RegisterDdeSession();

    QTimer::singleShot(300, &mw, &MainWindow::show);

    return app.exec();
}
