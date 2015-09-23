#ifndef LAUNCHERITEM_H
#define LAUNCHERITEM_H

#include <QTimer>
#include <QDebug>
#include <QWidget>
#include <QObject>
#include <QProcess>

#include "appicon.h"
#include "abstractdockitem.h"
#include "controller/dockmodedata.h"
#include "interfaces/dockconstants.h"

class LauncherItem : public AbstractDockItem
{
    Q_OBJECT
public:
    explicit LauncherItem(QWidget *parent = 0);
    ~LauncherItem();

    QString getTitle(){return "Launcher";}
    QWidget * getApplet(){return NULL;}
    bool moveable(){return false;}

protected:
    void enterEvent(QEvent *);
    void leaveEvent(QEvent *);
    void mousePressEvent(QMouseEvent *event);
    void mouseReleaseEvent(QMouseEvent *event);

private slots:
    void slotMousePress(QMouseEvent *event);
    void slotMouseRelease(QMouseEvent *event);
    void updateIcon();

private:
    void changeDockMode(Dock::DockMode newMode, Dock::DockMode oldMode);
    void reanchorIcon();

private:
    DockModeData * m_dockModeData = DockModeData::instance();
    AppIcon * m_appIcon = NULL;
    QProcess * m_launcherProcess = NULL;
    QString m_menuInterfacePath = "";
};

#endif // LAUNCHERITEM_H
