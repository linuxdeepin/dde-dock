#include "accessibledefine.h"

#include "mainwindow.h"
#include "mainpanelcontrol.h"
#include "tipswidget.h"
#include "dockpopupwindow.h"
#include "statebutton.h"

#include "launcheritem.h"
#include "appitem.h"
#include "components/previewcontainer.h"
#include "pluginsitem.h"
#include "traypluginitem.h"
#include "placeholderitem.h"
#include "components/appdragwidget.h"
#include "components/appsnapshot.h"
#include "components/floatingpreview.h"

#include "snitraywidget.h"
#include "abstracttraywidget.h"
#include "indicatortraywidget.h"
#include "xembedtraywidget.h"
#include "system-trays/systemtrayitem.h"
#include "fashiontray/fashiontrayitem.h"
#include "fashiontray/fashiontraywidgetwrapper.h"
#include "fashiontray/fashiontraycontrolwidget.h"
#include "fashiontray/containers/attentioncontainer.h"
#include "fashiontray/containers/holdcontainer.h"
#include "fashiontray/containers/normalcontainer.h"
#include "fashiontray/containers/spliteranimated.h"

// 这部分由sound插件单独维护,这样做是因为在标记volumeslider这个类时,需要用到其setValue的实现,
// 但插件的源文件dock这边并没有包含,不想引入复杂的包含关系,其实最好的做法就是像sound插件这样,谁维护谁的
//#include "../plugins/sound/sounditem.h"
//#include "../plugins/sound/soundapplet.h"
//#include "../plugins/sound/sinkinputwidget.h"
//#include "../plugins/sound/componments/volumeslider.h"
//#include "../plugins/sound/componments/horizontalseparator.h"

#include "showdesktopwidget.h"
#include "networkitem.h"
#include "item/applet/devicecontrolwidget.h"
#include "datetimewidget.h"
#include "onboarditem.h"
#include "trashwidget.h"
#include "popupcontrolwidget.h"
#include "shutdownwidget.h"
#include "multitaskingwidget.h"
#include "overlaywarningwidget.h"

#include <DIconButton>
#include <DSwitchButton>
#include <DPushButton>
#include <DListView>
#include <DSwitchButton>
#include <DSpinner>
#include <dloadingindicator.h>

#include <QScrollBar>

DWIDGET_USE_NAMESPACE
using namespace Dock;

// 添加accessible
SET_FORM_ACCESSIBLE(MainWindow, "mainwindow")
SET_BUTTON_ACCESSIBLE(MainPanelControl, "mainpanelcontrol")
SET_LABEL_ACCESSIBLE(TipsWidget, "tips")
SET_FORM_ACCESSIBLE(DockPopupWindow, "popupwindow")
SET_BUTTON_ACCESSIBLE(LauncherItem, "launcheritem")
SET_BUTTON_ACCESSIBLE(AppItem, "appitem")
SET_BUTTON_ACCESSIBLE(PreviewContainer, "previewcontainer")
SET_BUTTON_ACCESSIBLE(PluginsItem, m_w->pluginName())
SET_BUTTON_ACCESSIBLE(TrayPluginItem, m_w->pluginName())
SET_BUTTON_ACCESSIBLE(PlaceholderItem, "placeholderitem")
SET_BUTTON_ACCESSIBLE(AppDragWidget, "appdragwidget")
SET_BUTTON_ACCESSIBLE(AppSnapshot, "appsnapshot")
SET_BUTTON_ACCESSIBLE(FloatingPreview, "floatingpreview")
SET_BUTTON_ACCESSIBLE(XEmbedTrayWidget, m_w->itemKeyForConfig().replace("sni:", ""))
SET_BUTTON_ACCESSIBLE(IndicatorTrayWidget, m_w->itemKeyForConfig().replace("sni:", ""))
SET_BUTTON_ACCESSIBLE(SNITrayWidget, m_w->itemKeyForConfig().replace("sni:", ""))
SET_BUTTON_ACCESSIBLE(AbstractTrayWidget, m_w->itemKeyForConfig().replace("sni:", ""))
SET_BUTTON_ACCESSIBLE(SystemTrayItem, m_w->itemKeyForConfig().replace("sni:", ""))
SET_FORM_ACCESSIBLE(FashionTrayItem, "fashiontrayitem")
SET_FORM_ACCESSIBLE(FashionTrayWidgetWrapper, "fashiontraywrapper")
SET_BUTTON_ACCESSIBLE(FashionTrayControlWidget, "fashiontraycontrolwidget")
SET_FORM_ACCESSIBLE(AttentionContainer, "attentioncontainer")
SET_FORM_ACCESSIBLE(HoldContainer, "holdcontainer")
SET_FORM_ACCESSIBLE(NormalContainer, "normalcontainer")
SET_FORM_ACCESSIBLE(SpliterAnimated, "spliteranimated")
SET_FORM_ACCESSIBLE(DatetimeWidget, "plugin-datetime")
SET_FORM_ACCESSIBLE(OnboardItem, "plugin-onboard")
SET_FORM_ACCESSIBLE(TrashWidget, "plugin-trash")
SET_BUTTON_ACCESSIBLE(PopupControlWidget, "popupcontrolwidget")
SET_FORM_ACCESSIBLE(ShutdownWidget, "plugin-shutdown")
SET_FORM_ACCESSIBLE(MultitaskingWidget, "plugin-multitasking")
SET_FORM_ACCESSIBLE(ShowDesktopWidget, "plugin-showdesktop")
SET_FORM_ACCESSIBLE(OverlayWarningWidget, "plugin-overlaywarningwidget")
SET_FORM_ACCESSIBLE(QWidget, m_w->objectName().isEmpty() ? "widget" : m_w->objectName())
SET_LABEL_ACCESSIBLE(QLabel, m_w->objectName() == "notifications" ? m_w->objectName() : m_w->text().isEmpty() ? m_w->objectName().isEmpty() ? "text" : m_w->objectName() : m_w->text())
SET_BUTTON_ACCESSIBLE(DIconButton, m_w->objectName().isEmpty() ? "imagebutton" : m_w->objectName())
SET_BUTTON_ACCESSIBLE(DSwitchButton, m_w->text().isEmpty() ? "switchbutton" : m_w->text())
SET_BUTTON_ACCESSIBLE(DesktopWidget, "desktopWidget");
// 几个没什么用的标记，但为了提醒大家不要遗漏标记控件，还是不要去掉
SET_FORM_ACCESSIBLE(DBlurEffectWidget, "DBlurEffectWidget")
SET_FORM_ACCESSIBLE(DListView, "DListView")
SET_FORM_ACCESSIBLE(DLoadingIndicator, "DLoadingIndicator")
SET_FORM_ACCESSIBLE(DSpinner, "DSpinner")
SET_FORM_ACCESSIBLE(QMenu, "QMenu")
SET_FORM_ACCESSIBLE(QPushButton, "QPushButton")
SET_FORM_ACCESSIBLE(QSlider, "QSlider")
SET_FORM_ACCESSIBLE(QScrollBar, "QScrollBar")
SET_FORM_ACCESSIBLE(QScrollArea, "QScrollArea")
SET_FORM_ACCESSIBLE(QFrame, "QFrame")
SET_FORM_ACCESSIBLE(QGraphicsView, "QGraphicsView")
SET_FORM_ACCESSIBLE(DragWidget, "DragWidget")
SET_FORM_ACCESSIBLE(NetworkItem, "NetworkItem")
SET_FORM_ACCESSIBLE(StateButton, "StateButton")
SET_FORM_ACCESSIBLE(DeviceControlWidget, "DeviceControlWidget")

