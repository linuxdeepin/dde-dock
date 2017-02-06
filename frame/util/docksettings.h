#ifndef DOCKSETTINGS_H
#define DOCKSETTINGS_H

#include "constants.h"
#include "dbus/dbusdock.h"
#include "dbus/dbusmenumanager.h"
#include "dbus/dbusdisplay.h"
#include "controller/dockitemcontroller.h"

#include <QAction>
#include <QMenu>

#include <QObject>
#include <QSize>

#include <QStyleFactory>

extern "C"
{
#ifdef signals
#undef signals
#endif

#include <gtk/gtk.h>
#undef signals
#define signals public
}

DWIDGET_USE_NAMESPACE

using namespace Dock;

class WhiteMenu : public QMenu
{
    Q_OBJECT
public:
    WhiteMenu(QWidget * parent = nullptr) : QMenu(parent) {
        QStyle *style = QStyleFactory::create("dlight");
        if (style) setStyle(style);
    }

    virtual ~WhiteMenu() {}
};

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
    int screenWidth() const;
    int expandTimeout() const;
    int narrowTimeout() const;

    bool autoHide() const;
    const QRect primaryRect() const;
    const QSize windowSize() const;
    const QRect windowRect(const Position position, const bool hide = false) const;

    void showDockSettingsMenu();

signals:
    void dataChanged() const;
    void positionChanged(const Position prevPosition) const;
    void autoHideChanged(const bool autoHide) const;
    void windowVisibleChanged() const;
    void windowHideModeChanged() const;
    void windowGeometryChanged() const;

public slots:
    void updateGeometry();
    void setAutoHide(const bool autoHide);

private slots:
    void menuActionClicked(QAction *action);
    void onPositionChanged();
    void iconSizeChanged();
    void displayModeChanged();
    void hideModeChanged();
    void hideStateChanged();
    void dockItemCountChanged();
    void primaryScreenChanged();
    void resetFrontendGeometry();

private:
    void calculateWindowConfig();
    static void gtkIconThemeChanged(GtkSettings *gs, GParamSpec *pspec, gpointer udata);

private:
    int m_iconSize;
    bool m_autoHide;
    Position m_position;
    HideMode m_hideMode;
    HideState m_hideState;
    DisplayMode m_displayMode;
    QRect m_primaryRect;
    QSize m_mainWindowSize;

    WhiteMenu m_settingsMenu;
    QAction m_fashionModeAct;
    QAction m_efficientModeAct;
    QAction m_topPosAct;
    QAction m_bottomPosAct;
    QAction m_leftPosAct;
    QAction m_rightPosAct;
    QAction m_largeSizeAct;
    QAction m_mediumSizeAct;
    QAction m_smallSizeAct;
    QAction m_keepShownAct;
    QAction m_keepHiddenAct;
    QAction m_smartHideAct;

    DBusDisplay *m_displayInter;
    DBusDock *m_dockInter;
    DockItemController *m_itemController;
};

#endif // DOCKSETTINGS_H
