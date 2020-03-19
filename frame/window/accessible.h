/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     fpc_diesel <fanpengcheng@uniontech.com>
 *
 * Maintainer: fpc_diesel <fanpengcheng@uniontech.com>
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
#ifndef ACCESSIBLE_H
#define ACCESSIBLE_H

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
//#include "../plugins/overlay-warning/overlaywarningwidget.h"// TODO

#include <QAccessible>
#include <QAccessibleWidget>
#include <QEvent>
#include <QMouseEvent>
#include <QApplication>

#include <DImageButton>
#include <DSwitchButton>
#include <DPushButton>

DWIDGET_USE_NAMESPACE

/**************************************************************************************/
// 构造函数
#define FUNC_CREATE(classname,accessibletype,accessdescription)    Accessible##classname(classname *w) \
        : QAccessibleWidget(w,accessibletype,#classname)\
        , m_w(w)\
        , m_description(accessdescription)\
    {}\
    private:\
    classname *m_w;\
    QString m_description;\

// 左键点击
#define FUNC_PRESS(classobj)     QStringList actionNames() const override{\
        if(!classobj->isEnabled())\
            return QStringList();\
        return QStringList() << pressAction();}\
    void doAction(const QString &actionName) override{\
        if(actionName == pressAction())\
        {\
            QPointF localPos = classobj->geometry().center();\
            QMouseEvent event(QEvent::MouseButtonPress,localPos,Qt::LeftButton,Qt::LeftButton,Qt::NoModifier);\
            QMouseEvent event2(QEvent::MouseButtonRelease,localPos,Qt::LeftButton,Qt::LeftButton,Qt::NoModifier);\
            qApp->sendEvent(classobj,&event);\
            qApp->sendEvent(classobj,&event2);\
        }\
    }\

// 右键点击
#define FUNC_SHOWMENU(classobj)     QStringList actionNames() const override{\
        if(!classobj->isEnabled())\
            return QStringList();\
        return QStringList() << showMenuAction();}\
    void doAction(const QString &actionName) override{\
        if(actionName == showMenuAction())\
        {\
            QPointF localPos = classobj->geometry().center();\
            QMouseEvent event(QEvent::MouseButtonPress,localPos,Qt::RightButton,Qt::RightButton,Qt::NoModifier);\
            QMouseEvent event2(QEvent::MouseButtonRelease,localPos,Qt::RightButton,Qt::RightButton,Qt::NoModifier);\
            qApp->sendEvent(classobj,&event);\
            qApp->sendEvent(classobj,&event2);\
        }\
    }\

// 左键和右键点击
#define FUNC_PRESS_SHOWMENU(classobj)     QStringList actionNames() const override{\
        if(!classobj->isEnabled())\
            return QStringList();\
        return QStringList() << pressAction() << showMenuAction();}\
    void doAction(const QString &actionName) override{\
        if(actionName == pressAction())\
        {\
            QPointF localPos = classobj->geometry().center();\
            QMouseEvent event(QEvent::MouseButtonPress,localPos,Qt::LeftButton,Qt::LeftButton,Qt::NoModifier);\
            QMouseEvent event2(QEvent::MouseButtonRelease,localPos,Qt::LeftButton,Qt::LeftButton,Qt::NoModifier);\
            qApp->sendEvent(classobj,&event);\
            qApp->sendEvent(classobj,&event2);\
        }\
        else if(actionName == showMenuAction())\
        {\
            QPointF localPos = classobj->geometry().center();\
            QMouseEvent event(QEvent::MouseButtonPress,localPos,Qt::RightButton,Qt::RightButton,Qt::NoModifier);\
            QMouseEvent event2(QEvent::MouseButtonRelease,localPos,Qt::RightButton,Qt::RightButton,Qt::NoModifier);\
            qApp->sendEvent(classobj,&event);\
            qApp->sendEvent(classobj,&event2);\
        }\
    }\

// 实现rect接口
#define FUNC_RECT(classobj) QRect rect() const override{\
        if (!classobj->isVisible())\
            return QRect();\
        return classobj->geometry();\
    }\

// 启用accessible
#define GET_ACCESSIBLE(classnamestring,classname)    if (classnamestring == QLatin1String(#classname) && object && object->isWidgetType())\
    {\
        interface = new Accessible##classname(static_cast<classname *>(object));\
    }\

// 启用accessible[指定objectname]---适用同一个类，但objectname不同的情况
#define GET_ACCESSIBLE_BY_OBJECTNAME(classnamestring,classname,objectname)    if (classnamestring == QLatin1String(#classname) && object && (object->objectName() == objectname) && object->isWidgetType())\
    {\
        interface = new Accessible##classname(static_cast<classname *>(object));\
    }\

// 按钮类型的控件[仅有左键点击]
#define SET_BUTTON_ACCESSIBLE_PRESS_DESCRIPTION(classname,accessdescription)  class Accessible##classname : public QAccessibleWidget\
    {\
    public:\
        FUNC_CREATE(classname,QAccessible::Button,accessdescription)\
        QString text(QAccessible::Text t) const override;/*需要单独实现*/\
        FUNC_PRESS(m_w)\
    };\

// 按钮类型的控件[仅有右键点击]
#define SET_BUTTON_ACCESSIBLE_SHOWMENU_DESCRIPTION(classname,accessdescription)  class Accessible##classname : public QAccessibleWidget\
    {\
    public:\
        FUNC_CREATE(classname,QAccessible::Button,accessdescription)\
        QString text(QAccessible::Text t) const override;/*需要单独实现*/\
        FUNC_SHOWMENU(m_w)\
    };\