QAccessibleInterface *accessibleFactory(const QString &classname, QObject *object)
{
    // 自动化标记确定不需要的控件，方可加入忽略列表
    const static QStringList ignoreLst = {"WirelessItem", "WiredItem", "SsidButton", "WirelessList", "AccessPointWidget"};

    QAccessibleInterface *interface = nullptr;

    USE_ACCESSIBLE(classname, MainWindow);
    USE_ACCESSIBLE(classname, MainPanelControl);
    USE_ACCESSIBLE(QString(classname).replace("Dock::", ""), TipsWidget);
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
    USE_ACCESSIBLE(classname, AbstractTrayWidget);
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
    USE_ACCESSIBLE(classname, DesktopWidget);
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
    USE_ACCESSIBLE_BY_OBJECTNAME(QString(classname).replace("Dtk::Widget::", ""), DIconButton, "closebutton-2d");
    USE_ACCESSIBLE_BY_OBJECTNAME(QString(classname).replace("Dtk::Widget::", ""), DIconButton, "closebutton-3d");
    USE_ACCESSIBLE_BY_OBJECTNAME(QString(classname).replace("Dtk::Widget::", ""), DSwitchButton, "");
    USE_ACCESSIBLE(QString(classname).replace("Dtk::Widget::", ""), DBlurEffectWidget);
    USE_ACCESSIBLE(QString(classname).replace("Dtk::Widget::", ""), DListView);
    USE_ACCESSIBLE(QString(classname).replace("Dtk::Widget::", ""), DLoadingIndicator);
    USE_ACCESSIBLE(QString(classname).replace("Dtk::Widget::", ""), DSpinner);
    USE_ACCESSIBLE(QString(classname).replace("Dtk::Widget::", ""), DSwitchButton);
    USE_ACCESSIBLE(QString(classname).replace("Dtk::Widget::", ""), DIconButton);
    USE_ACCESSIBLE(classname, QMenu);
    USE_ACCESSIBLE(classname, QPushButton);
    USE_ACCESSIBLE(classname, QSlider);
    USE_ACCESSIBLE(classname, QScrollBar);
    USE_ACCESSIBLE(classname, QScrollArea);
    USE_ACCESSIBLE(classname, QFrame);
    USE_ACCESSIBLE(classname, QGraphicsView);
    USE_ACCESSIBLE(classname, DragWidget);
    USE_ACCESSIBLE(classname, NetworkItem);
    USE_ACCESSIBLE(classname, StateButton);
    USE_ACCESSIBLE(classname, DeviceControlWidget);

    if (!interface && object->inherits("QWidget") && !ignoreLst.contains(classname)) {
        QWidget *w = static_cast<QWidget *>(object);
        // 如果你看到这里的输出，说明代码中仍有控件未兼顾到accessible功能，请帮忙添加
        if (w->accessibleName().isEmpty())
            qWarning() << "accessibleFactory()" + QString("Class: " + classname + " cannot access");
    }

    return interface;
}
