/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef PANEL_H
#define PANEL_H

#include <QWidget>
#include <QLabel>
#include <QPushButton>
#include <QTimer>

#include "dbus/dbushidestatemanager.h"
#include "controller/dockmodedata.h"
#include "controller/old/appmanager.h"
#include "widgets/old/appitem.h"
#include "widgets/old/docklayout.h"
#include "widgets/old/screenmask.h"
#include "widgets/previewwindow.h"
#include "widgets/reflectioneffect.h"
#include "dbus/dbusdisplay.h"
#include "panelmenu.h"

class LayoutDropMask;
class PluginManager;
class Panel : public QLabel
{
    Q_OBJECT
    Q_PROPERTY(int y READ y WRITE setY)
    Q_PROPERTY(bool isFashionMode READ isFashionMode)
    Q_PROPERTY(int width READ width WRITE setFixedWidth)

public:
    explicit Panel(QWidget *parent = 0);
    ~Panel();

    bool isFashionMode();               //for qss setting background
    void showPanelMenu();
    void loadResources();

public slots:
    void resizeWithContent();

signals:
    void startShow();
    void startHide();
    void panelHasShown();
    void panelHasHidden();
    void sizeChanged();

protected:
    void mousePressEvent(QMouseEvent *event);
    void mouseReleaseEvent(QMouseEvent *);

private:
    void initShowHideAnimation();
    void initHideStateManager();
    void initWidthAnimation();
    void initPluginManager();
    void initPluginLayout();
    void initAppLayout();
    void initAppManager();
    void initReflection();
    void initScreenMask();
    void initGlobalPreview();

    void onItemDropped();
    void onItemDragStarted();
    void onAppItemAdd(AbstractDockItem *item, bool delayShow);
    void onAppItemAppend(AbstractDockItem *item, bool delayShow);
    void onAppItemRemove(const QString &id);
    void onDockModeChanged(Dock::DockMode newMode, Dock::DockMode);
    void onHideStateChanged(int dockState);
    void onShowPanelFinished();
    void onHidePanelFinished();
    void onNeedPreviewHide();
    void onNeedPreviewShow(QPoint pos);
    void onNeedPreviewImmediatelyHide();
    void onNeedPreviewUpdate();

    void changeItemHoverable(QAbstractAnimation::State state);
    void reanchorsLayout(Dock::DockMode mode);
    void updateRightReflection();
    void updateLeftReflection();
    void showPluginLayoutMask();
    void hidePluginLayoutMask();
    void reloadStyleSheet();
    void setY(int value);   //for hide and show animation

    DisplayRect getScreenRect();

private:
    QPoint m_lastPreviewPos;
    PreviewWindow *m_globalPreview = NULL;
    DBusDockedAppManager *m_ddam = new DBusDockedAppManager(this);
    DockModeData *m_dockModeData = DockModeData::instance();
    QPropertyAnimation *m_widthAnimation = NULL;
    DBusHideStateManager *m_HSManager = NULL;
    ReflectionEffect *m_pluginReflection = NULL;
    ReflectionEffect *m_appReflection = NULL;
    DockLayout *m_pluginLayout = NULL;
    ScreenMask * m_maskWidget = NULL;
    AppManager *m_appManager = NULL;
    QWidget *m_parentWidget = NULL;
    LayoutDropMask *m_pluginLayoutMask = NULL;
    DockLayout *m_appLayout = NULL;
    PluginManager *m_pluginManager = NULL;

    bool m_previewShown = false;
    bool m_isFashionMode = false;
    const int REFLECTION_HEIGHT = 15;
    const int FASHION_PANEL_LPADDING = 21;
    const int FASHION_PANEL_RPADDING = 21;
    const int WIDTH_ANIMATION_DURATION = 200;
    const int SHOW_ANIMATION_DURATION = 300;
    const int HIDE_ANIMATION_DURATION = 300;
    const int DELAY_HIDE_PREVIEW_INTERVAL = 200;
    const int DELAY_SHOW_PREVIEW_INTERVAL = 200;
    const QEasingCurve SHOW_EASINGCURVE = QEasingCurve::OutCubic;
    const QEasingCurve HIDE_EASINGCURVE = QEasingCurve::Linear;
};

class LayoutDropMask : public  QFrame
{
    Q_OBJECT
public:
    LayoutDropMask(QWidget *parent = 0);

signals:
    void itemMove();
    void itemEnter();
    void itemDrop();

protected:
    void dragEnterEvent(QDragEnterEvent *event);
    void dragMoveEvent(QDragMoveEvent *event);
    void dropEvent(QDropEvent *event);

};

#endif // PANEL_H
