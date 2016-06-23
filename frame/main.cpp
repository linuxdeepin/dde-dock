
#include "window/mainwindow.h"

#include <dapplication.h>

#include <unistd.h>

DWIDGET_USE_NAMESPACE

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

    MainWindow mw;
    QDBusConnection::sessionBus().registerService("com.deepin.dde.dock");
    QDBusConnection::sessionBus().registerObject("/com/deepin/dde/dock", "com.deepin.dde.dock", &mw);

    QTimer::singleShot(500, &mw, &MainWindow::show);

    return app.exec();
}
