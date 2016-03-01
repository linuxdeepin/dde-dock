/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

#ifndef DOCKPANEL_H
#define DOCKPANEL_H

#include <QLabel>

#include "dbus/dbushidestatemanager.h"
#include "controller/dockmodedata.h"
#include "widgets/launcher/docklauncheritem.h"
#include "widgets/plugin/dockpluginlayout.h"
#include "widgets/app/dockapplayout.h"
#include "widgets/reflectioneffect.h"
#include "widgets/previewwindow.h"
#include "dbus/dbusdisplay.h"
#include "panelmenu.h"

class LayoutDropMask;
class DockPanel : public QLabel
{
    Q_OBJECT
    Q_PROPERTY(int y READ y WRITE setY)
    Q_PROPERTY(bool isFashionMode READ isFashionMode)
    Q_PROPERTY(int width READ width WRITE setFixedWidth)

public:
    explicit DockPanel(QWidget *parent = 0);
    ~DockPanel();

    bool isFashionMode();               //for qss setting background
    void loadResources();
    QSize sizeHint() const;

signals:
    void startShow();
    void startHide();
    void panelHasShown();
    void panelHasHidden();
    void sizeChanged();

protected:
    void mousePressEvent(QMouseEvent *event);

private:
    void initShowHideAnimation();
    void initHideStateManager();
    void initPluginLayout();
    void initAppLayout();
    void initMainLayout();
    void initReflection();
    void initGlobalPreview();

    void onDockModeChanged(Dock::DockMode, Dock::DockMode);
    void onHideStateChanged(int dockState);
    void onShowPanelFinished();
    void onHidePanelFinished();
    void onNeedPreviewHide(bool immediately);
    void onNeedPreviewShow(DockItem *item, const QPoint &pos);
    void onNeedPreviewUpdate();
    void onContentsSizeChanged();

    void changeItemHoverable(QAbstractAnimation::State state);
    void reloadStyleSheet();
    void setY(int value);   //for hide and show animation
    void showPanelMenu();
    void updateReflection();

    DisplayRect getScreenRect();

private:
    QPoint m_lastPreviewPos;
    PreviewWindow *m_globalPreview = NULL;
    DBusHideStateManager *m_HSManager = NULL;
    DockModeData *m_dockModeData = DockModeData::instance();
    DockLauncherItem *m_launcherItem = NULL;
    DockPluginLayout *m_pluginLayout = NULL;
    DockAppLayout *m_appLayout = NULL;
    ReflectionEffect *m_launcherReflection;
    ReflectionEffect *m_pluginReflection;
    ReflectionEffect *m_appReflection;

    bool m_isFashionMode = false;
};

#endif // DOCKPANEL_H
