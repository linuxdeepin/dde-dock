#include "accessibledefine.h"

#include "mainwindow.h"
#include "../panel/mainpanelcontrol.h"
#include "../../widgets/tipswidget.h"
#include "../util/dockpopupwindow.h"

#include "../item/launcheritem.h"
#include "../item/appitem.h"
#include "../item/components/previewcontainer.h"
#include "../item/pluginsitem.h"
#include "../item/traypluginitem.h"
#include "../item/placeholderitem.h"
#include "../item/components/appdragwidget.h"
#include "../item/components/appsnapshot.h"
#include "../item/components/floatingpreview.h"

#include "../plugins/tray/snitraywidget.h"
#include "../plugins/tray/indicatortraywidget.h"
#include "../plugins/tray/xembedtraywidget.h"
#include "../plugins/tray/system-trays/systemtrayitem.h"
#include "../plugins/tray/fashiontray/fashiontrayitem.h"
#include "../plugins/tray/fashiontray/fashiontraywidgetwrapper.h"
#include "../plugins/tray/fashiontray/fashiontraycontrolwidget.h"
#include "../plugins/tray/fashiontray/containers/attentioncontainer.h"
#include "../plugins/tray/fashiontray/containers/holdcontainer.h"
#include "../plugins/tray/fashiontray/containers/normalcontainer.h"
#include "../plugins/tray/fashiontray/containers/spliteranimated.h"

#include "../plugins/show-desktop/showdesktopwidget.h"

#include "../plugins/sound/sounditem.h"
#include "../plugins/sound/soundapplet.h"
#include "../plugins/sound/sinkinputwidget.h"
#include "../plugins/sound/componments/volumeslider.h"
#include "../plugins/sound/componments/horizontalseparator.h"

//#include "../plugins/network/item/deviceitem.h"// TODO

#include "../plugins/datetime/datetimewidget.h"
#include "../plugins/onboard/onboarditem.h"
#include "../plugins/trash/trashwidget.h"
#include "../plugins/trash/popupcontrolwidget.h"
#include "../plugins/shutdown/shutdownwidget.h"
#include "../plugins/multitasking/multitaskingwidget.h"
#include "../plugins/overlay-warning/overlaywarningwidget.h"

#include <DImageButton>
#include <DSwitchButton>
#include <DPushButton>

DWIDGET_USE_NAMESPACE

// 添加accessible
SET_WIDGET_ACCESSIBLE(MainWindow, QAccessible::Form, "mainwindow")
SET_BUTTON_ACCESSIBLE(MainPanelControl, "mainpanelcontrol")
SET_LABEL_ACCESSIBLE(TipsWidget, "tips")
SET_WIDGET_ACCESSIBLE(DockPopupWindow, QAccessible::Form, "popupwindow")
SET_BUTTON_ACCESSIBLE(LauncherItem, "launcheritem")
SET_BUTTON_ACCESSIBLE(AppItem, "appitem")
SET_BUTTON_ACCESSIBLE(PreviewContainer, "previewcontainer")
SET_BUTTON_ACCESSIBLE(PluginsItem, m_w->pluginName())
SET_BUTTON_ACCESSIBLE(TrayPluginItem, m_w->pluginName())
SET_BUTTON_ACCESSIBLE(PlaceholderItem, "placeholderitem")
SET_BUTTON_ACCESSIBLE(AppDragWidget, "appdragwidget")
SET_BUTTON_ACCESSIBLE(AppSnapshot, "appsnapshot")
SET_BUTTON_ACCESSIBLE(FloatingPreview, "floatingpreview")
SET_BUTTON_ACCESSIBLE(XEmbedTrayWidget, m_w->itemKeyForConfig())
SET_BUTTON_ACCESSIBLE(IndicatorTrayWidget, m_w->itemKeyForConfig())
SET_BUTTON_ACCESSIBLE(SNITrayWidget, m_w->itemKeyForConfig())
SET_BUTTON_ACCESSIBLE(SystemTrayItem, m_w->itemKeyForConfig())
SET_WIDGET_ACCESSIBLE(FashionTrayItem, QAccessible::Form, "fashiontrayitem")
SET_WIDGET_ACCESSIBLE(FashionTrayWidgetWrapper, QAccessible::Form, "fashiontraywrapper")
SET_BUTTON_ACCESSIBLE(FashionTrayControlWidget, "fashiontraycontrolwidget")
SET_WIDGET_ACCESSIBLE(AttentionContainer, QAccessible::Form, "attentioncontainer")
SET_WIDGET_ACCESSIBLE(HoldContainer, QAccessible::Form, "holdcontainer")
SET_WIDGET_ACCESSIBLE(NormalContainer, QAccessible::Form, "normalcontainer")
SET_WIDGET_ACCESSIBLE(SpliterAnimated, QAccessible::Form, "spliteranimated")
SET_BUTTON_ACCESSIBLE(SoundItem, "plugin-sounditem")
SET_WIDGET_ACCESSIBLE(SoundApplet, QAccessible::Form, "soundapplet")
SET_WIDGET_ACCESSIBLE(SinkInputWidget, QAccessible::Form, "sinkinputwidget")
SET_WIDGET_ACCESSIBLE(VolumeSlider, QAccessible::Slider, "volumeslider")
SET_WIDGET_ACCESSIBLE(HorizontalSeparator, QAccessible::Separator, "horizontalseparator")
SET_WIDGET_ACCESSIBLE(DatetimeWidget, QAccessible::Form, "plugin-datetime")
SET_WIDGET_ACCESSIBLE(OnboardItem, QAccessible::Form, "plugin-onboard")
SET_WIDGET_ACCESSIBLE(TrashWidget, QAccessible::Form, "plugin-trash")
SET_BUTTON_ACCESSIBLE(PopupControlWidget, "popupcontrolwidget")
SET_WIDGET_ACCESSIBLE(ShutdownWidget, QAccessible::Form, "plugin-shutdown")
SET_WIDGET_ACCESSIBLE(MultitaskingWidget, QAccessible::Form, "plugin-multitasking")
SET_WIDGET_ACCESSIBLE(ShowDesktopWidget, QAccessible::Form, "plugin-showdesktop")
SET_WIDGET_ACCESSIBLE(OverlayWarningWidget, QAccessible::Form, "plugin-overlaywarningwidget")
SET_WIDGET_ACCESSIBLE(QWidget, QAccessible::Form, m_w->objectName().isEmpty() ? "widget" : m_w->objectName())
SET_LABEL_ACCESSIBLE(QLabel, m_w->text().isEmpty() ? m_w->objectName().isEmpty() ? "text" : m_w->objectName() : m_w->text())
SET_BUTTON_ACCESSIBLE(DImageButton, m_w->objectName().isEmpty() ? "imagebutton" : m_w->objectName())
SET_BUTTON_ACCESSIBLE(DSwitchButton, m_w->text().isEmpty() ? "switchbutton" : m_w->text())

