#ifndef DOCKLAUNCHERITEM_H
#define DOCKLAUNCHERITEM_H

#include <QWidget>

#include "../dockitem.h"
#include "../app/dockappicon.h"
#include "controller/dockmodedata.h"
#include "interfaces/dockconstants.h"

class QProcess;

class DockLauncherItem : public DockItem
{
    Q_OBJECT
public:
    explicit DockLauncherItem(QWidget *parent = 0);
    ~DockLauncherItem();

    QString getItemId() Q_DECL_OVERRIDE {return "dde-launcher";}
    QString getTitle() Q_DECL_OVERRIDE { return tr("Launcher"); }
    QWidget * getApplet() Q_DECL_OVERRIDE { return NULL; }

protected:
    void enterEvent(QEvent *) Q_DECL_OVERRIDE;
    void leaveEvent(QEvent *) Q_DECL_OVERRIDE;
    void mousePressEvent(QMouseEvent *event) Q_DECL_OVERRIDE;
    void mouseReleaseEvent(QMouseEvent *event) Q_DECL_OVERRIDE;

private slots:
    void slotMousePress(QMouseEvent *event);
    void slotMouseRelease(QMouseEvent *event);
    void updateIcon();

private:
    void changeDockMode(Dock::DockMode newMode, Dock::DockMode oldMode);
    void reanchorIcon();

private:
    DockAppIcon * m_appIcon;
    QProcess * m_launcherProcess;
    QString m_menuInterfacePath = "";
    DockModeData * m_dockModeData = DockModeData::instance();
};

#endif // DOCKLAUNCHERITEM_H
