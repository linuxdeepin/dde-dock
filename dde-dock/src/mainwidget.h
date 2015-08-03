#ifndef MAINWIDGET_H
#define MAINWIDGET_H

#include <QApplication>
#include <QDesktopWidget>
#include <QWidget>
#include <QScreen>
#include <QStateMachine>
#include <QState>
#include <QPropertyAnimation>
#include <QDBusConnection>
#include "DBus/dbushidestatemanager.h"
#include "DBus/dbusdocksetting.h"
#include "Controller/dockmodedata.h"
#include "Panel/panel.h"

const QString DBUS_PATH = "/com/deepin/dde/dock";
const QString DBUS_NAME = "com.deepin.dde.dock";

class DockUIDbus;
class MainWidget : public QWidget
{
    Q_OBJECT

public:
    MainWidget(QWidget *parent = 0);
    ~MainWidget();

protected:
    void enterEvent(QEvent *event);
    void leaveEvent(QEvent *);

private:
    void showDock();
    void hideDock();
    void changeDockMode(Dock::DockMode newMode,Dock::DockMode oldMode);
    void updateXcbStructPartial();
    void initHideStateManager();
    void initDockSetting();
private:
    Panel *mainPanel = NULL;
    bool hasHidden = false;
    DockModeData * m_dmd = DockModeData::instance();
    DBusHideStateManager *m_dhsm = NULL;
    DBusDockSetting *m_dds = NULL;
};

class DockUIDbus : public QDBusAbstractAdaptor {
    Q_OBJECT
    Q_CLASSINFO("D-Bus Interface", "com.deepin.dde.dock")

public:
    DockUIDbus(MainWidget* parent);
    ~DockUIDbus();

    Q_SLOT qulonglong Xid();
private:
    MainWidget* m_parent;
};

#endif // MAINWIDGET_H
