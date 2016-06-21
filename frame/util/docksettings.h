#ifndef DOCKSETTINGS_H
#define DOCKSETTINGS_H

#include "constants.h"
#include "dbus/dbusdock.h"
#include "dbus/dbusmenumanager.h"
#include "dbus/dbusdisplay.h"
#include "controller/dockitemcontroller.h"

#include <DAction>
#include <DMenu>

#include <QObject>
#include <QSize>

DWIDGET_USE_NAMESPACE

using namespace Dock;

class DockSettings : public QObject
{
    Q_OBJECT

public:
    explicit DockSettings(QObject *parent = 0);

    Position position() const;
    const QSize windowSize() const;

    void showDockSettingsMenu();

signals:
    void dataChanged() const;

public slots:
    void updateGeometry();

private slots:
    void menuActionClicked(DAction *action);
    void positionChanged();

private:
    void calculateWindowConfig();

private:
    int m_iconSize;
    Position m_position;
    HideMode m_hideMode;
    DisplayMode m_displayMode;
    QSize m_mainWindowSize;

    DMenu m_settingsMenu;
    DAction m_fashionModeAct;
    DAction m_efficientModeAct;
    DAction m_topPosAct;
    DAction m_bottomPosAct;
    DAction m_leftPosAct;
    DAction m_rightPosAct;
    DAction m_largeSizeAct;
    DAction m_mediumSizeAct;
    DAction m_smallSizeAct;
    DAction m_keepShownAct;
    DAction m_keepHiddenAct;
    DAction m_smartHideAct;

    DBusDisplay *m_displayInter;
    DBusDock *m_dockInter;
    DockItemController *m_itemController;
};

#endif // DOCKSETTINGS_H
