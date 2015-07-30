#include <QApplication>
#include <QFile>
#include <QDebug>
#include "mainwidget.h"

int main(int argc, char *argv[])
{
    QApplication a(argc, argv);

    if(QDBusConnection::sessionBus().registerService(DBUS_NAME)){
        QFile file("://Resources/qss/default.qss");
        if (file.open(QFile::ReadOnly))
        {
            QString styleSheet = QLatin1String(file.readAll());
            qApp->setStyleSheet(styleSheet);
            file.close();
        }
        else
        {
            qWarning() << "[Error:] Open  style file errr!";
        }

        MainWidget w;
        w.show();

        return a.exec();
    }
    else {
        qWarning() << "dde dock is running...";
        return 0;
    }
}
