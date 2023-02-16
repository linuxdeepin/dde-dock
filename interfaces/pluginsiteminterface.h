// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef PLUGINSITEMINTERFACE_H
#define PLUGINSITEMINTERFACE_H

#include "pluginproxyinterface.h"

#include <DGuiApplicationHelper>

#include <QIcon>
#include <QtCore>

DGUI_USE_NAMESPACE

// 任务栏的部件位置
enum class DockPart {
    QuickShow = 0,    // 快捷插件显示区域
    QuickPanel,       // 快捷面板区域
    SystemPanel,      // 系统插件显示区域
    DCCSetting        // 显示在控制中心个性化设置的图标
};

enum PluginFlag {
    Type_NoneFlag = 0x1,                 // 插件类型-没有任何的属性，不在任何地方显示
    Type_Common = 0x2,                   // 插件类型-快捷插件区
    Type_Tool = 0x4,                     // 插件类型-工具插件，例如回收站
    Type_System = 0x8,                   // 插件类型-系统插件，例如关机插件
    Type_Tray = 0x10,                    // 插件类型-托盘区，例如U盘插件
    Type_Fixed = 0x20,                   // 插件类型-固定区域,例如多任务视图和显示桌面

    Quick_Single = 0x40,                 // 当插件类型为Common时,快捷插件区域只有一列的那种插件
    Quick_Multi = 0x80,                  // 当插件类型为Common时,快捷插件区占两列的那种插件
    Quick_Full = 0x100,                  // 当插件类型为Common时,快捷插件区占用4列的那种插件，例如声音、亮度设置和音乐等

    Attribute_CanDrag = 0x200,           // 插件属性-是否支持拖动
    Attribute_CanInsert = 0x400,         // 插件属性-是否支持在其前面插入其他的插件，普通的快捷插件是支持的
    Attribute_CanSetting = 0x800,        // 插件属性-是否可以在控制中心设置显示或隐藏
    Attribute_ForceDock = 0x1000,        // 插件属性-强制显示在任务栏上

    FlagMask = 0xffffffff                // 掩码
};

Q_DECLARE_FLAGS(PluginFlags, PluginFlag)
Q_DECLARE_OPERATORS_FOR_FLAGS(PluginFlags)

// 快捷面板详情页面的itemWidget对应的itemKey
#define QUICK_ITEM_KEY "quick_item_key"
///
/// \brief The PluginsItemInterface class
/// the dock plugins item interface, all dock plugins should
/// inheirt this class and override all pure virtual function.
///

class PluginsItemInterface
{
public:
    enum PluginType {
        Normal,
        Fixed
    };

    /**
    * @brief Plugin size policy
    */
    enum PluginSizePolicy {
        System = 1 << 0, // Follow the system
        Custom = 1 << 1  // The custom
    };

    enum PluginMode {
        Deactive = 0,
        Active,
        Disabled
    };

    ///
    /// \brief ~PluginsItemInterface
    /// DON'T try to delete m_proxyInter.
    ///
    virtual ~PluginsItemInterface() {}

    ///
    /// \brief pluginName
    /// tell dock the unique plugin id
    /// \return
    ///
    virtual const QString pluginName() const = 0;
    virtual const QString pluginDisplayName() const { return QString(); }

    ///
    /// \brief init
    /// init your plugins, you need to save proxyInter to m_proxyInter
    /// member variable. but you shouldn't free this pointer.
    /// \param proxyInter
    /// DON'T try to delete this pointer.
    ///
    virtual void init(PluginProxyInterface *proxyInter) = 0;
    ///
    /// \brief itemWidget
    /// your plugin item widget, each item should have a unique key.
    /// \param itemKey
    /// your widget' unqiue key.
    /// \return
    ///
    virtual QWidget *itemWidget(const QString &itemKey) = 0;

    ///
    /// \brief itemTipsWidget
    /// override this function if your item want to have a tips.
    /// the tips will shown when user hover your item.
    /// nullptr will be ignored.
    /// \param itemKey
    /// \return
    ///
    virtual QWidget *itemTipsWidget(const QString &itemKey) {Q_UNUSED(itemKey); return nullptr;}
    ///
    /// \brief itemPopupApplet
    /// override this function if your item wants to have an popup applet.
    /// the popup applet will shown when user click your item.
    ///
    /// Tips:
    /// dock should receive mouse press/release event to check user mouse operate,
    /// if your item filter mouse event, this function will not be called.
    /// so if you override mouse event and want to use popup applet, you
    /// should pass event to your parent use QWidget::someEvent(e);
    /// \param itemKey
    /// \return
    ///
    virtual QWidget *itemPopupApplet(const QString &itemKey) {Q_UNUSED(itemKey); return nullptr;}
    ///
    /// \brief itemCommand
    /// execute spec command when user clicked your item.
    /// ensure your command do not get user input.
    ///
    /// empty string will be ignored.
    /// \param itemKey
    /// \return
    ///
    virtual const QString itemCommand(const QString &itemKey) {Q_UNUSED(itemKey); return QString();}

