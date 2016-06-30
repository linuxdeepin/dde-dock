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
    explicit DockSettings(QWidget *parent = 0);

    DisplayMode displayMode() const;
    HideMode hideMode() const;
    HideState hideState() const;
    Position position() const;
    int screenHeight() const;
    const QRect primaryRect() const;
    const QSize windowSize() const;

    void showDockSettingsMenu();

signals:
    void dataChanged() const;
    void windowVisibleChanegd() const;
    void windowHideModeChanged() const;
    void windowGeometryChanged() const;

public slots:
    void updateGeometry();

private slots:
    void menuActionClicked(DAction *action);
    void positionChanged();
    void iconSizeChanged();
    void displayModeChanged();
    void hideModeChanged();
    void hideStateChanegd();
    void dockItemCountChanged();
    void primaryScreenChanged();

    void resetFrontendWinId();

private:
    void calculateWindowConfig();

private:
    int m_iconSize;
    Position m_position;
    HideMode m_hideMode;
    HideState m_hideState;
    DisplayMode m_displayMode;
    QRect m_primaryRect;
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
