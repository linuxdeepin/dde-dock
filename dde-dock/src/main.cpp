#include <QApplication>
#include <QFile>
#include <QDebug>
#include <QTranslator>
#include <QDBusConnection>

#include "mainwidget.h"
#include "logmanager.h"
#include "Logger.h"
#include "controller/stylemanager.h"
#include "controller/signalmanager.h"

#include <sys/file.h>

#undef signals
extern "C" {
  #include <gtk/gtk.h>
}
#define signals public

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

int main(int argc, char *argv[])
{
    QApplication a(argc, argv);
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

    LogManager::instance()->debug_log_console_on();
    LOG_INFO()<< "LogFile:" << LogManager::instance()->getlogFilePath();

    int pidLock = open("/tmp/dde-dock.pid", O_CREAT | O_RDWR, 0666);
    int rc = flock(pidLock, LOCK_EX | LOCK_NB);
    if (!rc) {

        // TODO: 根据日志发现注册 dbus 来做单例不靠谱，所以改用 pid 锁来做，这里注册 dbus
        // 是为了兼容，等后端把服务监控改为 pid 监控后可以移除这部分代码
        QDBusConnection::sessionBus().registerService(DBUS_NAME);
        RegisterDdeSession();

        StyleManager::instance()->initStyleSheet();

        MainWidget w;
        w.show();
        qWarning() << "Start Dock, The main window has been shown.............................................................";
        w.loadResources();

        initGtkThemeWatcher();

        return a.exec();
    } else {
        qWarning() << "Dock is running!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!";
        return 0;
    }
}
