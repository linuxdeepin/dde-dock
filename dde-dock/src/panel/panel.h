#ifndef PANEL_H
#define PANEL_H

#include <QWidget>
#include <QLabel>
#include <QPushButton>
#include <QTimer>

#include "dbus/dbushidestatemanager.h"
#include "controller/dockmodedata.h"
#include "controller/appmanager.h"
#include "widgets/appitem.h"
#include "widgets/docklayout.h"
#include "widgets/screenmask.h"
#include "widgets/previewframe.h"
#include "widgets/reflectioneffect.h"
#include "panelmenu.h"

class LayoutDropMask;
class DockPluginManager;
class Panel : public QLabel
{
    Q_OBJECT
    Q_PROPERTY(int y READ y WRITE setY)
    Q_PROPERTY(bool isFashionMode READ isFashionMode)
    Q_PROPERTY(int width READ width WRITE setFixedWidth)

public:
    explicit Panel(QWidget *parent = 0);
    ~Panel();

    void setContainMouse(bool value);   //for smart-hide and keep-hide
    bool isFashionMode();               //for qss setting background
    void showPanelMenu();
    void loadResources();

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
    void onLayoutContentsWidthChanged();
    void onAppItemAdd(AbstractDockItem *item);
    void onAppItemRemove(const QString &id);
    void onDockModeChanged(Dock::DockMode newMode, Dock::DockMode);
    void onHideStateChanged(int dockState);
    void onShowPanelFinished();
    void onHidePanelFinished();
    void onNeedPreviewHide();
    void onNeedPreviewShow(QPoint pos);
    void onNeedPreviewImmediatelyHide();

    void reanchorsLayout(Dock::DockMode mode);
    void updateRightReflection();
    void updateLeftReflection();
    void showPluginLayoutMask();
    void hidePluginLayoutMask();
    void reloadStyleSheet();
    void setY(int value);   //for hide and show animation

private:
    PreviewFrame *m_globalPreview = NULL;
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
    DockPluginManager *m_pluginManager = NULL;

    bool m_containMouse = false;
    bool m_isFashionMode = false;
    const int REFLECTION_HEIGHT = 15;
    const int FASHION_PANEL_LPADDING = 21;
    const int FASHION_PANEL_RPADDING = 21;
    const int WIDTH_ANIMATION_DURATION = 200;
    const int SHOW_HIDE_ANIMATION_DURATION = 200;
    const int DELAY_HIDE_PREVIEW_INTERVAL = 200;
    const int DELAY_SHOW_PREVIEW_INTERVAL = 200;
    const QEasingCurve SHOW_HIDE_EASINGCURVE = QEasingCurve::InSine;
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
