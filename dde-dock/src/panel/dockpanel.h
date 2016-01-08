#ifndef DOCKPANEL_H
#define DOCKPANEL_H

#include <QLabel>

#include "dbus/dbushidestatemanager.h"
#include "controller/dockmodedata.h"
#include "widgets/old/appitem.h"
#include "widgets/plugin/dockpluginlayout.h"
#include "widgets/app/dockapplayout.h"
#include "widgets/screenmask.h"
#include "widgets/previewwindow.h"
#include "widgets/reflectioneffect.h"
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
    void initWidthAnimation();
    void initPluginLayout();
    void initAppLayout();
    void initReflection();
    void initScreenMask();
    void initGlobalPreview();

    void onItemDragStarted();
    void onDockModeChanged(Dock::DockMode newMode, Dock::DockMode);
    void onHideStateChanged(int dockState);
    void onShowPanelFinished();
    void onHidePanelFinished();
    void onNeedPreviewHide(bool immediately);
    void onNeedPreviewShow(QPoint pos);
    void onNeedPreviewUpdate();

    void reanchorsLayout(Dock::DockMode mode);
    void updateRightReflection();
    void updateLeftReflection();
    void resizeWithContent();
    void reloadStyleSheet();
    void setY(int value);   //for hide and show animation
    void showPanelMenu();

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
    DockPluginLayout *m_pluginLayout = NULL;
    ScreenMask * m_maskWidget = NULL;
    QWidget *m_parentWidget = NULL;
    LayoutDropMask *m_pluginLayoutMask = NULL;
    DockAppLayout *m_appLayout = NULL;

    bool m_previewShown = false;
    bool m_menuItemInvoked = false;
    bool m_isFashionMode = false;
};

#endif // DOCKPANEL_H