QAccessibleInterface *accessibleFactory(const QString &classname, QObject *object)
{
    QAccessibleInterface *interface = nullptr;

    USE_ACCESSIBLE(classname, MainWindow);
    USE_ACCESSIBLE(classname, MainPanelControl);
    USE_ACCESSIBLE(classname, TipsWidget);
    USE_ACCESSIBLE(classname, DockPopupWindow);
    USE_ACCESSIBLE(classname, LauncherItem);
    USE_ACCESSIBLE(classname, AppItem);
    USE_ACCESSIBLE(classname, PreviewContainer);
    USE_ACCESSIBLE(classname, PluginsItem);
    USE_ACCESSIBLE(classname, TrayPluginItem);
    USE_ACCESSIBLE(classname, PlaceholderItem);
    USE_ACCESSIBLE(classname, AppDragWidget);
    USE_ACCESSIBLE(classname, AppSnapshot);
    USE_ACCESSIBLE(classname, FloatingPreview);
    USE_ACCESSIBLE(classname, SNITrayWidget);
    USE_ACCESSIBLE(classname, SystemTrayItem);
    USE_ACCESSIBLE(classname, FashionTrayItem);
    USE_ACCESSIBLE(classname, FashionTrayWidgetWrapper);
    USE_ACCESSIBLE(classname, FashionTrayControlWidget);
    USE_ACCESSIBLE(classname, AttentionContainer);
    USE_ACCESSIBLE(classname, HoldContainer);
    USE_ACCESSIBLE(classname, NormalContainer);
    USE_ACCESSIBLE(classname, SpliterAnimated);
    USE_ACCESSIBLE(classname, IndicatorTrayWidget);
    USE_ACCESSIBLE(classname, XEmbedTrayWidget);
    USE_ACCESSIBLE(classname, SoundItem);
    USE_ACCESSIBLE(classname, SoundApplet);
    USE_ACCESSIBLE(classname, SinkInputWidget);
    USE_ACCESSIBLE(classname, VolumeSlider);
    USE_ACCESSIBLE(classname, HorizontalSeparator);
    USE_ACCESSIBLE(classname, DatetimeWidget);
    USE_ACCESSIBLE(classname, OnboardItem);
    USE_ACCESSIBLE(classname, TrashWidget);
    USE_ACCESSIBLE(classname, PopupControlWidget);
    USE_ACCESSIBLE(classname, ShutdownWidget);
    USE_ACCESSIBLE(classname, MultitaskingWidget);
    USE_ACCESSIBLE(classname, ShowDesktopWidget);
    USE_ACCESSIBLE(classname, OverlayWarningWidget);
    USE_ACCESSIBLE(classname, QWidget);
    USE_ACCESSIBLE_BY_OBJECTNAME(classname, QLabel, "spliter_fix");
    USE_ACCESSIBLE_BY_OBJECTNAME(classname, QLabel, "spliter_app");
    USE_ACCESSIBLE_BY_OBJECTNAME(classname, QLabel, "spliter_tray");
    USE_ACCESSIBLE(classname, QLabel);
    USE_ACCESSIBLE_BY_OBJECTNAME(QString(classname).replace("Dtk::Widget::", ""), DImageButton, "closebutton-2d");
    USE_ACCESSIBLE_BY_OBJECTNAME(QString(classname).replace("Dtk::Widget::", ""), DImageButton, "closebutton-3d");
    USE_ACCESSIBLE_BY_OBJECTNAME(QString(classname).replace("Dtk::Widget::", ""), DSwitchButton, "");

    return interface;
}
