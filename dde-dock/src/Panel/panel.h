#ifndef PANEL_H
#define PANEL_H

#include <QWidget>
#include <QLabel>
#include <QPushButton>
#include <QTimer>
#include <QDebug>
#include "DBus/dbushidestatemanager.h"
#include "Controller/dockmodedata.h"
#include "Controller/appmanager.h"
#include "Widgets/appitem.h"
#include "Widgets/docklayout.h"
#include "Widgets/screenmask.h"
#include "panelmenu.h"

class Panel : public QLabel
{
    Q_OBJECT
    Q_PROPERTY(QPoint pos READ pos WRITE move)

public:
    explicit Panel(QWidget *parent = 0);
    ~Panel();

    void showScreenMask();
    void hideScreenMask();

public slots:
    void slotDragStarted();
    void slotItemDropped();
    void slotEnteredMask();
    void slotExitedMask();
signals:
    void startShow();
    void startHide();
    void panelHasShown();
    void panelHasHidden();

private slots:
    void slotLayoutContentsWidthChanged();

    void slotAddAppItem(AppItem *item);
    void slotRemoveAppItem(const QString &id);

protected:
    void mousePressEvent(QMouseEvent *event);
    void mouseReleaseEvent(QMouseEvent *event);

private:
    void changeDockMode(Dock::DockMode newMode,Dock::DockMode oldMode);
    void reanchorsLayout(Dock::DockMode mode);

    void showMenu();

    void initAppManager();

    void hasShown();
    void hasHidden();
    void hideStateChanged(int value);
    void initHSManager();
    void initState();
private:
    DBusHideStateManager * m_HSManager = NULL;
    DockLayout * leftLayout = NULL;
    DockLayout *rightLayout = NULL;
    AppManager *m_appManager = NULL;
    QWidget * parentWidget = NULL;
    ScreenMask * maskWidget = NULL;
    DockModeData *dockCons = DockModeData::instance();
};

#endif // PANEL_H
