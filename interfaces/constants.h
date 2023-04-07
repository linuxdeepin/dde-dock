// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef CONSTANTS_H
#define CONSTANTS_H

#include <QtCore>

namespace Dock {

#define DOCK_PLUGIN_MIME    "dock/plugin"
#define DOCK_PLUGIN_API_VERSION    "2.0.0"

#define PROP_DISPLAY_MODE   "DisplayMode"
#define PROP_DOCK_DRAGING   "isDraging"

#define PLUGIN_BACKGROUND_MAX_SIZE 40
#define PLUGIN_BACKGROUND_MIN_SIZE 20

#define PLUGIN_ICON_MAX_SIZE 20
#define PLUGIN_ITEM_WIDTH 300

#define QUICK_PATH "/plugins/quick-trays"
#define PLUGIN_PATH "/plugins"
#define TRAY_PATH "/plugins/system-trays"

// 需求变更成插件图标始终保持20x20,但16x16的资源还在。所以暂时保留此宏
#define PLUGIN_ICON_MIN_SIZE 20

// 插件最小尺寸，图标采用深色
#define PLUGIN_MIN_ICON_NAME "-dark"

// dock最小尺寸
#define DOCK_MIN_SIZE 40
// dock最大尺寸
#define DOCK_MAX_SIZE 100
///
/// \brief The DisplayMode enum
/// spec dock display mode
///
enum DisplayMode {
    Fashion     = 0,
    Efficient   = 1,
    // deprecreated
//    Classic     = 2,
};

#define PROP_HIDE_MODE      "HideMode"
///
/// \brief The HideMode enum
/// spec dock hide behavior
///
enum HideMode {
    KeepShowing     = 0,
    KeepHidden      = 1,
    SmartHide       = 2
};

#define PROP_POSITION       "Position"
///
/// \brief The Position enum
/// spec dock position, dock always placed at primary screen,
/// so all position is the primary screen edge.
///
enum Position {
    Top         = 0,
    Right       = 1,
    Bottom      = 2,
    Left        = 3,
};

#define PROP_HIDE_STATE     "HideState"
///
/// \brief The HideState enum
/// spec current dock should hide or shown.
/// this argument works only HideMode is SmartHide
///
enum HideState {
    Unknown     = 0,
    Show        = 1,
    Hide        = 2,
};

enum class AniAction {
    Show = 0,
    Hide
};

#define IS_TOUCH_STATE "isTouchState"
#define POPUP_PADDING 10

}

Q_DECLARE_METATYPE(Dock::DisplayMode)
Q_DECLARE_METATYPE(Dock::Position)

#endif // CONSTANTS_H