    ///
    /// \brief itemContextMenu
    /// context menu is shown when RequestPopupMenu called.
    /// \param itemKey
    /// \return
    ///
    virtual const QString itemContextMenu(const QString &itemKey) {Q_UNUSED(itemKey); return QString();}
    ///
    /// \brief invokedMenuItem
    /// call if context menu item is clicked
    /// \param itemKey
    /// \param itemId
    /// menu item id
    /// \param checked
    ///
    virtual void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked) {Q_UNUSED(itemKey); Q_UNUSED(menuId); Q_UNUSED(checked);}

    ///
    /// \brief itemSortKey
    /// tell dock where your item wants to put on.
    ///
    /// this index is start from 1 and
    /// 0 for left side
    /// -1 for right side
    /// \param itemKey
    /// \return
    ///
    virtual int itemSortKey(const QString &itemKey) {Q_UNUSED(itemKey); return 1;}
    ///
    /// \brief setSortKey
    /// save your item new position
    /// sort key will be changed when plugins order
    /// changed(by user drag-drop)
    /// \param itemKey
    /// \param order
    ///
    virtual void setSortKey(const QString &itemKey, const int order) {Q_UNUSED(itemKey); Q_UNUSED(order);}

    ///
    /// \brief itemAllowContainer
    /// tell dock is your item allow to move into container
    ///
    /// if your item placed into container, popup tips and popup
    /// applet will be disabled.
    /// \param itemKey
    /// \return
    ///
    virtual bool itemAllowContainer(const QString &itemKey) {Q_UNUSED(itemKey); return false;}
    ///
    /// \brief itemIsInContainer
    /// tell dock your item is in container, this function
    /// called at item init and if your item enable container.
    /// \param itemKey
    /// \return
    ///
    virtual bool itemIsInContainer(const QString &itemKey) {Q_UNUSED(itemKey); return false;}
    ///
    /// \brief setItemIsInContainer
    /// save your item new state.
    /// this function called when user drag out your item from
    /// container or user drop item into container(if your item
    /// allow drop into container).
    /// \param itemKey
    /// \param container
    ///
    virtual void setItemIsInContainer(const QString &itemKey, const bool container) {Q_UNUSED(itemKey); Q_UNUSED(container);}

    virtual bool pluginIsAllowDisable() { return false; }
    virtual bool pluginIsDisable() { return false; }
    virtual void pluginStateSwitched() {}

    ///
    /// \brief displayModeChanged
    /// override this function to receive display mode changed signal
    /// \param displayMode
    ///
    virtual void displayModeChanged(const Dock::DisplayMode displayMode) {Q_UNUSED(displayMode);}
    ///
    /// \brief positionChanged
    /// override this function to receive dock position changed signal
    /// \param position
    ///
    virtual void positionChanged(const Dock::Position position) {Q_UNUSED(position);}

    ///
    /// \brief refreshIcon
    /// refresh item icon, its triggered when system icon theme changed.
    /// \param itemKey
    /// item key
    ///
    virtual void refreshIcon(const QString &itemKey) { Q_UNUSED(itemKey); }

    ///
    /// \brief displayMode
    /// get current dock display mode
    /// \return
    ///
    inline Dock::DisplayMode displayMode() const
    {
        return qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();
    }

    ///
    /// \brief position
    /// get current dock position
    /// \return
    ///
    inline Dock::Position position() const
    {
        return qApp->property(PROP_POSITION).value<Dock::Position>();
    }

    ///
    /// \brief settingsChanged
    /// override this function to receive plugin settings changed signal(DeepinSync)
    ///
    virtual void pluginSettingsChanged() {}

    ///
    /// \brief type
    /// default plugin add dock right,fixed plugin add to dock fixed area
    ///
    virtual PluginType type() { return Normal; }

    ///
    /// \brief plugin size policy
    /// default plugin size policy
    ///
    virtual PluginSizePolicy pluginSizePolicy() const { return System; }

    ///
    /// the plugin status
    ///
    ///
    virtual PluginMode status() const { return PluginMode::Deactive; }

    ///
    /// return the detail value, it will display in the center
    ///
    ///
    virtual QString description() const { return QString(); }

    ///
    /// the icon for the plugin
    /// themeType {0:UnknownType 1:LightType 2:DarkType}
    ///
    virtual QIcon icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType = DGuiApplicationHelper::instance()->themeType()) { return QIcon(); }

    ///
    /// \brief m_proxyInter
    /// return the falgs for current plugin
    ///
    virtual PluginFlags flags() const { return PluginFlag::Type_Common | PluginFlag::Quick_Single | PluginFlag::Attribute_CanDrag | PluginFlag::Attribute_CanInsert | PluginFlag::Attribute_CanSetting; }

    ///
    /// \brief m_proxyInter
    ///
    ///
    virtual bool eventHandler(QEvent *event) { return false; }

protected:
    ///
    /// \brief m_proxyInter
    /// NEVER delete this object.
    ///
    PluginProxyInterface *m_proxyInter = nullptr;
};

QT_BEGIN_NAMESPACE

#define ModuleInterface_iid "com.deepin.dock.PluginsItemInterface_2_0_0"

Q_DECLARE_INTERFACE(PluginsItemInterface, ModuleInterface_iid)
QT_END_NAMESPACE

#endif // PLUGINSITEMINTERFACE_H
