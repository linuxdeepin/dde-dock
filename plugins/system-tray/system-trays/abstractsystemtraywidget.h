/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     listenerri <listenerri@gmail.com>
 *
 * Maintainer: listenerri <listenerri@gmail.com>
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

#ifndef ABSTRACTSYSTEMTRAYWIDGET_H
#define ABSTRACTSYSTEMTRAYWIDGET_H

#include "constants.h"
#include "abstracttraywidget.h"
#include "util/dockpopupwindow.h"
#include "dbus/dbusmenumanager.h"

class AbstractSystemTrayWidget : public AbstractTrayWidget
{
    Q_OBJECT

public:
    AbstractSystemTrayWidget(QWidget *parent = nullptr);
    virtual ~AbstractSystemTrayWidget();

    void sendClick(uint8_t mouseButton, int x, int y) Q_DECL_OVERRIDE;
    inline TrayType trayTyep() const Q_DECL_OVERRIDE { return TrayType::SystemTray; }

    virtual inline QWidget *trayTipsWidget() { return nullptr; }
    virtual inline QWidget *trayPopupApplet() { return nullptr; }
    virtual inline const QString trayClickCommand() { return QString(); }
    virtual inline const QString contextMenu() const {return QString(); }
    virtual inline void invokedMenuItem(const QString &menuId, const bool checked) { Q_UNUSED(menuId); Q_UNUSED(checked); }

    static void setDockPostion(const Dock::Position pos) { DockPosition = pos; }

Q_SIGNALS:
    void requestWindowAutoHide(const bool autoHide) const;
    void requestRefershWindowVisible() const;

protected:
    bool event(QEvent *event) Q_DECL_OVERRIDE;
    void enterEvent(QEvent *event) Q_DECL_OVERRIDE;
    void leaveEvent(QEvent *event) Q_DECL_OVERRIDE;
    void mousePressEvent(QMouseEvent *event) Q_DECL_OVERRIDE;

protected:
    const QPoint popupMarkPoint() const;
    const QPoint topleftPoint() const;

    void hidePopup();
    void hideNonModel();
    void popupWindowAccept();
    void showPopupApplet(QWidget * const applet);

    virtual void showPopupWindow(QWidget * const content, const bool model = false);
    virtual void showHoverTips();

protected Q_SLOTS:
    void showContextMenu();
    void onContextMenuAccepted();

private:
    void updatePopupPosition();

private:
    bool m_popupShown;

    QPointer<QWidget> m_lastPopupWidget;

    QTimer *m_popupTipsDelayTimer;
    QTimer *m_popupAdjustDelayTimer;

    DBusMenuManager *m_menuManagerInter;

    static Dock::Position DockPosition;
    static QPointer<DockPopupWindow> PopupWindow;
};

#endif // ABSTRACTSYSTEMTRAYWIDGET_H
