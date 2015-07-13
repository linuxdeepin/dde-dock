#ifndef PANEL_H
#define PANEL_H

#include <QWidget>
#include <QLabel>
#include <QPushButton>
#include <QDebug>
#include "Controller/dockmodedata.h"
#include "Widgets/appitem.h"
#include "Widgets/docklayout.h"
#include "Widgets/screenmask.h"
#include "panelmenu.h"
#include "Controller/appmanager.h"

class Panel : public QLabel
{
    Q_OBJECT
public:
    explicit Panel(QWidget *parent = 0);
    ~Panel();

    void resize(const QSize &size);
    void resize(int width,int height);

    void showScreenMask();
    void hideScreenMask();

public slots:
    void slotDragStarted();
    void slotItemDropped();
    void slotEnteredMask();
    void slotExitedMask();

private slots:
    void slotDockModeChanged(Dock::DockMode newMode,Dock::DockMode oldMode);
    void slotLayoutContentsWidthChanged();

    void slotAddAppItem(AppItem *item);
    void slotRemoveAppItem(const QString &id);

protected:
    void mousePressEvent(QMouseEvent *event);
    void mouseReleaseEvent(QMouseEvent *event);

private:
    void reanchorsLayout(Dock::DockMode mode);

    void showMenu();

    void initAppManager();
private:
    DockLayout * leftLayout;
    DockLayout *rightLayout;
    AppManager *m_appManager = NULL;
    QWidget * parentWidget = NULL;
    ScreenMask * maskWidget = NULL;
    DockModeData *dockCons = DockModeData::instance();
};

#endif // PANEL_H
