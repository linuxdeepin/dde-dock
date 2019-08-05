/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#ifndef DOCKSETTINGS_H
#define DOCKSETTINGS_H

#include "constants.h"
#include "dbus/dbusmenumanager.h"
#include "dbus/dbusdisplay.h"
#include "controller/dockitemcontroller.h"

#include <com_deepin_dde_daemon_dock.h>

#include <QAction>
#include <QMenu>

#include <QObject>
#include <QSize>

#include <QStyleFactory>

DWIDGET_USE_NAMESPACE

using namespace Dock;
using DBusDock = com::deepin::dde::daemon::Dock;

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
    static DockSettings& Instance();

    inline DisplayMode displayMode() const { return m_displayMode; }
    inline HideMode hideMode() const { return m_hideMode; }
    inline HideState hideState() const { return m_hideState; }
    inline Position position() const { return m_position; }
    inline int screenRawHeight() const { return m_screenRawHeight; }
    inline int screenRawWidth() const { return m_screenRawWidth; }
    inline int expandTimeout() const { return m_dockInter->showTimeout(); }
    inline int narrowTimeout() const { return 100; }
    inline bool autoHide() const { return m_autoHide; }
    inline bool isMaxSize() const { return m_isMaxSize; }
    const QRect primaryRect() const;
    inline const QRect primaryRawRect() const { return m_primaryRawRect; }
    inline const QRect frontendWindowRect() const { return m_frontendRect; }
    inline const QSize windowSize() const { return m_mainWindowSize; }
    inline const quint8 Opacity() const { return m_opacity * 255; }
    inline const QSize fashionTraySize() const { return m_fashionTraySize; }

    const QSize panelSize() const;
    const QRect windowRect(const Position position, const bool hide = false) const;
    qreal dockRatio() const;

    void showDockSettingsMenu();

signals:
    void dataChanged() const;
    void positionChanged(const Position prevPosition) const;
    void autoHideChanged(const bool autoHide) const;
    void displayModeChanegd() const;
    void windowVisibleChanged() const;
    void windowHideModeChanged() const;
    void windowGeometryChanged() const;
    void opacityChanged(const quint8 value) const;

public slots:
    void updateGeometry();
    void setAutoHide(const bool autoHide);

private slots:
    void menuActionClicked(QAction *action);
    void onPositionChanged();
    void iconSizeChanged();
    void onDisplayModeChanged();
    void hideModeChanged();
    void hideStateChanged();
    void dockItemCountChanged();
    void primaryScreenChanged();
    void resetFrontendGeometry();
    void updateForbidPostions();
    void onOpacityChanged(const double value);
    void onFashionTraySizeChanged(const QSize &traySize);

private:
    DockSettings(QWidget *parent = 0);
    DockSettings(DockSettings const &) = delete;
    DockSettings operator =(DockSettings const &) = delete;

    bool test(const Position pos, const QList<QRect> &otherScreens) const;
    void calculateWindowConfig();
    void gtkIconThemeChanged();

private:
    int m_iconSize;
    bool m_autoHide;
    bool m_isMaxSize;
    int m_screenRawHeight;
    int m_screenRawWidth;
    double m_opacity;
    QSet<Position> m_forbidPositions;
    Position m_position;
    HideMode m_hideMode;
    HideState m_hideState;
    DisplayMode m_displayMode;
    QRect m_primaryRawRect;
    QRect m_frontendRect;
    QSize m_mainWindowSize;
    QSize m_fashionTraySize;

    WhiteMenu m_settingsMenu;
    WhiteMenu *m_hideSubMenu;
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
