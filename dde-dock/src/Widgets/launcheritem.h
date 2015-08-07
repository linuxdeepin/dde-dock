#ifndef LAUNCHERITEM_H
#define LAUNCHERITEM_H

#include <QObject>
#include <QWidget>
#include <QTimer>
#include <QProcess>
#include <QDebug>
#include "Controller/dockmodedata.h"
#include "abstractdockitem.h"
#include "appicon.h"
#include "arrowrectangle.h"
#include "../dockconstants.h"

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

private slots:
    void slotMousePress(QMouseEvent *event);
    void slotMouseRelease(QMouseEvent *event);

private:
    void changeDockMode(Dock::DockMode newMode, Dock::DockMode oldMode);
    void updateIcon();
    void reanchorIcon();

private:
    DockModeData * m_dmd = DockModeData::instance();
    AppIcon * m_appIcon = NULL;
    QProcess * m_launcherProcess = NULL;
    QString m_menuInterfacePath = "";
};

#endif // LAUNCHERITEM_H
