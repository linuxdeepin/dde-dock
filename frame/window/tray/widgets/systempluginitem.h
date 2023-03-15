// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SYSTEMTRAYITEM_H
#define SYSTEMTRAYITEM_H

#include "constants.h"
#include "basetraywidget.h"
#include "pluginsiteminterface.h"

#include <QGestureEvent>

class QGSettings;
class QMenu;
class DockPopupWindow;
class QHBoxLayout;

class SystemPluginItem : public BaseTrayWidget
{
    Q_OBJECT

public:
    SystemPluginItem(PluginsItemInterface *const pluginInter, const QString &itemKey, QWidget *parent = nullptr);
    ~SystemPluginItem() override;

public:
    QString itemKeyForConfig() override;
    void updateIcon() override;
    void sendClick(uint8_t mouseButton, int x, int y) override;
    inline TrayType trayType() const override { return TrayType::SystemTray; }
    QPixmap icon() override;

    QWidget *trayTipsWidget();
    QWidget *trayPopupApplet();
    const QString trayClickCommand();
    const QString contextMenu() const;
    void invokedMenuItem(const QString &menuId, const bool checked);

    static void setDockPostion(const Dock::Position pos) { DockPosition = pos; }

    QWidget *centralWidget() const;
    void detachPluginWidget();
    void showContextMenu();

    void showPopupApplet(QWidget * const applet);
    void hidePopup();
    bool containsPoint(const QPoint& mouse) override;

signals:
    void itemVisibleChanged(bool visible);
    void execActionFinished();

protected:
    bool event(QEvent *event) override;
    bool eventFilter(QObject *watched, QEvent *event) override;
    void enterEvent(QEvent *event) override;
    void leaveEvent(QEvent *event) override;
    void mousePressEvent(QMouseEvent *event) override;
    void mouseReleaseEvent(QMouseEvent *event) override;
    void showEvent(QShowEvent *event) override;
    void hideEvent(QHideEvent *event) override;
    void paintEvent(QPaintEvent *event) override;

protected:
    const QPoint popupMarkPoint() const;
    const QPoint topleftPoint() const;

    void hideNonModel();
    void popupWindowAccept();

    virtual void showPopupWindow(QWidget * const content, const bool model = false);
    virtual void showHoverTips();

    bool checkAndResetTapHoldGestureState();
    virtual void gestureEvent(QGestureEvent *event);

protected Q_SLOTS:
    void onContextMenuAccepted();

private:
    void updatePopupPosition();
    void onGSettingsChanged(const QString &key);
    bool checkGSettingsControl() const;
    void menuActionClicked(QAction *action);
    void showCentralWidget();

private:
    bool m_popupShown;
    bool m_tapAndHold;
    QMenu *m_contextMenu;

    PluginsItemInterface* m_pluginInter;
    QWidget *m_centralWidget;

    QTimer *m_popupTipsDelayTimer;
    QTimer *m_popupAdjustDelayTimer;

    QPointer<QWidget> m_lastPopupWidget;
    QString m_itemKey;

    static Dock::Position DockPosition;
    static QPointer<DockPopupWindow> PopupWindow;
    const QGSettings* m_gsettings;
};

#endif // SYSTEMTRAYITEM_H