// 按钮类型的控件[有左键点击和右键点击]
#define SET_BUTTON_ACCESSIBLE_PRESS_SHOEMENU_DESCRIPTION(classname,accessdescription)  class Accessible##classname : public QAccessibleWidget\
    {\
    public:\
        FUNC_CREATE(classname,QAccessible::Button,accessdescription)\
        QString text(QAccessible::Text t) const override;/*需要单独实现*/\
        FUNC_PRESS_SHOWMENU(m_w)\
    };\

// 标签类型的控件
#define SET_LABEL_ACCESSIBLE_WITH_DESCRIPTION(classname,aaccessibletype,accessdescription)  class Accessible##classname : public QAccessibleWidget\
    {\
    public:\
        FUNC_CREATE(classname,aaccessibletype,accessdescription)\
        QString text(QAccessible::Text t) const override;/*需要单独实现*/\
        FUNC_RECT(m_w)\
    };\

// 简化使用
#define SET_BUTTON_ACCESSIBLE_PRESS_SHOWMENU(classname)         SET_BUTTON_ACCESSIBLE_PRESS_SHOEMENU_DESCRIPTION(classname,"")
#define SET_BUTTON_ACCESSIBLE_SHOWMENU(classname)               SET_BUTTON_ACCESSIBLE_SHOWMENU_DESCRIPTION(classname,"")
#define SET_BUTTON_ACCESSIBLE(classname)                        SET_BUTTON_ACCESSIBLE_PRESS_DESCRIPTION(classname,"")

#define SET_LABEL_ACCESSIBLE(classname)                         SET_LABEL_ACCESSIBLE_WITH_DESCRIPTION(classname,QAccessible::StaticText,"")
#define SET_FORM_ACCESSIBLE(classname)                          SET_LABEL_ACCESSIBLE_WITH_DESCRIPTION(classname,QAccessible::Form,"")
#define SET_SLIDER_ACCESSIBLE(classname)                        SET_LABEL_ACCESSIBLE_WITH_DESCRIPTION(classname,QAccessible::Slider,"")
#define SET_SEPARATOR_ACCESSIBLE(classname)                     SET_LABEL_ACCESSIBLE_WITH_DESCRIPTION(classname,QAccessible::Separator,"")
/**************************************************************************************/

// 添加accessible
SET_FORM_ACCESSIBLE(MainWindow)
SET_BUTTON_ACCESSIBLE_SHOWMENU(MainPanelControl)
SET_LABEL_ACCESSIBLE(TipsWidget)
SET_FORM_ACCESSIBLE(DockPopupWindow)

SET_BUTTON_ACCESSIBLE(LauncherItem)
SET_BUTTON_ACCESSIBLE_PRESS_SHOWMENU(AppItem)
SET_BUTTON_ACCESSIBLE(PreviewContainer)
SET_BUTTON_ACCESSIBLE_PRESS_SHOWMENU(PluginsItem)
SET_BUTTON_ACCESSIBLE(TrayPluginItem)
SET_BUTTON_ACCESSIBLE(PlaceholderItem)
SET_BUTTON_ACCESSIBLE(AppDragWidget)
SET_BUTTON_ACCESSIBLE(AppSnapshot)
SET_BUTTON_ACCESSIBLE(FloatingPreview)

// tray plugin
SET_BUTTON_ACCESSIBLE_PRESS_SHOWMENU(XEmbedTrayWidget)
SET_BUTTON_ACCESSIBLE_PRESS_SHOWMENU(IndicatorTrayWidget)
SET_BUTTON_ACCESSIBLE_PRESS_SHOWMENU(SNITrayWidget)
SET_BUTTON_ACCESSIBLE_PRESS_SHOWMENU(SystemTrayItem)
SET_FORM_ACCESSIBLE(FashionTrayItem)
SET_FORM_ACCESSIBLE(FashionTrayWidgetWrapper)
SET_BUTTON_ACCESSIBLE(FashionTrayControlWidget)
SET_FORM_ACCESSIBLE(AttentionContainer)
SET_FORM_ACCESSIBLE(HoldContainer)
SET_FORM_ACCESSIBLE(NormalContainer)
SET_FORM_ACCESSIBLE(SpliterAnimated)

// sound plugin
SET_BUTTON_ACCESSIBLE(SoundItem)
SET_FORM_ACCESSIBLE(SoundApplet)
SET_FORM_ACCESSIBLE(SinkInputWidget)
SET_SLIDER_ACCESSIBLE(VolumeSlider)
SET_SEPARATOR_ACCESSIBLE(HorizontalSeparator)

// fixed plugin
SET_FORM_ACCESSIBLE(DatetimeWidget)
SET_FORM_ACCESSIBLE(OnboardItem)
SET_FORM_ACCESSIBLE(TrashWidget)
SET_BUTTON_ACCESSIBLE(PopupControlWidget)
SET_FORM_ACCESSIBLE(ShutdownWidget)
SET_FORM_ACCESSIBLE(MultitaskingWidget)
SET_FORM_ACCESSIBLE(ShowDesktopWidget)

// special class for other useage
SET_BUTTON_ACCESSIBLE(QWidget)
SET_BUTTON_ACCESSIBLE(DImageButton)
SET_BUTTON_ACCESSIBLE(DSwitchButton)

#endif // ACCESSIBLE_H
