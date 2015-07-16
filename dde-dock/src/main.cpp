#include <QApplication>
#include <QFile>
#include <QDebug>
#include "mainwidget.h"

#include "xcb_misc.h"

int main(int argc, char *argv[])
{
    QApplication a(argc, argv);

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

    XcbMisc::instance()->set_window_type(w.winId(),
                                         XcbMisc::Dock);

    XcbMisc::instance()->set_strut_partial(w.winId(),
                                           XcbMisc::OrientationBottom,
                                           w.height(),
                                           w.x(),
                                           w.x() + w.width());

    return a.exec();
}
